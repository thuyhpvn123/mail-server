#include <iostream>
#include <tuple>
#include "mvm_linker.hpp"
#include <string>
#include "mvm/opcode.h"
#include "mvm/processor.h"
#include "my_global_state.h"
#include "my_logger.h"
#include "my_extension.h"

#include <stdlib.h>
#include <cassert>
#include <fmt/format_header_only.h>
#include <fstream>
#include <iostream>
#include <nlohmann/json.hpp>
#include <random>
#include <sstream>
#include <vector>
#include <fstream>

ExecuteResult* pendingResult;

nlohmann::json vectorLogsHandlerToJson(mvm::VectorLogHandler logHandler)
{
  auto json_logs = nlohmann::json::array();
  for (const auto &log : logHandler.logs)
  {
    nlohmann::json json_log;
    mvm::to_json(json_log, log);
    json_logs.push_back(json_log);
  }

  return json_logs;
}

void append_argument(std::vector<uint8_t> &code, const uint256_t &arg)
{
  // To ABI encode a function call with a uint256_t (or Address) argument,
  // simply append the big-endian byte representation to the code (function
  // selector, or bin). ABI-encoding for more complicated types is more
  // complicated, so not shown in this sample.
  const auto pre_size = code.size();
  code.resize(pre_size + 32u);
  mvm::to_big_endian(arg, code.data() + pre_size);
}

// Run input as an EVM transaction, check the result and return the output
mvm::ExecResult run(
  mvm::MyGlobalState &gs,
  bool deploy,
  const mvm::Address &from,
  const mvm::Address &to,
  const uint256_t &amount,
  uint64_t gas_price,
  uint64_t gas_limit,
  mvm::VectorLogHandler &log_handler,
  const mvm::Code &input)
{
  mvm::Transaction tx(
    from,
    amount, 
    gas_price, 
    gas_limit
  );
  MyLogger logger = MyLogger();
  MyExtension extension = MyExtension();
  mvm::Processor p(gs, log_handler, extension, logger);
  // Run the transaction
  const auto exec_result = p.run(tx, deploy, from, gs.get(to), input, amount);
  return exec_result;
}

mvm::BlockContext CreateBlockContext(
  uint64_t prevrandao,
  uint64_t gas_limit, 
  uint64_t time,      
  uint64_t base_fee,
  uint256_t number,   
  uint256_t coinbase
)
{
  mvm::BlockContext block_context;
  block_context.prevrandao = prevrandao;
  block_context.gas_limit = gas_limit;
  block_context.time = time;
  block_context.base_fee = base_fee;
  block_context.number = number;
  block_context.coinbase = coinbase;
  return block_context;
}

ExecuteResult* processResult (mvm::ExecResult result, mvm::MyGlobalState &gs, mvm::VectorLogHandler &log_handler) {
  // storage
  char *b_output = new char[result.output.size()];
  int length_output = static_cast<int>(result.output.size());
  
  std::vector<std::vector<uint8_t>> add_balance_change = gs.get_add_balance_change();
  int length_add_balance_change = add_balance_change.size();
  uint8_t **b_add_balance_change = new uint8_t *[length_add_balance_change];
  for (int i = 0; i < length_add_balance_change; ++i) {
     b_add_balance_change[i] = new uint8_t[add_balance_change[i].size()];
     memcpy(b_add_balance_change[i], add_balance_change[i].data(), add_balance_change[i].size());
  }
  
  std::vector<std::vector<uint8_t>> sub_balance_change = gs.get_sub_balance_change();
  int length_sub_balance_change = sub_balance_change.size();
  uint8_t **b_sub_balance_change = new uint8_t *[length_sub_balance_change];
  for (int i = 0; i < length_sub_balance_change; ++i) {
     b_sub_balance_change[i] = new uint8_t[sub_balance_change[i].size()];
     memcpy(b_sub_balance_change[i], sub_balance_change[i].data(), sub_balance_change[i].size());
  }
  
  std::vector<std::vector<uint8_t>> code_change = gs.get_newly_deploy();
  int length_code_change = code_change.size();
  int *length_codes = new int[length_code_change];
  uint8_t **b_code_change = new uint8_t *[length_code_change];
  for (size_t i = 0; i < length_code_change; ++i) {
     length_codes[i] = code_change[i].size();
     b_code_change[i] = new uint8_t[code_change[i].size()];
     memcpy(b_code_change[i], code_change[i].data(), code_change[i].size());
  }
  
  std::vector<std::vector<uint8_t>> storage_change = gs.get_storage_change();
  int length_storage_change = storage_change.size();
  int *length_storages = new int[length_storage_change];
  uint8_t **b_storage_change = new uint8_t *[length_storage_change];
  for (int i = 0; i < length_storage_change; ++i) {
     length_storages[i] = storage_change[i].size();
     b_storage_change[i] = new uint8_t[storage_change[i].size()];
     memcpy(b_storage_change[i], storage_change[i].data(), storage_change[i].size());
  }
  
  // logs
  auto json_logs = vectorLogsHandlerToJson(log_handler);
  std::string str_logs = json_logs.dump();

  pendingResult = new ExecuteResult{
    b_exitReason : (char)result.er,
    b_exception : (char)result.ex,
    b_exmsg : new char[result.exmsg.size()],
    length_exmsg : (int)result.exmsg.size(),

    b_output : b_output,
    length_output : length_output,
    b_add_balance_change : (char **)b_add_balance_change,
    length_add_balance_change : length_add_balance_change,

    b_sub_balance_change : (char **)b_sub_balance_change,
    length_sub_balance_change : length_sub_balance_change,

    b_code_change : (char **)b_code_change,
    length_code_change : length_code_change,
    length_codes : length_codes,
    
    b_storage_change : (char **)b_storage_change,
    length_storage_change : length_storage_change,
    length_storages : length_storages,

    b_logs : new char[str_logs.size()],
    length_logs : (int)str_logs.size(),

    gas_used: result.gas_used
  };

  memcpy(pendingResult->b_exmsg, (char *)result.exmsg.c_str(), result.exmsg.size());
  memcpy(pendingResult->b_output, result.output.data(), result.output.size());
  memcpy(pendingResult->b_logs, str_logs.c_str(), str_logs.size());

  return pendingResult;
}

ExecuteResult* deploy(
  // transaction data
  unsigned char* b_caller_address,
  unsigned char* b_caller_last_hash,
  unsigned char* b_contract_constructor,
  int contract_constructor_length,
  unsigned char* b_amount,
  unsigned long long gas_price,
  unsigned long long gas_limit,
  // block context data
  unsigned long long block_prevrandao,
  unsigned long long block_gas_limit,
  unsigned long long block_time,
  unsigned long long block_base_fee,
  unsigned char* b_block_number,
  unsigned char* b_block_coinbase
)
{

  // format argument to right data type
  uint256_t caller_address = mvm::from_big_endian((uint8_t *)b_caller_address, 20u);
  uint256_t caller_last_hash = mvm::from_big_endian((uint8_t *)b_caller_last_hash, 32u);
  std::vector<uint8_t> contract_constructor((uint8_t *)b_contract_constructor, (uint8_t *)b_contract_constructor + contract_constructor_length);
  uint256_t amount = mvm::from_big_endian((uint8_t *)b_amount, 32u);

  uint256_t block_number = mvm::from_big_endian((uint8_t *)b_block_number, 32u);
  uint256_t block_coinbase = mvm::from_big_endian((uint8_t *)b_block_coinbase, 20u);
  

  mvm::BlockContext blockContext = CreateBlockContext(
    block_prevrandao,
    block_gas_limit,
    block_time, 
    block_base_fee,
    block_number,   
    block_coinbase
  );
  mvm::MyGlobalState gs(blockContext);
  //  init env
  mvm::VectorLogHandler log_handler;
  const auto contract_address = mvm::generate_contract_address(caller_address, caller_last_hash);
  // Set this constructor as the contract's code body
  auto contract = gs.create(contract_address, 0u, contract_constructor);

  auto result = run(
    gs,
    true,
    caller_address,
    contract_address,
    amount,
    gas_price,
    gas_limit,
    log_handler,
    {});
  auto code = result.output;
  contract.acc.set_code(std::move(code));

  gs.add_addresses_newly_deploy(contract_address, code);
  // update output to contract address
  std::vector<uint8_t> deployed_address(32);

  mvm::to_big_endian(contract_address, deployed_address.data());
  std::vector<uint8_t> truncated_address(20);
  memcpy(truncated_address.data(), deployed_address.data()+12, 20);

  result.output = truncated_address;
  ExecuteResult *rs = processResult(result, gs, log_handler);
  
  gs.Clear();

  return rs;
}

ExecuteResult *call(
  // transaction data
  unsigned char* b_caller_address,
  unsigned char* b_contract_address,
  unsigned char* b_input,
  int   length_input,
  unsigned char* b_amount,
  unsigned long long gas_price,
  unsigned long long gas_limit,
  // block context data
  unsigned long long block_prevrandao,
  unsigned long long block_gas_limit,
  unsigned long long block_time,
  unsigned long long block_base_fee,
  unsigned char* b_block_number,
  unsigned char* b_block_coinbase
)
{
  // format argument to right data type
  uint256_t caller_address = mvm::from_big_endian((uint8_t *)b_caller_address, 20u);
  uint256_t contract_address = mvm::from_big_endian((uint8_t *)b_contract_address, 20u);
  std::vector<uint8_t> input((uint8_t *)b_input, (uint8_t *)b_input + length_input);
  uint256_t amount = mvm::from_big_endian((uint8_t *)b_amount, 32u);


  uint256_t block_number = mvm::from_big_endian((uint8_t *)b_block_number, 32u);
  uint256_t block_coinbase = mvm::from_big_endian((uint8_t *)b_block_coinbase, 20u);

  mvm::BlockContext blockContext = CreateBlockContext(
    block_prevrandao,
    block_gas_limit,
    block_time, 
    block_base_fee,
    block_number,   
    block_coinbase
  );

  mvm::MyGlobalState gs(blockContext);
  //  init env
  mvm::VectorLogHandler log_handler;

  auto result = run(
    gs,
    false,
    caller_address,
    contract_address,
    amount,
    gas_price,
    gas_limit,
    log_handler,
    input
  );

  ExecuteResult *rs = processResult(result, gs, log_handler);

  gs.Clear();

  return rs;
}

void freeResult(ExecuteResult* ptr) {
    delete[] ptr->b_exmsg;
    delete[] ptr->b_output;
    for (int i = 0; i < ptr->length_add_balance_change; i++) {
        delete[] ptr->b_add_balance_change[i];
    }
    delete[] ptr->b_add_balance_change;
    for (int i = 0; i < ptr->length_sub_balance_change; i++) {
        delete[] ptr->b_sub_balance_change[i];
    }
    delete[] ptr->b_sub_balance_change;
    for (int i = 0; i < ptr->length_code_change; i++) {
        delete[] ptr->b_code_change[i];
    }
    delete[] ptr->b_code_change;
    delete[] ptr->length_codes;
    for (int i = 0; i < ptr->length_storage_change; i++) {
        delete[] ptr->b_storage_change[i];
    }
    delete[] ptr->b_storage_change;
    delete[] ptr->length_storages;
    delete[] ptr->b_logs;
    delete ptr;  
    std::cout << "freed Result" << std::endl;
}

void freePendingResult() {
    delete[] pendingResult->b_exmsg;
    delete[] pendingResult->b_output;
    for (int i = 0; i < pendingResult->length_add_balance_change; i++) {
        delete[] pendingResult->b_add_balance_change[i];
    }
    delete[] pendingResult->b_add_balance_change;
    for (int i = 0; i < pendingResult->length_sub_balance_change; i++) {
        delete[] pendingResult->b_sub_balance_change[i];
    }
    delete[] pendingResult->b_sub_balance_change;
    for (int i = 0; i < pendingResult->length_code_change; i++) {
        delete[] pendingResult->b_code_change[i];
    }
    delete[] pendingResult->b_code_change;
    delete[] pendingResult->length_codes;
    for (int i = 0; i < pendingResult->length_storage_change; i++) {
        delete[] pendingResult->b_storage_change[i];
    }
    delete[] pendingResult->b_storage_change;
    delete[] pendingResult->length_storages;
    delete[] pendingResult->b_logs;
    pendingResult = NULL;
    std::cout << "freed pending Result" << std::endl;
}

ExecuteResult* testMemLeak() {
  // storage
  char *b_output = reinterpret_cast<char *>(malloc(32 * sizeof(char)));
  int length_output = static_cast<int>(32);
  
  int length_add_balance_change = 0;
  int length_sub_balance_change = 0;
  int length_code_change = 0;

  mvm::MyGlobalState gs;
  std::vector<std::vector<uint8_t>> storage_change;
  
  int length_storage_change = 100;
  int *length_storages = new int[length_storage_change];
  uint8_t **b_storage_change = new uint8_t *[length_storage_change];
  for (int i = 0; i < length_storage_change; ++i) {
     int allocSize = 64*10 + 32;
     length_storages[i] = allocSize;
     b_storage_change[i] = new uint8_t[allocSize];
     for(int u=0; u<allocSize; u++) {
      b_storage_change[i][u] = i;
    }
  }
  
  // logs
  std::string str_logs = "{}";

  pendingResult = new ExecuteResult{
    b_exitReason : (char)1,
    b_exception : (char)1,
    b_exmsg : new char[0],
    length_exmsg : 0,

    b_output : b_output,
    length_output : length_output,
    length_add_balance_change : length_add_balance_change,
    length_sub_balance_change : length_sub_balance_change,
    length_code_change : length_code_change,
    b_storage_change : (char **)b_storage_change,
    length_storage_change : length_storage_change,
    length_storages : length_storages,

    b_logs : (char *)malloc((int)str_logs.size() * sizeof(char)),
    length_logs : (int)str_logs.size(),

    gas_used: 0
  };

  memcpy(pendingResult->b_logs, str_logs.c_str(), str_logs.size());

  return pendingResult;
}


void testMemLeakGS(
  int total_address,
  unsigned char* b_contract_addresses
) {
  mvm::MyGlobalState gs;
  for(int i=0; i  < total_address; i++ ){
    uint256_t contract_address = mvm::from_big_endian(b_contract_addresses+ (i*20), 20u);
    gs.get(contract_address);
    std::cout << "gs.get(contract_address)" << std::endl;
  }
}


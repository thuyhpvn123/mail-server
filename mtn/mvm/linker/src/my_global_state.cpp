#include "my_global_state.h"
#include "mvm_linker.hpp"
#include "mvm/exception.h"
#include "mvm/gas.h"

struct GlobalStateGet_return
{
  int status;
  unsigned char *balance_p;
  unsigned char *code_p;
  int code_length;
};

namespace mvm
{
  using ET = Exception::Type;
  void MyGlobalState::remove(const Address &addr)
  {
    accounts.erase(addr);
  }

  AccountState MyGlobalState::get(const Address &addr, GasTracker *gas_tracker)
  {
    uint8_t b_address[32];
    mvm::to_big_endian(addr, b_address);

    const auto acc = accounts.find(addr);
    if (acc != accounts.cend())
    {
      if (gas_tracker != nullptr)
      {
        gas_tracker->add_gas_used(getTouchedAddressGasCost());
      }
      return acc->second;
    }
    else
    {
      GlobalStateGet_return accountQueryData = GlobalStateGet(b_address + 12);
      if (accountQueryData.status == 2)
      {
        throw Exception(
            ET::addressNotInRelated,
            "Address not in related addresses: " + mvm::address_to_hex_string(addr));
      }
      if (accountQueryData.status == 1)
      {
        uint256_t balance = from_big_endian(accountQueryData.balance_p, 32u);
        std::vector<uint8_t> code(accountQueryData.code_p, accountQueryData.code_p + accountQueryData.code_length);

        insert({MyAccount(addr, balance, code), MyStorage(addr)});
        const auto acc = accounts.find(addr);
        if (gas_tracker != nullptr)
        {
          gas_tracker->add_gas_used(getUnTouchedAddressGasCost());
        }
        ClearProcessingPointers();
        return acc->second;
      }
      return create(addr, 0, {});
    }
  }

  AccountState MyGlobalState::create(
      const Address &addr, const uint256_t &balance, const Code &code)
  {
    insert({MyAccount(addr, balance, code), MyStorage(addr)});

    return get(addr);
  }

  bool MyGlobalState::exists(const Address &addr)
  {
    return accounts.find(addr) != accounts.end();
  }

  size_t MyGlobalState::num_accounts()
  {
    return accounts.size();
  }

  const BlockContext &MyGlobalState::get_block_context()
  {
    return blockContext;
  }

  uint256_t MyGlobalState::get_block_hash(uint8_t offset)
  {
    return 0u;
  }

  uint256_t MyGlobalState::get_chain_id()
  {
    // TODO: may load from config
    return 0u;
  }

  void MyGlobalState::insert(const StateEntry &p)
  {
    const auto ib = accounts.insert(std::make_pair(p.first.get_address(), p));

    assert(ib.second);
  }

  bool operator==(const MyGlobalState &l, const MyGlobalState &r)
  {
    // TODO
    return true;
    // return (l.accounts == r.accounts) && (l.currentBlock == r.currentBlock);
  }

  // add changes functions
  void MyGlobalState::add_addresses_newly_deploy(const Address &addr, const Code &code)
  {
    addresses_newly_deploy[addr] = code;
  };

  void MyGlobalState::add_addresses_storage_change(const Address &addr, const uint256_t &key, const uint256_t &value)
  {
    addresses_storage_change[addr][key] = value;
  };

  void MyGlobalState::add_addresses_add_balance_change(const Address &addr, const uint256_t &amount)
  {
    addresses_add_balance_change[addr] += amount;
  };

  void MyGlobalState::add_addresses_sub_balance_change(const Address &addr, const uint256_t &amount)
  {
    addresses_sub_balance_change[addr] += amount;
  };

const std::vector<std::vector<uint8_t>> MyGlobalState::get_newly_deploy() {
    std::vector<std::vector<uint8_t>> result;
    
    for (const auto &p : addresses_newly_deploy) {
        std::vector<uint8_t> address_with_code(32 + p.second.size());
        mvm::to_big_endian(p.first, address_with_code.data());
        std::memcpy(address_with_code.data() + 32, p.second.data(), p.second.size());
        result.push_back(address_with_code);
    }
    
    return result;
}

const std::vector<std::vector<uint8_t>> MyGlobalState::get_storage_change() {
    int size = addresses_storage_change.size();
    std::vector<std::vector<uint8_t>> result(size);
    int count = 0;
    
    for (const auto &p : addresses_storage_change) {
        int storage_size = 64 * p.second.size();
        std::vector<uint8_t> storage(storage_size);
        int storage_count = 0;
        
        for (const auto &s : p.second) {
            int idx = storage_count * 64;
            mvm::to_big_endian(s.first, storage.data() + idx);
            mvm::to_big_endian(s.second, storage.data() + idx + 32);
            storage_count++;
        }
        
        std::vector<uint8_t> address_with_storage_change(32 + storage_size);
        mvm::to_big_endian(p.first, address_with_storage_change.data());
        std::memcpy(address_with_storage_change.data() + 32, storage.data(), storage_size);
        result[count] = address_with_storage_change;
        count++;
    }
    
    return result;
}

const std::vector<std::vector<uint8_t>> MyGlobalState::get_add_balance_change() {
    int size = addresses_add_balance_change.size();
    std::vector<std::vector<uint8_t>> result(size);
    int count = 0;
    
    for (const auto &p : addresses_add_balance_change) {
        std::vector<uint8_t> address_with_add_balance_change(64);
        mvm::to_big_endian(p.first, address_with_add_balance_change.data());
        mvm::to_big_endian(p.second, address_with_add_balance_change.data() + 32);
        result[count] = address_with_add_balance_change;
        count++;
    }
    
    return result;
}

const std::vector<std::vector<uint8_t>> MyGlobalState::get_sub_balance_change() {
    int size = addresses_sub_balance_change.size();
    std::vector<std::vector<uint8_t>> result(size);
    int count = 0;
    
    for (const auto &p : addresses_sub_balance_change) {
        std::vector<uint8_t> address_with_sub_balance_change(64);
        mvm::to_big_endian(p.first, address_with_sub_balance_change.data());
        mvm::to_big_endian(p.second, address_with_sub_balance_change.data() + 32);
        result[count] = address_with_sub_balance_change;
        count++;
    }
    
    return result;
}

void MyGlobalState::Clear() {
  for (const auto &a : accounts) {
        MyStorage storage = accounts[a.first].second;
        storage.Clear();
    }
};


} // namespace mvm

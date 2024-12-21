#pragma once

#include "mvm/globalstate.h"
#include "my_storage.h"
#include "my_account.h"

namespace mvm
{
  /**
   * My implementation of GlobalState
   */
  class MyGlobalState : public GlobalState
  {
  public:
    using StateEntry = std::pair<MyAccount, MyStorage>;

    std::map<Address, Code> addresses_newly_deploy;
    std::map<Address, std::map<uint256_t,uint256_t>> addresses_storage_change;
    std::map<Address, uint256_t> addresses_add_balance_change;
    std::map<Address, uint256_t> addresses_sub_balance_change;

  private:
    BlockContext blockContext;

    std::map<Address, StateEntry> accounts;
  

  public:
    MyGlobalState() = default;
    explicit MyGlobalState(BlockContext blockContext) : blockContext(std::move(blockContext)) {}

    virtual void remove(const Address& addr) override;

    AccountState get(const Address& addr, GasTracker *gas_tracker = NULL) override;
    AccountState create(
      const Address& addr, const uint256_t& balance, const Code& code) override;

    bool exists(const Address& addr);
    size_t num_accounts();

    virtual const BlockContext& get_block_context() override;
    virtual uint256_t get_block_hash(uint8_t offset) override;
    virtual uint256_t get_chain_id() override;

    /**
     * Add and Extract changes data from global state
     * to create result
     */
    virtual void add_addresses_newly_deploy(const Address& addr, const Code& code) override;
    virtual void add_addresses_storage_change(const Address& addr, const uint256_t& key, const uint256_t& value) override;
    virtual void add_addresses_add_balance_change(const Address& addr, const uint256_t& amount) override;
    virtual void add_addresses_sub_balance_change(const Address& addr, const uint256_t& amount) override;

    const std::vector<std::vector<uint8_t>> get_newly_deploy();
    const std::vector<std::vector<uint8_t>> get_storage_change();
    const std::vector<std::vector<uint8_t>> get_add_balance_change();
    const std::vector<std::vector<uint8_t>> get_sub_balance_change();
    
    /**
     * For tests which require some initial state, allow manual insertion of
     * pre-constructed accounts
     */
    void insert(const StateEntry& e);
    friend bool operator==(const MyGlobalState&, const MyGlobalState&);
    void Clear();
  };
} // namespace mvm

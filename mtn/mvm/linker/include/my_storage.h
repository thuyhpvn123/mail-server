// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

#pragma once

#include "mvm/storage.h"
#include "my_account.h"

#include <map>
#include <nlohmann/json.hpp>

namespace mvm
{
  /**
   * merkle patricia trie implementation of Storage
   */
  class MyStorage : public Storage
  {

  private:
    Address address = {};
    std::map<uint256_t, uint256_t> cache;

  public:
    MyStorage(){};
    MyStorage(const Address &a) : address(a)
    {
    }
    ~MyStorage();
    // MyStorage(const std::vector<std::vector<uint8_t>>&storage);

    void Clear();
    void store(const uint256_t& key, const uint256_t& value, GasTracker* gas_tracker = NULL) override;
    uint256_t load(const uint256_t& key, GasTracker* gas_tracker = NULL) override;
    bool remove(const uint256_t& key) override;
    bool exists(const uint256_t& key);
  };
} // namespace mvm

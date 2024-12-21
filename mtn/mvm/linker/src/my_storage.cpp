// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

#include "my_storage.h"

#include "mvm/util.h"
#include "mvm/gas.h"
#include "mvm_linker.hpp"
#include "mvm/exception.h"

#include <ostream>

struct GetStorageValue_return
{
  unsigned char *value;
  bool success;
};

namespace mvm
{
  using ET = Exception::Type;
  MyStorage::~MyStorage()
  {
  }

  void MyStorage::Clear()
  {
    cache.clear();
  }

  void MyStorage::store(const uint256_t &key, const uint256_t &value, GasTracker *gas_tracker)
  {
    if (gas_tracker != NULL)
    {
      uint256_t old_value = load(key);
      if (value == old_value)
      {
        gas_tracker->add_gas_used(getSstoreGasCost(old_value, value));
      }
    }
    cache[key] = value;
  }

  uint256_t MyStorage::load(const uint256_t &key, GasTracker *gas_tracker)
  {
    if (gas_tracker != NULL)
    {
      // TODO: check touched storage
      gas_tracker->add_gas_used(getTouchedStorageGasCost());
    }
    // find in cache
    auto e = cache.find(key);
    if (e == cache.end())
    {
      // if not found in cache, smart contract db to find
      uint8_t b_address[32];
      mvm::to_big_endian(address, b_address);
      uint8_t b_key[32];
      mvm::to_big_endian(key, b_key);
      auto get_rs = GetStorageValue(
          b_address + 12,
          b_key);

      if (!get_rs.success)
      {
        throw Exception(
            ET::ErrExecutionReverted,
            "Get Storage Value Error");
      }

      uint256_t value = mvm::from_big_endian(get_rs.value, 32u);
      // delete[] get_rs; // ! Need to confirm
      // add value to cache
      cache[key] = value;
      return value;
    }
    return e->second;
  }

  bool MyStorage::exists(const uint256_t &key)
  {
    // I dont see this function is used
    // This may wrong because not check in remote storage
    return cache.find(key) != cache.end();
  }

  bool MyStorage::remove(const uint256_t &key)
  {
    // load to check if key in remote storage
    load(key, NULL);
    auto e = cache.find(key);
    if (e == cache.end())
      return false;

    cache[key] = 0;
    return true;
  }

  inline std::ostream &operator<<(std::ostream &os, const MyStorage &s)
  {
    // os << nlohmann::json(s).dump(2);
    // TODO
    return os;
  }

} // namespace mvm

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

#include "my_extension.h"
#include "mvm_linker.hpp"

struct ExtensionCallGetApi_return {
  unsigned char* data_p;
  int data_size;
};

struct ExtensionExtractJsonField_return {
  unsigned char* data_p;
  int data_size;
};

struct ExtensionBlst_return{
  unsigned char* data_p;
  int data_size;
};

mvm::Code MyExtension::CallGetApi(mvm::Code input)
{
  ExtensionCallGetApi_return data = ExtensionCallGetApi(input.data(), input.size());
  std::vector<uint8_t> vec(data.data_p, data.data_p + data.data_size);
  return vec;
}

mvm::Code MyExtension::ExtractJsonField(mvm::Code input)
{
  ExtensionExtractJsonField_return data = ExtensionExtractJsonField(input.data(), input.size());
  std::vector<uint8_t> vec(data.data_p, data.data_p + data.data_size);
  return vec;
}


mvm::Code MyExtension::Blst(mvm::Code input)
{
  ExtensionBlst_return data = ExtensionBlst(input.data(), input.size());
  std::vector<uint8_t> vec(data.data_p, data.data_p + data.data_size);
  return vec;
}

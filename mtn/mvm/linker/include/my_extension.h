// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

#pragma once

#include "mvm/extension.h"

class MyExtension : public mvm::Extension
{
    public:
    MyExtension() = default;
    virtual mvm::Code CallGetApi(mvm::Code input) override;
    virtual mvm::Code ExtractJsonField(mvm::Code input) override;
    virtual mvm::Code Blst(mvm::Code input) override;
};

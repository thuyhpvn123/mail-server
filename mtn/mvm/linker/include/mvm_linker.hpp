#pragma once

#include <stdbool.h>

#ifdef __cplusplus
extern "C"
{
#endif
    struct ExecuteResult
    {
        char b_exitReason;
        char b_exception;
        char *b_exmsg;
        int length_exmsg;

        char *b_output;
        int length_output;

        char **b_add_balance_change;
        int length_add_balance_change;

        char **b_sub_balance_change;
        int length_sub_balance_change;

        char **b_code_change;
        int length_code_change;
        int *length_codes;

        char **b_storage_change;
        int length_storage_change;
        int *length_storages;

        char *b_logs;
        int length_logs;

        unsigned long long gas_used;
    };

    struct ExecuteResult *deploy(
        // transaction data
        unsigned char *b_caller_address,
        unsigned char *b_caller_last_hash,
        unsigned char *b_contract_constructor,
        int contract_constructor_length,
        unsigned char *b_amount,
        unsigned long long gas_price,
        unsigned long long gas_limit,
        // block context data
        unsigned long long block_prevrandao,
        unsigned long long block_gas_limit,
        unsigned long long block_time,
        unsigned long long block_base_fee,
        unsigned char *b_block_number,
        unsigned char *b_block_coinbase);

    struct ExecuteResult *call(
        // transaction data
        unsigned char *b_caller_address,
        unsigned char *b_contract_address,
        unsigned char *b_input,
        int length_input,
        unsigned char *b_amount,
        unsigned long long gas_price,
        unsigned long long gas_limit,
        // block context data
        unsigned long long block_prevrandao,
        unsigned long long block_gas_limit,
        unsigned long long block_time,
        unsigned long long block_base_fee,
        unsigned char *b_block_number,
        unsigned char *b_block_coinbase);

    void freeResult(struct ExecuteResult *);
    void freePendingResult();

    extern struct GlobalStateGet_return GlobalStateGet(unsigned char *);
    extern void ClearProcessingPointers();
    extern struct ExtensionCallGetApi_return ExtensionCallGetApi(unsigned char *, int);
    extern struct ExtensionExtractJsonField_return ExtensionExtractJsonField(unsigned char *, int);
    extern struct ExtensionBlst_return ExtensionBlst(unsigned char *, int);
    extern struct GetStorageValue_return GetStorageValue(unsigned char *, unsigned char *);

    extern void GoLogString(int, char *);
    extern void GoLogBytes(int, unsigned char *, int);

    struct ExecuteResult *testMemLeak();
    void testMemLeakGS(
        int total_address,
        unsigned char *b_addresses);

#ifdef __cplusplus
}
#endif

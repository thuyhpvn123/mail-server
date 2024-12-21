package mvm

/*
#cgo CFLAGS: -w
#cgo CXXFLAGS: -std=c++17 -w
#cgo LDFLAGS: -L./linker/build/lib/static -lmvm_linker -L./c_mvm/build/lib/static -lmvm -lstdc++
#cgo CPPFLAGS: -I./linker/build/include
#include "mvm_linker.hpp"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/hex"
	"encoding/json"
	"unsafe"

	"github.com/ethereum/go-ethereum/crypto"
	"gomail/mtn/logger"
	pb "gomail/mtn/proto"
)

func extractExecuteResult(cExecuteResult *C.struct_ExecuteResult) *MVMExecuteResult {
	status := pb.RECEIPT_STATUS(cExecuteResult.b_exitReason)
	var exception pb.EXCEPTION
	if status == pb.RECEIPT_STATUS_THREW {
		exception = pb.EXCEPTION(cExecuteResult.b_exception)
	} else {
		exception = pb.EXCEPTION_NONE
	}

	// extract add balance
	mapAddBalance := extractAddBalance(cExecuteResult)
	mapSubBalance := extractSubBalance(cExecuteResult)
	mapCodeChange, mapCodeHash := extractCodeChange(cExecuteResult)
	mapStorageChange := extractStorageChange(cExecuteResult)
	jEventLogs := extractEventLogs(cExecuteResult)

	uptr := unsafe.Pointer(cExecuteResult.b_exmsg)
	exmsg := string(C.GoBytes(uptr, cExecuteResult.length_exmsg))
	// C.free(uptr)

	uptr = unsafe.Pointer(cExecuteResult.b_output)
	rt := C.GoBytes(uptr, cExecuteResult.length_output)
	// C.free(uptr)

	gasUsed := (uint64)(cExecuteResult.gas_used)

	return &MVMExecuteResult{
		mapAddBalance,
		mapSubBalance,
		mapCodeChange, mapCodeHash,
		mapStorageChange, 
		jEventLogs,
		status,
		exception,
		exmsg,
		rt,
		gasUsed,
	}
}

// extract funcs
func extractAddBalance(
	cExecuteResult *C.struct_ExecuteResult,
) (
	mapAddBalance map[string][]byte,
) {
	// extract add balance
	bAddBalanceChange := unsafe.Slice(cExecuteResult.b_add_balance_change, cExecuteResult.length_add_balance_change)
	mapAddBalance = make(map[string][]byte, len(bAddBalanceChange))
	for _, v := range bAddBalanceChange {
		uptr := unsafe.Pointer(v)
		addrWithAddBalanceChange := C.GoBytes(uptr, (C.int)(64))
		// C.free(uptr)
		mapAddBalance[hex.EncodeToString(addrWithAddBalanceChange[12:32])] = addrWithAddBalanceChange[32:]
	}
	// C.free(unsafe.Pointer(cExecuteResult.b_add_balance_change))
	return
}

func extractSubBalance(
	cExecuteResult *C.struct_ExecuteResult,
) (
	mapSubBalance map[string][]byte,
) {
	bSubBalanceChange := unsafe.Slice(cExecuteResult.b_sub_balance_change, cExecuteResult.length_sub_balance_change)
	mapSubBalance = make(map[string][]byte, len(bSubBalanceChange))
	for _, v := range bSubBalanceChange {
		uptr := unsafe.Pointer(v)
		addrWithSubBalanceChange := C.GoBytes(uptr, (C.int)(64))
		// C.free(uptr)
		mapSubBalance[hex.EncodeToString(addrWithSubBalanceChange[12:32])] = addrWithSubBalanceChange[32:]
	}
	// C.free(unsafe.Pointer(cExecuteResult.b_sub_balance_change))
	return
}

func extractCodeChange(
	cExecuteResult *C.struct_ExecuteResult,
) (
	mapCodeChange map[string][]byte,
	mapCodeHash map[string][]byte,
) {
	mapCodeChange = make(map[string][]byte, cExecuteResult.length_code_change)
	mapCodeHash = make(map[string][]byte, cExecuteResult.length_code_change)

	bCodeChange := unsafe.Slice(cExecuteResult.b_code_change, cExecuteResult.length_code_change)
	cLengthCodes := unsafe.Slice(cExecuteResult.length_codes, cExecuteResult.length_code_change)
	lengthCodes := make([]int, cExecuteResult.length_code_change)
	for i, v := range cLengthCodes {
		lengthCodes[i] = int(v)
	}

	for i, v := range lengthCodes {
    sptr := unsafe.Pointer(bCodeChange[i])
		uptr := unsafe.Pointer(sptr)
		addrWithCode := C.GoBytes(uptr, (C.int)(v))
		address := hex.EncodeToString(addrWithCode[12:32])
		code := addrWithCode[32:]
		mapCodeChange[address] = code
		mapCodeHash[address] = crypto.Keccak256(code)
    logger.DebugP("Code for hash: ", code)
	}

	return
}

func extractStorageChange(
	cExecuteResult *C.struct_ExecuteResult,
) (
	mapStorageChange map[string]map[string][]byte,
) {
	// extract storage changes
	mapStorageChange = make(map[string]map[string][]byte, cExecuteResult.length_storage_change)

	bStorageChange := unsafe.Slice(cExecuteResult.b_storage_change, cExecuteResult.length_storage_change)
	cLengthStorages := unsafe.Slice(cExecuteResult.length_storages, cExecuteResult.length_storage_change)
	lengthStorages := make([]int, cExecuteResult.length_storage_change)
	for i, v := range cLengthStorages {
		lengthStorages[i] = int(v)
	}

	for i, v := range lengthStorages {
    sprt := unsafe.Pointer(bStorageChange[i])
		addrWithStorageChanges := C.GoBytes(sprt, (C.int)(v+32))
		// C.free(sprt)
		address := hex.EncodeToString(addrWithStorageChanges[12:32])
		addrWithStorageChanges = addrWithStorageChanges[32:]
		storageCount := v / 64
		mapStorageChange[address] = make(map[string][]byte, storageCount)
		for j := 0; j < storageCount; j++ {
      // 32 bytes for key, 32 bytes for value
      key := hex.EncodeToString(addrWithStorageChanges[j*64 : j*64 + 32])
      value := addrWithStorageChanges[j*64 + 32 : (j+1)*64]
			mapStorageChange[address][key] = value
		}
	}

	return
}

func extractEventLogs(
	cExecuteResult *C.struct_ExecuteResult,
) (
	logJson LogsJson,
) {
  sptr :=  unsafe.Pointer(cExecuteResult.b_logs)
	rawLogs := C.GoBytes(sptr, cExecuteResult.length_logs)
 	// C.free(sptr) 
	json.Unmarshal(rawLogs, &logJson.Logs)
	return
}

package mvm

/*
#cgo CFLAGS: -w -O3 -march=native -mtune=native
#cgo CXXFLAGS: -std=c++17 -w -O3 -march=native -mtune=native
#cgo LDFLAGS: -L./linker/build/lib/static -lmvm_linker -L./c_mvm/build/lib/static -lmvm -lstdc++
#cgo CPPFLAGS: -I./linker/build/include
#include "mvm_linker.hpp"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/hex"
	"math/big"
	"time"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"gomail/mtn/logger"
	"gomail/mtn/types"
)

var apiInstance *MVMApi

type AccountStateDB interface {
	AccountState(address common.Address) (types.AccountState, error)
}

type SmartContractDB interface {
	Code(address common.Address) []byte
	StorageValue(address common.Address, key []byte) ([]byte, bool)
}

type MVMApi struct {
	smartContractDb         SmartContractDB
	accountStateDb          AccountStateDB
	currentRelatedAddresses map[common.Address]struct{}
}

func MVMApiInstance() *MVMApi {
	return apiInstance
}

func Clear() {
	apiInstance = nil
}

func InitMVMApi(
	smartContractDb SmartContractDB,
	accountStateDb AccountStateDB,
) {
	if apiInstance == nil {
		apiInstance = &MVMApi{
			smartContractDb,
			accountStateDb,
			make(map[common.Address]struct{}),
		}
	}
}

func (a *MVMApi) SetSmartContractDb(smartContractDb SmartContractDB) {
	a.smartContractDb = smartContractDb
}

func (a *MVMApi) SmartContractDatas() SmartContractDB {
	return a.smartContractDb
}

func (a *MVMApi) SetAccountStateDb(accountStateDb AccountStateDB) {
	a.accountStateDb = accountStateDb
}

func (a *MVMApi) AccountStateDb() AccountStateDB {
	return a.accountStateDb
}

func (a *MVMApi) SetRelatedAddresses(addresses []common.Address) {
	a.currentRelatedAddresses = make(map[common.Address]struct{}, len(addresses))
	for _, v := range addresses {
		a.currentRelatedAddresses[v] = struct{}{}
	}
}

func (a *MVMApi) InRelatedAddress(address common.Address) bool {
	_, ok := a.currentRelatedAddresses[address]
	return ok
}

func (a *MVMApi) Call(
	// transaction data
	bSender []byte,
	bContractAddress []byte,
	bInput []byte,
	amount *big.Int,
	gasPrice uint64,
	gasLimit uint64,
	// block context data
	blockPrevrandao uint64,
	blockGasLimit uint64,
	blockTime uint64,
	blockBaseFee uint64,
	blockNumber uint64,
	blockCoinbase common.Address,
) *MVMExecuteResult {
	// transaction data
	bAmount := [32]byte{}
	amount.FillBytes(bAmount[:])
	cBSender := C.CBytes(bSender)
	cBContractAddress := C.CBytes(bContractAddress)
	cBInput := C.CBytes(bInput)
	cBAmount := C.CBytes(bAmount[:])

	// block context data
	bBlockNumber := [32]byte{}
	bigBlockNumber := big.NewInt(int64(blockNumber))
	bigBlockNumber.FillBytes(bBlockNumber[:])
	bBlockCoinbase := blockCoinbase.Bytes()
	cBBlockNumber := C.CBytes(bBlockNumber[:])
	cBBlockCoinbase := C.CBytes(bBlockCoinbase)

	defer C.free(unsafe.Pointer(cBSender))
	defer C.free(unsafe.Pointer(cBContractAddress))
	defer C.free(unsafe.Pointer(cBInput))
	defer C.free(unsafe.Pointer(cBAmount))
	defer C.free(unsafe.Pointer(cBBlockNumber))
	defer C.free(unsafe.Pointer(cBBlockCoinbase))
	startTime := time.Now()
	cRs := C.call(
		// transaction data
		(*C.uchar)(cBSender),
		(*C.uchar)(cBContractAddress),
		(*C.uchar)(cBInput),
		(C.int)(len(bInput)),
		(*C.uchar)(cBAmount),
		(C.ulonglong)(gasPrice),
		(C.ulonglong)(gasLimit),
		// block context data
		(C.ulonglong)(blockPrevrandao),
		(C.ulonglong)(blockGasLimit),
		(C.ulonglong)(blockTime),
		(C.ulonglong)(blockBaseFee),
		(*C.uchar)(cBBlockNumber),
		(*C.uchar)(cBBlockCoinbase),
	)
	logger.DebugP("Call took: ", time.Since(startTime))
	rs := extractExecuteResult(cRs)
	// C.freeResult(cRs)
	C.freePendingResult()
	return rs
}

func (a *MVMApi) Deploy(
	// transaction data
	bSender []byte,
	bLastHash []byte,
	bContractConstructor []byte,
	amount *big.Int,
	gasPrice uint64,
	gasLimit uint64,
	// block context data
	blockPrevrandao uint64,
	blockGasLimit uint64,
	blockTime uint64,
	blockBaseFee uint64,
	blockNumber uint64,
	blockCoinbase common.Address,
) *MVMExecuteResult {
	// transaction data
	bAmount := [32]byte{}
	amount.FillBytes(bAmount[:])
	constructorLength := len(bContractConstructor)
	cBSender := C.CBytes(bSender)
	cBLastHash := C.CBytes(bLastHash)
	cBContractConstructor := C.CBytes(bContractConstructor)
	cBAmount := C.CBytes(bAmount[:])
	// block context data
	bBlockNumber := [32]byte{}
	bigBlockNumber := big.NewInt(int64(blockNumber))
	bigBlockNumber.FillBytes(bBlockNumber[:])
	bBlockCoinbase := blockCoinbase.Bytes()

	cBBlockNumber := C.CBytes(bBlockNumber[:])
	cBBlockCoinbase := C.CBytes(bBlockCoinbase)

	defer C.free(unsafe.Pointer(cBSender))
	defer C.free(unsafe.Pointer(cBLastHash))
	defer C.free(unsafe.Pointer(cBContractConstructor))
	defer C.free(unsafe.Pointer(cBAmount))

	defer C.free(unsafe.Pointer(cBBlockNumber))
	defer C.free(unsafe.Pointer(cBBlockCoinbase))

	startTime := time.Now()
	cRs := C.deploy(
		// transaction data
		(*C.uchar)(cBSender),
		(*C.uchar)(cBLastHash),
		(*C.uchar)(cBContractConstructor),
		(C.int)(constructorLength),
		(*C.uchar)(cBAmount),
		(C.ulonglong)(gasPrice),
		(C.ulonglong)(gasLimit),
		// block context data
		(C.ulonglong)(blockPrevrandao),
		(C.ulonglong)(blockGasLimit),
		(C.ulonglong)(blockTime),
		(C.ulonglong)(blockBaseFee),
		(*C.uchar)(cBBlockNumber),
		(*C.uchar)(cBBlockCoinbase),
	)
	logger.DebugP("Deploy took: ", time.Since(startTime))
	rs := extractExecuteResult(cRs)
	// C.freeResult(cRs)
	C.freePendingResult()
	return rs
}

// GLOBAL STATE Functions
var (
	processingPointers []unsafe.Pointer
)

//export GlobalStateGet
func GlobalStateGet(
	address *C.uchar,
) (
	status C.int, // 0 not found, 1 found, 2 not in related
	balance_p *C.uchar,
	code_p *C.uchar,
	code_length C.int,
) {
	startTime := time.Now()
	mvmApi := MVMApiInstance()
	bAddress := C.GoBytes(unsafe.Pointer(address), 20)
	fAddress := common.BytesToAddress(bAddress)
	// extensions
	logger.DebugP("GlobalStateGet address: ", fAddress)
	if fAddress == common.HexToAddress("0x0000000000000000000000000000000000000101") ||
		fAddress == common.HexToAddress("0x0000000000000000000000000000000000000102") ||
		fAddress == common.HexToAddress("0x0000000000000000000000000000000000000103") {
		balance := uint256.NewInt(0).Bytes32()
		cBBalance := C.CBytes(balance[:])
		code := []byte{0x01}
		lenCode := len(code)
		cBCode := C.CBytes(code)
		processingPointers = append(processingPointers, cBBalance)
		processingPointers = append(processingPointers, cBCode)
		logger.DebugP("GlobalStateGet extension took: ", time.Since(startTime))
		return C.int(1), (*C.uchar)(cBBalance), (*C.uchar)(cBCode), (C.int)(lenCode)
	}

	//
	inRelatedAddresses := mvmApi.InRelatedAddress(fAddress)
	if !inRelatedAddresses {
		return C.int(2), nil, nil, 0
	}
	accountState, err := mvmApi.accountStateDb.AccountState(fAddress)
	if err != nil {
		logger.Error("Failed to get account state", "err", err)
	}
	logger.Trace("Geted account state", accountState)

	if accountState == nil {
		return C.int(0), nil, nil, 0
	}

	bigBalance := big.NewInt(0).Add(
		accountState.Balance(),
		accountState.PendingBalance(),
	)
	b32Balance := [32]byte{}
	bigBalance.FillBytes(b32Balance[:])
	bCode := []byte{}

	smartContractState := accountState.SmartContractState()
	if smartContractState != nil {
		bCode = mvmApi.smartContractDb.Code(fAddress)
	}

	cBBalance := C.CBytes(b32Balance[:])
	lenCode := len(bCode)
	cBCode := C.CBytes(bCode)
	processingPointers = append(processingPointers, cBBalance)
	processingPointers = append(processingPointers, cBCode)
	logger.DebugP("GlobalStateGet took: ", time.Since(startTime))

	return C.int(1), (*C.uchar)(cBBalance), (*C.uchar)(cBCode), (C.int)(lenCode)
}

// go functions
//
//export ClearProcessingPointers
func ClearProcessingPointers() {
	for _, p := range processingPointers {
		C.free(p)
	}
	processingPointers = []unsafe.Pointer{}
}

func TestMemLeak() {
	cRs := C.testMemLeak()
	rs := extractExecuteResult(cRs)
	logger.DebugP("TestMemLeak: ", rs)
	C.freePendingResult()
}

func TestMemLeakGs(addresses []common.Address) {
	totalAddress := len(addresses)
	bAddress := []byte{}
	for i := range totalAddress {
		bAddress = append(bAddress, addresses[i].Bytes()...)
	}
	cAddress := C.CBytes(bAddress)
	logger.DebugP("TotalAddress", totalAddress)
	logger.DebugP("bAddress", hex.EncodeToString(bAddress))

	C.testMemLeakGS(
		C.int(totalAddress),
		(*C.uchar)(cAddress),
	)
}

//export GetStorageValue
func GetStorageValue(
	address *C.uchar,
	key *C.uchar,
) (value *C.uchar, success bool) {
	mvmApi := MVMApiInstance()
	bAddress := C.GoBytes(unsafe.Pointer(address), 20)
	bKey := C.GoBytes(unsafe.Pointer(key), 32)
	fAddress := common.BytesToAddress(bAddress)
	logger.Debug("GetStorageValue address: ", fAddress, hex.EncodeToString(bKey))
	bValue, success := mvmApi.smartContractDb.StorageValue(fAddress, bKey)
	cValue := C.CBytes(bValue)
	return (*C.uchar)(cValue), success
}

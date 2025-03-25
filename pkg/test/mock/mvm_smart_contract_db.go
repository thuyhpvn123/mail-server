package mock

import (
	"github.com/ethereum/go-ethereum/common"
)

type TestMVMSmartContractDB struct {
	CodeFn         func(address common.Address) []byte
	StorageValueFn func(address common.Address, key []byte) []byte
}

func (t *TestMVMSmartContractDB) Code(address common.Address) []byte {
	return t.CodeFn(address)
}

func (t *TestMVMSmartContractDB) StorageValue(address common.Address, key []byte) []byte {
	return t.StorageValueFn(address, key)
}

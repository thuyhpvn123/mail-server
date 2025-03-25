package mock

import (
	"github.com/ethereum/go-ethereum/common"

	"gomail/types"
)

type TestMVMAccountStateDB struct {
	AccountStateFn func(address common.Address) types.AccountState
}

func (t *TestMVMAccountStateDB) AccountState(address common.Address) types.AccountState {
	return t.AccountStateFn(address)
}

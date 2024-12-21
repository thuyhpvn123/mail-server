package mock

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	p_common "gomail/mtn/common"
	"gomail/mtn/state"
	"gomail/mtn/types"
)

type TestAccountStateDB struct {
	AccountStates map[common.Address]types.AccountState
	DirtyAccounts map[common.Address]types.AccountState
}

func NewTestAccountStateDB() *TestAccountStateDB {
	return &TestAccountStateDB{
		AccountStates: make(map[common.Address]types.AccountState),
		DirtyAccounts: make(map[common.Address]types.AccountState),
	}
}

func (db *TestAccountStateDB) AccountState(address common.Address) types.AccountState {
	return db.AccountStates[address]
}

func (db *TestAccountStateDB) SubPendingBalance(address common.Address, amount *big.Int) error {
	as := db.getOrCreateAccountState(address)
	err := as.SubPendingBalance(amount)
	if err != nil {
		return err
	}
	return nil
}

func (db *TestAccountStateDB) AddPendingBalance(address common.Address, amount *big.Int) {
	as := db.getOrCreateAccountState(address)
	as.AddPendingBalance(amount)
}

func (db *TestAccountStateDB) AddBalance(address common.Address, amount *big.Int) {
	as := db.getOrCreateAccountState(address)
	as.AddBalance(amount)
}

func (db *TestAccountStateDB) SubBalance(address common.Address, amount *big.Int) error {
	as := db.getOrCreateAccountState(address)
	err := as.SubBalance(amount)
	if err != nil {
		return err
	}
	return nil
}

func (db *TestAccountStateDB) SubTotalBalance(address common.Address, amount *big.Int) error {
	as := db.getOrCreateAccountState(address)
	err := as.SubTotalBalance(amount)
	if err != nil {
		return err
	}
	return nil
}

func (db *TestAccountStateDB) SetLastHash(address common.Address, hash common.Hash) {
	as := db.getOrCreateAccountState(address)
	as.SetLastHash(hash)
}

func (db *TestAccountStateDB) SetNewDeviceKey(address common.Address, newDeviceKey common.Hash) {
	as := db.getOrCreateAccountState(address)
	as.SetNewDeviceKey(newDeviceKey)
}

func (db *TestAccountStateDB) SetState(as types.AccountState) {
	db.DirtyAccounts[as.Address()] = as
}

// smart contract state
func (db *TestAccountStateDB) SetCreatorPublicKey(
	address common.Address,
	creatorPublicKey p_common.PublicKey,
) {
	as := db.getOrCreateAccountState(address)
	as.SetCreatorPublicKey(creatorPublicKey)
}

func (db *TestAccountStateDB) SetCodeHash(address common.Address, codeHash common.Hash) {
	as := db.getOrCreateAccountState(address)
	as.SetCodeHash(codeHash)
}

func (db *TestAccountStateDB) SetStorageRoot(address common.Address, storageRoot common.Hash) {
	as := db.getOrCreateAccountState(address)
	as.SetStorageRoot(storageRoot)
}

func (db *TestAccountStateDB) SetStorageAddress(
	address common.Address,
	storageAddress common.Address,
) {
	as := db.getOrCreateAccountState(address)
	as.SetStorageAddress(storageAddress)
}

func (db *TestAccountStateDB) getOrCreateAccountState(address common.Address) types.AccountState {
	// get from dirty
	// if not exist in dirty then get from account states
	// if not exist in account states then create new
	as := db.DirtyAccounts[address]
	if as == nil {
		as = db.AccountStates[address]
		if as == nil {
			as = state.NewAccountState(address)
			db.AccountStates[address] = as
		}
		db.DirtyAccounts[address] = as.Copy()
		as = db.DirtyAccounts[address]
	}
	// return copy of account state
	return as
}

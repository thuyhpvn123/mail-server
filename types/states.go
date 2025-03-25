package types

// type AccountStatesManager interface {
// 	// general
// 	Hash() common.Hash
// 	Commit() (common.Hash, error)
// 	CommitTransaction() (common.Hash, error)
// 	RollbackToLastTransaction() error
// 	RollbackToLastBlock() error
//
// 	Copy() AccountStatesManager
// 	// getter
// 	AccountState(address common.Address) AccountState
// 	Exist(common.Address) bool
// 	OriginAccountState(address common.Address) AccountState
// 	GetStorageSnapShot() storage.SnapShot
// 	Storage() storage.Storage
// 	AccountStateChanges() map[common.Address]AccountState
// 	// setter
// 	SetState(AccountState)
// 	SetSmartContractState(address common.Address, smState SmartContractState)
// 	SetNewDeviceKey(address common.Address, newDeviceKey common.Hash)
// 	SetLastHash(address common.Address, newLastHash common.Hash)
// 	AddPendingBalance(address common.Address, amount *uint256.Int)
// 	SubPendingBalance(address common.Address, amount *uint256.Int) error
// 	SubBalance(address common.Address, amount *uint256.Int) error
// 	AddBalance(address common.Address, amount *uint256.Int)
// 	SubTotalBalance(address common.Address, amount *uint256.Int) error
// 	SetCodeHash(address common.Address, hash common.Hash)
// 	SetStorageHost(address common.Address, storageHost string)
// 	SetStorageAddress(address common.Address, storagAddress common.Address)
// 	SetStorageRoot(address common.Address, hash common.Hash)
// 	SetStateChannelState(address common.Address, scState StateChannelState)
//
// 	CopyFrom(AccountStatesManager)
// }

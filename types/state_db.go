package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	p_common "gomail/pkg/common"
	"gomail/pkg/storage"
)

type AccountStateDB interface {
	AccountState(common.Address) (AccountState, error)

	SubPendingBalance(common.Address, *big.Int) error
	AddPendingBalance(common.Address, *big.Int) error

	AddBalance(common.Address, *big.Int) error
	SubBalance(common.Address, *big.Int) error

	SubTotalBalance(common.Address, *big.Int) error

	SetLastHash(common.Address, common.Hash)
	SetNewDeviceKey(common.Address, common.Hash)

	SetState(AccountState)

	IntermediateRoot() (common.Hash, error)
	Commit() (common.Hash, error)
	Discard() error
	Storage() storage.Storage

	// smart contract state
	SetCreatorPublicKey(address common.Address, creatorPublicKey p_common.PublicKey)
	SetCodeHash(address common.Address, codeHash common.Hash)
	SetStorageRoot(address common.Address, storageRoot common.Hash)
	SetStorageAddress(address common.Address, storageAddress common.Address)
	AddLogHash(address common.Address, logsHash common.Hash)

	CopyFrom(as AccountStateDB) error
}

type SmartContractDB interface {
	Code(address common.Address) []byte
	StorageValue(address common.Address, key []byte) ([]byte, bool)
	SetAccountStateDB(asdb AccountStateDB)
	SetBlockNumber(blockNumber uint64)
	SetCode(
		address common.Address,
		codeHash common.Hash,
		code []byte,
	)
	SetStorageValue(address common.Address, key []byte, value []byte) error
	AddEventLogs(eventLogs []EventLog)
	StorageRoot(
		address common.Address,
	) common.Hash
	NewTrieStorage(
		address common.Address,
	) common.Hash
	DeleteAddress(address common.Address)

	GetSmartContractUpdateDatas() map[common.Address]SmartContractUpdateData
	ClearSmartContractUpdateDatas()
}

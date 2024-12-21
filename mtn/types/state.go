package types

import (
	"math/big"

	e_common "github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/reflect/protoreflect"

	"gomail/mtn/common"
	p_common "gomail/mtn/common"
	pb "gomail/mtn/proto"
)

type AccountState interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() *pb.AccountState
	FromProto(*pb.AccountState)
	Copy() AccountState
	String() string

	// getter
	Address() e_common.Address
	LastHash() e_common.Hash
	Balance() *big.Int
	PendingBalance() *big.Int
	TotalBalance() *big.Int
	SmartContractState() SmartContractState
	DeviceKey() e_common.Hash

	SubPendingBalance(*big.Int) error
	AddPendingBalance(*big.Int)

	AddBalance(*big.Int)
	SubBalance(*big.Int) error

	SubTotalBalance(*big.Int) error

	SetLastHash(e_common.Hash)
	SetNewDeviceKey(e_common.Hash)

	// smart contract state
	SetCreatorPublicKey(creatorPublicKey p_common.PublicKey)
	SetCodeHash(codeHash e_common.Hash)
	SetStorageRoot(storageRoot e_common.Hash)
	SetStorageAddress(storageAddress e_common.Address)
	AddLogHash(logsHash e_common.Hash)
}

type SmartContractState interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	String() string

	// getter
	Proto() *pb.SmartContractState
	FromProto(*pb.SmartContractState)

	CreatorPublicKey() common.PublicKey
	CreatorAddress() e_common.Address
	StorageAddress() e_common.Address
	CodeHash() e_common.Hash
	StorageRoot() e_common.Hash
	LogsHash() e_common.Hash

	// setter
	SetCreatorPublicKey(p_common.PublicKey)
	SetStorageAddress(storageAddress e_common.Address)
	SetCodeHash(e_common.Hash)
	SetStorageRoot(e_common.Hash)
	SetLogsHash(e_common.Hash)

	Copy() SmartContractState
}

type UpdateField interface {
	Field() pb.UPDATE_STATE_FIELD
	Value() []byte
}

type UpdateStateFields interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	String() string
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	Address() e_common.Address
	Fields() []UpdateField
	AddField(field pb.UPDATE_STATE_FIELD, value []byte)
}

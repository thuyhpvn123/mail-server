package types

import (
	"math/big"

	e_common "github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/reflect/protoreflect"

	"gomail/mtn/common"
	pb "gomail/mtn/proto"
)

type Transaction interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	String() string

	// getter
	CalculateHash() e_common.Hash
	Hash() e_common.Hash
	NewDeviceKey() e_common.Hash
	LastDeviceKey() e_common.Hash
	FromAddress() e_common.Address
	ToAddress() e_common.Address
	Pubkey() common.PublicKey
	LastHash() e_common.Hash
	Sign() common.Sign
	Amount() *big.Int
	PendingUse() *big.Int
	Action() pb.ACTION
	BRelatedAddresses() [][]byte
	RelatedAddresses() []e_common.Address
	Data() []byte
	Fee(currentGasPrice uint64) *big.Int
	DeployData() DeployData
	CallData() CallData
	OpenStateChannelData() OpenStateChannelData
	UpdateStorageHostData() UpdateStorageHostData
	CommissionSign() common.Sign
	MaxGas() uint64
	MaxGasPrice() uint64
	MaxTimeUse() uint64
	MaxFee() *big.Int

	// setter
	SetSign(privateKey common.PrivateKey)
	SetCommissionSign(privateKey common.PrivateKey)
	SetHash(e_common.Hash)

	// verifiers
	ValidTransactionHash() bool
	ValidSign() bool
	ValidLastHash(fromAccountState AccountState) bool
	ValidDeviceKey(fromAccountState AccountState) bool
	ValidMaxGas() bool
	ValidMaxGasPrice(currentGasPrice uint64) bool
	ValidAmount(fromAccountState AccountState, currentGasPrice uint64) bool
	ValidPendingUse(fromAccountState AccountState) bool
	ValidDeploySmartContractToAccount(fromAccountState AccountState) bool
	ValidCallSmartContractToAccount(toAccountState AccountState) bool
}

type CallData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	// geter
	Input() []byte
}

type DeployData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	// getter
	Code() []byte
	StorageAddress() e_common.Address
}

type OpenStateChannelData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	// geter
	ValidatorAddresses() []e_common.Address
}

type CommitAccountStateChannelData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	// geter
	Address() e_common.Address
	CloseSmartContract() bool
	Amount() *big.Int
}

type VerifyTransactionSignResult interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	TransactionHash() e_common.Hash
	Valid() bool
	Proto() *pb.VerifyTransactionSignResult
}

type VerifyTransactionSignRequest interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	TransactionHash() e_common.Hash
	SenderPublicKey() common.PublicKey
	SenderSign() common.Sign
	Proto() *pb.VerifyTransactionSignRequest
}

type UpdateStorageHostData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	// geter
	StorageHost() string
	StorageAddress() e_common.Address
}

type TransactionError interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	String() string
}

type FromNodeTransactionsResult interface {
	ValidTransactionHashes() []e_common.Hash
	TransactionErrors() map[e_common.Hash]int64
	BlockNumber() uint64
}

type ToNodeTransactionsResult interface {
	ValidTransactionHashes() []e_common.Hash
	BlockNumber() uint64
}

type ExecuteSCTransactions interface {
	Transactions() []Transaction
	BlockNumber() uint64
	GroupId() uint64
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

type TransactionsFromLeader interface {
	Transactions() []Transaction
	BlockNumber() uint64
	TimeStamp() uint64
	AggSign() []byte
	IsValidSign() bool
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

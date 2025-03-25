package types

import (
	"math/big"

	e_common "github.com/ethereum/go-ethereum/common"
	e_types "github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/protobuf/reflect/protoreflect"

	"gomail/pkg/common"
	pb "gomail/pkg/proto"
	//e_types "github.com/ethereum/go-ethereum/core/types"
)

type Transaction interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	String() string

	// getter
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
	GetNonce() uint64
	GetChainID() uint64
	GetEthTx() *e_types.Transaction
	ValidEthSign() bool
	// setter
	SetSign(privateKey common.PrivateKey)
	SetSignBytes(bytes []byte)
	SetNonce(uint64)
	SetFromAddress(address e_common.Address)
	SetToAddress(address e_common.Address)
	CopyTransaction() Transaction

	SetCommissionSign(privateKey common.PrivateKey)
	UpdateRelatedAddresses(relatedAddresses [][]byte)
	SetReadOnly(readOnly bool)
	GetReadOnly() bool
	// verifiers
	ValidTx0(fromAccountState AccountState, chainId string) (bool, int64)
	ValidChainID(chainId uint64) bool
	ValidSign() bool
	ValidLastHash(fromAccountState AccountState) bool
	ValidDeviceKey(fromAccountState AccountState) bool
	ValidMaxGas() bool
	ValidMaxGasPrice(currentGasPrice uint64) bool
	ValidAmount(fromAccountState AccountState, currentGasPrice uint64) bool
	ValidAmountSpend(fromAccountState AccountState, spendAmount *big.Int) bool
	ValidPendingUse(fromAccountState AccountState) bool
	ValidDeploySmartContractToAccount(fromAccountState AccountState) bool
	ValidCallSmartContractToAccount(toAccountState AccountState) bool
	ValidDeployData() bool
	ValidCallData() bool
	//
	RawSignatureValues() (v, r, s *big.Int)
	SetSignatureValues(chainID, v, r, s *big.Int)
	IsDeployContract() bool
	IsCallContract() bool
	IsRegularTransaction() bool
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
	Valid() bool
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

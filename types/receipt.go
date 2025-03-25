package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	e_common "github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "gomail/pkg/proto"
)

type Receipt interface {
	// general
	FromProto(proto *pb.Receipt)
	Proto() protoreflect.ProtoMessage
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	String() string
	Json() ([]byte, error)
	MarshalReceiptToMap() (map[string]interface{}, error)
	// getter
	TransactionHash() e_common.Hash
	FromAddress() e_common.Address
	ToAddress() e_common.Address
	Amount() *big.Int
	Status() pb.RECEIPT_STATUS
	GasUsed() uint64
	GastFee() uint64

	Return() []byte
	EventLogs() []*pb.EventLog
	SetReturn([]byte)

	TransactionIndex() uint64
	BlockNumber() uint64
	BlockHash() common.Hash
	// setter
	UpdateExecuteResult(
		status pb.RECEIPT_STATUS,
		output []byte,
		exception pb.EXCEPTION,
		gasUsed uint64,
		eventLogs []EventLog,
	)
	UpdateInfoTx(
		transactionIndex uint64,
		blockNumber uint64,
	)
}

type Receipts interface {
	// getter
	ReceiptsRoot() (e_common.Hash, error)
	Commit() (e_common.Hash, error)
	IntermediateRoot() (e_common.Hash, error)
	ReceiptsMap() map[e_common.Hash]Receipt
	GasUsed() uint64
	Discard() error
	// setter
	AddReceipt(Receipt) error
	UpdateExecuteResultToReceipt(
		e_common.Hash,
		pb.RECEIPT_STATUS,
		[]byte,
		pb.EXCEPTION,
		uint64,
		[]EventLog,
	) error
}

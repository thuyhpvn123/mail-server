package types

import (
	"math/big"

	e_common "github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "gomail/mtn/proto"
)

type Receipt interface {
	// general
	FromProto(proto *pb.Receipt)
	Proto() protoreflect.ProtoMessage
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	String() string
	Json() ([]byte, error)

	// getter
	TransactionHash() e_common.Hash
	FromAddress() e_common.Address
	ToAddress() e_common.Address
	Amount() *big.Int
	Status() pb.RECEIPT_STATUS
	Action() pb.ACTION
	GasUsed() uint64
	GastFee() uint64

	Return() []byte
	EventLogs() []*pb.EventLog
	// setter
	UpdateExecuteResult(
		status pb.RECEIPT_STATUS,
		output []byte,
		exception pb.EXCEPTION,
		gasUsed uint64,
		eventLogs []EventLog,
	)
}

type Receipts interface {
	// getter
	ReceiptsRoot() (e_common.Hash, error)
	ReceiptsMap() map[e_common.Hash]Receipt
	GasUsed() uint64

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

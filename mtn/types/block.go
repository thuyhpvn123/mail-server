package types

import (
	e_common "github.com/ethereum/go-ethereum/common"

	pb "gomail/mtn/proto"
)

type Block interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() *pb.Block
	FromProto(*pb.Block)

	Header() BlockHeader
	Transactions() []Transaction
	ExecuteSCResults() []ExecuteSCResult
}

type BlockHeader interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() *pb.BlockHeader
	FromProto(*pb.BlockHeader)

	Hash() e_common.Hash
	LastBlockHash() e_common.Hash
	BlockNumber() uint64
	AccountStatesRoot() e_common.Hash
	ReceiptRoot() e_common.Hash
	TimeStamp() uint64
	LeaderAddress() e_common.Address
	AggregateSignature() []byte
	String() string
}

type ConfirmedBlockData interface {
	Header() BlockHeader
	Receipts() []Receipt
	BranchStateRoot() e_common.Hash
	ValidatorSigns() map[e_common.Address][]byte
	SetBranchStateRoot(e_common.Hash)
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() *pb.ConfirmedBlockData
	FromProto(*pb.ConfirmedBlockData)
}

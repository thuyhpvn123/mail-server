package block

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/proto"

	pb "gomail/mtn/proto"
)

type BlockHeader struct {
	lastBlockHash      common.Hash
	blockNumber        uint64
	accountStatesRoot  common.Hash
	receiptRoot        common.Hash
	leaderAddress      common.Address
	timeStamp          uint64
	aggregateSignature []byte
}

func NewBlockHeader(
	lastBlockHash common.Hash,
	blockNumber uint64,
	accountStatesRoot common.Hash,
	receiptRoot common.Hash,
	leaderAddress common.Address,
	timeStamp uint64,
) *BlockHeader {
	return &BlockHeader{
		lastBlockHash:     lastBlockHash,
		blockNumber:       blockNumber,
		accountStatesRoot: accountStatesRoot,
		receiptRoot:       receiptRoot,
		leaderAddress:     leaderAddress,
		timeStamp:         timeStamp,
	}
}

func (b *BlockHeader) LastBlockHash() common.Hash {
	return b.lastBlockHash
}

func (b *BlockHeader) BlockNumber() uint64 {
	return b.blockNumber
}

func (b *BlockHeader) AccountStatesRoot() common.Hash {
	return b.accountStatesRoot
}

func (b *BlockHeader) ReceiptRoot() common.Hash {
	return b.receiptRoot
}

func (b *BlockHeader) LeaderAddress() common.Address {
	return b.leaderAddress
}

func (b *BlockHeader) TimeStamp() uint64 {
	return b.timeStamp
}

func (b *BlockHeader) AggregateSignature() []byte {
	return b.aggregateSignature
}

func (b *BlockHeader) Marshal() ([]byte, error) {
	return proto.Marshal(b.Proto())
}

func (b *BlockHeader) Unmarshal(bData []byte) error {
	pbBlockHeader := &pb.BlockHeader{}
	if err := proto.Unmarshal(bData, pbBlockHeader); err != nil {
		return err
	}
	b.FromProto(pbBlockHeader)
	return nil
}

func (b *BlockHeader) Hash() common.Hash {
	bData, _ := b.Marshal()
	return crypto.Keccak256Hash(bData)
}

func (b *BlockHeader) Proto() *pb.BlockHeader {
	return &pb.BlockHeader{
		LastBlockHash:     b.lastBlockHash.Bytes(),
		BlockNumber:       b.blockNumber,
		AccountStatesRoot: b.accountStatesRoot.Bytes(),
		ReceiptRoot:       b.receiptRoot.Bytes(),
		LeaderAddress:     b.leaderAddress.Bytes(),
		TimeStamp:         b.timeStamp,
	}
}

func (b *BlockHeader) FromProto(pbBlockHeader *pb.BlockHeader) {
	b.lastBlockHash = common.BytesToHash(pbBlockHeader.LastBlockHash)
	b.blockNumber = pbBlockHeader.BlockNumber
	b.accountStatesRoot = common.BytesToHash(pbBlockHeader.AccountStatesRoot)
	b.receiptRoot = common.BytesToHash(pbBlockHeader.ReceiptRoot)
	b.leaderAddress = common.BytesToAddress(pbBlockHeader.LeaderAddress)
	b.timeStamp = pbBlockHeader.TimeStamp
	b.aggregateSignature = pbBlockHeader.AggregateSignature
}

func (b *BlockHeader) String() string {
	str := fmt.Sprintf(`
BlockHeader{
  LastBlockHash: %v,
  BlockNumber: %d,
  AccountStatesRoot: %v,
  ReceiptRoot: %v,
  LeaderAddress: %v,
  TimeStamp: %d,
  AggregateSignature: %v
}
`, b.lastBlockHash, b.blockNumber, b.accountStatesRoot, b.receiptRoot, b.leaderAddress, b.timeStamp, b.aggregateSignature)
	return str
}

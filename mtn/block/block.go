package block

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	pb "gomail/mtn/proto"
	"gomail/mtn/smart_contract"
	"gomail/mtn/transaction"
	"gomail/mtn/types"
)

type Block struct {
	header           types.BlockHeader
	transactions     []types.Transaction
	executeSCResults []types.ExecuteSCResult
}

func NewBlock(
	header types.BlockHeader,
	transactions []types.Transaction,
	executeSCResults []types.ExecuteSCResult,
) *Block {
	return &Block{
		header:           header,
		transactions:     transactions,
		executeSCResults: executeSCResults,
	}
}

func (b *Block) Header() types.BlockHeader {
	return b.header
}

func (b *Block) Transactions() []types.Transaction {
	return b.transactions
}

func (b *Block) ExecuteSCResults() []types.ExecuteSCResult {
	return b.executeSCResults
}

func (b *Block) Proto() *pb.Block {
	return &pb.Block{
		Header:           b.header.Proto(),
		Transactions:     transaction.TransactionsToProto(b.transactions),
		ExecuteSCResults: smart_contract.ExecuteSCResultsToProto(b.executeSCResults),
	}
}

func (b *Block) FromProto(pbBlock *pb.Block) {
	b.header = &BlockHeader{}
	b.header.FromProto(pbBlock.Header)
	b.transactions = transaction.TransactionsFromProto(pbBlock.Transactions)
	b.executeSCResults = smart_contract.ExecuteSCResultsFromProto(pbBlock.ExecuteSCResults)
}

func (b *Block) Marshal() ([]byte, error) {
	return proto.Marshal(b.Proto())
}

func (b *Block) Unmarshal(bData []byte) error {
	pbBlock := &pb.Block{}
	err := proto.Unmarshal(bData, pbBlock)
	if err != nil {
		return err
	}
	b.FromProto(pbBlock)
	return nil
}

func (b *Block) String() string {
	return fmt.Sprintf("Block: Header: %v", b.header.String())
}

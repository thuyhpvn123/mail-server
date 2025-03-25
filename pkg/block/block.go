package block

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	pb "gomail/pkg/proto"
	"gomail/pkg/shard_storage"
	"gomail/pkg/smart_contract"
	"gomail/types"

	"github.com/ethereum/go-ethereum/common"
)

const (
	maxBlocksPerShard = 1000
	lineByte          = 66
)

type Block struct {
	header           types.BlockHeader
	transactions     []common.Hash
	executeSCResults []types.ExecuteSCResult
}

func NewBlock(
	header types.BlockHeader,
	transactions []common.Hash,
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

func (b *Block) Transactions() []common.Hash {
	return b.transactions
}

func (b *Block) ExecuteSCResults() []types.ExecuteSCResult {
	return b.executeSCResults
}

func (b *Block) Proto() *pb.Block {
	transactionsHash := make([][]byte, len(b.transactions))
	for i, txHash := range b.transactions {
		transactionsHash[i] = txHash.Bytes()
	}
	return &pb.Block{
		Header:       b.header.Proto(),
		Transactions: transactionsHash,
	}
}

func (b *Block) FromProto(pbBlock *pb.Block) {
	b.header = &BlockHeader{}
	b.header.FromProto(pbBlock.Header)
	transactions := make([]common.Hash, len(pbBlock.GetTransactions()))
	for i, txBytes := range pbBlock.GetTransactions() {
		transactions[i] = common.BytesToHash(txBytes)
	}
	b.transactions = transactions
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

// Lưu block vào file shard, đảm bảo dòng cuối cùng khớp
func SaveBlock(shardDir string, blockNumber int, blockHash string) error {
	// Tạo một đối tượng ShardStorage
	shardStorage, err := shard_storage.NewShardStorage(maxBlocksPerShard, shardDir, lineByte)
	if err != nil {
		return fmt.Errorf("failed to create shard_storage: %w", err)
	}
	shardStorage.SetIndexValue(blockNumber, blockHash)
	return nil
}

// Tìm blockHash dựa trên blockNumber
func FindBlockHashByBlockNumber(shardDir string, blockNumber int) (string, error) {
	shardStorage, err := shard_storage.NewShardStorage(maxBlocksPerShard, shardDir, lineByte)
	if err != nil {
		return "", fmt.Errorf("failed to create shard_storage: %w", err)
	}
	return shardStorage.FindValueByIndex(blockNumber)
}

package block

import (
	"sync"

	"gomail/pkg/storage"
	"gomail/types"

	"github.com/ethereum/go-ethereum/common"
)

var lastBlockHashKey common.Hash

func init() {
	lastBlockHashKey = common.Hash{}
}

type BlockDatabase struct {
	db          *storage.ShardelDB
	dirtyBlocks sync.Map // Use sync.Map instead of map[common.Hash]types.Block
}

func NewBlockDatabase(
	db *storage.ShardelDB,
) *BlockDatabase {
	return &BlockDatabase{
		db:          db,
		dirtyBlocks: sync.Map{}, // Initialize the sync.Map

	}
}

func (blockDatabase *BlockDatabase) SaveBlock(block types.Block) error {
	// Encode block to bytes
	blockBytes, err := block.Marshal()
	if err != nil {
		return err
	}

	// Save block to database
	blockHash := block.Header().Hash()
	if err := blockDatabase.db.Put(blockHash.Bytes(), blockBytes); err != nil {
		return err
	}

	return nil
}

// SaveLastBlock saves the last block's hash to the database.
func (blockDatabase *BlockDatabase) SaveLastBlock(block types.Block) error {
	// Encode block to bytes
	blockBytes, err := block.Marshal()
	if err != nil {
		return err
	}

	if err := blockDatabase.db.Put(lastBlockHashKey.Bytes(), blockBytes); err != nil {
		return err
	}

	return nil
}

// func (blockDatabase *BlockDatabase) GetBlockByNumber(blockNumber uint64) (types.Block, error) {
// 	shardDir := filepath.Join(blockDatabase.blockStorageDir, fmt.Sprintf("%d", blockDatabase.shardID))
// 	blockHash, err := FindBlockHashByBlockNumber(shardDir, int(blockNumber))
// 	if err != nil {
// 		return nil, err
// 	}

// 	blockFilePath := filepath.Join(blockDatabase.blockStorageDir, fmt.Sprintf("%d", blockDatabase.shardID), fmt.Sprintf("%s.dat", blockHash))
// 	block, err := LoadBlockFromFile(blockFilePath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return block, nil
// }

func (blockDatabase *BlockDatabase) GetBlockByHash(blockHash common.Hash) (types.Block, error) {
	// Check if the block is in dirtyBlocks
	if block, ok := blockDatabase.dirtyBlocks.Load(blockHash); ok {
		return block.(types.Block), nil
	}

	// Try to load the block from the database
	blockBytes, err := blockDatabase.db.Get(blockHash.Bytes())
	if err != nil {
		return nil, err
	}
	block := &Block{}

	err = block.Unmarshal(blockBytes)

	if err != nil {
		return nil, err
	}
	// Save the block to dirtyBlocks
	blockDatabase.dirtyBlocks.Store(blockHash, block)

	return block, nil
}

// GetLastBlock retrieves the last block from the database.
func (blockDatabase *BlockDatabase) GetLastBlock() (types.Block, error) {
	return blockDatabase.GetBlockByHash(lastBlockHashKey)
}

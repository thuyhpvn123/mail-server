package mock

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"gomail/cmd/validator/pkg/env"
	v_types "gomail/cmd/validator/types"
	"gomail/types"
)

// Chain struct holds the chain data
type Chain struct {
	sync.RWMutex
	lastBlock             types.Block
	maximumInMemoryBlocks int

	mHashBlocks   map[common.Hash]types.Block
	mNumberBlocks map[uint64]types.Block
}

// NewChain creates a new Chain instance
func NewChain() *Chain {
	return &Chain{
		mHashBlocks:   make(map[common.Hash]types.Block),
		mNumberBlocks: make(map[uint64]types.Block),
	}
}

// AddBlock adds a new block to the chain
func (c *Chain) AddBlock(block types.Block) {
	c.Lock()
	defer c.Unlock()

	// Add block to memory
	c.lastBlock = block
	c.mHashBlocks[block.Header().Hash()] = block
	c.mNumberBlocks[block.Header().BlockNumber()] = block
}

// LastBlock returns the last block
func (c *Chain) LastBlock() types.Block {
	c.RLock()
	defer c.RUnlock()
	return c.lastBlock
}

// BlockByHash returns the block by its hash
func (c *Chain) BlockByHash(hash common.Hash) types.Block {
	c.RLock()
	defer c.RUnlock()
	return c.mHashBlocks[hash]
}

// BlockByNumber returns the block by its number
func (c *Chain) BlockByNumber(number uint64) types.Block {
	c.RLock()
	v, ok := c.mNumberBlocks[number]
	c.RUnlock()
	if !ok {
		return nil
	}
	return v
}

func (c *Chain) Env() v_types.Env {
	return &env.DevEnv{}
}

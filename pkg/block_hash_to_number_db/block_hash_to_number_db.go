package block_hash_to_number_db

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"gomail/pkg/state" // Import types
	"gomail/pkg/storage"
	p_trie "gomail/pkg/trie"
)

type BlockHashToNumberDB struct {
	trie              *p_trie.MerklePatriciaTrie
	originRootHash    common.Hash
	db                storage.Storage
	dirtyTransactions map[common.Hash]state.BlockHashToNumber
	mu                sync.Mutex
}

func NewBlockHashToNumberDB(
	trie *p_trie.MerklePatriciaTrie,
	db storage.Storage,
) *BlockHashToNumberDB {
	return &BlockHashToNumberDB{
		trie:              trie,
		db:                db,
		originRootHash:    trie.Hash(),
		dirtyTransactions: make(map[common.Hash]state.BlockHashToNumber),
	}
}

func (db *BlockHashToNumberDB) Get(hash common.Hash) (state.BlockHashToNumber, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// get from dirty
	tx, ok := db.dirtyTransactions[hash]
	if ok {
		return tx, nil
	}
	// if not exist in dirty then get from trie
	bData, _ := db.trie.Get(hash.Bytes())
	if len(bData) == 0 {
		return state.BlockHashToNumber{}, errors.New("transaction not found") // if not exist in trie create new once
	}
	// exist in trie, unmarshal
	var txData state.BlockHashToNumber // Declare variable outside the if statement
	err := txData.Unmarshal(bData)
	if err != nil {
		return state.BlockHashToNumber{}, err
	}
	return txData, nil

}

func (db *BlockHashToNumberDB) Set(tx state.BlockHashToNumber) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.setDirtyTransaction(tx)

}

func (db *BlockHashToNumberDB) Commit() (common.Hash, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	rootHash, err := db.IntermediateRoot()
	if err != nil {
		return common.Hash{}, err
	}

	hash, nodeSet, _, err := db.trie.Commit(true)
	if err != nil {
		return common.Hash{}, err
	}

	if rootHash != hash {
		return common.Hash{}, errors.New("inconsistent root hash")
	}

	if nodeSet != nil {
		batch := [][2][]byte{}
		for _, node := range nodeSet.Nodes {

			batch = append(batch, [2][]byte{node.Hash.Bytes(), node.Blob})
		}

		err := db.db.BatchPut(batch)
		if err != nil {
			return common.Hash{}, err
		}
	} else {
	}

	db.dirtyTransactions = make(map[common.Hash]state.BlockHashToNumber)
	db.trie, err = p_trie.New(hash, db.db)
	if err != nil {
		return common.Hash{}, err
	}
	db.originRootHash = hash

	return hash, nil
}

func (db *BlockHashToNumberDB) IntermediateRoot() (common.Hash, error) {

	for hash, tx := range db.dirtyTransactions {
		b, err := tx.Marshal()
		if err != nil {
			return common.Hash{}, err
		}

		err = db.trie.Update(hash.Bytes(), b)
		if err != nil {
			return common.Hash{}, err
		}
	}
	return db.trie.Hash(), nil
}

func (db *BlockHashToNumberDB) setDirtyTransaction(tx state.BlockHashToNumber) {
	db.dirtyTransactions[tx.Hash] = tx

}

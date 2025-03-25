package transaction_blocknumber_db

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"

	"gomail/pkg/state" // Import types
	"gomail/pkg/storage"
	p_trie "gomail/pkg/trie"
)

type TransactionBlockNumberDB struct {
	trie              *p_trie.MerklePatriciaTrie
	originRootHash    common.Hash
	db                storage.Storage
	dirtyTransactions map[common.Hash]state.TransactionBlockNumber
}

func NewTransactionBlockNumberDB(
	trie *p_trie.MerklePatriciaTrie,
	db storage.Storage,
) *TransactionBlockNumberDB {
	return &TransactionBlockNumberDB{
		trie:              trie,
		db:                db,
		originRootHash:    trie.Hash(),
		dirtyTransactions: make(map[common.Hash]state.TransactionBlockNumber),
	}
}

func (db *TransactionBlockNumberDB) GetTransaction(hash common.Hash) (state.TransactionBlockNumber, error) {

	// get from dirty
	tx, ok := db.dirtyTransactions[hash]
	if ok {
		return tx, nil
	}
	// if not exist in dirty then get from trie
	bData, _ := db.trie.Get(hash.Bytes())
	if len(bData) == 0 {
		return state.TransactionBlockNumber{}, errors.New("transaction not found") // if not exist in trie create new once
	}
	// exist in trie, unmarshal
	var txData state.TransactionBlockNumber // Declare variable outside the if statement
	err := txData.Unmarshal(bData)
	if err != nil {
		return state.TransactionBlockNumber{}, err
	}
	return txData, nil

}

func (db *TransactionBlockNumberDB) SetTransaction(tx state.TransactionBlockNumber) {

	db.setDirtyTransaction(tx)

}

func (db *TransactionBlockNumberDB) Commit() (common.Hash, error) {

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

	db.dirtyTransactions = make(map[common.Hash]state.TransactionBlockNumber)
	db.trie, err = p_trie.New(hash, db.db, true)
	if err != nil {
		return common.Hash{}, err
	}
	db.originRootHash = hash

	return hash, nil
}

func (db *TransactionBlockNumberDB) IntermediateRoot() (common.Hash, error) {

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

func (db *TransactionBlockNumberDB) setDirtyTransaction(tx state.TransactionBlockNumber) {
	db.dirtyTransactions[tx.Hash] = tx

}

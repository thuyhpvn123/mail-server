package transaction_state_db

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"

	// Import types
	"gomail/pkg/logger"
	"gomail/pkg/storage"
	"gomail/pkg/transaction"
	p_trie "gomail/pkg/trie"
	"gomail/types"
)

type TransactionStateDB struct {
	trie              *p_trie.MerklePatriciaTrie
	originRootHash    common.Hash
	db                storage.Storage
	dirtyTransactions map[common.Hash]types.Transaction
}

func NewTransactionStateDB(
	trie *p_trie.MerklePatriciaTrie,
	db storage.Storage,
) *TransactionStateDB {
	return &TransactionStateDB{
		trie:              trie,
		db:                db,
		originRootHash:    trie.Hash(),
		dirtyTransactions: make(map[common.Hash]types.Transaction),
	}
}
func (db *TransactionStateDB) GetTransaction(hash common.Hash) (types.Transaction, error) {
	tx, ok := db.dirtyTransactions[hash]
	if ok {
		return tx, nil // Trả về con trỏ thay vì giá trị trực tiếp
	}

	// if not exist in dirty then get from trie
	bData, _ := db.trie.Get(hash.Bytes())
	if len(bData) == 0 {
		logger.Error("transaction not found")
		return nil, errors.New("transaction not found") // Trả về nil hợp lệ cho con trỏ
	}

	// exist in trie, unmarshal
	txData := &transaction.Transaction{} // Bạn cần implement hàm này để tạo transaction phù hợp

	err := txData.Unmarshal(bData)

	if err != nil {
		logger.Error("err: ", err)
		return nil, err
	}

	return txData, nil // Trả về con trỏ đến struct
}

func (db *TransactionStateDB) SetTransaction(tx types.Transaction) {
	// fileLoggerPT, _ := loggerfile.NewFileLogger("txStatedb.log")
	// fileLoggerPT.Info("Set: ", tx)
	db.setDirtyTransaction(tx)

}

func (db *TransactionStateDB) Commit() (common.Hash, error) {

	rootHash, err := db.IntermediateRoot()
	if err != nil {
		return common.Hash{}, err
	}
	trieCommit := db.trie.Copy()

	hash, nodeSet, _, err := trieCommit.Commit(true)
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
		logger.Info("tx BatchPut")
	} else {
	}

	// db.dirtyTransactions = make(map[common.Hash]types.Transaction)
	db.trie, err = p_trie.New(hash, db.db, true)
	if err != nil {
		// panic(fmt.Sprintf("set db.trie is nil vvv %v", err)) // Sử dụng fmt.Sprintf để nối chuỗi

		return common.Hash{}, err
	}
	db.originRootHash = hash

	return hash, nil
}

func (db *TransactionStateDB) IntermediateRoot() (common.Hash, error) {
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
	db.dirtyTransactions = make(map[common.Hash]types.Transaction)

	return db.trie.Hash(), nil
}

func (db *TransactionStateDB) setDirtyTransaction(tx types.Transaction) {
	db.dirtyTransactions[tx.Hash()] = tx

}

func (db *TransactionStateDB) Discard() (err error) {
	db.dirtyTransactions = make(map[common.Hash]types.Transaction)
	db.trie, err = p_trie.New(db.originRootHash, db.db, true)
	return err
}

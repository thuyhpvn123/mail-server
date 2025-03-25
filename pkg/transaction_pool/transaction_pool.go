package transaction_pool

import (
	"strconv"
	"sync"

	blst "gomail/pkg/bls/blst/bindings/go"
	"gomail/pkg/logger"
	"gomail/types"
)

type TransactionPool struct {
	transactions    []types.Transaction
	transactionKeys map[string]bool // Use a map to store transaction keys for quick existence checks

	aggSign *blst.P2Aggregate
	mutex   sync.Mutex
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		transactions:    make([]types.Transaction, 0), // Initialize the transactions slice
		transactionKeys: make(map[string]bool),        // Initialize the transactionKeys map
		aggSign:         new(blst.P2Aggregate)}
}

func (tp *TransactionPool) CountTransactions() int {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	return len(tp.transactions)
}

func (tp *TransactionPool) AddTransaction(tx types.Transaction) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	tp.addTransaction(tx)
}

func (tp *TransactionPool) AddTransactions(txs []types.Transaction) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	for _, tx := range txs {
		tp.addTransaction(tx)
	}
}

func (tp *TransactionPool) addTransaction(tx types.Transaction) {
	key := tx.FromAddress().String() + strconv.FormatUint(tx.GetNonce(), 10) // Combine FromAddress and Nonce for a unique key

	// Check if the transaction already exists in the pool
	if _, exists := tp.transactionKeys[key]; exists {
		logger.Info("Transaction already exists in pool, skipping", "key", key)
		return // Skip adding the transaction if it already exists
	}

	// Add the transaction to the pool
	tp.transactions = append(tp.transactions, tx)
	tp.transactionKeys[key] = true // Add the key to the transactionKeys map

	// Aggregate the signature
	p := new(blst.P2Affine)
	err := p.Uncompress(tx.Sign().Bytes())
	if err != nil {
		logger.Error("Error uncompressing signature ", err)
		return // Handle the error, possibly remove the transaction from the pool
	}
	tp.aggSign.Add(p, false)
}

// TransactionsWithAggSign returns transactions and aggregate sign
// and clear transactions
func (tp *TransactionPool) TransactionsWithAggSign() ([]types.Transaction, []byte) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	tx := tp.transactions
	sign := tp.aggSign.ToAffine().Compress()

	// Clear the transaction pool
	tp.transactions = make([]types.Transaction, 0)
	tp.transactionKeys = make(map[string]bool) // Clear the transaction keys as well
	tp.aggSign = new(blst.P2Aggregate)

	return tx, sign
}

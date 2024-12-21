package transaction_pool

import (
	"sync"

	blst "gomail/mtn/bls/blst/bindings/go"
	"gomail/mtn/types"
)

type TransactionPool struct {
	transactions []types.Transaction
	aggSign      *blst.P2Aggregate
	mutex        sync.Mutex
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		aggSign: new(blst.P2Aggregate),
	}
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

// TransactionsWithAggSign returns transactions and aggregate sign
// and clear transactions
func (tp *TransactionPool) TransactionsWithAggSign() ([]types.Transaction, []byte) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	tx := tp.transactions
	sign := tp.aggSign.ToAffine().Compress()
	tp.transactions = make([]types.Transaction, 0)
	tp.aggSign = new(blst.P2Aggregate)
	return tx, sign
}

func (tp *TransactionPool) addTransaction(tx types.Transaction) {
	tp.transactions = append(tp.transactions, tx)
	p := new(blst.P2Affine)
	p.Uncompress(tx.Sign().Bytes())
	tp.aggSign.Add(p, false)
}

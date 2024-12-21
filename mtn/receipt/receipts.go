package receipt

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"

	pb "gomail/mtn/proto"
	"gomail/mtn/storage"
	"gomail/mtn/trie"
	"gomail/mtn/types"
)

var ErrorReceiptNotFound = errors.New("receipt not found")

type Receipts struct {
	trie     *trie.MerklePatriciaTrie
	receipts map[common.Hash]types.Receipt
}

func NewReceipts() types.Receipts {
	trie, _ := trie.New(trie.EmptyRootHash, storage.NewMemoryDb())
	return &Receipts{
		trie:     trie,
		receipts: make(map[common.Hash]types.Receipt),
	}
}

func (r *Receipts) ReceiptsRoot() (common.Hash, error) {
	hash := r.trie.Hash()
	return hash, nil
}

func (r *Receipts) AddReceipt(receipt types.Receipt) error {
	b, err := receipt.Marshal()
	if err != nil {
		return err
	}
	r.receipts[receipt.TransactionHash()] = receipt
	r.trie.Update(receipt.TransactionHash().Bytes(), b)
	return nil
}

func (r *Receipts) ReceiptsMap() map[common.Hash]types.Receipt {
	return r.receipts
}

func (r *Receipts) UpdateExecuteResultToReceipt(
	hash common.Hash,
	status pb.RECEIPT_STATUS,
	returnValue []byte,
	exception pb.EXCEPTION,
	gasUsed uint64,
	eventLogs []types.EventLog,
) error {
	receipt := r.receipts[hash]
	if receipt == nil {
		return ErrorReceiptNotFound
	}
	receipt.UpdateExecuteResult(
		status,
		returnValue,
		exception,
		gasUsed,
		eventLogs,
	)
	err := r.AddReceipt(receipt)
	return err
}

func (r *Receipts) GasUsed() uint64 {
	gasUsed := uint64(0)
	if r.receipts == nil {
		return gasUsed
	} else {
		for _, v := range r.receipts {
			gasUsed += v.GasUsed()
		}
	}
	return gasUsed
}

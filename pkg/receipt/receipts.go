package receipt

import (
	"errors"
	"sort"

	"github.com/ethereum/go-ethereum/common"

	pb "gomail/pkg/proto"
	"gomail/pkg/storage"
	"gomail/pkg/trie"
	"gomail/types"
)

var ErrorReceiptNotFound = errors.New("receipt not found")

type Receipts struct {
	trie           *trie.MerklePatriciaTrie
	db             storage.Storage
	originRootHash common.Hash
	dirtyReceipts  map[common.Hash]types.Receipt
}

func NewReceipts(db storage.Storage) types.Receipts {
	trie, err := trie.New(trie.EmptyRootHash, db, true)
	if err != nil {
		panic(err) // Hoặc có thể trả về lỗi thay vì panic
	}
	return &Receipts{
		trie:           trie,
		db:             db,
		originRootHash: trie.Hash(),
		dirtyReceipts:  make(map[common.Hash]types.Receipt),
	}
}

func (r *Receipts) ReceiptsRoot() (common.Hash, error) {
	return r.trie.Hash(), nil
}

func (r *Receipts) AddReceipt(receipt types.Receipt) error {
	b, err := receipt.Marshal()
	if err != nil {
		return err
	}
	err = r.trie.Update(receipt.TransactionHash().Bytes(), b)
	if err != nil {
		return err
	}
	r.dirtyReceipts[receipt.TransactionHash()] = receipt // Cập nhật map để tránh lỗi truy xuất
	r.setDirtyReceipt(receipt)
	r.db.Put(receipt.TransactionHash().Bytes(), b)

	return nil
}

func (r *Receipts) ReceiptsMap() map[common.Hash]types.Receipt {
	return r.dirtyReceipts
}

func (r *Receipts) UpdateExecuteResultToReceipt(
	hash common.Hash,
	status pb.RECEIPT_STATUS,
	returnValue []byte,
	exception pb.EXCEPTION,
	gasUsed uint64,
	eventLogs []types.EventLog,
) error {
	receipt, exists := r.dirtyReceipts[hash]
	if !exists {
		return ErrorReceiptNotFound
	}
	receipt.UpdateExecuteResult(
		status,
		returnValue,
		exception,
		gasUsed,
		eventLogs,
	)
	return r.AddReceipt(receipt)
}

func (r *Receipts) GasUsed() uint64 {
	gasUsed := uint64(0)
	for _, v := range r.dirtyReceipts {
		gasUsed += v.GasUsed()
	}
	return gasUsed
}

func (r *Receipts) GetReceipt(hash common.Hash) (types.Receipt, error) {
	data, err := r.trie.Get(hash.Bytes())
	if err != nil {
		return nil, ErrorReceiptNotFound
	}
	var receipt types.Receipt
	err = receipt.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (r *Receipts) Discard() error {
	r.dirtyReceipts = make(map[common.Hash]types.Receipt)
	trie, err := trie.New(r.originRootHash, r.db, true)
	if err != nil {
		return err
	}
	r.trie = trie
	return nil
}

func (r *Receipts) Commit() (common.Hash, error) {
	rootHash, err := r.IntermediateRoot()
	if err != nil {
		return common.Hash{}, err
	}
	hash, _, _, err := r.trie.Commit(true)
	if err != nil {
		return common.Hash{}, err
	}
	if rootHash != hash {
		return common.Hash{}, errors.New("root hash mismatch")
	}
	r.dirtyReceipts = make(map[common.Hash]types.Receipt)
	trie, err := trie.New(hash, r.db, true)
	if err != nil {
		return common.Hash{}, err
	}
	r.trie = trie
	r.originRootHash = hash
	return hash, nil
}

func (r *Receipts) IntermediateRoot() (common.Hash, error) {
	receiptItems := make([]types.Receipt, 0, len(r.dirtyReceipts))
	for _, receipt := range r.dirtyReceipts {
		receiptItems = append(receiptItems, receipt)
	}

	sort.Slice(receiptItems, func(i, j int) bool {
		return receiptItems[i].TransactionHash().Cmp(receiptItems[j].TransactionHash()) < 0
	})

	for _, receiptItem := range receiptItems {
		b, err := receiptItem.Marshal()
		if err != nil {
			return common.Hash{}, err
		}
		err = r.trie.Update(receiptItem.TransactionHash().Bytes(), b)
		if err != nil {
			return common.Hash{}, err
		}
	}
	return r.trie.Hash(), nil
}

func (r *Receipts) setDirtyReceipt(receipt types.Receipt) {
	r.dirtyReceipts[receipt.TransactionHash()] = receipt
}

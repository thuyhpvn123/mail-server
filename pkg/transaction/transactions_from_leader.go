package transaction

import (
	"google.golang.org/protobuf/proto"

	"gomail/pkg/bls"
	pb "gomail/pkg/proto"
	"gomail/types"
)

type TransactionsFromLeader struct {
	transactions []types.Transaction
	blockNumber  uint64
	aggSign      []byte
	timeStamp    uint64
}

func NewTransactionsFromLeader(
	transactions []types.Transaction,
	blockNumber uint64,
	aggSign []byte,
	timeStamp uint64,
) *TransactionsFromLeader {
	return &TransactionsFromLeader{
		transactions: transactions,
		blockNumber:  blockNumber,
		aggSign:      aggSign,
		timeStamp:    timeStamp,
	}
}

func (t *TransactionsFromLeader) Transactions() []types.Transaction {
	return t.transactions
}

func (t *TransactionsFromLeader) BlockNumber() uint64 {
	return t.blockNumber
}

func (t *TransactionsFromLeader) AggSign() []byte {
	return t.aggSign
}

func (t *TransactionsFromLeader) TimeStamp() uint64 {
	return t.timeStamp
}

func (t *TransactionsFromLeader) Marshal() ([]byte, error) {
	return proto.Marshal(t.Proto())
}

func (t *TransactionsFromLeader) Unmarshal(b []byte) error {
	pbData := &pb.TransactionsFromLeader{}
	if err := proto.Unmarshal(b, pbData); err != nil {
		return err
	}
	t.FromProto(pbData)
	return nil
}

func (t *TransactionsFromLeader) IsValidSign() bool {
	if len(t.transactions) == 0 {
		return true
	}
	publicKeys := make([][]byte, len(t.transactions))
	hashes := make([][]byte, len(t.transactions))
	for i, tx := range t.transactions {
		publicKeys[i] = tx.Pubkey().Bytes()
		hashes[i] = tx.Hash().Bytes()
	}
	return bls.VerifyAggregateSign(publicKeys, t.aggSign, hashes)
}

func (t *TransactionsFromLeader) Proto() *pb.TransactionsFromLeader {
	transactions := make([]*pb.Transaction, 0, len(t.transactions))
	for _, transaction := range t.transactions {
		transactions = append(transactions, transaction.Proto().(*pb.Transaction))
	}
	return &pb.TransactionsFromLeader{
		Transactions: transactions,
		BlockNumber:  t.blockNumber,
		AggSign:      t.aggSign,
		TimeStamp:    t.timeStamp,
	}
}

func (t *TransactionsFromLeader) FromProto(pbData *pb.TransactionsFromLeader) {
	transactions := make([]types.Transaction, len(pbData.Transactions))
	for i, pbTransaction := range pbData.Transactions {
		transaction := &Transaction{}
		transaction.FromProto(pbTransaction)
		transactions[i] = transaction
	}
	t.transactions = transactions
	t.blockNumber = pbData.BlockNumber
	t.aggSign = pbData.AggSign
	t.timeStamp = pbData.TimeStamp
}

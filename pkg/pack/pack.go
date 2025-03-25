package pack

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"gomail/pkg/bls"
	pb "gomail/pkg/proto"
	"gomail/pkg/transaction"
	"gomail/types"
)

type Pack struct {
	id           string
	transactions []types.Transaction
	aggSign      []byte
	timeStamp    uint64
}

func NewPack(
	transactions []types.Transaction,
	aggSign []byte,
	timeStamp uint64,
) types.Pack {
	return &Pack{
		transactions: transactions,
		aggSign:      aggSign,
		timeStamp:    timeStamp,
		id:           uuid.New().String(),
	}
}

func (p *Pack) NewVerifyPackSignRequest() types.VerifyPackSignRequest {
	publicKeys := make([][]byte, len(p.transactions))
	hashes := make([][]byte, len(p.transactions))
	for i, tx := range p.transactions {
		publicKeys[i] = tx.Pubkey().Bytes()
		hashes[i] = tx.Hash().Bytes()
	}
	return &VerifyPackSignRequest{
		packId:        p.Id(),
		publicKeys:    publicKeys,
		hashes:        hashes,
		aggregateSign: p.AggregateSign(),
	}
}

// general
func (p *Pack) Unmarshal(b []byte) error {
	protoPack := &pb.Pack{}
	err := proto.Unmarshal(b, protoPack)
	if err != nil {
		return err
	}
	p.FromProto(protoPack)
	return nil
}

func (p *Pack) Marshal() ([]byte, error) {
	return proto.Marshal(p.Proto())
}

func (p *Pack) Proto() *pb.Pack {
	pbTransactions := transaction.TransactionsToProto(p.transactions)
	return &pb.Pack{
		Transactions:  pbTransactions,
		AggregateSign: p.aggSign,
		TimeStamp:     p.timeStamp,
		Id:            p.id,
	}
}

func (p *Pack) FromProto(pbMessage *pb.Pack) {
	p.transactions = transaction.TransactionsFromProto(pbMessage.Transactions)
	p.aggSign = pbMessage.AggregateSign
	p.timeStamp = pbMessage.TimeStamp
	p.id = pbMessage.Id
}

// getter
func (p *Pack) Transactions() []types.Transaction {
	return p.transactions
}

func (p *Pack) Timestamp() uint64 {
	return p.timeStamp
}

func (p *Pack) Id() string {
	return p.id
}

func (p *Pack) AggregateSign() []byte {
	return p.aggSign
}

func (p *Pack) ValidSign() bool {
	publicKeys := make([][]byte, len(p.transactions))
	hashes := make([][]byte, len(p.transactions))
	for i, tx := range p.transactions {
		publicKeys[i] = tx.Pubkey().Bytes()
		hashes[i] = tx.Hash().Bytes()
	}
	validSign := bls.VerifyAggregateSign(publicKeys, p.aggSign, hashes)
	return validSign
}

func PacksToProto(packs []types.Pack) []*pb.Pack {
	rs := make([]*pb.Pack, len(packs))
	for i, v := range packs {
		rs[i] = v.Proto()
	}
	return rs
}

func PackFromProto(pbPack *pb.Pack) types.Pack {
	return &Pack{
		transactions: transaction.TransactionsFromProto(pbPack.Transactions),
		aggSign:      pbPack.AggregateSign,
		timeStamp:    pbPack.TimeStamp,
		id:           pbPack.Id,
	}
}

func PacksFromProto(pbPacks []*pb.Pack) []types.Pack {
	rs := make([]types.Pack, len(pbPacks))
	for i, v := range pbPacks {
		rs[i] = PackFromProto(v)
	}
	return rs
}

func MarshalPacks(packs []types.Pack) ([]byte, error) {
	return proto.Marshal(&pb.Packs{Packs: PacksToProto(packs)})
}

func UnmarshalTransactions(b []byte) ([]types.Pack, error) {
	pbPacks := &pb.Packs{}
	err := proto.Unmarshal(b, pbPacks)
	if err != nil {
		return nil, err
	}
	return PacksFromProto(pbPacks.Packs), nil
}

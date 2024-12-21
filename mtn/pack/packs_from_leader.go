package pack

import (
	"google.golang.org/protobuf/proto"

	pb "gomail/mtn/proto"
	"gomail/mtn/types"
)

type PacksFromLeader struct {
	packs       []types.Pack
	blockNumber uint64
	timeStamp   uint64
}

func NewPacksFromLeader(
	packs []types.Pack,
	blockNumber uint64,
	timeStamp uint64,
) *PacksFromLeader {
	return &PacksFromLeader{
		packs:       packs,
		blockNumber: blockNumber,
		timeStamp:   timeStamp,
	}
}

func (t *PacksFromLeader) Packs() []types.Pack {
	return t.packs
}

func (t *PacksFromLeader) BlockNumber() uint64 {
	return t.blockNumber
}

func (t *PacksFromLeader) TimeStamp() uint64 {
	return t.timeStamp
}

func (t *PacksFromLeader) Marshal() ([]byte, error) {
	return proto.Marshal(t.Proto())
}

func (t *PacksFromLeader) Unmarshal(b []byte) error {
	pbData := &pb.PacksFromLeader{}
	if err := proto.Unmarshal(b, pbData); err != nil {
		return err
	}
	t.FromProto(pbData)
	return nil
}

func (t *PacksFromLeader) IsValidSign() bool {
	for _, pack := range t.packs {
		if !pack.ValidSign() {
			return false
		}
	}

	return true
}

func (t *PacksFromLeader) Transactions() []types.Transaction {
	txs := []types.Transaction{}
	for _, pack := range t.packs {
		txs = append(txs, pack.Transactions()...)
	}

	return txs
}

func (t *PacksFromLeader) Proto() *pb.PacksFromLeader {
	packs := make([]*pb.Pack, 0, len(t.packs))
	for _, pack := range t.packs {
		packs = append(packs, pack.Proto())
	}
	return &pb.PacksFromLeader{
		Packs:       packs,
		BlockNumber: t.blockNumber,
		TimeStamp:   t.timeStamp,
	}
}

func (t *PacksFromLeader) FromProto(pbData *pb.PacksFromLeader) {
	packs := make([]types.Pack, len(pbData.Packs))
	for i, pbPack := range pbData.Packs {
		transaction := &Pack{}
		transaction.FromProto(pbPack)
		packs[i] = transaction
	}
	t.packs = packs
	t.blockNumber = pbData.BlockNumber
	t.timeStamp = pbData.TimeStamp
}

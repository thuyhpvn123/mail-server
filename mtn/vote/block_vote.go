package vote

import (
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"

	"gomail/mtn/bls"
	cm "gomail/mtn/common"
	pb "gomail/mtn/proto"
	"gomail/mtn/types"
)

type BlockVote struct {
	blockHash common.Hash
	number    uint64
	publicKey cm.PublicKey
	sign      cm.Sign
}

func NewBlockVote(
	blockHash common.Hash,
	number uint64,
	publicKey cm.PublicKey,
	sign cm.Sign,
) types.BlockVote {
	return &BlockVote{
		blockHash: blockHash,
		number:    number,
		publicKey: publicKey,
		sign:      sign,
	}
}

func (b *BlockVote) BlockHash() common.Hash {
	return b.blockHash
}

func (b *BlockVote) Number() uint64 {
	return b.number
}

func (b *BlockVote) PublicKey() cm.PublicKey {
	return b.publicKey
}

func (b *BlockVote) Address() common.Address {
	return cm.AddressFromPubkey(b.publicKey)
}

func (b *BlockVote) Sign() cm.Sign {
	return b.sign
}

func (b *BlockVote) Valid() bool {
	return bls.VerifySign(
		b.publicKey,
		b.sign,
		b.blockHash.Bytes(),
	)
}

func (b *BlockVote) Marshal() ([]byte, error) {
	return proto.Marshal(b.Proto())
}

func (b *BlockVote) Unmarshal(bData []byte) error {
	pbBlockVote := &pb.BlockVote{}
	if err := proto.Unmarshal(bData, pbBlockVote); err != nil {
		return err
	}
	b.FromProto(pbBlockVote)
	return nil
}

func (b *BlockVote) Proto() *pb.BlockVote {
	return &pb.BlockVote{
		BlockHash: b.blockHash.Bytes(),
		Number:    b.number,
		PublicKey: b.publicKey.Bytes(),
		Sign:      b.sign.Bytes(),
	}
}

func (b *BlockVote) FromProto(v *pb.BlockVote) {
	b.blockHash = common.BytesToHash(v.BlockHash)
	b.number = v.Number
	b.publicKey = cm.PubkeyFromBytes(v.PublicKey)
	b.sign = cm.SignFromBytes(v.Sign)
}

package transaction

import (
	"fmt"

	e_common "github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "gomail/pkg/proto"
)

type ToNodeTransactionResult struct {
	validTransactionHashes []e_common.Hash
	blockNumber            uint64
}

func NewToNodeTransactionResult(
	validTransactionHashes []e_common.Hash,
	blockNumber uint64,
) *ToNodeTransactionResult {
	return &ToNodeTransactionResult{
		validTransactionHashes: validTransactionHashes,
		blockNumber:            blockNumber,
	}
}

func (f *ToNodeTransactionResult) ValidTransactionHashes() []e_common.Hash {
	return f.validTransactionHashes
}

func (f *ToNodeTransactionResult) BlockNumber() uint64 {
	return f.blockNumber
}

func (f *ToNodeTransactionResult) Marshal() ([]byte, error) {
	return proto.Marshal(f.Proto())
}

func (f *ToNodeTransactionResult) Unmarshal(
	b []byte,
) error {
	pbData := &pb.ToNodeTransactionsResult{}
	if err := proto.Unmarshal(b, pbData); err != nil {
		return err
	}
	f.FromProto(pbData)
	return nil
}

func (f *ToNodeTransactionResult) String() string {
	return fmt.Sprintf(
		"ValidTransactionHashes: %v, BlockNumber: %v",
		f.validTransactionHashes,
		f.blockNumber,
	)
}

func (f *ToNodeTransactionResult) Proto() protoreflect.ProtoMessage {
	bHashes := make([][]byte, len(f.validTransactionHashes))
	for i, hash := range f.validTransactionHashes {
		bHashes[i] = hash.Bytes()
	}

	return &pb.ToNodeTransactionsResult{
		ValidTransactionHashes: bHashes,
		BlockNumber:            f.blockNumber,
	}
}

func (f *ToNodeTransactionResult) FromProto(
	proto *pb.ToNodeTransactionsResult,
) {
	hashes := make([]e_common.Hash, len(proto.ValidTransactionHashes))
	for i, hash := range proto.ValidTransactionHashes {
		hashes[i] = e_common.BytesToHash(hash)
	}
	f.validTransactionHashes = hashes
	f.blockNumber = proto.BlockNumber
}

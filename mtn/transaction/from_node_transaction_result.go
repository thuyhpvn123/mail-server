package transaction

import (
	"fmt"

	e_common "github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "gomail/mtn/proto"
)

type FromNodeTransactionResult struct {
	validTransactionHashes []e_common.Hash
	transactionErrors      map[e_common.Hash]int64
	blockNumber            uint64
}

func NewFromNodeTransactionResult(
	validTransactionHashes []e_common.Hash,
	transactionErrors map[e_common.Hash]int64,
	blockNumber uint64,
) *FromNodeTransactionResult {
	return &FromNodeTransactionResult{
		validTransactionHashes: validTransactionHashes,
		transactionErrors:      transactionErrors,
		blockNumber:            blockNumber,
	}
}

func (f *FromNodeTransactionResult) ValidTransactionHashes() []e_common.Hash {
	return f.validTransactionHashes
}

func (f *FromNodeTransactionResult) TransactionErrors() map[e_common.Hash]int64 {
	return f.transactionErrors
}

func (f *FromNodeTransactionResult) BlockNumber() uint64 {
	return f.blockNumber
}

func (f *FromNodeTransactionResult) Marshal() ([]byte, error) {
	return proto.Marshal(f.Proto())
}

func (f *FromNodeTransactionResult) Unmarshal(
	b []byte,
) error {
	pbData := &pb.FromNodeTransactionsResult{}
	if err := proto.Unmarshal(b, pbData); err != nil {
		return err
	}
	f.FromProto(pbData)
	return nil
}

func (f *FromNodeTransactionResult) String() string {
	return fmt.Sprintf(
		"ValidTransactionHashes: %v, TransactionErrors: %v",
		f.validTransactionHashes,
		f.transactionErrors,
	)
}

func (f *FromNodeTransactionResult) Proto() protoreflect.ProtoMessage {
	bHashes := make([][]byte, len(f.validTransactionHashes))
	for i, hash := range f.validTransactionHashes {
		bHashes[i] = hash.Bytes()
	}
	errorsCode := make([]*pb.TransactionHashWithErrorCode, len(f.transactionErrors))
	i := 0
	for txHash, errCode := range f.transactionErrors {
		errorsCode[i] = &pb.TransactionHashWithErrorCode{
			TransactionHash: txHash.Bytes(),
			Code:            errCode,
		}
		i++
	}

	return &pb.FromNodeTransactionsResult{
		ValidTransactionHashes: bHashes,
		TransactionErrors:      errorsCode,
		BlockNumber:            f.blockNumber,
	}
}

func (f *FromNodeTransactionResult) FromProto(
	proto *pb.FromNodeTransactionsResult,
) {
	hashes := make([]e_common.Hash, len(proto.ValidTransactionHashes))
	for i, hash := range proto.ValidTransactionHashes {
		hashes[i] = e_common.BytesToHash(hash)
	}
	errors := make(map[e_common.Hash]int64, len(proto.TransactionErrors))
	for _, txErr := range proto.TransactionErrors {
		errors[e_common.BytesToHash(txErr.TransactionHash)] = txErr.Code
	}
	f.validTransactionHashes = hashes
	f.transactionErrors = errors
	f.blockNumber = proto.BlockNumber
}

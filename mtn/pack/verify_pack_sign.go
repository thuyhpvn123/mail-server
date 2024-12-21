package pack

import (
	"google.golang.org/protobuf/proto"

	"gomail/mtn/bls"
	pb "gomail/mtn/proto"
	"gomail/mtn/types"
)

type VerifyPackSignRequest struct {
	packId        string
	publicKeys    [][]byte
	hashes        [][]byte
	aggregateSign []byte
}

func NewVerifyPackSignRequest(
	packId string,
	publicKeys [][]byte,
	hashes [][]byte,
	aggregateSign []byte,
) types.VerifyPackSignRequest {
	return &VerifyPackSignRequest{
		packId:        packId,
		publicKeys:    publicKeys,
		hashes:        hashes,
		aggregateSign: aggregateSign,
	}
}

func (request *VerifyPackSignRequest) Unmarshal(b []byte) error {
	pbRequest := &pb.VerifyPackSignRequest{}
	err := proto.Unmarshal(b, pbRequest)
	if err != nil {
		return err
	}
	request.packId = pbRequest.PackId
	request.publicKeys = pbRequest.PublicKeys
	request.hashes = pbRequest.Hashes
	request.aggregateSign = pbRequest.Sign
	return nil
}

func (request *VerifyPackSignRequest) Marshal() ([]byte, error) {
	pbRequest := &pb.VerifyPackSignRequest{
		PackId:     request.packId,
		PublicKeys: request.publicKeys,
		Hashes:     request.hashes,
		Sign:       request.aggregateSign,
	}
	return proto.Marshal(pbRequest)
}

func (r *VerifyPackSignRequest) Id() string {
	return r.packId
}

func (r *VerifyPackSignRequest) PublicKeys() [][]byte {
	return r.publicKeys
}

func (r *VerifyPackSignRequest) Hashes() [][]byte {
	return r.hashes
}

func (r *VerifyPackSignRequest) AggregateSign() []byte {
	return r.aggregateSign
}

func (r *VerifyPackSignRequest) Valid() bool {
	validSign := bls.VerifyAggregateSign(r.PublicKeys(), r.AggregateSign(), r.Hashes())
	return validSign
}

func (rs *VerifyPackSignRequest) Proto() *pb.VerifyPackSignRequest {
	return &pb.VerifyPackSignRequest{
		PackId:     rs.packId,
		PublicKeys: rs.publicKeys,
		Hashes:     rs.hashes,
		Sign:       rs.aggregateSign,
	}
}

type VerifyPackSignResult struct {
	packId string
	valid  bool
}

// ===========
func NewVerifyPackSignResult(
	packId string,
	valid bool,
) types.VerifyPackSignResult {
	return &VerifyPackSignResult{
		packId: packId,
		valid:  valid,
	}
}

func (rs *VerifyPackSignResult) Unmarshal(b []byte) error {
	rsPb := &pb.VerifyPackSignResult{}
	err := proto.Unmarshal(b, rsPb)
	if err != nil {
		return err
	}
	rs.packId = rsPb.PackId
	rs.valid = rsPb.Valid
	return nil
}

func (rs *VerifyPackSignResult) Marshal() ([]byte, error) {
	rsPb := &pb.VerifyPackSignResult{
		PackId: rs.packId,
		Valid:  rs.valid,
	}
	return proto.Marshal(rsPb)
}

func (rs *VerifyPackSignResult) PackId() string {
	return rs.packId
}

func (rs *VerifyPackSignResult) Valid() bool {
	return rs.valid
}

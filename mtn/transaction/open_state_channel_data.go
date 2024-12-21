package transaction

import (
	e_common "github.com/ethereum/go-ethereum/common"
	pb "gomail/mtn/proto"
	"gomail/mtn/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type OpenStateChannelData struct {
	validatorAddresses []e_common.Address
}

func NewOpenStateChannelData(validatorAddresses []e_common.Address) types.OpenStateChannelData {
	return &OpenStateChannelData{
		validatorAddresses: validatorAddresses,
	}
}

func (d *OpenStateChannelData) Unmarshal(b []byte) error {
	cdPb := &pb.OpenStateChannelData{}
	err := proto.Unmarshal(b, cdPb)
	if err != nil {
		return err
	}
	d.FromProto(cdPb)
	return nil
}

func (d *OpenStateChannelData) Marshal() ([]byte, error) {
	return proto.Marshal(d.Proto())
}

func (d *OpenStateChannelData) Proto() protoreflect.ProtoMessage {
	bAddresses := make([][]byte, len(d.validatorAddresses))
	for i, v := range d.validatorAddresses {
		bAddresses[i] = v.Bytes()
	}
	return &pb.OpenStateChannelData{
		ValidatorAddresses: bAddresses,
	}
}

func (d *OpenStateChannelData) FromProto(pbMessage protoreflect.ProtoMessage) {
	cdPb := pbMessage.(*pb.OpenStateChannelData)
	addresses := make([]e_common.Address, len(cdPb.ValidatorAddresses))
	for i, v := range cdPb.ValidatorAddresses {
		addresses[i] = e_common.BytesToAddress(v)
	}
	d.validatorAddresses = addresses
}

// geter
func (d *OpenStateChannelData) ValidatorAddresses() []e_common.Address {
	return d.validatorAddresses
}

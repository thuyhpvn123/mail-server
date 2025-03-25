package transaction

import (
	pb "gomail/pkg/proto"
	"gomail/types"

	e_common "github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type UpdateStorageHostData struct {
	storageHost    string
	storageAddress e_common.Address
}

func NewUpdateStorageHostData(
	storageHost string,
	storageAddress e_common.Address,
) types.UpdateStorageHostData {
	return &UpdateStorageHostData{
		storageHost:    storageHost,
		storageAddress: storageAddress,
	}
}

func (d *UpdateStorageHostData) Unmarshal(b []byte) error {
	cdPb := &pb.UpdateStorageHostData{}
	err := proto.Unmarshal(b, cdPb)
	if err != nil {
		return err
	}
	d.FromProto(cdPb)
	return nil
}

func (d *UpdateStorageHostData) Marshal() ([]byte, error) {
	return proto.Marshal(d.Proto())
}

func (d *UpdateStorageHostData) Proto() protoreflect.ProtoMessage {
	return &pb.UpdateStorageHostData{
		StorageHost:    d.storageHost,
		StorageAddress: d.storageAddress.Bytes(),
	}
}

func (d *UpdateStorageHostData) FromProto(pbMessage protoreflect.ProtoMessage) {
	dPb := pbMessage.(*pb.UpdateStorageHostData)
	d.storageAddress = e_common.BytesToAddress(dPb.StorageAddress)
	d.storageHost = dPb.StorageHost
}

// geter
func (d *UpdateStorageHostData) StorageHost() string {
	return d.storageHost
}

func (d *UpdateStorageHostData) StorageAddress() e_common.Address {
	return d.storageAddress
}

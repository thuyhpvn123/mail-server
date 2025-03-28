// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: checkpoint.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ValidatorWithStakeAmount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address     []byte `protobuf:"bytes,1,opt,name=Address,proto3" json:"Address,omitempty"`
	StakeAmount []byte `protobuf:"bytes,2,opt,name=StakeAmount,proto3" json:"StakeAmount,omitempty"`
}

func (x *ValidatorWithStakeAmount) Reset() {
	*x = ValidatorWithStakeAmount{}
	if protoimpl.UnsafeEnabled {
		mi := &file_checkpoint_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidatorWithStakeAmount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidatorWithStakeAmount) ProtoMessage() {}

func (x *ValidatorWithStakeAmount) ProtoReflect() protoreflect.Message {
	mi := &file_checkpoint_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidatorWithStakeAmount.ProtoReflect.Descriptor instead.
func (*ValidatorWithStakeAmount) Descriptor() ([]byte, []int) {
	return file_checkpoint_proto_rawDescGZIP(), []int{0}
}

func (x *ValidatorWithStakeAmount) GetAddress() []byte {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *ValidatorWithStakeAmount) GetStakeAmount() []byte {
	if x != nil {
		return x.StakeAmount
	}
	return nil
}

type ValidatorCheckPoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockNumber                 uint64                      `protobuf:"varint,1,opt,name=BlockNumber,proto3" json:"BlockNumber,omitempty"`
	PreValidatorWithStakeAmount []*ValidatorWithStakeAmount `protobuf:"bytes,2,rep,name=PreValidatorWithStakeAmount,proto3" json:"PreValidatorWithStakeAmount,omitempty"`
	ValidatorWithStakeAmount    []*ValidatorWithStakeAmount `protobuf:"bytes,3,rep,name=ValidatorWithStakeAmount,proto3" json:"ValidatorWithStakeAmount,omitempty"`
	ScheduleSeed                []byte                      `protobuf:"bytes,4,opt,name=ScheduleSeed,proto3" json:"ScheduleSeed,omitempty"`
	NextScheduleSeed            []byte                      `protobuf:"bytes,5,opt,name=NextScheduleSeed,proto3" json:"NextScheduleSeed,omitempty"`
	StakeStorageRoot            []byte                      `protobuf:"bytes,6,opt,name=StakeStorageRoot,proto3" json:"StakeStorageRoot,omitempty"`
}

func (x *ValidatorCheckPoint) Reset() {
	*x = ValidatorCheckPoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_checkpoint_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidatorCheckPoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidatorCheckPoint) ProtoMessage() {}

func (x *ValidatorCheckPoint) ProtoReflect() protoreflect.Message {
	mi := &file_checkpoint_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidatorCheckPoint.ProtoReflect.Descriptor instead.
func (*ValidatorCheckPoint) Descriptor() ([]byte, []int) {
	return file_checkpoint_proto_rawDescGZIP(), []int{1}
}

func (x *ValidatorCheckPoint) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *ValidatorCheckPoint) GetPreValidatorWithStakeAmount() []*ValidatorWithStakeAmount {
	if x != nil {
		return x.PreValidatorWithStakeAmount
	}
	return nil
}

func (x *ValidatorCheckPoint) GetValidatorWithStakeAmount() []*ValidatorWithStakeAmount {
	if x != nil {
		return x.ValidatorWithStakeAmount
	}
	return nil
}

func (x *ValidatorCheckPoint) GetScheduleSeed() []byte {
	if x != nil {
		return x.ScheduleSeed
	}
	return nil
}

func (x *ValidatorCheckPoint) GetNextScheduleSeed() []byte {
	if x != nil {
		return x.NextScheduleSeed
	}
	return nil
}

func (x *ValidatorCheckPoint) GetStakeStorageRoot() []byte {
	if x != nil {
		return x.StakeStorageRoot
	}
	return nil
}

var File_checkpoint_proto protoreflect.FileDescriptor

var file_checkpoint_proto_rawDesc = []byte{
	0x0a, 0x10, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x22, 0x56,
	0x0a, 0x18, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x57, 0x69, 0x74, 0x68, 0x53,
	0x74, 0x61, 0x6b, 0x65, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x41, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x53, 0x74, 0x61, 0x6b, 0x65,
	0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0xfd, 0x02, 0x0a, 0x13, 0x56, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x6f, 0x72, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x20,
	0x0a, 0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x12, 0x66, 0x0a, 0x1b, 0x50, 0x72, 0x65, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72,
	0x57, 0x69, 0x74, 0x68, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x57, 0x69, 0x74, 0x68,
	0x53, 0x74, 0x61, 0x6b, 0x65, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x1b, 0x50, 0x72, 0x65,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x57, 0x69, 0x74, 0x68, 0x53, 0x74, 0x61,
	0x6b, 0x65, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x60, 0x0a, 0x18, 0x56, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x6f, 0x72, 0x57, 0x69, 0x74, 0x68, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x41, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f,
	0x72, 0x57, 0x69, 0x74, 0x68, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
	0x52, 0x18, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x57, 0x69, 0x74, 0x68, 0x53,
	0x74, 0x61, 0x6b, 0x65, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x53, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x53, 0x65, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0c, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x53, 0x65, 0x65, 0x64, 0x12, 0x2a,
	0x0a, 0x10, 0x4e, 0x65, 0x78, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x53, 0x65,
	0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x4e, 0x65, 0x78, 0x74, 0x53, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x53, 0x65, 0x65, 0x64, 0x12, 0x2a, 0x0a, 0x10, 0x53, 0x74,
	0x61, 0x6b, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x42, 0x32, 0x0a, 0x28, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x65,
	0x74, 0x61, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x63,
	0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x65, 0x64, 0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x5a, 0x06, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_checkpoint_proto_rawDescOnce sync.Once
	file_checkpoint_proto_rawDescData = file_checkpoint_proto_rawDesc
)

func file_checkpoint_proto_rawDescGZIP() []byte {
	file_checkpoint_proto_rawDescOnce.Do(func() {
		file_checkpoint_proto_rawDescData = protoimpl.X.CompressGZIP(file_checkpoint_proto_rawDescData)
	})
	return file_checkpoint_proto_rawDescData
}

var file_checkpoint_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_checkpoint_proto_goTypes = []interface{}{
	(*ValidatorWithStakeAmount)(nil), // 0: checkpoint.ValidatorWithStakeAmount
	(*ValidatorCheckPoint)(nil),      // 1: checkpoint.ValidatorCheckPoint
}
var file_checkpoint_proto_depIdxs = []int32{
	0, // 0: checkpoint.ValidatorCheckPoint.PreValidatorWithStakeAmount:type_name -> checkpoint.ValidatorWithStakeAmount
	0, // 1: checkpoint.ValidatorCheckPoint.ValidatorWithStakeAmount:type_name -> checkpoint.ValidatorWithStakeAmount
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_checkpoint_proto_init() }
func file_checkpoint_proto_init() {
	if File_checkpoint_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_checkpoint_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidatorWithStakeAmount); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_checkpoint_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidatorCheckPoint); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_checkpoint_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_checkpoint_proto_goTypes,
		DependencyIndexes: file_checkpoint_proto_depIdxs,
		MessageInfos:      file_checkpoint_proto_msgTypes,
	}.Build()
	File_checkpoint_proto = out.File
	file_checkpoint_proto_rawDesc = nil
	file_checkpoint_proto_goTypes = nil
	file_checkpoint_proto_depIdxs = nil
}

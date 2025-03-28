// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: block_vote.proto

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

type BlockVote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockHash []byte `protobuf:"bytes,1,opt,name=BlockHash,proto3" json:"BlockHash,omitempty"`
	Number    uint64 `protobuf:"varint,2,opt,name=Number,proto3" json:"Number,omitempty"`
	PublicKey []byte `protobuf:"bytes,3,opt,name=PublicKey,proto3" json:"PublicKey,omitempty"`
	Sign      []byte `protobuf:"bytes,4,opt,name=Sign,proto3" json:"Sign,omitempty"`
}

func (x *BlockVote) Reset() {
	*x = BlockVote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_vote_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockVote) ProtoMessage() {}

func (x *BlockVote) ProtoReflect() protoreflect.Message {
	mi := &file_block_vote_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockVote.ProtoReflect.Descriptor instead.
func (*BlockVote) Descriptor() ([]byte, []int) {
	return file_block_vote_proto_rawDescGZIP(), []int{0}
}

func (x *BlockVote) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *BlockVote) GetNumber() uint64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *BlockVote) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (x *BlockVote) GetSign() []byte {
	if x != nil {
		return x.Sign
	}
	return nil
}

var File_block_vote_proto protoreflect.FileDescriptor

var file_block_vote_proto_rawDesc = []byte{
	0x0a, 0x10, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x22, 0x73,
	0x0a, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x1c, 0x0a, 0x09, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12,
	0x12, 0x0a, 0x04, 0x53, 0x69, 0x67, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x53,
	0x69, 0x67, 0x6e, 0x42, 0x32, 0x0a, 0x28, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x5f,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x70,
	0x69, 0x6c, 0x65, 0x64, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x5a,
	0x06, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_block_vote_proto_rawDescOnce sync.Once
	file_block_vote_proto_rawDescData = file_block_vote_proto_rawDesc
)

func file_block_vote_proto_rawDescGZIP() []byte {
	file_block_vote_proto_rawDescOnce.Do(func() {
		file_block_vote_proto_rawDescData = protoimpl.X.CompressGZIP(file_block_vote_proto_rawDescData)
	})
	return file_block_vote_proto_rawDescData
}

var file_block_vote_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_block_vote_proto_goTypes = []interface{}{
	(*BlockVote)(nil), // 0: block_vote.BlockVote
}
var file_block_vote_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_block_vote_proto_init() }
func file_block_vote_proto_init() {
	if File_block_vote_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_block_vote_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockVote); i {
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
			RawDescriptor: file_block_vote_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_block_vote_proto_goTypes,
		DependencyIndexes: file_block_vote_proto_depIdxs,
		MessageInfos:      file_block_vote_proto_msgTypes,
	}.Build()
	File_block_vote_proto = out.File
	file_block_vote_proto_rawDesc = nil
	file_block_vote_proto_goTypes = nil
	file_block_vote_proto_depIdxs = nil
}

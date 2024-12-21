// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.3
// source: poh.proto

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

type PohHashData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PreHash    []byte   `protobuf:"bytes,1,opt,name=PreHash,proto3" json:"PreHash,omitempty"`
	PackHashes [][]byte `protobuf:"bytes,2,rep,name=PackHashes,proto3" json:"PackHashes,omitempty"`
}

func (x *PohHashData) Reset() {
	*x = PohHashData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_poh_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PohHashData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PohHashData) ProtoMessage() {}

func (x *PohHashData) ProtoReflect() protoreflect.Message {
	mi := &file_poh_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PohHashData.ProtoReflect.Descriptor instead.
func (*PohHashData) Descriptor() ([]byte, []int) {
	return file_poh_proto_rawDescGZIP(), []int{0}
}

func (x *PohHashData) GetPreHash() []byte {
	if x != nil {
		return x.PreHash
	}
	return nil
}

func (x *PohHashData) GetPackHashes() [][]byte {
	if x != nil {
		return x.PackHashes
	}
	return nil
}

type PohEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NumHashes uint64  `protobuf:"varint,1,opt,name=NumHashes,proto3" json:"NumHashes,omitempty"`
	Hash      []byte  `protobuf:"bytes,2,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Packs     []*Pack `protobuf:"bytes,3,rep,name=Packs,proto3" json:"Packs,omitempty"`
	Time      uint64  `protobuf:"varint,4,opt,name=Time,proto3" json:"Time,omitempty"`
}

func (x *PohEntry) Reset() {
	*x = PohEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_poh_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PohEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PohEntry) ProtoMessage() {}

func (x *PohEntry) ProtoReflect() protoreflect.Message {
	mi := &file_poh_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PohEntry.ProtoReflect.Descriptor instead.
func (*PohEntry) Descriptor() ([]byte, []int) {
	return file_poh_proto_rawDescGZIP(), []int{1}
}

func (x *PohEntry) GetNumHashes() uint64 {
	if x != nil {
		return x.NumHashes
	}
	return 0
}

func (x *PohEntry) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *PohEntry) GetPacks() []*Pack {
	if x != nil {
		return x.Packs
	}
	return nil
}

func (x *PohEntry) GetTime() uint64 {
	if x != nil {
		return x.Time
	}
	return 0
}

type LeaderSchedule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Seed     []byte            `protobuf:"bytes,1,opt,name=Seed,proto3" json:"Seed,omitempty"`
	FromSlot []byte            `protobuf:"bytes,2,opt,name=FromSlot,proto3" json:"FromSlot,omitempty"`
	ToSlot   []byte            `protobuf:"bytes,3,opt,name=ToSlot,proto3" json:"ToSlot,omitempty"`
	Slots    map[string][]byte `protobuf:"bytes,4,rep,name=Slots,proto3" json:"Slots,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *LeaderSchedule) Reset() {
	*x = LeaderSchedule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_poh_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaderSchedule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaderSchedule) ProtoMessage() {}

func (x *LeaderSchedule) ProtoReflect() protoreflect.Message {
	mi := &file_poh_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeaderSchedule.ProtoReflect.Descriptor instead.
func (*LeaderSchedule) Descriptor() ([]byte, []int) {
	return file_poh_proto_rawDescGZIP(), []int{2}
}

func (x *LeaderSchedule) GetSeed() []byte {
	if x != nil {
		return x.Seed
	}
	return nil
}

func (x *LeaderSchedule) GetFromSlot() []byte {
	if x != nil {
		return x.FromSlot
	}
	return nil
}

func (x *LeaderSchedule) GetToSlot() []byte {
	if x != nil {
		return x.ToSlot
	}
	return nil
}

func (x *LeaderSchedule) GetSlots() map[string][]byte {
	if x != nil {
		return x.Slots
	}
	return nil
}

type EntryWithBlockNumber struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockNumber []byte    `protobuf:"bytes,1,opt,name=BlockNumber,proto3" json:"BlockNumber,omitempty"`
	Entry       *PohEntry `protobuf:"bytes,2,opt,name=Entry,proto3" json:"Entry,omitempty"`
}

func (x *EntryWithBlockNumber) Reset() {
	*x = EntryWithBlockNumber{}
	if protoimpl.UnsafeEnabled {
		mi := &file_poh_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntryWithBlockNumber) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntryWithBlockNumber) ProtoMessage() {}

func (x *EntryWithBlockNumber) ProtoReflect() protoreflect.Message {
	mi := &file_poh_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntryWithBlockNumber.ProtoReflect.Descriptor instead.
func (*EntryWithBlockNumber) Descriptor() ([]byte, []int) {
	return file_poh_proto_rawDescGZIP(), []int{3}
}

func (x *EntryWithBlockNumber) GetBlockNumber() []byte {
	if x != nil {
		return x.BlockNumber
	}
	return nil
}

func (x *EntryWithBlockNumber) GetEntry() *PohEntry {
	if x != nil {
		return x.Entry
	}
	return nil
}

var File_poh_proto protoreflect.FileDescriptor

var file_poh_proto_rawDesc = []byte{
	0x0a, 0x09, 0x70, 0x6f, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x70, 0x6f, 0x68,
	0x1a, 0x0a, 0x70, 0x61, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x47, 0x0a, 0x0b,
	0x50, 0x6f, 0x68, 0x48, 0x61, 0x73, 0x68, 0x44, 0x61, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x50,
	0x72, 0x65, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x50, 0x72,
	0x65, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1e, 0x0a, 0x0a, 0x50, 0x61, 0x63, 0x6b, 0x48, 0x61, 0x73,
	0x68, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0a, 0x50, 0x61, 0x63, 0x6b, 0x48,
	0x61, 0x73, 0x68, 0x65, 0x73, 0x22, 0x72, 0x0a, 0x08, 0x50, 0x6f, 0x68, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x1c, 0x0a, 0x09, 0x4e, 0x75, 0x6d, 0x48, 0x61, 0x73, 0x68, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x4e, 0x75, 0x6d, 0x48, 0x61, 0x73, 0x68, 0x65, 0x73, 0x12,
	0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x48,
	0x61, 0x73, 0x68, 0x12, 0x20, 0x0a, 0x05, 0x50, 0x61, 0x63, 0x6b, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x52, 0x05,
	0x50, 0x61, 0x63, 0x6b, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x22, 0xc8, 0x01, 0x0a, 0x0e, 0x4c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x53, 0x65, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x53, 0x65, 0x65, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x46, 0x72, 0x6f, 0x6d, 0x53, 0x6c, 0x6f, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x08, 0x46, 0x72, 0x6f, 0x6d, 0x53, 0x6c, 0x6f, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x54, 0x6f, 0x53, 0x6c, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x54, 0x6f,
	0x53, 0x6c, 0x6f, 0x74, 0x12, 0x34, 0x0a, 0x05, 0x53, 0x6c, 0x6f, 0x74, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x70, 0x6f, 0x68, 0x2e, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x2e, 0x53, 0x6c, 0x6f, 0x74, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x05, 0x53, 0x6c, 0x6f, 0x74, 0x73, 0x1a, 0x38, 0x0a, 0x0a, 0x53, 0x6c,
	0x6f, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0x5d, 0x0a, 0x14, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x57, 0x69, 0x74,
	0x68, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0b,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x23,
	0x0a, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x70, 0x6f, 0x68, 0x2e, 0x50, 0x6f, 0x68, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x42, 0x2b, 0x0a, 0x21, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x5f,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x70,
	0x69, 0x6c, 0x65, 0x64, 0x2e, 0x70, 0x6f, 0x68, 0x5a, 0x06, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_poh_proto_rawDescOnce sync.Once
	file_poh_proto_rawDescData = file_poh_proto_rawDesc
)

func file_poh_proto_rawDescGZIP() []byte {
	file_poh_proto_rawDescOnce.Do(func() {
		file_poh_proto_rawDescData = protoimpl.X.CompressGZIP(file_poh_proto_rawDescData)
	})
	return file_poh_proto_rawDescData
}

var file_poh_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_poh_proto_goTypes = []interface{}{
	(*PohHashData)(nil),          // 0: poh.PohHashData
	(*PohEntry)(nil),             // 1: poh.PohEntry
	(*LeaderSchedule)(nil),       // 2: poh.LeaderSchedule
	(*EntryWithBlockNumber)(nil), // 3: poh.EntryWithBlockNumber
	nil,                          // 4: poh.LeaderSchedule.SlotsEntry
	(*Pack)(nil),                 // 5: pack.Pack
}
var file_poh_proto_depIdxs = []int32{
	5, // 0: poh.PohEntry.Packs:type_name -> pack.Pack
	4, // 1: poh.LeaderSchedule.Slots:type_name -> poh.LeaderSchedule.SlotsEntry
	1, // 2: poh.EntryWithBlockNumber.Entry:type_name -> poh.PohEntry
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_poh_proto_init() }
func file_poh_proto_init() {
	if File_poh_proto != nil {
		return
	}
	file_pack_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_poh_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PohHashData); i {
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
		file_poh_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PohEntry); i {
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
		file_poh_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeaderSchedule); i {
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
		file_poh_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntryWithBlockNumber); i {
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
			RawDescriptor: file_poh_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_poh_proto_goTypes,
		DependencyIndexes: file_poh_proto_depIdxs,
		MessageInfos:      file_poh_proto_msgTypes,
	}.Build()
	File_poh_proto = out.File
	file_poh_proto_rawDesc = nil
	file_poh_proto_goTypes = nil
	file_poh_proto_depIdxs = nil
}

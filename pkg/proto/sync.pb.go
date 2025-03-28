// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: sync.proto

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

type SyncAccoutStates struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StorageData map[string][]byte `protobuf:"bytes,1,rep,name=StorageData,proto3" json:"StorageData,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Finished    bool              `protobuf:"varint,2,opt,name=Finished,proto3" json:"Finished,omitempty"`
}

func (x *SyncAccoutStates) Reset() {
	*x = SyncAccoutStates{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sync_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncAccoutStates) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncAccoutStates) ProtoMessage() {}

func (x *SyncAccoutStates) ProtoReflect() protoreflect.Message {
	mi := &file_sync_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncAccoutStates.ProtoReflect.Descriptor instead.
func (*SyncAccoutStates) Descriptor() ([]byte, []int) {
	return file_sync_proto_rawDescGZIP(), []int{0}
}

func (x *SyncAccoutStates) GetStorageData() map[string][]byte {
	if x != nil {
		return x.StorageData
	}
	return nil
}

func (x *SyncAccoutStates) GetFinished() bool {
	if x != nil {
		return x.Finished
	}
	return false
}

type SyncNodeConsensusConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PacksPerEntry           uint64 `protobuf:"varint,1,opt,name=PacksPerEntry,proto3" json:"PacksPerEntry,omitempty"`
	EntriesPerSlot          uint64 `protobuf:"varint,2,opt,name=EntriesPerSlot,proto3" json:"EntriesPerSlot,omitempty"`
	EntriesPerSecond        uint64 `protobuf:"varint,3,opt,name=EntriesPerSecond,proto3" json:"EntriesPerSecond,omitempty"`
	HashesPerEntry          uint64 `protobuf:"varint,4,opt,name=HashesPerEntry,proto3" json:"HashesPerEntry,omitempty"`
	ValidatorMinStakeAmount []byte `protobuf:"bytes,5,opt,name=ValidatorMinStakeAmount,proto3" json:"ValidatorMinStakeAmount,omitempty"`
	StartRewardAmount       []byte `protobuf:"bytes,6,opt,name=StartRewardAmount,proto3" json:"StartRewardAmount,omitempty"`
	HalvingAfterBlockCount  []byte `protobuf:"bytes,7,opt,name=HalvingAfterBlockCount,proto3" json:"HalvingAfterBlockCount,omitempty"`
}

func (x *SyncNodeConsensusConfig) Reset() {
	*x = SyncNodeConsensusConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sync_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncNodeConsensusConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncNodeConsensusConfig) ProtoMessage() {}

func (x *SyncNodeConsensusConfig) ProtoReflect() protoreflect.Message {
	mi := &file_sync_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncNodeConsensusConfig.ProtoReflect.Descriptor instead.
func (*SyncNodeConsensusConfig) Descriptor() ([]byte, []int) {
	return file_sync_proto_rawDescGZIP(), []int{1}
}

func (x *SyncNodeConsensusConfig) GetPacksPerEntry() uint64 {
	if x != nil {
		return x.PacksPerEntry
	}
	return 0
}

func (x *SyncNodeConsensusConfig) GetEntriesPerSlot() uint64 {
	if x != nil {
		return x.EntriesPerSlot
	}
	return 0
}

func (x *SyncNodeConsensusConfig) GetEntriesPerSecond() uint64 {
	if x != nil {
		return x.EntriesPerSecond
	}
	return 0
}

func (x *SyncNodeConsensusConfig) GetHashesPerEntry() uint64 {
	if x != nil {
		return x.HashesPerEntry
	}
	return 0
}

func (x *SyncNodeConsensusConfig) GetValidatorMinStakeAmount() []byte {
	if x != nil {
		return x.ValidatorMinStakeAmount
	}
	return nil
}

func (x *SyncNodeConsensusConfig) GetStartRewardAmount() []byte {
	if x != nil {
		return x.StartRewardAmount
	}
	return nil
}

func (x *SyncNodeConsensusConfig) GetHalvingAfterBlockCount() []byte {
	if x != nil {
		return x.HalvingAfterBlockCount
	}
	return nil
}

type ValidatorSyncData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LatestBlock           *Block               `protobuf:"bytes,1,opt,name=LatestBlock,proto3" json:"LatestBlock,omitempty"`
	LatestCheckPoint      *ValidatorCheckPoint `protobuf:"bytes,2,opt,name=LatestCheckPoint,proto3" json:"LatestCheckPoint,omitempty"`
	LatestCheckPointBlock *Block               `protobuf:"bytes,3,opt,name=LatestCheckPointBlock,proto3" json:"LatestCheckPointBlock,omitempty"`
	StakeStorageData      []*StorageData       `protobuf:"bytes,4,rep,name=StakeStorageData,proto3" json:"StakeStorageData,omitempty"`
	PendingBlockVotes     []*BlockVote         `protobuf:"bytes,5,rep,name=PendingBlockVotes,proto3" json:"PendingBlockVotes,omitempty"`
}

func (x *ValidatorSyncData) Reset() {
	*x = ValidatorSyncData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sync_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidatorSyncData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidatorSyncData) ProtoMessage() {}

func (x *ValidatorSyncData) ProtoReflect() protoreflect.Message {
	mi := &file_sync_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidatorSyncData.ProtoReflect.Descriptor instead.
func (*ValidatorSyncData) Descriptor() ([]byte, []int) {
	return file_sync_proto_rawDescGZIP(), []int{2}
}

func (x *ValidatorSyncData) GetLatestBlock() *Block {
	if x != nil {
		return x.LatestBlock
	}
	return nil
}

func (x *ValidatorSyncData) GetLatestCheckPoint() *ValidatorCheckPoint {
	if x != nil {
		return x.LatestCheckPoint
	}
	return nil
}

func (x *ValidatorSyncData) GetLatestCheckPointBlock() *Block {
	if x != nil {
		return x.LatestCheckPointBlock
	}
	return nil
}

func (x *ValidatorSyncData) GetStakeStorageData() []*StorageData {
	if x != nil {
		return x.StakeStorageData
	}
	return nil
}

func (x *ValidatorSyncData) GetPendingBlockVotes() []*BlockVote {
	if x != nil {
		return x.PendingBlockVotes
	}
	return nil
}

type GetNodeSyncData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LatestCheckPointBlockNumber uint64 `protobuf:"varint,1,opt,name=LatestCheckPointBlockNumber,proto3" json:"LatestCheckPointBlockNumber,omitempty"`
	ValidatorAddress            []byte `protobuf:"bytes,2,opt,name=ValidatorAddress,proto3" json:"ValidatorAddress,omitempty"`
	NodeStatesIndex             int64  `protobuf:"varint,3,opt,name=NodeStatesIndex,proto3" json:"NodeStatesIndex,omitempty"`
}

func (x *GetNodeSyncData) Reset() {
	*x = GetNodeSyncData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sync_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetNodeSyncData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetNodeSyncData) ProtoMessage() {}

func (x *GetNodeSyncData) ProtoReflect() protoreflect.Message {
	mi := &file_sync_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetNodeSyncData.ProtoReflect.Descriptor instead.
func (*GetNodeSyncData) Descriptor() ([]byte, []int) {
	return file_sync_proto_rawDescGZIP(), []int{3}
}

func (x *GetNodeSyncData) GetLatestCheckPointBlockNumber() uint64 {
	if x != nil {
		return x.LatestCheckPointBlockNumber
	}
	return 0
}

func (x *GetNodeSyncData) GetValidatorAddress() []byte {
	if x != nil {
		return x.ValidatorAddress
	}
	return nil
}

func (x *GetNodeSyncData) GetNodeStatesIndex() int64 {
	if x != nil {
		return x.NodeStatesIndex
	}
	return 0
}

type NodeSyncData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ValidatorAddress []byte         `protobuf:"bytes,1,opt,name=ValidatorAddress,proto3" json:"ValidatorAddress,omitempty"`
	NodeStatesIndex  int64          `protobuf:"varint,2,opt,name=NodeStatesIndex,proto3" json:"NodeStatesIndex,omitempty"`
	AccountStateRoot []byte         `protobuf:"bytes,3,opt,name=AccountStateRoot,proto3" json:"AccountStateRoot,omitempty"`
	StorageData      []*StorageData `protobuf:"bytes,4,rep,name=StorageData,proto3" json:"StorageData,omitempty"`
	Finished         bool           `protobuf:"varint,5,opt,name=Finished,proto3" json:"Finished,omitempty"`
}

func (x *NodeSyncData) Reset() {
	*x = NodeSyncData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sync_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeSyncData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeSyncData) ProtoMessage() {}

func (x *NodeSyncData) ProtoReflect() protoreflect.Message {
	mi := &file_sync_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeSyncData.ProtoReflect.Descriptor instead.
func (*NodeSyncData) Descriptor() ([]byte, []int) {
	return file_sync_proto_rawDescGZIP(), []int{4}
}

func (x *NodeSyncData) GetValidatorAddress() []byte {
	if x != nil {
		return x.ValidatorAddress
	}
	return nil
}

func (x *NodeSyncData) GetNodeStatesIndex() int64 {
	if x != nil {
		return x.NodeStatesIndex
	}
	return 0
}

func (x *NodeSyncData) GetAccountStateRoot() []byte {
	if x != nil {
		return x.AccountStateRoot
	}
	return nil
}

func (x *NodeSyncData) GetStorageData() []*StorageData {
	if x != nil {
		return x.StorageData
	}
	return nil
}

func (x *NodeSyncData) GetFinished() bool {
	if x != nil {
		return x.Finished
	}
	return false
}

var File_sync_proto protoreflect.FileDescriptor

var file_sync_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x73, 0x79,
	0x6e, 0x63, 0x1a, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x10, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xb9, 0x01, 0x0a, 0x10, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x73, 0x12, 0x49, 0x0a, 0x0b, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x73, 0x79, 0x6e,
	0x63, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x0b, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x1a, 0x3e, 0x0a, 0x10,
	0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xdb, 0x02, 0x0a,
	0x17, 0x53, 0x79, 0x6e, 0x63, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x73,
	0x75, 0x73, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x24, 0x0a, 0x0d, 0x50, 0x61, 0x63, 0x6b,
	0x73, 0x50, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0d, 0x50, 0x61, 0x63, 0x6b, 0x73, 0x50, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x26,
	0x0a, 0x0e, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x50, 0x65, 0x72, 0x53, 0x6c, 0x6f, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x50,
	0x65, 0x72, 0x53, 0x6c, 0x6f, 0x74, 0x12, 0x2a, 0x0a, 0x10, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65,
	0x73, 0x50, 0x65, 0x72, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x10, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x50, 0x65, 0x72, 0x53, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x12, 0x26, 0x0a, 0x0e, 0x48, 0x61, 0x73, 0x68, 0x65, 0x73, 0x50, 0x65, 0x72, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x48, 0x61, 0x73, 0x68,
	0x65, 0x73, 0x50, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x38, 0x0a, 0x17, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x4d, 0x69, 0x6e, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x41,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x17, 0x56, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x4d, 0x69, 0x6e, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x41, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2c, 0x0a, 0x11, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x77,
	0x61, 0x72, 0x64, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x11, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x41, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x36, 0x0a, 0x16, 0x48, 0x61, 0x6c, 0x76, 0x69, 0x6e, 0x67, 0x41, 0x66, 0x74,
	0x65, 0x72, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x16, 0x48, 0x61, 0x6c, 0x76, 0x69, 0x6e, 0x67, 0x41, 0x66, 0x74, 0x65, 0x72,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0xe1, 0x02, 0x0a, 0x11, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x79, 0x6e, 0x63, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x2e, 0x0a, 0x0b, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x52, 0x0b, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x12, 0x4b, 0x0a, 0x10, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x50,
	0x6f, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f,
	0x72, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x10, 0x4c, 0x61, 0x74,
	0x65, 0x73, 0x74, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x42, 0x0a,
	0x15, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x50, 0x6f, 0x69, 0x6e,
	0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x15, 0x4c, 0x61, 0x74, 0x65,
	0x73, 0x74, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x12, 0x46, 0x0a, 0x10, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x10, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x53, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x43, 0x0a, 0x11, 0x50, 0x65, 0x6e,
	0x64, 0x69, 0x6e, 0x67, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x56, 0x6f, 0x74, 0x65, 0x73, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x76, 0x6f, 0x74,
	0x65, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x11, 0x50, 0x65, 0x6e,
	0x64, 0x69, 0x6e, 0x67, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x56, 0x6f, 0x74, 0x65, 0x73, 0x22, 0xa9,
	0x01, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x40, 0x0a, 0x1b, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x1b, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x10, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f,
	0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x28, 0x0a, 0x0f, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x4e, 0x6f, 0x64, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x73, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x22, 0xea, 0x01, 0x0a, 0x0c, 0x4e,
	0x6f, 0x64, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x44, 0x61, 0x74, 0x61, 0x12, 0x2a, 0x0a, 0x10, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x28, 0x0a, 0x0f, 0x4e, 0x6f, 0x64, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x73, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0f, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x12, 0x2a, 0x0a, 0x10, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x3c, 0x0a,
	0x0b, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0b,
	0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x46,
	0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x46,
	0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x42, 0x2c, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e,
	0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x65, 0x64, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x5a, 0x06, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sync_proto_rawDescOnce sync.Once
	file_sync_proto_rawDescData = file_sync_proto_rawDesc
)

func file_sync_proto_rawDescGZIP() []byte {
	file_sync_proto_rawDescOnce.Do(func() {
		file_sync_proto_rawDescData = protoimpl.X.CompressGZIP(file_sync_proto_rawDescData)
	})
	return file_sync_proto_rawDescData
}

var file_sync_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_sync_proto_goTypes = []interface{}{
	(*SyncAccoutStates)(nil),        // 0: sync.SyncAccoutStates
	(*SyncNodeConsensusConfig)(nil), // 1: sync.SyncNodeConsensusConfig
	(*ValidatorSyncData)(nil),       // 2: sync.ValidatorSyncData
	(*GetNodeSyncData)(nil),         // 3: sync.GetNodeSyncData
	(*NodeSyncData)(nil),            // 4: sync.NodeSyncData
	nil,                             // 5: sync.SyncAccoutStates.StorageDataEntry
	(*Block)(nil),                   // 6: block.Block
	(*ValidatorCheckPoint)(nil),     // 7: checkpoint.ValidatorCheckPoint
	(*StorageData)(nil),             // 8: account_state.StorageData
	(*BlockVote)(nil),               // 9: block_vote.BlockVote
}
var file_sync_proto_depIdxs = []int32{
	5, // 0: sync.SyncAccoutStates.StorageData:type_name -> sync.SyncAccoutStates.StorageDataEntry
	6, // 1: sync.ValidatorSyncData.LatestBlock:type_name -> block.Block
	7, // 2: sync.ValidatorSyncData.LatestCheckPoint:type_name -> checkpoint.ValidatorCheckPoint
	6, // 3: sync.ValidatorSyncData.LatestCheckPointBlock:type_name -> block.Block
	8, // 4: sync.ValidatorSyncData.StakeStorageData:type_name -> account_state.StorageData
	9, // 5: sync.ValidatorSyncData.PendingBlockVotes:type_name -> block_vote.BlockVote
	8, // 6: sync.NodeSyncData.StorageData:type_name -> account_state.StorageData
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_sync_proto_init() }
func file_sync_proto_init() {
	if File_sync_proto != nil {
		return
	}
	file_block_proto_init()
	file_checkpoint_proto_init()
	file_state_proto_init()
	file_block_vote_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_sync_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncAccoutStates); i {
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
		file_sync_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncNodeConsensusConfig); i {
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
		file_sync_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidatorSyncData); i {
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
		file_sync_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetNodeSyncData); i {
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
		file_sync_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeSyncData); i {
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
			RawDescriptor: file_sync_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sync_proto_goTypes,
		DependencyIndexes: file_sync_proto_depIdxs,
		MessageInfos:      file_sync_proto_msgTypes,
	}.Build()
	File_sync_proto = out.File
	file_sync_proto_rawDesc = nil
	file_sync_proto_goTypes = nil
	file_sync_proto_depIdxs = nil
}

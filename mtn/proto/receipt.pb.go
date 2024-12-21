// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.27.1
// source: receipt.proto

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

type RECEIPT_STATUS int32

const (
	RECEIPT_STATUS_RETURNED          RECEIPT_STATUS = 0
	RECEIPT_STATUS_HALTED            RECEIPT_STATUS = 1
	RECEIPT_STATUS_THREW             RECEIPT_STATUS = 2
	RECEIPT_STATUS_TRANSACTION_ERROR RECEIPT_STATUS = -1
)

// Enum value maps for RECEIPT_STATUS.
var (
	RECEIPT_STATUS_name = map[int32]string{
		0:  "RETURNED",
		1:  "HALTED",
		2:  "THREW",
		-1: "TRANSACTION_ERROR",
	}
	RECEIPT_STATUS_value = map[string]int32{
		"RETURNED":          0,
		"HALTED":            1,
		"THREW":             2,
		"TRANSACTION_ERROR": -1,
	}
)

func (x RECEIPT_STATUS) Enum() *RECEIPT_STATUS {
	p := new(RECEIPT_STATUS)
	*p = x
	return p
}

func (x RECEIPT_STATUS) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RECEIPT_STATUS) Descriptor() protoreflect.EnumDescriptor {
	return file_receipt_proto_enumTypes[0].Descriptor()
}

func (RECEIPT_STATUS) Type() protoreflect.EnumType {
	return &file_receipt_proto_enumTypes[0]
}

func (x RECEIPT_STATUS) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RECEIPT_STATUS.Descriptor instead.
func (RECEIPT_STATUS) EnumDescriptor() ([]byte, []int) {
	return file_receipt_proto_rawDescGZIP(), []int{0}
}

type EXCEPTION int32

const (
	EXCEPTION_ERR_OUT_OF_GAS                 EXCEPTION = 0
	EXCEPTION_ERR_CODE_STORE_OUT_OF_GAS      EXCEPTION = 1
	EXCEPTION_ERR_DEPTH                      EXCEPTION = 2
	EXCEPTION_ERR_INSUFFICIENT_BALANCE       EXCEPTION = 3
	EXCEPTION_ERR_CONTRACT_ADDRESS_COLLISION EXCEPTION = 4
	EXCEPTION_ERR_EXECUTION_REVERTED         EXCEPTION = 5
	EXCEPTION_ERR_MAX_CODE_SIZE_EXCEEDED     EXCEPTION = 6
	EXCEPTION_ERR_INVALID_JUMP               EXCEPTION = 7
	EXCEPTION_ERR_WRITE_PROTECTION           EXCEPTION = 8
	EXCEPTION_ERR_RETURN_DATA_OUT_OF_BOUNDS  EXCEPTION = 9
	EXCEPTION_ERR_GAS_UINT_OVERFLOW          EXCEPTION = 10
	EXCEPTION_ERR_INVALID_CODE               EXCEPTION = 11
	EXCEPTION_ERR_NONCE_UINT_OVERFLOW        EXCEPTION = 12
	EXCEPTION_ERR_OUT_OF_BOUNDS              EXCEPTION = 13
	EXCEPTION_ERR_OVERFLOW                   EXCEPTION = 14
	EXCEPTION_ERR_ADDRESS_NOT_IN_RELATED     EXCEPTION = 15
	EXCEPTION_NONE                           EXCEPTION = -1
)

// Enum value maps for EXCEPTION.
var (
	EXCEPTION_name = map[int32]string{
		0:  "ERR_OUT_OF_GAS",
		1:  "ERR_CODE_STORE_OUT_OF_GAS",
		2:  "ERR_DEPTH",
		3:  "ERR_INSUFFICIENT_BALANCE",
		4:  "ERR_CONTRACT_ADDRESS_COLLISION",
		5:  "ERR_EXECUTION_REVERTED",
		6:  "ERR_MAX_CODE_SIZE_EXCEEDED",
		7:  "ERR_INVALID_JUMP",
		8:  "ERR_WRITE_PROTECTION",
		9:  "ERR_RETURN_DATA_OUT_OF_BOUNDS",
		10: "ERR_GAS_UINT_OVERFLOW",
		11: "ERR_INVALID_CODE",
		12: "ERR_NONCE_UINT_OVERFLOW",
		13: "ERR_OUT_OF_BOUNDS",
		14: "ERR_OVERFLOW",
		15: "ERR_ADDRESS_NOT_IN_RELATED",
		-1: "NONE",
	}
	EXCEPTION_value = map[string]int32{
		"ERR_OUT_OF_GAS":                 0,
		"ERR_CODE_STORE_OUT_OF_GAS":      1,
		"ERR_DEPTH":                      2,
		"ERR_INSUFFICIENT_BALANCE":       3,
		"ERR_CONTRACT_ADDRESS_COLLISION": 4,
		"ERR_EXECUTION_REVERTED":         5,
		"ERR_MAX_CODE_SIZE_EXCEEDED":     6,
		"ERR_INVALID_JUMP":               7,
		"ERR_WRITE_PROTECTION":           8,
		"ERR_RETURN_DATA_OUT_OF_BOUNDS":  9,
		"ERR_GAS_UINT_OVERFLOW":          10,
		"ERR_INVALID_CODE":               11,
		"ERR_NONCE_UINT_OVERFLOW":        12,
		"ERR_OUT_OF_BOUNDS":              13,
		"ERR_OVERFLOW":                   14,
		"ERR_ADDRESS_NOT_IN_RELATED":     15,
		"NONE":                           -1,
	}
)

func (x EXCEPTION) Enum() *EXCEPTION {
	p := new(EXCEPTION)
	*p = x
	return p
}

func (x EXCEPTION) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EXCEPTION) Descriptor() protoreflect.EnumDescriptor {
	return file_receipt_proto_enumTypes[1].Descriptor()
}

func (EXCEPTION) Type() protoreflect.EnumType {
	return &file_receipt_proto_enumTypes[1]
}

func (x EXCEPTION) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EXCEPTION.Descriptor instead.
func (EXCEPTION) EnumDescriptor() ([]byte, []int) {
	return file_receipt_proto_rawDescGZIP(), []int{1}
}

type Receipt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TransactionHash []byte         `protobuf:"bytes,1,opt,name=TransactionHash,proto3" json:"TransactionHash,omitempty"`
	FromAddress     []byte         `protobuf:"bytes,2,opt,name=FromAddress,proto3" json:"FromAddress,omitempty"`
	ToAddress       []byte         `protobuf:"bytes,3,opt,name=ToAddress,proto3" json:"ToAddress,omitempty"`
	Amount          []byte         `protobuf:"bytes,4,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Action          ACTION         `protobuf:"varint,5,opt,name=Action,proto3,enum=transaction.ACTION" json:"Action,omitempty"`
	Status          RECEIPT_STATUS `protobuf:"varint,6,opt,name=Status,proto3,enum=receipt.RECEIPT_STATUS" json:"Status,omitempty"`
	Return          []byte         `protobuf:"bytes,7,opt,name=Return,proto3" json:"Return,omitempty"`
	Exception       EXCEPTION      `protobuf:"varint,8,opt,name=Exception,proto3,enum=receipt.EXCEPTION" json:"Exception,omitempty"`
	GasUsed         uint64         `protobuf:"varint,9,opt,name=GasUsed,proto3" json:"GasUsed,omitempty"`
	GasFee          uint64         `protobuf:"varint,10,opt,name=GasFee,proto3" json:"GasFee,omitempty"`
	EventLogs       []*EventLog    `protobuf:"bytes,11,rep,name=EventLogs,proto3" json:"EventLogs,omitempty"`
}

func (x *Receipt) Reset() {
	*x = Receipt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_receipt_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Receipt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Receipt) ProtoMessage() {}

func (x *Receipt) ProtoReflect() protoreflect.Message {
	mi := &file_receipt_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Receipt.ProtoReflect.Descriptor instead.
func (*Receipt) Descriptor() ([]byte, []int) {
	return file_receipt_proto_rawDescGZIP(), []int{0}
}

func (x *Receipt) GetTransactionHash() []byte {
	if x != nil {
		return x.TransactionHash
	}
	return nil
}

func (x *Receipt) GetFromAddress() []byte {
	if x != nil {
		return x.FromAddress
	}
	return nil
}

func (x *Receipt) GetToAddress() []byte {
	if x != nil {
		return x.ToAddress
	}
	return nil
}

func (x *Receipt) GetAmount() []byte {
	if x != nil {
		return x.Amount
	}
	return nil
}

func (x *Receipt) GetAction() ACTION {
	if x != nil {
		return x.Action
	}
	return ACTION_EMPTY
}

func (x *Receipt) GetStatus() RECEIPT_STATUS {
	if x != nil {
		return x.Status
	}
	return RECEIPT_STATUS_RETURNED
}

func (x *Receipt) GetReturn() []byte {
	if x != nil {
		return x.Return
	}
	return nil
}

func (x *Receipt) GetException() EXCEPTION {
	if x != nil {
		return x.Exception
	}
	return EXCEPTION_ERR_OUT_OF_GAS
}

func (x *Receipt) GetGasUsed() uint64 {
	if x != nil {
		return x.GasUsed
	}
	return 0
}

func (x *Receipt) GetGasFee() uint64 {
	if x != nil {
		return x.GasFee
	}
	return 0
}

func (x *Receipt) GetEventLogs() []*EventLog {
	if x != nil {
		return x.EventLogs
	}
	return nil
}

var File_receipt_proto protoreflect.FileDescriptor

var file_receipt_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x1a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x98, 0x03, 0x0a,
	0x07, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x12, 0x28, 0x0a, 0x0f, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0f, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x61,
	0x73, 0x68, 0x12, 0x20, 0x0a, 0x0b, 0x46, 0x72, 0x6f, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x46, 0x72, 0x6f, 0x6d, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x06, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2b, 0x0a, 0x06, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x52,
	0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2f, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70,
	0x74, 0x2e, 0x52, 0x45, 0x43, 0x45, 0x49, 0x50, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53,
	0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x74, 0x75,
	0x72, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e,
	0x12, 0x30, 0x0a, 0x09, 0x45, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x2e, 0x45, 0x58,
	0x43, 0x45, 0x50, 0x54, 0x49, 0x4f, 0x4e, 0x52, 0x09, 0x45, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x47, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x07, 0x47, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x12, 0x16, 0x0a, 0x06,
	0x47, 0x61, 0x73, 0x46, 0x65, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x47, 0x61,
	0x73, 0x46, 0x65, 0x65, 0x12, 0x31, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67,
	0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f,
	0x6c, 0x6f, 0x67, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x52, 0x09, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x73, 0x2a, 0x55, 0x0a, 0x0e, 0x52, 0x45, 0x43, 0x45, 0x49,
	0x50, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x54,
	0x55, 0x52, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x48, 0x41, 0x4c, 0x54, 0x45,
	0x44, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x54, 0x48, 0x52, 0x45, 0x57, 0x10, 0x02, 0x12, 0x1e,
	0x0a, 0x11, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x10, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x2a, 0xc8,
	0x03, 0x0a, 0x09, 0x45, 0x58, 0x43, 0x45, 0x50, 0x54, 0x49, 0x4f, 0x4e, 0x12, 0x12, 0x0a, 0x0e,
	0x45, 0x52, 0x52, 0x5f, 0x4f, 0x55, 0x54, 0x5f, 0x4f, 0x46, 0x5f, 0x47, 0x41, 0x53, 0x10, 0x00,
	0x12, 0x1d, 0x0a, 0x19, 0x45, 0x52, 0x52, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x53, 0x54, 0x4f,
	0x52, 0x45, 0x5f, 0x4f, 0x55, 0x54, 0x5f, 0x4f, 0x46, 0x5f, 0x47, 0x41, 0x53, 0x10, 0x01, 0x12,
	0x0d, 0x0a, 0x09, 0x45, 0x52, 0x52, 0x5f, 0x44, 0x45, 0x50, 0x54, 0x48, 0x10, 0x02, 0x12, 0x1c,
	0x0a, 0x18, 0x45, 0x52, 0x52, 0x5f, 0x49, 0x4e, 0x53, 0x55, 0x46, 0x46, 0x49, 0x43, 0x49, 0x45,
	0x4e, 0x54, 0x5f, 0x42, 0x41, 0x4c, 0x41, 0x4e, 0x43, 0x45, 0x10, 0x03, 0x12, 0x22, 0x0a, 0x1e,
	0x45, 0x52, 0x52, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x52, 0x41, 0x43, 0x54, 0x5f, 0x41, 0x44, 0x44,
	0x52, 0x45, 0x53, 0x53, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x49, 0x53, 0x49, 0x4f, 0x4e, 0x10, 0x04,
	0x12, 0x1a, 0x0a, 0x16, 0x45, 0x52, 0x52, 0x5f, 0x45, 0x58, 0x45, 0x43, 0x55, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x52, 0x45, 0x56, 0x45, 0x52, 0x54, 0x45, 0x44, 0x10, 0x05, 0x12, 0x1e, 0x0a, 0x1a,
	0x45, 0x52, 0x52, 0x5f, 0x4d, 0x41, 0x58, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x53, 0x49, 0x5a,
	0x45, 0x5f, 0x45, 0x58, 0x43, 0x45, 0x45, 0x44, 0x45, 0x44, 0x10, 0x06, 0x12, 0x14, 0x0a, 0x10,
	0x45, 0x52, 0x52, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x4a, 0x55, 0x4d, 0x50,
	0x10, 0x07, 0x12, 0x18, 0x0a, 0x14, 0x45, 0x52, 0x52, 0x5f, 0x57, 0x52, 0x49, 0x54, 0x45, 0x5f,
	0x50, 0x52, 0x4f, 0x54, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x08, 0x12, 0x21, 0x0a, 0x1d,
	0x45, 0x52, 0x52, 0x5f, 0x52, 0x45, 0x54, 0x55, 0x52, 0x4e, 0x5f, 0x44, 0x41, 0x54, 0x41, 0x5f,
	0x4f, 0x55, 0x54, 0x5f, 0x4f, 0x46, 0x5f, 0x42, 0x4f, 0x55, 0x4e, 0x44, 0x53, 0x10, 0x09, 0x12,
	0x19, 0x0a, 0x15, 0x45, 0x52, 0x52, 0x5f, 0x47, 0x41, 0x53, 0x5f, 0x55, 0x49, 0x4e, 0x54, 0x5f,
	0x4f, 0x56, 0x45, 0x52, 0x46, 0x4c, 0x4f, 0x57, 0x10, 0x0a, 0x12, 0x14, 0x0a, 0x10, 0x45, 0x52,
	0x52, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x10, 0x0b,
	0x12, 0x1b, 0x0a, 0x17, 0x45, 0x52, 0x52, 0x5f, 0x4e, 0x4f, 0x4e, 0x43, 0x45, 0x5f, 0x55, 0x49,
	0x4e, 0x54, 0x5f, 0x4f, 0x56, 0x45, 0x52, 0x46, 0x4c, 0x4f, 0x57, 0x10, 0x0c, 0x12, 0x15, 0x0a,
	0x11, 0x45, 0x52, 0x52, 0x5f, 0x4f, 0x55, 0x54, 0x5f, 0x4f, 0x46, 0x5f, 0x42, 0x4f, 0x55, 0x4e,
	0x44, 0x53, 0x10, 0x0d, 0x12, 0x10, 0x0a, 0x0c, 0x45, 0x52, 0x52, 0x5f, 0x4f, 0x56, 0x45, 0x52,
	0x46, 0x4c, 0x4f, 0x57, 0x10, 0x0e, 0x12, 0x1e, 0x0a, 0x1a, 0x45, 0x52, 0x52, 0x5f, 0x41, 0x44,
	0x44, 0x52, 0x45, 0x53, 0x53, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x49, 0x4e, 0x5f, 0x52, 0x45, 0x4c,
	0x41, 0x54, 0x45, 0x44, 0x10, 0x0f, 0x12, 0x11, 0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x42, 0x2f, 0x0a, 0x25, 0x63, 0x6f, 0x6d,
	0x2e, 0x6d, 0x65, 0x74, 0x61, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x65, 0x64, 0x2e, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x70, 0x74, 0x5a, 0x06, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_receipt_proto_rawDescOnce sync.Once
	file_receipt_proto_rawDescData = file_receipt_proto_rawDesc
)

func file_receipt_proto_rawDescGZIP() []byte {
	file_receipt_proto_rawDescOnce.Do(func() {
		file_receipt_proto_rawDescData = protoimpl.X.CompressGZIP(file_receipt_proto_rawDescData)
	})
	return file_receipt_proto_rawDescData
}

var file_receipt_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_receipt_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_receipt_proto_goTypes = []interface{}{
	(RECEIPT_STATUS)(0), // 0: receipt.RECEIPT_STATUS
	(EXCEPTION)(0),      // 1: receipt.EXCEPTION
	(*Receipt)(nil),     // 2: receipt.Receipt
	(ACTION)(0),         // 3: transaction.ACTION
	(*EventLog)(nil),    // 4: event_log.EventLog
}
var file_receipt_proto_depIdxs = []int32{
	3, // 0: receipt.Receipt.Action:type_name -> transaction.ACTION
	0, // 1: receipt.Receipt.Status:type_name -> receipt.RECEIPT_STATUS
	1, // 2: receipt.Receipt.Exception:type_name -> receipt.EXCEPTION
	4, // 3: receipt.Receipt.EventLogs:type_name -> event_log.EventLog
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_receipt_proto_init() }
func file_receipt_proto_init() {
	if File_receipt_proto != nil {
		return
	}
	file_transaction_proto_init()
	file_event_log_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_receipt_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Receipt); i {
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
			RawDescriptor: file_receipt_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_receipt_proto_goTypes,
		DependencyIndexes: file_receipt_proto_depIdxs,
		EnumInfos:         file_receipt_proto_enumTypes,
		MessageInfos:      file_receipt_proto_msgTypes,
	}.Build()
	File_receipt_proto = out.File
	file_receipt_proto_rawDesc = nil
	file_receipt_proto_goTypes = nil
	file_receipt_proto_depIdxs = nil
}

package receipt

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "gomail/mtn/proto"
	"gomail/mtn/types"
)

type Receipt struct {
	proto *pb.Receipt
}

func NewReceipt(
	transactionHash common.Hash,
	fromAddress common.Address,
	toAddress common.Address,
	amount *big.Int,
	action pb.ACTION,
	status pb.RECEIPT_STATUS,
	returnValue []byte,
	exception pb.EXCEPTION,
	gastFee uint64,
	gasUsed uint64,
	eventLogs []types.EventLog,
) types.Receipt {
	pbEventlogs := make([]*pb.EventLog, len(eventLogs))
	for idx, eventLog := range eventLogs {
		pbEventlogs[idx] = eventLog.Proto()
	}
	proto := &pb.Receipt{
		TransactionHash: transactionHash.Bytes(),
		FromAddress:     fromAddress.Bytes(),
		ToAddress:       toAddress.Bytes(),
		Amount:          amount.Bytes(),
		Action:          action,
		Status:          status,
		Return:          returnValue,
		Exception:       exception,
		GasUsed:         gasUsed,
		GasFee:          gastFee,
		EventLogs:       pbEventlogs,
	}
	return ReceiptFromProto(proto)
}

func ReceiptFromProto(proto *pb.Receipt) types.Receipt {
	return &Receipt{
		proto: proto,
	}
}

// general
func (r *Receipt) FromProto(proto *pb.Receipt) {
	r.proto = proto
}

func (r *Receipt) Unmarshal(b []byte) error {
	receiptPb := &pb.Receipt{}
	err := proto.Unmarshal(b, receiptPb)
	if err != nil {
		return err
	}
	r.proto = receiptPb
	return nil
}

func (r *Receipt) Marshal() ([]byte, error) {
	return proto.Marshal(r.proto)
}

func (r *Receipt) Proto() protoreflect.ProtoMessage {
	return r.proto
}

func (r *Receipt) String() string {
	str := fmt.Sprintf(`
	Transaction hash: %v
	From address: %v
	To address: %v
	Amount: %v
	Action: %v
	Status: %v
	Return: %v
	Exception: %v
	GasUsed: %v
	GasFee: %v
`,
		common.BytesToHash(r.proto.TransactionHash),
		common.BytesToAddress(r.proto.FromAddress),
		common.BytesToAddress(r.proto.ToAddress),
		uint256.NewInt(0).SetBytes(r.proto.Amount),
		r.proto.Action,
		r.proto.Status,
		common.Bytes2Hex(r.proto.Return),
		r.proto.Exception,
		r.proto.GasUsed,
		r.proto.GasFee,
	)
	return str
}

// getter
func (r *Receipt) TransactionHash() common.Hash {
	return common.BytesToHash(r.proto.TransactionHash)
}

func (r *Receipt) FromAddress() common.Address {
	return common.BytesToAddress(r.proto.FromAddress)
}

func (r *Receipt) ToAddress() common.Address {
	return common.BytesToAddress(r.proto.ToAddress)
}

func (r *Receipt) GasUsed() uint64 {
	return r.proto.GasUsed
}

func (r *Receipt) GastFee() uint64 {
	return r.proto.GasFee
}

func (r *Receipt) Amount() *big.Int {
	return big.NewInt(0).SetBytes(r.proto.Amount)
}

func (r *Receipt) Return() []byte {
	return r.proto.Return
}

func (r *Receipt) Status() pb.RECEIPT_STATUS {
	return r.proto.Status
}

func (r *Receipt) Action() pb.ACTION {
	return r.proto.Action
}

func (r *Receipt) EventLogs() []*pb.EventLog {
	return r.proto.EventLogs
}

// setter
func (r *Receipt) UpdateExecuteResult(
	status pb.RECEIPT_STATUS,
	returnValue []byte,
	exception pb.EXCEPTION,
	gasUsed uint64,
	eventLogs []types.EventLog,
) {
	r.proto.Status = status
	r.proto.Return = returnValue
	r.proto.Exception = exception
	r.proto.GasUsed = gasUsed
	pbEventlogs := make([]*pb.EventLog, len(eventLogs))
	for idx, eventLog := range eventLogs {
		pbEventlogs[idx] = eventLog.Proto()
	}
	r.proto.EventLogs = pbEventlogs
}

func (r *Receipt) Json() ([]byte, error) {
	mapReceipt := map[string]interface{}{
		"transaction_hash": hex.EncodeToString(r.TransactionHash().Bytes()),
		"from_address":     hex.EncodeToString(r.FromAddress().Bytes()),
		"to_address":       hex.EncodeToString(r.ToAddress().Bytes()),
		"amount":           hex.EncodeToString(r.Amount().Bytes()),
		"action":           r.Action().Enum(),
		"status":           r.Status().Enum(),
		"return_value":     hex.EncodeToString(r.Return()),
		"exception":        r.proto.Exception.Enum(),
		"gas_fee":          r.proto.GasFee,
		"gas_used":         r.GasUsed(),
	}
	return json.Marshal(mapReceipt)
}

func ReceiptsToProto(receipts []types.Receipt) []*pb.Receipt {
	protoReceipts := make([]*pb.Receipt, len(receipts))
	for i, receipt := range receipts {
		protoReceipts[i] = receipt.Proto().(*pb.Receipt)
	}
	return protoReceipts
}

func ProtoToReceipts(protoReceipts []*pb.Receipt) []types.Receipt {
	receipts := make([]types.Receipt, len(protoReceipts))
	for i, protoReceipt := range protoReceipts {
		receipts[i] = ReceiptFromProto(protoReceipt)
	}
	return receipts
}

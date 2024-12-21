package smart_contract

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"google.golang.org/protobuf/proto"

	pb "gomail/mtn/proto"
	"gomail/mtn/types"
)

type EventLog struct {
	proto *pb.EventLog
}

func NewEventLog(
	blockNumber uint64,
	transactionHash common.Hash,
	address common.Address,
	data []byte,
	topics [][]byte,
) types.EventLog {
	return &EventLog{
		proto: &pb.EventLog{
			BlockNumber:     big.NewInt(int64(blockNumber)).Bytes(),
			TransactionHash: transactionHash.Bytes(),
			Address:         address.Bytes(),
			Data:            data,
			Topics:          topics,
		},
	}
}

// general
func (l *EventLog) Proto() *pb.EventLog {
	return l.proto
}

func (l *EventLog) FromProto(logPb *pb.EventLog) {
	l.proto = logPb
}

func (l *EventLog) Unmarshal(b []byte) error {
	logPb := &pb.EventLog{}
	err := proto.Unmarshal(b, logPb)
	if err != nil {
		return err
	}
	l.FromProto(logPb)
	return nil
}

func (l *EventLog) Marshal() ([]byte, error) {
	return proto.Marshal(l.proto)
}

func (l *EventLog) Copy() types.EventLog {
	cp := &EventLog{}
	cp.FromProto(proto.Clone(l.proto).(*pb.EventLog))
	return cp
}

// getter
func (l *EventLog) Hash() common.Hash {
	b, _ := l.Marshal()
	return crypto.Keccak256Hash(b)
}

func (l *EventLog) Address() common.Address {
	return common.BytesToAddress(l.proto.Address)
}

func (l *EventLog) BlockNumber() string {
	return common.Bytes2Hex(l.proto.BlockNumber)
}

func (l *EventLog) TransactionHash() string {
	return common.Bytes2Hex(l.proto.TransactionHash)
}

func (l *EventLog) Data() string {
	return common.Bytes2Hex(l.proto.Data)
}

func (l *EventLog) Topics() []string {
	topics := make([]string, 0)
	for _, topic := range l.proto.Topics {
		topics = append(topics, common.Bytes2Hex(topic))
	}
	return topics
}

func (l *EventLog) String() string {
	str := fmt.Sprintf(`
	Block Count: %v
	Transaction Hash: %v
	Address: %v
	Data: %v
	Topics: 
	`,
		uint256.NewInt(0).SetBytes(l.proto.BlockNumber),
		common.BytesToHash(l.proto.TransactionHash),
		common.BytesToAddress(l.proto.Address),
		common.Bytes2Hex(l.proto.Data),
	)

	for i, t := range l.proto.Topics {
		str += fmt.Sprintf("\tTopic %v: %v\n", i, common.Bytes2Hex(t))
	}
	return str
}

func EventLogsToProto(eventLogs []types.EventLog) []*pb.EventLog {
	pbEventLogs := make([]*pb.EventLog, len(eventLogs))
	for idx, eventLog := range eventLogs {
		pbEventLogs[idx] = eventLog.Proto()
	}
	return pbEventLogs
}

func EventLogsFromProto(eventLogsPb []*pb.EventLog) []types.EventLog {
	eventLogs := make([]types.EventLog, len(eventLogsPb))
	for idx, eventLog := range eventLogsPb {
		eventLogs[idx] = &EventLog{}
		eventLogs[idx].FromProto(eventLog)
	}
	return eventLogs
}

func GetLogsHash(
	eventLogs []types.EventLog,
) common.Hash {
	hashData := &pb.EventLogsHashData{
		Hashes: make([][]byte, len(eventLogs)),
	}
	for i, v := range eventLogs {
		hashData.Hashes[i] = v.Hash().Bytes()
	}
	b, _ := proto.Marshal(hashData)
	return crypto.Keccak256Hash(b)
}

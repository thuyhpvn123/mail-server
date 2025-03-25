package smart_contract

import (
	"google.golang.org/protobuf/proto"

	pb "gomail/pkg/proto"
	"gomail/types"
)

type EventLogs struct {
	proto *pb.EventLogs
}

func NewEventLogs(eventLogs []types.EventLog) types.EventLogs {
	pbEventLogs := make([]*pb.EventLog, len(eventLogs))
	for idx, eventLog := range eventLogs {
		pbEventLogs[idx] = eventLog.Proto()
	}
	return &EventLogs{
		proto: &pb.EventLogs{
			EventLogs: pbEventLogs,
		},
	}
}

// general
func (l *EventLogs) FromProto(logPb *pb.EventLogs) {
	l.proto = logPb
}

func (l *EventLogs) Unmarshal(b []byte) error {
	logsPb := &pb.EventLogs{}
	err := proto.Unmarshal(b, logsPb)
	if err != nil {
		return err
	}
	l.FromProto(logsPb)
	return nil
}

func (l *EventLogs) Marshal() ([]byte, error) {
	return proto.Marshal(l.proto)
}

func (l *EventLogs) Proto() *pb.EventLogs {
	return l.proto
}

// getter
func (l *EventLogs) EventLogList() []types.EventLog {
	eventLogsPb := l.proto.EventLogs
	eventLogList := make([]types.EventLog, len(eventLogsPb))
	for idx, eventLog := range eventLogsPb {
		eventLogList[idx] = &EventLog{}
		eventLogList[idx].FromProto(eventLog)
	}
	return eventLogList
}

func (l *EventLogs) Copy() types.EventLogs {
	cp := &EventLogs{}
	cp.FromProto(proto.Clone(l.proto).(*pb.EventLogs))
	return cp
}

package network

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	cm "gomail/mtn/common"
	pb "gomail/mtn/proto"
	"gomail/mtn/types/network"
)

type Message struct {
	proto *pb.Message
}

func NewMessage(pbMessage *pb.Message) network.Message {
	return &Message{
		proto: pbMessage,
	}
}

func (m *Message) Marshal() ([]byte, error) {
	return proto.Marshal(m.proto)
}

func (m *Message) Unmarshal(protoStruct protoreflect.ProtoMessage) error {
	err := proto.Unmarshal(m.proto.Body, protoStruct)
	return err
}

func (m *Message) String() string {
	str := fmt.Sprintf(`
	Header:
		Command: %v
		Pubkey: %v
		ToAddress: %v
		Sign: %v
		Version: %v
	Body: %v
`,
		m.proto.Header.Command,
		hex.EncodeToString(m.proto.Header.Pubkey),
		hex.EncodeToString(m.proto.Header.ToAddress),
		hex.EncodeToString(m.proto.Header.Sign),
		m.proto.Header.Version,
		hex.EncodeToString(m.proto.Body),
	)
	return str
}

// getter
func (m *Message) Command() string {
	if m == nil || m.proto == nil || m.proto.Header == nil {
		return ""
	}
	return m.proto.Header.Command
}

func (m *Message) Body() []byte {
	return m.proto.Body
}

func (m *Message) Pubkey() cm.PublicKey {
	return cm.PubkeyFromBytes(m.proto.Header.Pubkey)
}

func (m *Message) Sign() cm.Sign {
	return cm.SignFromBytes(m.proto.Header.Sign)
}

func (m *Message) ToAddress() common.Address {
	return common.BytesToAddress(m.proto.Header.ToAddress)
}

func (m *Message) ID() string {
	return m.proto.Header.ID
}

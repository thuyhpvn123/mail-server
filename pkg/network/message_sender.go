package network

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"gomail/pkg/logger"
	pb "gomail/pkg/proto"
	"gomail/types/network"
)

type MessageSender struct {
	version string
}

func NewMessageSender(
	version string,
) network.MessageSender {
	return &MessageSender{
		version: version,
	}
}

func (s *MessageSender) SendMessage(
	connection network.Connection,
	command string,
	pbMessage protoreflect.ProtoMessage,
) error {
	return SendMessage(
		connection,
		command,
		pbMessage,
		s.version,
	)
}

func (s *MessageSender) SendBytes(
	connection network.Connection,
	command string,
	b []byte,
) error {
	return SendBytes(
		connection,
		command,
		b,
		s.version,
	)
}

func (s *MessageSender) BroadcastMessage(
	mapAddressConnections map[common.Address]network.Connection,
	command string,
	marshaler network.Marshaler,
) error {
	wg := sync.WaitGroup{}
	bytes, err := marshaler.Marshal()
	if err != nil {
		logger.Error(err)
		return err
	}
	for _, con := range mapAddressConnections {
		if con != nil {
			wg.Add(1)
			go func(conn network.Connection, wg *sync.WaitGroup, b []byte) {
				defer wg.Done()
				err := s.SendBytes(conn, command, b)
				if err != nil {
					logger.Error(err)
				}
			}(con, &wg, bytes)
		}
	}
	wg.Wait()
	return nil
}

func getHeaderForCommand(
	command string,
	toAddress common.Address,
	version string,
) *pb.Header {
	return &pb.Header{
		Command:   command,
		Version:   version,
		ToAddress: toAddress.Bytes(),
		ID:        uuid.New().String(),
	}
}

func generateMessage(
	toAddress common.Address,
	command string,
	body []byte,
	version string,
) network.Message {
	messageProto := &pb.Message{
		Header: getHeaderForCommand(
			command,
			toAddress,
			version,
		),
		Body: body,
	}
	message := NewMessage(messageProto)
	return message
}

func SendMessage(
	connection network.Connection,
	command string,
	pbMessage proto.Message,
	version string,
) (err error) {
	body := []byte{}
	if pbMessage != nil {
		body, err = proto.Marshal(pbMessage)
		if err != nil {
			return err
		}
	}
	return SendBytes(connection, command, body, version)
}

func SendBytes(
	connection network.Connection,
	command string,
	bytes []byte,
	version string,
) error {
	if connection == nil {
		return errors.New("nil connection")
	}
	message := generateMessage(connection.Address(), command, bytes, version)
	return connection.SendMessage(message)
}

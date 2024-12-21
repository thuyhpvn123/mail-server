package network

import (
	"github.com/ethereum/go-ethereum/common"
)

type Connection interface {
	// getter
	Address() common.Address
	ConnectionAddress() (string, error)

	RequestChan() (chan Request, chan error)
	Type() string
	String() string

	RemoteAddr() string
	// setter
	Init(common.Address, string)
	SetRealConnAddr(realConnAddr string)

	// other
	SendMessage(message Message) error
	Connect() error
	Disconnect() error
	IsConnect() bool
	ReadRequest()
	Clone() Connection
}

package mock

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"

	"gomail/types/network"
)

type TestConnection struct {
	connected         bool
	address           common.Address
	connectionAddress string
	connectionType    string
	requestChan       chan network.Request
	errorChan         chan error
	handler           func(message network.Message) (network.Request, error)
}

func NewTestConnection(
	address common.Address,
	connectionAddress string,
	connectionType string,
	handler func(message network.Message) (network.Request, error),
) *TestConnection {
	return &TestConnection{
		connected:         true,
		handler:           handler,
		address:           address,
		connectionAddress: connectionAddress,
		connectionType:    connectionType,
	}
}

func (tc *TestConnection) Address() common.Address {
	return tc.address
}

func (tc *TestConnection) ConnectionAddress() (string, error) {
	return tc.connectionAddress, nil
}

func (tc *TestConnection) RequestChan() (chan network.Request, chan error) {
	return tc.requestChan, tc.errorChan
}

func (tc *TestConnection) Type() string {
	return tc.connectionType
}

func (tc *TestConnection) String() string {
	return tc.address.String()
}

func (tc *TestConnection) RemoteAddr() string {
	return "test connection"
}

func (tc *TestConnection) Init(address common.Address, connectionAddress string) {
	tc.address = address
	tc.connectionAddress = connectionAddress
}

func (tc *TestConnection) SendMessage(message network.Message) error {
	go func() {
		request, err := tc.handler(message)
		if err != nil {
			tc.errorChan <- err
		}
		tc.requestChan <- request
	}()
	return nil
}

func (tc *TestConnection) Connect() error {
	tc.connected = true
	return nil
}

func (tc *TestConnection) Disconnect() error {
	tc.errorChan <- errors.New("disconnect")
	tc.connected = false
	return nil
}

func (tc *TestConnection) IsConnect() bool {
	return true
}

func (tc *TestConnection) ReadRequest() {
}

func (tc *TestConnection) Clone() network.Connection {
	return &TestConnection{
		connected:         tc.connected,
		address:           tc.address,
		connectionAddress: tc.connectionAddress,
		connectionType:    tc.connectionType,
		requestChan:       tc.requestChan,
		errorChan:         tc.errorChan,
		handler:           tc.handler,
	}
}

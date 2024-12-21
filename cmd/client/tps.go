package client

import (
	"github.com/ethereum/go-ethereum/common"

	"gomail/cmd/client/pkg/config"
	p_network "gomail/cmd/client/pkg/network"
	"gomail/mtn/bls"
	"gomail/mtn/network"
	"gomail/mtn/types"
	t_network "gomail/mtn/types/network"
)

type TpsClient struct {
	Keypair            *bls.KeyPair
	connectionsManager t_network.ConnectionsManager
	socketServer       t_network.SocketServer
	messageSender      t_network.MessageSender
	AccountStateChan   chan types.AccountState
	receiptChan        chan types.Receipt
}

func NewTpsClient(config *config.TpsConfig) *TpsClient {
	client := &TpsClient{
		Keypair:            bls.NewKeyPair(common.FromHex(config.PrivateKey_)),
		connectionsManager: network.NewConnectionsManager(),
		messageSender:      network.NewMessageSender(""),
		AccountStateChan:   make(chan types.AccountState, 1),
		receiptChan:        make(chan types.Receipt, 1),
	}

	handler := p_network.NewTpsHandler(client.AccountStateChan, client.receiptChan)

	client.socketServer = network.NewSockerServer(
		client.Keypair,
		client.connectionsManager,
		handler,
		"client",
		"",
		config.DnsLink_,
	)

	parentConnection := network.NewConnection(
		common.HexToAddress(config.NodeAddress),
		"",
		config.DnsLink_,
	)

	err := parentConnection.Connect()
	if err != nil {
		panic(err)
	}
	client.connectionsManager.AddParentConnection(parentConnection)
	go client.socketServer.HandleConnection(parentConnection)
	// start socket server
	client.socketServer.OnConnect(parentConnection)
	return client
}

func (c *TpsClient) ParentConnection() t_network.Connection {
	return c.connectionsManager.ParentConnection()
}

func (c *TpsClient) SendBytes(conn t_network.Connection, command string, bData []byte) (err error) {
	return c.messageSender.SendBytes(c.connectionsManager.ParentConnection(), command, bData)
}

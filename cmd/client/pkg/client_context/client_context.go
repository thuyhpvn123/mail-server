package client_context

import (
	"gomail/cmd/client/pkg/config"
	client_types "gomail/cmd/client/types"
	"gomail/pkg/bls"
	"gomail/types/network"
)

type ClientContext struct {
	// config
	Config  *config.ClientConfig
	KeyPair *bls.KeyPair

	// network
	ConnectionsManager network.ConnectionsManager
	MessageSender      network.MessageSender
	SocketServer       network.SocketServer
	Handler            client_types.ClientHandler
}

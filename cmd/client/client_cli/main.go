package main

import (
	"flag"
	"fmt"
	"gomail/cmd/client"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"gomail/cmd/client/pkg/cli"
	"gomail/cmd/client/pkg/client_context"
	c_config "gomail/cmd/client/pkg/config"
	"gomail/cmd/client/pkg/controllers"
	// c_network "gomail/cmd/client/pkg/network"
	"gomail/mtn/bls"
	"gomail/mtn/logger"
	"gomail/mtn/network"
	// "gomail/cmd/client/types"
)

const (
	defaultConfigPath = "config.json"
	defaultLogLevel   = logger.FLAG_DEBUG
)

var (
	// flags
	CONFIG_FILE_PATH string
	LOG_LEVEL        int
)

func main() {
	flag.StringVar(&CONFIG_FILE_PATH, "config", defaultConfigPath, "Config path")
	flag.StringVar(&CONFIG_FILE_PATH, "c", defaultConfigPath, "Config path (shorthand)")

	flag.IntVar(&LOG_LEVEL, "log-level", defaultLogLevel, "Log level")
	flag.IntVar(&LOG_LEVEL, "ll", defaultLogLevel, "Log level (shorthand)")
	flag.Parse()
	clientContext := &client_context.ClientContext{}
	//
	// set logger config
	var loggerConfig = &logger.LoggerConfig{
		Flag:    LOG_LEVEL,
		Outputs: []*os.File{os.Stdout},
	}
	logger.SetConfig(loggerConfig)
	//
	


	//
	config, err := c_config.LoadConfig(CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	clientContext.Config = config.(*c_config.ClientConfig)
	clientContext.KeyPair = bls.NewKeyPair(config.PrivateKey())
	logger.Debug("Running with key pair: " + "\n" + clientContext.KeyPair.String())
	// init message sender
	clientContext.MessageSender = network.NewMessageSender( config.Version())
	// connect to parent
	clientContext.ConnectionsManager = network.NewConnectionsManager()
	// connection to parent

	parentConn := network.NewConnection(
		common.HexToAddress(clientContext.Config.ParentAddress),
		clientContext.Config.ParentConnectionType,
		clientContext.Config.DnsLink(),
	)

	// accountStateChan := make(chan types.AccountState, 1)
	// clientContext.Handler = c_network.NewHandler(accountStateChan, nil)
	// clientContext.SocketServer = network.NewSockerServer(
	// 	config,
	// 	clientContext.KeyPair,
	// 	clientContext.ConnectionsManager,
	// 	clientContext.Handler,
	// )
	err = parentConn.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent %v", err))
		// panic(fmt.Sprintf("error when connect to parent %v", err))
	} else {
		// init connection
		// clientContext.ConnectionsManager.AddParentConnection(parentConn)
		// clientContext.SocketServer.OnConnect(parentConn)
		// go clientContext.SocketServer.HandleConnection(parentConn)
	}
	// init controller
	transactionCtl := controllers.NewTransactionController(clientContext)
	// init and start cli
	c, err := client.NewClient(
		clientContext.Config,
		// accountStateChan,
	)
	if err != nil {
		logger.Error("Error when init client", err)
		panic(err)
	}


	cli := cli.NewCli(
		c,
		clientContext,
		transactionCtl,
		// accountStateChan,
	)
	cli.Start()
}

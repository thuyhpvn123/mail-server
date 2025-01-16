package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"gomail/cmd/client/command"
	"gomail/cmd/client/pkg/client_context"
	c_config "gomail/cmd/client/pkg/config"
	"gomail/cmd/client/pkg/controllers"
	c_network "gomail/cmd/client/pkg/network"
	client_types "gomail/cmd/client/types"
	"gomail/mtn/bls"
	p_common "gomail/mtn/common"
	"gomail/mtn/logger"
	p_network "gomail/mtn/network"
	pb "gomail/mtn/proto"
	"gomail/mtn/types"
	t_network "gomail/mtn/types/network"
	"github.com/ethereum/go-ethereum/crypto"
)

type Client struct {
	clientContext *client_context.ClientContext

	mu                    sync.Mutex
	accountStateChan      chan types.AccountState
	receiptChan           chan types.Receipt
	deviceKeyChan           chan types.LastDeviceKey
	transactionErrorChan  chan types.TransactionError
	transactionController client_types.TransactionController
	subscribeSCAddresses  []common.Address
}

// var client = Client{}

func NewClient(
	config *c_config.ClientConfig,
) (*Client, error) {
	clientContext := &client_context.ClientContext{
		Config: config,
	}
	client := Client{
		clientContext:        clientContext,
		accountStateChan:     make(chan types.AccountState, 1),
		receiptChan:          make(chan types.Receipt, 1),
		deviceKeyChan:          make(chan types.LastDeviceKey, 1),
		transactionErrorChan: make(chan types.TransactionError, 1),
	}

	clientContext.KeyPair = bls.NewKeyPair(config.PrivateKey())
	clientContext.MessageSender = p_network.NewMessageSender(
		config.Version(),
	)
	clientContext.ConnectionsManager = p_network.NewConnectionsManager()
	parentConn := p_network.NewConnection(
		common.HexToAddress(config.ParentAddress),
		config.ParentConnectionType,
		config.DnsLink(),
	)
	clientContext.Handler = c_network.NewHandler(
		client.accountStateChan,
		client.receiptChan,
		client.deviceKeyChan,
		client.transactionErrorChan,
	)
	clientContext.SocketServer = p_network.NewSockerServer(
		clientContext.KeyPair,
		clientContext.ConnectionsManager,
		clientContext.Handler,
		config.NodeType(),
		config.Version(),
		config.DnsLink(),
	)
	err := parentConn.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent: %v", err))
		return nil, err
	} else {
		// init connection
		clientContext.ConnectionsManager.AddParentConnection(parentConn)
		clientContext.SocketServer.OnConnect(parentConn)
		go clientContext.SocketServer.HandleConnection(parentConn)
	}
	client.transactionController = controllers.NewTransactionController(
		clientContext,
	)
	return &client, nil
}

func (client *Client) ReconnectToParent(conn t_network.Connection) error {
	err := conn.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent. %v", err))
		return err
	} else {
		client.clientContext.ConnectionsManager.AddParentConnection(conn)
		client.clientContext.SocketServer.OnConnect(conn)
		go client.clientContext.SocketServer.HandleConnection(conn)
	}
	return nil
}

func (client *Client) SendTransaction(
	toAddress common.Address,
	amount *big.Int,
	action pb.ACTION,
	data []byte,
	relatedAddress []common.Address,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
) (types.Receipt, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	// get account state
	parentConn := client.clientContext.ConnectionsManager.ParentConnection()
	if !parentConn.IsConnect() {
		err := client.ReconnectToParent(parentConn)
		if err != nil {
			return nil, err
		}
	}

	client.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetAccountState,
		client.clientContext.KeyPair.Address().Bytes(),
	)
	lastDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	newDeviceKey := common.HexToHash(
		"290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563",
	)

	select {
	case as := <-client.accountStateChan:
		lastHash := as.LastHash()
		pendingBalance := as.PendingBalance()

		bRelatedAddresses := make([][]byte, len(relatedAddress))
		for i, v := range relatedAddress {
			bRelatedAddresses[i] = v.Bytes()
		}
		_, err := client.transactionController.SendTransaction(
			lastHash,
			toAddress,
			pendingBalance,
			amount,
			maxGas,
			maxGasPrice,
			maxTimeUse,
			action,
			data,
			bRelatedAddresses,
			lastDeviceKey,
			newDeviceKey,
			nil,
		)
		if err != nil {
			return nil, err
		}

		select {
		case receipt := <-client.receiptChan:
			return receipt, nil
		case <-time.After(10 * time.Second):
			logger.DebugP("Timeout recive receipt")
			return nil, errors.New("err timeout when SendTransaction")
		}
	case <-time.After(10 * time.Second):
		logger.DebugP("Timeout get account state")
		return nil, errors.New("err timeout when SendTransaction")
	}
}

func (client *Client) SendTransactionWithCommission(
	toAddress common.Address,
	amount *big.Int,
	action pb.ACTION,
	data []byte,
	relatedAddress []common.Address,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
	commissionPrivateKey []byte,
) (types.Receipt, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	// get account state
	parentConn := client.clientContext.ConnectionsManager.ParentConnection()
	if !parentConn.IsConnect() {
		err := client.ReconnectToParent(parentConn)
		if err != nil {
			return nil, err
		}
	}

	client.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetAccountState,
		client.clientContext.KeyPair.Address().Bytes(),
	)

	lastDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	newDeviceKey := common.HexToHash(
		"290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563",
	)
	select {
	case as := <-client.accountStateChan:
		lastHash := as.LastHash()
		pendingBalance := as.PendingBalance()

		bRelatedAddresses := make([][]byte, len(relatedAddress))
		for i, v := range relatedAddress {
			bRelatedAddresses[i] = v.Bytes()
		}
		_, err := client.transactionController.SendTransaction(
			lastHash,
			toAddress,
			pendingBalance,
			amount,
			maxGas,
			maxGasPrice,
			maxTimeUse,
			action,
			data,
			bRelatedAddresses,
			lastDeviceKey,
			newDeviceKey,
			commissionPrivateKey,
		)
		if err != nil {
			return nil, err
		}

		select {
		case receipt := <-client.receiptChan:
			return receipt, nil
		case <-time.After(10 * time.Second):
			logger.DebugP("Timeout recive receipt")
			return nil, errors.New("err timeout when SendTransaction")
		}
	case <-time.After(10 * time.Second):
		logger.DebugP("Timeout get account state")
		return nil, errors.New("err timeout when SendTransaction")
	}
}


func (client *Client) SendTransactionWithDeviceKey(
	toAddress common.Address,
	amount *big.Int,
	action pb.ACTION,
	data []byte,
	relatedAddress []common.Address,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
) (types.Receipt, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	// get account state
	parentConn := client.clientContext.ConnectionsManager.ParentConnection()
	if !parentConn.IsConnect() {
		err := client.ReconnectToParent(parentConn)
		if err != nil {
			return nil, err
		}
	}

	client.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetAccountState,
		client.clientContext.KeyPair.Address().Bytes(),
	)

	// lastDeviceKey := common.HexToHash(
	// 	"0000000000000000000000000000000000000000000000000000000000000000",
	// )
	// newDeviceKey := common.HexToHash(
	// 	"290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563",
	// )

	select {
	case as := <-client.accountStateChan:

		lastHash := as.LastHash()
		pendingBalance := as.PendingBalance()
		
		// lastDeviceKey := common.HexToHash(string(as.Proto().GetDeviceKey()))
		fmt.Println("lastDeviceKey:",lastHash)

		err := client.clientContext.MessageSender.SendBytes(
			parentConn,
			"GetDeviceKey",
			lastHash.Bytes(),
			// make([]byte, 0),
		)


		if err != nil {
			return nil, err
		}

		select {
		case receiveDeviceKey := <-client.deviceKeyChan:
				

			TransactionHash := receiveDeviceKey.TransactionHash
			
			lastDeviceKey := common.HexToHash(
				hex.EncodeToString(receiveDeviceKey.LastDeviceKeyFromServer),
			)

			// logger.DebugP("TransactionHash", hex.EncodeToString(TransactionHash))
			// logger.DebugP("lastDeviceKey", lastDeviceKey)


			rawNewDeviceKeyBytes := []byte( fmt.Sprintf("%s-%d", hex.EncodeToString(TransactionHash), time.Now().Unix()) )

			rawNewDeviceKey := crypto.Keccak256(rawNewDeviceKeyBytes)

			// logger.DebugP("rawNewDeviceKey", hex.EncodeToString(rawNewDeviceKey))

			newDeviceKey := crypto.Keccak256Hash(rawNewDeviceKey)
			

			bRelatedAddresses := make([][]byte, len(relatedAddress))
			for i, v := range relatedAddress {
				bRelatedAddresses[i] = v.Bytes()
			}
			_, err = client.transactionController.SendTransactionWithDeviceKey(
				lastHash,
				toAddress,
				pendingBalance,
				amount,
				maxGas,
				maxGasPrice,
				maxTimeUse,
				action,
				data,
				bRelatedAddresses,
				lastDeviceKey,
				newDeviceKey,
				nil,
				rawNewDeviceKey,
			)
			if err != nil {
				return nil, err
			}

			select {
			case receipt := <-client.receiptChan:
				return receipt, nil
			case <-time.After(10 * time.Second):
				logger.DebugP("Timeout recive receipt")
				return nil, errors.New("err timeout when SendTransaction")
			}

		case <-time.After(10 * time.Second):
			logger.DebugP("Timeout recive receipt GetDeviceKey")
			return nil, errors.New("err timeout when GetDeviceKey")
		}




		
	case <-time.After(10 * time.Second):
		logger.DebugP("Timeout get account state")
		return nil, errors.New("err timeout when SendTransaction")
	}
}

func (client *Client) AccountState(address common.Address) (types.AccountState, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	// get account state
	parentConn := client.clientContext.ConnectionsManager.ParentConnection()
	client.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetAccountState,
		address.Bytes(),
	)
	as := <-client.accountStateChan
	return as, nil
}

func (client *Client) Get(address common.Address) (types.AccountState, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	// get account state
	parentConn := client.clientContext.ConnectionsManager.ParentConnection()
	client.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetAccountState,
		address.Bytes(),
	)
	as := <-client.accountStateChan
	return as, nil
}

func NewStorageClient(
	config *c_config.ClientConfig,
	listSCAddress []common.Address,
) (*Client, error) {
	clientContext := &client_context.ClientContext{
		Config: config,
	}

	client := Client{
		clientContext:        clientContext,
		accountStateChan:     make(chan types.AccountState, 1),
		receiptChan:          make(chan types.Receipt, 1),
		transactionErrorChan: make(chan types.TransactionError, 1),
		subscribeSCAddresses: listSCAddress,
	}

	clientContext.KeyPair = bls.NewKeyPair(config.PrivateKey())
	clientContext.MessageSender = p_network.NewMessageSender(
		config.Version(),
	)
	clientContext.ConnectionsManager = p_network.NewConnectionsManager()
	parentConn := p_network.NewConnection(
		common.HexToAddress(config.ParentAddress),
		config.ParentConnectionType,
		config.DnsLink(),
	)
	clientContext.Handler = c_network.NewHandler(
		client.accountStateChan,
		client.receiptChan,
		client.deviceKeyChan,
		client.transactionErrorChan,
	)
	clientContext.SocketServer = p_network.NewSockerServer(
		clientContext.KeyPair,
		clientContext.ConnectionsManager,
		clientContext.Handler,
		config.NodeType(),
		config.Version(),
		config.DnsLink(),
	)
	// #debug
	// parentConn.SetRealConnAddr("192.168.1.253:7101")

	err := parentConn.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent! %v", err))
		return nil, err
	} else {
		// init connection
		clientContext.ConnectionsManager.AddParentConnection(parentConn)
		clientContext.SocketServer.OnConnect(parentConn)
		go clientContext.SocketServer.HandleConnection(parentConn)
	}

	for _, address := range listSCAddress {
		err = client.clientContext.MessageSender.SendBytes(parentConn, command.SubscribeToAddress, address.Bytes())
		if err != nil {
			return nil, fmt.Errorf("unable to send subscribe")
		}
	}

	client.transactionController = controllers.NewTransactionController(
		clientContext,
	)
	evenLogsChan := make(chan types.EventLogs)
	client.clientContext.Handler.(*c_network.Handler).SetEventLogsChan(evenLogsChan)
	client.clientContext.SocketServer.AddOnDisconnectedCallBack(client.RetryConnectToStorage)

	return &client, nil
}

func (client *Client) Subcribe(
	storageAddress common.Address,
	smartContractAddress common.Address,
) (chan types.EventLogs, error) {
	storageConnection := p_network.NewConnection(
		storageAddress,
		p_common.STORAGE_CONNECTION_TYPE,
		client.clientContext.Config.DnsLink(),
	)
	err := storageConnection.Connect()
	if err != nil {
		logger.Error("Unable to connect to storage", err)
		return nil, fmt.Errorf("unable to connect to storage")
	}
	go client.clientContext.SocketServer.HandleConnection(storageConnection)

	err = client.clientContext.MessageSender.SendBytes(
		storageConnection,
		command.SubscribeToAddress,
		smartContractAddress.Bytes(),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to send subscribe")
	}
	evenLogsChan := make(chan types.EventLogs)
	client.clientContext.Handler.(*c_network.Handler).SetEventLogsChan(evenLogsChan)
	return evenLogsChan, nil
}

func (client *Client) Subcribes(
	storageAddress common.Address,
	listSCAddress []common.Address,
) (chan types.EventLogs, error) {
	storageConnection := p_network.NewConnection(
		storageAddress,
		p_common.STORAGE_CONNECTION_TYPE,
		client.clientContext.Config.DnsLink(),
	)
	err := storageConnection.Connect()
	if err != nil {
		logger.Error("Unable to connect to storage", err)
		return nil, fmt.Errorf("unable to connect to storage")
	}
	go client.clientContext.SocketServer.HandleConnection(storageConnection)

	for _, address := range listSCAddress {
		err = client.clientContext.MessageSender.SendBytes(
			storageConnection,
			command.SubscribeToAddress,
			address.Bytes(),
		)
		if err != nil {
			return nil, fmt.Errorf("unable to send subscribe")
		}
	}

	evenLogsChan := make(chan types.EventLogs)
	client.clientContext.Handler.(*c_network.Handler).SetEventLogsChan(evenLogsChan)
	return evenLogsChan, nil
}

func (client *Client) ParentSubcribes(
	listSCAddress []common.Address,
) (chan types.EventLogs, error) {
	for _, address := range listSCAddress {
		err := client.clientContext.MessageSender.SendBytes(
			client.clientContext.ConnectionsManager.ParentConnection(),
			command.SubscribeToAddress,
			address.Bytes(),
		)
		if err != nil {
			return nil, fmt.Errorf("unable to send subscribe")
		}
	}

	evenLogsChan := make(chan types.EventLogs)
	client.clientContext.Handler.(*c_network.Handler).SetEventLogsChan(evenLogsChan)
	client.clientContext.SocketServer.AddOnDisconnectedCallBack(
		client.clientContext.SocketServer.StopAndRetryConnectToParent,
	)

	return evenLogsChan, nil
}

func (client *Client) RetryConnectToStorage(conn t_network.Connection) {
	for {
		<-time.After(5 * time.Second)
		parentConn := client.clientContext.ConnectionsManager.ParentConnection()
		if !parentConn.IsConnect() {
			err := client.ReconnectToParent(parentConn)
			if err != nil {
				logger.Warn(fmt.Sprintf("error when retry connect to parent %v", err))
				continue
			}
		}
		panic("panic when retry connect")
	}
}
func (client *Client) GetEventLogsChan() chan types.EventLogs {
	return client.clientContext.Handler.(*c_network.Handler).GetEventLogsChan()
}
func (client *Client) Close() {
	// remove parent connection to avoid reconnect
	client.clientContext.ConnectionsManager.AddParentConnection(nil)
	client.clientContext.SocketServer.Stop()
}

func (client *Client) SendQueryLogs(bQuery []byte) {
	client.mu.Lock()
	defer client.mu.Unlock()
	// get account state

	parentConn := client.clientContext.ConnectionsManager.ParentConnection()
	client.clientContext.MessageSender.SendBytes(
		parentConn,
		command.QueryLogs,
		bQuery,
	)
}

func (client *Client) NewEventLogsChan() chan types.EventLogs {
	evenLogsChan := make(chan types.EventLogs)
	client.clientContext.Handler.(*c_network.Handler).SetEventLogsChan(evenLogsChan)
	return evenLogsChan
}

func (client *Client) SendTransactionWithFullInfo(
	toAddress common.Address,
	amount *big.Int,
	maxGas uint64,
	maxGasFee uint64,
	maxTimeUse uint64,
	action pb.ACTION,
	data []byte,
	relatedAddress []common.Address,
	lastDeviceKey common.Hash,
	newDeviceKey common.Hash,
	commissionPrivateKey []byte,
) (types.Receipt, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	// get account state
	parentConn := client.clientContext.ConnectionsManager.ParentConnection()
	if !parentConn.IsConnect() {
		err := client.ReconnectToParent(parentConn)
		if err != nil {
			return nil, err
		}
	}

	client.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetAccountState,
		client.clientContext.KeyPair.Address().Bytes(),
	)

	select {
	case as := <-client.accountStateChan:
		lastHash := as.LastHash()
		pendingBalance := as.PendingBalance()

		bRelatedAddresses := make([][]byte, len(relatedAddress))
		for i, v := range relatedAddress {
			bRelatedAddresses[i] = v.Bytes()
		}
		_, err := client.transactionController.SendTransaction(
			lastHash,
			toAddress,
			pendingBalance,
			amount,
			maxGas,
			maxGasFee,
			maxTimeUse,
			action,
			data,
			bRelatedAddresses,
			lastDeviceKey,
			newDeviceKey,
			nil,
		)
		if err != nil {
			return nil, err
		}

		select {
		case receipt := <-client.receiptChan:
			return receipt, nil
		case <-time.After(10 * time.Second):
			logger.DebugP("Timeout recive receipt")
			return nil, errors.New("err timeout when SendTransaction")
		}
	case <-time.After(10 * time.Second):
		logger.DebugP("Timeout get account state")
		return nil, errors.New("err timeout when SendTransaction")
	}
}

func (s *Client) GetMtnAddress() common.Address {
	return s.clientContext.KeyPair.Address()
}

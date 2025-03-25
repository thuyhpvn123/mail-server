package client

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"gomail/cmd/client/command"
	"gomail/cmd/client/pkg/client_context"
	c_config "gomail/cmd/client/pkg/config"
	"gomail/cmd/client/pkg/controllers"
	c_network "gomail/cmd/client/pkg/network"
	client_types "gomail/cmd/client/types"
	"gomail/pkg/bls"
	p_common "gomail/pkg/common"
	"gomail/pkg/logger"
	p_network "gomail/pkg/network"
	pb "gomail/pkg/proto"
	"gomail/pkg/transaction"
	"gomail/types"
	t_network "gomail/types/network"
)

type Client struct {
	clientContext *client_context.ClientContext

	// mu               sync.Mutex
	accountStateChan chan types.AccountState
	receiptChan      chan types.Receipt
	deviceKeyChan    chan types.LastDeviceKey

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
		clientContext:    clientContext,
		accountStateChan: make(chan types.AccountState, 2000),
		receiptChan:      make(chan types.Receipt, 1),
		deviceKeyChan:    make(chan types.LastDeviceKey, 1),

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
	parentConn.SetRealConnAddr(config.ParentConnectionAddress)
	clientContext.Handler = c_network.NewHandler(
		client.accountStateChan,
		client.receiptChan,
		client.deviceKeyChan,
		client.transactionErrorChan,
	)
	clientContext.SocketServer = p_network.NewSocketServer(
		clientContext.KeyPair,
		clientContext.ConnectionsManager,
		clientContext.Handler,
		config.NodeType(),
		config.Version(),
		config.DnsLink(),
	)
	err := parentConn.Connect()
	logger.Info(err)
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent %v", err))
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

func (client *Client) GetClientContext() *client_context.ClientContext {
	return client.clientContext
}

func (client *Client) GetTransactionController() client_types.TransactionController {
	return client.transactionController
}
func (client *Client) GetAccountStateChan() chan types.AccountState {
	return client.accountStateChan
}

func (client *Client) ReconnectToParent(conn t_network.Connection) error {
	err := conn.Connect()
	if err != nil {
		return err
	} else {
		client.clientContext.ConnectionsManager.AddParentConnection(conn)
		client.clientContext.SocketServer.OnConnect(conn)
		go client.clientContext.SocketServer.HandleConnection(conn)
	}
	return nil
}

func (client *Client) SendTransaction(
	fromAddress common.Address,
	toAddress common.Address,
	amount *big.Int,
	action pb.ACTION,
	data []byte,
	relatedAddress []common.Address,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,

) (types.Receipt, error) {
	// client.mu.Lock()
	// defer client.mu.Unlock()
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
		fromAddress.Bytes(),
	)
	lastDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	newDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
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
			fromAddress,
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
			as.Nonce(),
			client.clientContext.Config.ChainId,
		)
		if err != nil {
			return nil, err
		}

		select {
		case receipt := <-client.receiptChan:
			return receipt, nil
			// case <-time.After(10 * time.Second):
			// 	logger.DebugP("Timeout recive receipt")
			// 	return nil, errors.New("err timeout when SendTransaction")
		}
		// case <-time.After(10 * time.Second):
		// 	logger.DebugP("Timeout get account state")
		// 	return nil, errors.New("err timeout when SendTransaction")
	}
}

func (client *Client) ReadTransaction(
	fromAddress common.Address,
	toAddress common.Address,
	amount *big.Int,
	action pb.ACTION,
	data []byte,
	relatedAddress []common.Address,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,

) (types.Receipt, error) {
	// client.mu.Lock()
	// defer client.mu.Unlock()
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
		fromAddress.Bytes(),
	)
	lastDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	newDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	select {
	case as := <-client.accountStateChan:
		lastHash := as.LastHash()
		pendingBalance := as.PendingBalance()

		bRelatedAddresses := make([][]byte, len(relatedAddress))
		for i, v := range relatedAddress {
			bRelatedAddresses[i] = v.Bytes()
		}
		_, err := client.transactionController.ReadTransaction(
			lastHash,
			fromAddress,
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
			as.Nonce(),
			client.clientContext.Config.ChainId,
		)
		if err != nil {
			return nil, err
		}

		select {
		case receipt := <-client.receiptChan:
			return receipt, nil
			// case <-time.After(10 * time.Second):
			// 	logger.DebugP("Timeout recive receipt")
			// 	return nil, errors.New("err timeout when SendTransaction")
		}
		// case <-time.After(10 * time.Second):
		// 	logger.DebugP("Timeout get account state")
		// 	return nil, errors.New("err timeout when SendTransaction")
	}
}

func (client *Client) AddAccountForClient(privateKey string, chainId string) (types.Receipt, error) {
	// Giải mã private key từ chuỗi hex sang dạng byte
	prk, _ := hex.DecodeString(privateKey)
	// Lấy public key BLS từ client context
	blsPublicKey := client.clientContext.KeyPair.PublicKey().Bytes()
	// blsPublicKey[len(blsPublicKey)-1] ^= 0x01
	// Đảo bit cuối để thay đổi giá trị

	// Ghép nối public key BLS và chain ID (là string) thành một mảng byte
	message := append(blsPublicKey, []byte(chainId)...)
	// Tính toán hash của message sử dụng Keccak256
	hash := crypto.Keccak256(message)
	// Ký hash sử dụng private key và thuật toán secp256k1
	sig, err := secp256k1.Sign(hash, prk)
	if err != nil {
		logger.Error("Error Sign", err)
	}

	// Khôi phục public key từ chữ ký và hash sử dụng thuật toán secp256k1
	pbk, err := secp256k1.RecoverPubkey(hash, sig)
	if err != nil {
		logger.Error("Error RecoverPubkey", err)
	}
	// Ghép nối public key BLS và chữ ký thành một mảng byte
	combined := make([]byte, len(blsPublicKey)+len(sig)) // Tạo một slice mới có kích thước đủ lớn
	copy(combined, blsPublicKey)                         // Sao chép blsPublicKey vào slice mới
	copy(combined[len(blsPublicKey):], sig)              // Sao chép chữ ký vào slice mới

	// Tính toán địa chỉ ví từ public key đã khôi phục
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pbk[1:])[12:])
	// Khởi tạo giá trị amount là 1
	bigI := big.NewInt(1)
	// Tạo dữ liệu gọi hàm (call data) từ mảng byte combined
	callData := transaction.NewCallData(combined)
	// Marshal call data sang dạng byte
	cData, _ := callData.Marshal()

	// Gửi giao dịch với dữ liệu đã chuẩn bị
	receipt, err := client.SendTransactionWithDeviceKey(
		addr, // Địa chỉ ví
		addr, // Địa chỉ nhận
		bigI,
		pb.ACTION_EMPTY,
		cData,
		nil,
		100000,
		1000000000,
		6000,
	)
	return receipt, err
}

func (client *Client) SendTransactionWithDeviceKey(
	fromAddress common.Address,
	toAddress common.Address,
	amount *big.Int,
	action pb.ACTION,
	data []byte,
	relatedAddress []common.Address,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
) (types.Receipt, error) {
	// Lấy kết nối tới Parent Node
	parentConn := client.clientContext.ConnectionsManager.ParentConnection()
	if !parentConn.IsConnect() {
		err := client.ReconnectToParent(parentConn)
		if err != nil {
			return nil, err
		}
	}

	// Gửi yêu cầu lấy trạng thái tài khoản
	client.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetAccountState,
		fromAddress.Bytes(),
	)

	// Lắng nghe tài khoản trong kênh accountStateChan bằng for range
	for as := range client.accountStateChan {
		// Nếu không phải tài khoản mong muốn, tiếp tục lắng nghe mà không bỏ dữ liệu
		if as.Address() != fromAddress {
			// Gửi lại dữ liệu cho luồng khác đọc (không bỏ dữ liệu)
			client.accountStateChan <- as
			time.Sleep(100 * time.Millisecond) // Delay trước khi tiếp tục lặp
			continue
		}

		logger.Info("as.Address: 2 ", as)

		// Nếu tìm thấy tài khoản phù hợp, xử lý giao dịch
		lastHash := as.LastHash()
		pendingBalance := as.PendingBalance()

		err := client.clientContext.MessageSender.SendBytes(
			parentConn,
			"GetDeviceKey",
			lastHash.Bytes(),
		)

		if err != nil {
			return nil, err
		}

		// Lắng nghe deviceKey từ server
		receiveDeviceKey := <-client.deviceKeyChan
		TransactionHash := receiveDeviceKey.TransactionHash
		lastDeviceKey := common.HexToHash(
			hex.EncodeToString(receiveDeviceKey.LastDeviceKeyFromServer),
		)
		logger.Info(lastHash, common.BytesToHash(receiveDeviceKey.TransactionHash))
		// Tạo khóa thiết bị mới
		rawNewDeviceKeyBytes := []byte(fmt.Sprintf("%s-%d", hex.EncodeToString(TransactionHash), time.Now().Unix()))
		rawNewDeviceKey := crypto.Keccak256(rawNewDeviceKeyBytes)
		newDeviceKey := crypto.Keccak256Hash(rawNewDeviceKey)

		// Chuyển đổi danh sách địa chỉ liên quan sang mảng byte
		bRelatedAddresses := make([][]byte, len(relatedAddress))
		for i, v := range relatedAddress {
			bRelatedAddresses[i] = v.Bytes()
		}
		logger.Info("fromAddress: toAddress ", fromAddress, toAddress, as.Nonce())

		// Gửi giao dịch với device key
		_, err = client.transactionController.SendTransactionWithDeviceKey(
			lastHash,
			fromAddress,
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
			as.Nonce(),
			rawNewDeviceKey,
			client.clientContext.Config.ChainId,
		)
		if err != nil {
			return nil, err
		}

		// Chờ biên lai giao dịch (receipt)
		receipt := <-client.receiptChan
		return receipt, nil
	}

	// Nếu kênh accountStateChan bị đóng, trả lỗi
	return nil, fmt.Errorf("account state channel closed unexpectedly")
}

func (client *Client) AccountState(address common.Address) (types.AccountState, error) {
	// client.mu.Lock()
	// defer client.mu.Unlock()
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
	// client.mu.Lock()
	// defer client.mu.Unlock()
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
	clientContext.SocketServer = p_network.NewSocketServer(
		clientContext.KeyPair,
		clientContext.ConnectionsManager,
		clientContext.Handler,
		config.NodeType(),
		config.Version(),
		config.DnsLink(),
	)
	err := parentConn.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent %v", err))
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
	// client.mu.Lock()
	// defer client.mu.Unlock()
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
	fromAddress common.Address,

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
	// client.mu.Lock()
	// defer client.mu.Unlock()
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
		fromAddress.Bytes(),
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
			fromAddress,
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
			as.Nonce(),
			client.clientContext.Config.ChainId,
		)
		if err != nil {
			return nil, err
		}

		select {
		case receipt := <-client.receiptChan:
			return receipt, nil
		}
	}
}

func (s *Client) GetMtnAddress() common.Address {
	return s.clientContext.KeyPair.Address()
}

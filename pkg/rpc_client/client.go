package rpc_client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"gomail/pkg/bls"
	"gomail/pkg/state"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	mt_common "gomail/pkg/common"
	"gomail/pkg/logger"
	mt_proto "gomail/pkg/proto"

	mt_transaction "gomail/pkg/transaction"
	mt_types "gomail/types"
)

// Client struct chứa các kết nối HTTP và WebSocket
type ClientRPC struct {
	HttpConn *http.Client
	WsConn   *websocket.Conn
	UrlHTTP  string
	UrlWS    string
	KeyPair  *bls.KeyPair
	ChainId  *big.Int
}
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type JSONRPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      interface{}   `json:"id"`
}

type JSONRPCResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"` // Sử dụng con trỏ và thêm `omitempty`
	Error   *JSONRPCError `json:"error,omitempty"`  // Sử dụng con trỏ và thêm `omitempty`
	Id      interface{}   `json:"id"`
}

// NewClient tạo một đối tượng Client mới
func NewClientRPC(urlHTTP, urlWS, privateKey string, chainId *big.Int) (*ClientRPC, error) {
	// Khởi tạo kết nối HTTP
	httpConn := &http.Client{}

	// Khởi tạo kết nối WebSocket
	// WsConn, _, err := websocket.DefaultDialer.Dial(urlWS, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("không thể kết nối WebSocket: %w", err)
	// }

	keyPair := bls.NewKeyPair(common.FromHex(privateKey))
	return &ClientRPC{
		HttpConn: httpConn,
		// WsConn:   WsConn,
		UrlHTTP: urlHTTP,
		UrlWS:   urlWS,
		KeyPair: keyPair,
		ChainId: chainId,
	}, nil
}

// SendHTTPRequest gửi yêu cầu HTTP đến server
func (c *ClientRPC) SendHTTPRequest(request *JSONRPCRequest) *JSONRPCResponse {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return &JSONRPCResponse{
			Jsonrpc: "2.0",
			Error: &JSONRPCError{
				Code:    -1,
				Message: err.Error(),
			}, // Chuyển đổi lỗi thành JSONRPCError
			Id: request.Id,
		}
	}

	resp, err := c.HttpConn.Post(c.UrlHTTP, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return &JSONRPCResponse{
			Jsonrpc: "2.0",
			Error: &JSONRPCError{
				Code:    -1,
				Message: err.Error(),
			}, // Chuyển đổi lỗi thành JSONRPCError
			Id: request.Id,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &JSONRPCResponse{
			Jsonrpc: "2.0",
			Error: &JSONRPCError{
				Code:    -1,
				Message: err.Error(),
			}, // Chuyển đổi lỗi thành JSONRPCError
			Id: request.Id,
		}
	}

	var response JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return &JSONRPCResponse{
			Jsonrpc: "2.0",
			Error: &JSONRPCError{
				Code:    -1,
				Message: err.Error(),
			}, // Chuyển đổi lỗi thành JSONRPCError
			Id: request.Id,
		}
	}

	return &response
}

// SendWSRequest gửi yêu cầu WebSocket đến server
func (c *ClientRPC) SendWSRequest(request *JSONRPCRequest) *JSONRPCResponse {
	if err := c.WsConn.WriteJSON(request); err != nil {
		return &JSONRPCResponse{
			Jsonrpc: "2.0",
			Error: &JSONRPCError{
				Code:    -1,
				Message: err.Error(),
			}, // Chuyển đổi lỗi thành JSONRPCError
			Id: request.Id,
		}
	}

	var response JSONRPCResponse
	if err := c.WsConn.ReadJSON(&response); err != nil {
		return &JSONRPCResponse{
			Jsonrpc: "2.0",
			Error: &JSONRPCError{
				Code:    -1,
				Message: err.Error(),
			}, // Chuyển đổi lỗi thành JSONRPCError
			Id: request.Id,
		}
	}

	return &response
}

func (c *ClientRPC) GetAccountState(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (mt_types.AccountState, error) {
	logger.Info("GetAccountState")

	request := &JSONRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getAccountState",
		Params:  []interface{}{address.String(), blockNrOrHash.String()}, // Thay đổi thành []interface{}
		Id:      1,
	}

	response := c.SendHTTPRequest(request)
	if response.Error != nil {
		return nil, fmt.Errorf("lỗi từ server: code=%d, message=%s", response.Error.Code, response.Error.Message)
	}
	if response.Result != nil {
		resultValue, ok := (response.Result).(map[string]interface{}) // Ép kiểu an toàn
		if !ok {
			return nil, fmt.Errorf("kết quả không phải là map: %v", response.Result)
		}
		// Khởi tạo JsonAccountState từ map JSON
		jsonAccountState := &state.JsonAccountState{
			Address:        resultValue["address"].(string),
			Balance:        resultValue["balance"].(string),
			PendingBalance: resultValue["pendingBalance"].(string),
			LastHash:       resultValue["lastHash"].(string),
			DeviceKey:      resultValue["deviceKey"].(string),
			Nonce:          uint64(resultValue["nonce"].(float64)), // Chuyển đổi từ float64 sang uint64
			PublicKeyBls:   resultValue["publicKeyBls"].(string),
			AccountType:    int32(resultValue["accountType"].(float64)), // Chuyển đổi từ float64 sang int32
		}
		accountState := jsonAccountState.ToAccountState()
		return accountState, nil
	}
	return nil, fmt.Errorf("kết quả không hợp lệ: %v", response.Result)
}

func (c *ClientRPC) GetDeviceKey(hash common.Hash) (common.Hash, error) {
	logger.Info("GetDeviceKey")

	request := &JSONRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getDeviceKey",
		Params:  []interface{}{hash.String()}, // Thay đổi thành []interface{}
		Id:      1,
	}

	response := c.SendHTTPRequest(request)
	logger.Info(response.Result)
	if response.Error != nil {
		return common.Hash{}, fmt.Errorf("lỗi từ server: code=%d, message=%s", response.Error.Code, response.Error.Message)
	}
	if response.Result != nil {
		return common.HexToHash(response.Result.(string)), nil
	}
	return common.Hash{}, fmt.Errorf("kết quả không hợp lệ: %v", response.Result)
}

func (c *ClientRPC) SendRawTransaction(input hexutil.Bytes, ethInput hexutil.Bytes, pubKeyBls hexutil.Bytes) JSONRPCResponse {
	logger.Info("SendRawTransaction")

	request := &JSONRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_sendRawTransactionWithDeviceKey",
		Params:  []interface{}{input.String(), ethInput.String(), pubKeyBls.String()}, // Thay đổi thành []interface{}		Id:      1,
	}

	response := c.SendHTTPRequest(request)
	return *response

}

func (c *ClientRPC) SendCallTransaction(input hexutil.Bytes) JSONRPCResponse {
	logger.Info("SendCallTransaction")

	request := &JSONRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_call",
		Params:  []interface{}{input.String()}, // Thay đổi thành []interface{}		Id:      1,
	}

	response := c.SendHTTPRequest(request)
	return *response

}

func (c *ClientRPC) SendEstimateGas(input hexutil.Bytes) JSONRPCResponse {
	logger.Info("SendEstimateGas")
	request := &JSONRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_estimateGas",
		Params:  []interface{}{input.String()}, // Thay đổi thành []interface{}		Id:      1,
	}

	response := c.SendHTTPRequest(request)
	return *response
}

func (c *ClientRPC) BuildTransaction(
	lastHash common.Hash,
	fromAddress common.Address,
	toAddress common.Address,
	pendingUse *big.Int,
	amount *big.Int,
	maxGas uint64,
	maxGasFee uint64,
	maxTimeUse uint64,
	action mt_proto.ACTION,
	data []byte,
	relatedAddress [][]byte,
	lastDeviceKey common.Hash,
	newDeviceKey common.Hash,
	commissionPrivateKey []byte,
	nonce uint64,
) ([]byte, error) {
	transaction := mt_transaction.NewTransaction(
		lastHash,
		fromAddress,
		toAddress,
		pendingUse,
		amount,
		maxGas,
		maxGasFee,
		maxTimeUse,
		data,
		relatedAddress,
		lastDeviceKey,
		newDeviceKey,
		nonce,
		c.ChainId.Uint64(),
	)
	transaction.SetSign(c.KeyPair.PrivateKey())
	if commissionPrivateKey != nil {
		transaction.SetCommissionSign(mt_common.PrivateKeyFromBytes(commissionPrivateKey))
	}
	bTransaction, err := transaction.Marshal()
	return bTransaction, err
}

func (c *ClientRPC) BuildTransactionEmpty(
	fromAddress common.Address,
	toAddress common.Address,
	amount *big.Int,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
) ([]byte, error) {

	lastDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	newDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	as, err := c.GetAccountState(fromAddress, rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))
	if err != nil {
		return nil, fmt.Errorf("lỗi khi get acccount state: %w", err) // Cập nhật thông báo lỗi
	}
	bRelatedAddresses := make([][]byte, 0)

	bTransaction, err := c.BuildTransaction(
		as.LastHash(),
		fromAddress,
		toAddress,
		as.PendingBalance(),
		amount,
		maxGas,
		maxGasPrice,
		maxTimeUse,
		mt_proto.ACTION_EMPTY,
		[]byte{},
		bRelatedAddresses,
		lastDeviceKey,
		newDeviceKey,
		nil,
		as.Nonce(),
	)
	return bTransaction, err
}

func (c *ClientRPC) BuildTransactionDeploy(
	fromAddress common.Address,
	amount *big.Int,
	data []byte,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
) ([]byte, error) {

	lastDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	newDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	as, err := c.GetAccountState(fromAddress, rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))
	if err != nil {
		return nil, fmt.Errorf("lỗi khi get acccount state: %w", err) // Cập nhật thông báo lỗi
	}
	bRelatedAddresses := make([][]byte, 0)
	var bData []byte

	// toAddress := common.BytesToAddress(
	// 	crypto.Keccak256(
	// 		append(
	// 			as.Address().Bytes(),
	// 			byte(as.Nonce())),
	// 	)[12:],
	// )
	toAddress := common.Address{}

	deployData := mt_transaction.NewDeployData(
		data,
		common.HexToAddress("0xda7284fac5e804f8b9d71aa39310f0f86776b51d"),
	)
	bData, err = deployData.Marshal()
	if err != nil {
		panic(err)
	}
	bTransaction, err := c.BuildTransaction(
		as.LastHash(),
		fromAddress,
		toAddress,
		as.PendingBalance(),
		amount,
		maxGas,
		maxGasPrice,
		maxTimeUse,
		mt_proto.ACTION_DEPLOY_SMART_CONTRACT,
		bData,
		bRelatedAddresses,
		lastDeviceKey,
		newDeviceKey,
		nil,
		as.Nonce(),
	)
	return bTransaction, err
}

func (c *ClientRPC) BuildTransactionCallContract(
	fromAddress common.Address,
	toAddress common.Address,
	amount *big.Int,
	data []byte,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
) ([]byte, error) {

	lastDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	newDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	as, err := c.GetAccountState(fromAddress, rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))
	if err != nil {
		return nil, fmt.Errorf("lỗi khi get acccount state: %w", err) // Cập nhật thông báo lỗi
	}
	bRelatedAddresses := make([][]byte, 0)
	var bData []byte

	callData := mt_transaction.NewCallData(data)

	bData, err = callData.Marshal()
	if err != nil {
		panic(err)
	}
	bTransaction, err := c.BuildTransaction(
		as.LastHash(),
		fromAddress,
		toAddress,
		as.PendingBalance(),
		amount,
		maxGas,
		maxGasPrice,
		maxTimeUse,
		mt_proto.ACTION_CALL_SMART_CONTRACT,
		bData,
		bRelatedAddresses,
		lastDeviceKey,
		newDeviceKey,
		nil,
		as.Nonce(),
	)
	return bTransaction, err
}

func (c *ClientRPC) BuildTransactionFromEthTx(
	ethTx *types.Transaction,
) ([]byte, error) {
	logger.Info("BuildTransactionFromEthTx")

	maxGas := uint64(10000000)
	maxGasPrice := uint64(mt_common.MINIMUM_BASE_FEE)
	sg := types.NewCancunSigner(ethTx.ChainId())
	fromAddress, err := sg.Sender(ethTx)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi get fromAddress : %w", err) // Cập nhật thông báo lỗi
	}
	// lastDeviceKey := common.HexToHash(
	// 	"0000000000000000000000000000000000000000000000000000000000000000",
	// )
	// newDeviceKey := common.HexToHash(
	// 	"0000000000000000000000000000000000000000000000000000000000000000",
	// )
	as, err := c.GetAccountState(fromAddress, rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))

	if err != nil {
		return nil, fmt.Errorf("lỗi khi get acccount state: %w", err) // Cập nhật thông báo lỗi
	}

	transactionHash := common.Hash{}
	deviceKey, err := c.GetDeviceKey(as.LastHash())
	if err != nil {
		logger.Info("lỗi khi get deviceKey", err)
	} else {
		transactionHash = deviceKey
	}
	logger.Info(transactionHash)
	logger.Info(hex.EncodeToString(crypto.Keccak256(transactionHash.Bytes())))

	rawNewDeviceKeyBytes := []byte(fmt.Sprintf("%s-%d", hex.EncodeToString(as.LastHash().Bytes()), time.Now().Unix()))

	rawNewDeviceKey := crypto.Keccak256(rawNewDeviceKeyBytes)

	newDeviceKey := crypto.Keccak256Hash(rawNewDeviceKey)

	logger.Info("as.LastHash()", as.LastHash())
	logger.Info("transactionHash", transactionHash)

	bRelatedAddresses := make([][]byte, 0)

	logger.Info("ethTx Nonce", ethTx.Nonce())

	var toAddress common.Address
	var bData []byte
	if len(ethTx.Data()) > 0 && ethTx.To() == nil {
		// toAddress = common.BytesToAddress(
		// 	crypto.Keccak256(
		// 		append(
		// 			as.Address().Bytes(),
		// 			as.LastHash().Bytes()...),
		// 	)[12:],
		// )
		toAddress = common.Address{}

		logger.Info("toAddress deploy: ", toAddress)
		deployData := mt_transaction.NewDeployData(
			ethTx.Data(),
			common.HexToAddress("0xda7284fac5e804f8b9d71aa39310f0f86776b51d"),
		)
		logger.Info(hexutil.Encode(ethTx.Data())) // Chuyển đổi thành chuỗi hex

		bData, err = deployData.Marshal()
		if err != nil {
			return nil, fmt.Errorf("lỗi khi create deployData : %w", err) // Cập nhật thông báo lỗi
		}
	}

	if len(ethTx.Data()) > 0 && ethTx.To() != nil {
		toAddress = common.BytesToAddress(ethTx.To().Bytes())
		logger.Info("toAddress deploy: ", toAddress)
		callData := mt_transaction.NewCallData(ethTx.Data())
		logger.Info(ethTx.Data())

		bData, err = callData.Marshal()
		if err != nil {
			panic(err)
		}
	}

	if len(ethTx.Data()) == 0 && ethTx.To() != nil {
		toAddress = common.BytesToAddress(ethTx.To().Bytes())
		maxGas = uint64(100000)

	}

	transaction := mt_transaction.NewTransaction(
		as.LastHash(),
		fromAddress,
		toAddress,
		as.PendingBalance(),
		ethTx.Value(),
		maxGas,
		maxGasPrice,
		// ethTx.Gas(),
		// ethTx.GasPrice().Uint64(),
		600,
		bData,
		bRelatedAddresses,
		as.LastHash(),
		newDeviceKey,
		ethTx.Nonce(),
		c.ChainId.Uint64(),
	)
	transaction.SetSign(c.KeyPair.PrivateKey())
	bTransaction, err := transaction.Marshal()
	return bTransaction, err
}

func (c *ClientRPC) BuildTransactionObjectFromEthTx(
	ethTx *types.Transaction,
) (mt_types.Transaction, error) {
	sg := types.NewCancunSigner(ethTx.ChainId())
	fromAddress, err := sg.Sender(ethTx)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi get fromAddress : %w", err) // Cập nhật thông báo lỗi
	}
	lastDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	newDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	as, err := c.GetAccountState(fromAddress, rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))
	if err != nil {
		return nil, fmt.Errorf("lỗi khi get acccount state: %w", err) // Cập nhật thông báo lỗi
	}
	bRelatedAddresses := make([][]byte, 0)

	var toAddress common.Address
	var bData []byte
	if len(ethTx.Data()) > 0 && ethTx.To() == nil {
		// toAddress = common.BytesToAddress(
		// 	crypto.Keccak256(
		// 		append(
		// 			as.Address().Bytes(),
		// 			as.LastHash().Bytes()...),
		// 	)[12:],
		// )
		toAddress = common.Address{}

		deployData := mt_transaction.NewDeployData(
			ethTx.Data(),
			common.HexToAddress("0xda7284fac5e804f8b9d71aa39310f0f86776b51d"),
		)
		bData, err = deployData.Marshal()
		if err != nil {
			return nil, fmt.Errorf("lỗi khi create deployData : %w", err) // Cập nhật thông báo lỗi
		}
	}

	if len(ethTx.Data()) > 0 && ethTx.To() != nil {
		toAddress = common.BytesToAddress(ethTx.To().Bytes())
		callData := mt_transaction.NewCallData(ethTx.Data())

		bData, err = callData.Marshal()
		if err != nil {
			panic(err)
		}
	}

	if len(ethTx.Data()) == 0 && ethTx.To() != nil {
		toAddress = common.BytesToAddress(ethTx.To().Bytes())
	}

	maxGas := uint64(10000000)
	maxGasPrice := uint64(mt_common.MINIMUM_BASE_FEE)

	transaction := mt_transaction.NewTransaction(
		as.LastHash(),
		fromAddress,
		toAddress,
		as.PendingBalance(),
		ethTx.Value(),
		maxGas,
		maxGasPrice,
		// ethTx.Gas(),
		// ethTx.GasPrice().Uint64(),
		600,
		bData,
		bRelatedAddresses,
		lastDeviceKey,
		newDeviceKey,
		ethTx.Nonce(),
		c.ChainId.Uint64(),
	)
	return transaction, err
}

func (c *ClientRPC) BuildCallTransaction(callDataT []byte, toAddress common.Address, fromAddress common.Address) ([]byte, error) {
	maxGas := uint64(10000000)
	maxGasPrice := uint64(mt_common.MINIMUM_BASE_FEE)
	lastDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	newDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)

	as, err := c.GetAccountState(fromAddress, rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))

	if err != nil {
		return nil, fmt.Errorf("lỗi khi get acccount state: %w", err) // Cập nhật thông báo lỗi
	}
	bRelatedAddresses := make([][]byte, 0)

	var bData []byte

	callData := mt_transaction.NewCallData(callDataT)

	bData, err = callData.Marshal()
	if err != nil {
		return nil, fmt.Errorf("lỗi convert callData: %w", err) // Cập nhật thông báo lỗi
	}

	txx := mt_transaction.NewTransaction(
		as.LastHash(),
		fromAddress,
		toAddress,
		as.PendingBalance(),
		big.NewInt(0),
		maxGas,
		maxGasPrice,
		600,
		bData,
		bRelatedAddresses,
		lastDeviceKey,
		newDeviceKey,
		as.Nonce(),
		c.ChainId.Uint64(),
	)
	txx.SetSign(c.KeyPair.PrivateKey())
	bTransaction, err := txx.Marshal()
	return bTransaction, err
}

func (c *ClientRPC) BuildDeployTransaction(callDataT []byte, from common.Address) ([]byte, error) {
	fromAddress := from
	maxGas := uint64(10000000)
	maxGasPrice := uint64(mt_common.MINIMUM_BASE_FEE)
	lastDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)
	newDeviceKey := common.HexToHash(
		"0000000000000000000000000000000000000000000000000000000000000000",
	)

	as, err := c.GetAccountState(fromAddress, rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))

	if err != nil {
		return nil, fmt.Errorf("lỗi khi get acccount state: %w", err) // Cập nhật thông báo lỗi
	}
	bRelatedAddresses := make([][]byte, 0)

	var bData []byte

	callData := mt_transaction.NewCallData(callDataT)

	bData, err = callData.Marshal()
	if err != nil {
		return nil, fmt.Errorf("lỗi convert callData: %w", err) // Cập nhật thông báo lỗi
	}
	toAddress := common.Address{}

	txx := mt_transaction.NewTransaction(
		as.LastHash(),
		fromAddress,
		toAddress,
		as.PendingBalance(),
		big.NewInt(0),
		maxGas,
		maxGasPrice,
		6000000,
		bData,
		bRelatedAddresses,
		lastDeviceKey,
		newDeviceKey,
		as.Nonce(),
		c.ChainId.Uint64(),
	)
	txx.SetSign(c.KeyPair.PrivateKey())
	bTransaction, err := txx.Marshal()
	return bTransaction, err
}

func (c *ClientRPC) BuildTransactionWithDeviceKeyFromEthTx(
	ethTx *types.Transaction,
) ([]byte, error) {

	sg := types.NewCancunSigner(ethTx.ChainId())
	fromAddress, err := sg.Sender(ethTx)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi get fromAddress : %w", err) // Cập nhật thông báo lỗi
	}
	as, err := c.GetAccountState(fromAddress, rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))

	if err != nil {
		return nil, fmt.Errorf("lỗi khi get acccount state: %w", err) // Cập nhật thông báo lỗi
	}

	deviceKey, err := c.GetDeviceKey(as.LastHash())
	if err != nil {
		logger.Info("lỗi khi get deviceKey", err)
	}

	rawNewDeviceKeyBytes := []byte(fmt.Sprintf("%s-%d", hex.EncodeToString(as.LastHash().Bytes()), time.Now().Unix()))

	rawNewDeviceKey := crypto.Keccak256(rawNewDeviceKeyBytes)

	newDeviceKey := crypto.Keccak256Hash(rawNewDeviceKey)

	bRelatedAddresses := make([][]byte, 0)

	logger.Info("ethTx Nonce", ethTx.Nonce())

	var toAddress common.Address
	var bData []byte
	if len(ethTx.Data()) > 0 && ethTx.To() == nil {
		toAddress = common.Address{}

		logger.Info("toAddress deploy: ", toAddress)
		deployData := mt_transaction.NewDeployData(
			ethTx.Data(),
			common.HexToAddress("0xda7284fac5e804f8b9d71aa39310f0f86776b51d"),
		)
		logger.Info(hexutil.Encode(ethTx.Data())) // Chuyển đổi thành chuỗi hex

		bData, err = deployData.Marshal()
		if err != nil {
			return nil, fmt.Errorf("lỗi khi create deployData : %w", err) // Cập nhật thông báo lỗi
		}
	}

	if len(ethTx.Data()) > 0 && ethTx.To() != nil {
		toAddress = common.BytesToAddress(ethTx.To().Bytes())
		logger.Info("toAddress deploy: ", toAddress)
		callData := mt_transaction.NewCallData(ethTx.Data())
		logger.Info(ethTx.Data())

		bData, err = callData.Marshal()
		if err != nil {
			panic(err)
		}
	}

	if len(ethTx.Data()) == 0 && ethTx.To() != nil {
		toAddress = common.BytesToAddress(ethTx.To().Bytes())

	}

	transaction := mt_transaction.NewTransaction(
		as.LastHash(),
		fromAddress,
		toAddress,
		as.PendingBalance(),
		ethTx.Value(),
		ethTx.Gas(),
		ethTx.GasPrice().Uint64(),
		0,
		bData,
		bRelatedAddresses,
		deviceKey,
		newDeviceKey,
		ethTx.Nonce(),
		c.ChainId.Uint64(),
	)
	v, r, s := ethTx.RawSignatureValues()
	transaction.SetSignatureValues(c.ChainId, v, r, s)
	transaction.SetSign(c.KeyPair.PrivateKey())

	// Create TransactionWithDeviceKey
	transactionWithDeviceKey := &mt_proto.TransactionWithDeviceKey{
		Transaction: transaction.Proto().(*mt_proto.Transaction),
		DeviceKey:   rawNewDeviceKey,
	}

	// Serialize to bytes
	bTransactionWithDeviceKey, err := proto.Marshal(transactionWithDeviceKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TransactionWithDeviceKey: %w", err)
	}
	return bTransactionWithDeviceKey, err
}

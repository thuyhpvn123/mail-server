package client

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/stretchr/testify/assert"

	c_config "gomail/cmd/client/pkg/config"
	"gomail/pkg/logger"
	pb "gomail/pkg/proto"
	"gomail/pkg/transaction"
	"gomail/types"
)

// test account
// Private key: 1c923d7764cb712f2f007a53f5c14f898fc3fcc3f00c609c55ec5cb4d7443211
// Public key: 897b1946f73218fb5cc1df08a7debdcdbfe19a404f923cfcbe0b9f64a9eec6b5b6967040a77a1d48cc4af1fac59960ae
// Address: ae357fc27436ed8aeeb2df11cbd745a12b2e2093
var client *Client

func init() {
	clientConfigJson := `
	{
		"version": "0.0.1.0",
		"private_key": "2b3aa0f620d2d73c046cd93eb64f2eb687a95b22e278500aa251c8c9dda1203b",
		"parent_address": "9065e2caaa4bf6533ded0b2a763e4ee81cd31bd2",
		"parent_connection_address": "0.0.0.0:4200",
		"parent_connection_type": "node",
		"dns_link": "http://127.0.0.1:7080/api/dns/connection-address/"
	}
	`
	config := &c_config.ClientConfig{}
	json.Unmarshal([]byte(clientConfigJson), config)
	client, _ = NewClient(
		config,
	)
}

func TestGetAccountState(t *testing.T) {
	as, err := client.AccountState(
		common.HexToAddress("97126B71376F7e55fBA904FdaA9dF0dBd396612f"),
	)
	assert.Nil(t, err)
	logger.Info(as)
}

func TestGetAccountState2(t *testing.T) {
	as, err := client.AccountState(
		common.HexToAddress("0x0000c3f90FE4788Bd0a4B31C5bc47810C044198F"),
	)
	assert.Nil(t, err)
	logger.Info(as)
	as2, err2 := client.AccountState(
		common.HexToAddress("97126B71376F7e55fBA904FdaA9dF0dBd396612f"),
	)
	assert.Nil(t, err2)
	logger.Info(as2)
}

func generateRandomHexAddress() string {
	// Generate a random 32-byte (256-bit) value.
	randomBytes := make([]byte, 20)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err) // Handle the error appropriately in a real application
	}

	// Convert the bytes to a hexadecimal string, prepending "0x".
	hexString := common.Bytes2Hex(randomBytes)
	return "0x" + hexString
}

func TestSendTransaction(t *testing.T) {
	bigI := big.NewInt(5000000000)
	bigI.Mul(bigI, big.NewInt(100000000))
	receipt, err := client.SendTransaction(
		common.HexToAddress("0x97126B71376F7e55fBA904FdaA9dF0dBd396612f"),
		common.HexToAddress("0x7ad5e1388e2fc94Fdb248914af7d300B0D794194"),
		bigI,
		pb.ACTION_EMPTY,
		[]byte{},
		nil,
		100000,
		1000000000,
		60,
	)
	assert.Nil(t, err)
	logger.Info(receipt)
}

func TestUpdateTypeAccont(t *testing.T) {
	callData := transaction.NewCallData(common.FromHex("0x1"))
	bData, _ := callData.Marshal()
	bigI := big.NewInt(1)
	bigI.Mul(bigI, big.NewInt(1))
	receipt, err := client.SendTransactionWithDeviceKey(
		common.HexToAddress("0x97126B71376F7e55fBA904FdaA9dF0dBd396612f"),
		common.HexToAddress("0x0000000000000000000000000000000000000000"),
		bigI,
		pb.ACTION_EMPTY,
		bData,
		nil,
		100000,
		1000000000,
		6000,
	)
	assert.Nil(t, err)
	logger.Info(receipt)
}

// Tạo liên kết đại chỉ ví với key của bls client
func TestUpdateAddBLSPublicKey(t *testing.T) {
	// Đây là private key của đia chị ví metamask
	prk, _ := hex.DecodeString("0cf85c6fdccc097927733b722b6ca4f158a846cc43d5d41642feb878b5e23239")
	message := append(client.clientContext.KeyPair.PublicKey().Bytes(), []byte("120")...)
	hash := crypto.Keccak256(message)

	sig, _ := secp256k1.Sign(hash, prk)
	logger.Info(len(client.clientContext.KeyPair.PublicKey().Bytes()))
	pbk, err := secp256k1.RecoverPubkey(hash, sig)
	// combined := append(client.clientContext.KeyPair.PublicKey().Bytes(), sig...)

	if err != nil {
		logger.Error("Error ValidSign", err)
	}
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pbk[1:])[12:])

	logger.Info(addr)
	// // 0x7ad5e1388e2fc94Fdb248914af7d300B0D794194
	// // address := common.BytesToAddress(pbk)
	// bigI := big.NewInt(1)
	// receipt, err := client.SendTransactionWithDeviceKey(
	// 	// Đây nhập đia chỉ ví tương ứng với private key phía trên
	// 	common.HexToAddress("0x0D465eD9b9b8dD9aca00692fC6eBf4c36c83e68A"),
	// 	common.HexToAddress("0x0D465eD9b9b8dD9aca00692fC6eBf4c36c83e68A"),
	// 	bigI,
	// 	pb.ACTION_EMPTY,
	// 	combined,
	// 	nil,
	// 	100000,
	// 	1000000000,
	// 	6000,
	// )
	// assert.Nil(t, err)
}
func TestAlwaysTrue(t *testing.T) {
	toAddress := crypto.CreateAddress(common.HexToAddress("0x92CC5d5127cd12715F18b531Eb73805e0678b9eC"), 12)
	logger.Info(toAddress)

}
func TestSendDeployTxs(t *testing.T) {
	clientConfigJson := `
	{
		"version": "0.0.1.0",
		"private_key": "2b3aa0f620d2d73c046cd93eb64f2eb687a95b22e278500aa251c8c9dda1203b",
		"parent_address": "9065e2caaa4bf6533ded0b2a763e4ee81cd31bd2",
		"parent_connection_address": "0.0.0.0:4200",
		"parent_connection_type": "node",
		"dns_link": "http://127.0.0.1:7080/api/dns/connection-address/"
	}
	`
	config := &c_config.ClientConfig{}
	json.Unmarshal([]byte(clientConfigJson), config)
	client, _ := NewClient(
		config,
	)
	txBytesList, err := LoadTransactionsFromFile("txs")
	if err != nil {
		logger.Error("error loading transactions from file: %w", err)
	}
	txBytesList1, err := LoadTransactionsFromFile("txs1")

	if err != nil {
		logger.Error("error loading transactions from file: %w", err)
	}

	txs2 := make([]types.Transaction, 0, len(txBytesList))
	txs := make([]types.Transaction, 0, len(txBytesList1))
	for _, txBytes := range txBytesList {
		tx := &transaction.Transaction{}
		err := tx.Unmarshal(txBytes)
		if err != nil {
			logger.Error("error unmarshaling transaction: %w", err)
		}
		txs2 = append(txs2, tx)

	}

	for _, txBytes1s := range txBytesList1 {
		tx := &transaction.Transaction{}
		err := tx.Unmarshal(txBytes1s)
		if err != nil {
			logger.Error("error unmarshaling transaction: %w", err)
		}
		txs = append(txs, tx)

	}
	allTransactions := append(txs, txs2...)
	logger.Info(len(allTransactions))
	client.transactionController.SendTransactions(allTransactions[:1000])

}

func TestReadAndSendTransactionsOption(t *testing.T) {

	txBytesList, err := LoadTransactionsFromFile("txs")
	if err != nil {
		logger.Error("error loading transactions from file: %w", err)
	}
	txs2 := make([]types.Transaction, 0, len(txBytesList))
	var txE types.Transaction
	for _, txBytes := range txBytesList {
		tx := &transaction.Transaction{}
		err := tx.Unmarshal(txBytes)
		if err != nil {
			logger.Error("error unmarshaling transaction: %w", err)
		}
		txs2 = append(txs2, tx)
		if tx.ToAddress() == common.HexToAddress("0xda73035a3Db008764DC5C7734045aD403E2CC9be") {
			logger.Info(tx)
			txE = tx
		}

	}
	txBytesList1, err := LoadTransactionsFromFile("txsdll")

	if err != nil {
		logger.Error("error loading transactions from file: %w", err)
	}

	txs := make([]types.Transaction, 0, len(txBytesList1))

	for _, txBytes1s := range txBytesList1 {
		tx := &transaction.Transaction{}
		err := tx.Unmarshal(txBytes1s)
		if err != nil {
			logger.Error("error unmarshaling transaction: %w", err)
		}
		txs = append(txs, tx)

	}
	allTransactions := append(txs, txs2...)
	logger.Info(txs[0].GetNonce())
	logger.Info(len(allTransactions))
	// client.transactionController.SendTransactions(txs2[:1000])
	logger.Info("ToAddress: ", txE.FromAddress())

	logger.Info("ToAddress: ", txE.ToAddress())
	logger.Info("Nonce: ", txE.GetNonce())
	client.transactionController.SendTransactions([]types.Transaction{txE, txs2[1]})

}

func TestCheckDeploy(t *testing.T) {
	// Đường dẫn đến file log
	filePath := "txs.og"

	// Đọc nội dung file
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Lỗi khi đọc file: %v", err)
	}

	// Tách các dòng trong file
	lines := bytes.Split(data, []byte("\n"))

	// Mảng để lưu trữ các địa chỉ
	addresses := []string{}

	// Duyệt qua từng dòng
	for _, line := range lines {
		// Kiểm tra xem dòng có chứa dấu hai chấm không
		if bytes.Contains(line, []byte(":")) {
			// Tìm vị trí của dấu hai chấm cuối cùng
			lastColonIndex := bytes.LastIndex(line, []byte(":"))

			// Nếu tìm thấy dấu hai chấm, trích xuất địa chỉ
			if lastColonIndex != -1 {
				address := string(bytes.TrimSpace(line[lastColonIndex+1:]))
				addresses = append(addresses, address)
			}
		}
	}

	// Kiểm tra kết quả
	assert.Greater(t, len(addresses), 0, "Không tìm thấy địa chỉ nào trong file")
	fmt.Println("Số địa chỉ được tìm thấy:", len(addresses))
	count := 0
	countCheck := 0
	listErrorAddress := []common.Address{} // Khởi tạo listErrorAddress
	for _, address := range addresses {
		as, err := client.AccountState(
			common.HexToAddress(address),
		)
		if err != nil {
			logger.Info(err)
		} else {

			response, err := makeRequest(as.Address().String())
			if err != nil {
				logger.Error("Get erro: ", as.Address().String())
			} else {
				result := strings.TrimSpace(response.Result)                    // Loại bỏ khoảng trắng thừa
				result = regexp.MustCompile(`\s+`).ReplaceAllString(result, "") // Loại bỏ khoảng trắng thừa
				result = strings.ReplaceAll(result, "\n", "")
				// Loại bỏ ký tự xuống dòng
				//result = sanitizeString(result) // Hàm để làm sạch chuỗi (nếu cần)
				expected := "0x606060405236156100ab576000357c01000000000000000000000000000000000000000000000000000000009004806306fdde03146100bd578063095ea7b31461013d57806318160ddd1461017957806323b872dd146101a1578063313ce567146101e657806354fd4d501461021157806370a082311461029157806395d89b41146102c2578063a9059cbb14610342578063cae9ca511461037e578063dd62ed3e14610401576100ab565b34610002576100bb5b610002565b565b005b34610002576100cf600480505061043b565b60405180806020018281038252838181518152602001915080519060200190808383829060006004602084601f0104600302600f01f150905090810190601f16801561012f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b346100025761016160048080359060200190919080359060200190919050506104dc565b60405180821515815260200191505060405180910390f35b346100025761018b60048050506105b0565b6040518082815260200191505060405180910390f35b34610002576101ce60048080359060200190919080359060200190919080359060200190919050506105b9565b60405180821515815260200191505060405180910390f35b34610002576101f860048050506107c5565b604051808260ff16815260200191505060405180910390f35b346100025761022360048050506107d8565b60405180806020018281038252838181518152602001915080519060200190808383829060006004602084601f0104600302600f01f150905090810190601f1680156102835780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34610002576102ac6004808035906020019091905050610879565b6040518082815260200191505060405180910390f35b34610002576102d460048050506108b7565b60405180806020018281038252838181518152602001915080519060200190808383829060006004602084601f0104600302600f01f150905090810190601f1680156103345780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34610002576103666004808035906020019091908035906020019091905050610958565b60405180821515815260200191505060405180910390f35b34610002576103e96004808035906020019091908035906020019091908035906020019082018035906020019191908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050909091905050610a98565b60405180821515815260200191505060405180910390f35b34610002576104256004808035906020019091908035906020019091905050610cdf565b6040518082815260200191505060405180910390f35b60036000508054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156104d45780601f106104a9576101008083540402835291602001916104d4565b820191906000526020600020905b8154815290600101906020018083116104b757829003601f168201915b505050505081565b600081600160005060003373ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005060008573ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040518082815260200191505060405180910390a3600190506105aa565b92915050565b60026000505481565b600081600060005060008673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000505410158015610653575081600160005060008673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005060003373ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000505410155b801561065f5750600082115b156107b45781600060005060008573ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282825054019250508190555081600060005060008673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282825054039250508190555081600160005060008673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005060003373ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828282505403925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a3600190506107be566107bd565b600090506107be565b5b9392505050565b600460009054906101000a900460ff1681565b60066000508054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108715780601f1061084657610100808354040283529160200191610871565b820191906000526020600020905b81548152906001019060200180831161085457829003601f168201915b505050505081565b6000600060005060008373ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000505490506108b2565b919050565b60056000508054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156109505780601f1061092557610100808354040283529160200191610950565b820191906000526020600020905b81548152906001019060200180831161093357829003601f168201915b505050505081565b600081600060005060003373ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005054101580156109995750600082115b15610a885781600060005060003373ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282825054039250508190555081600060005060008573ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828282505401925050819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a360019050610a9256610a91565b60009050610a92565b5b92915050565b600082600160005060003373ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005060008673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050819055508373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925856040518082815260200191505060405180910390a38373ffffffffffffffffffffffffffffffffffffffff1660405180807f72656365697665417070726f76616c28616464726573732c75696e743235362c81526020017f616464726573732c627974657329000000000000000000000000000000000000815260200150602e01905060405180910390207c0100000000000000000000000000000000000000000000000000000000900433853086604051857c0100000000000000000000000000000000000000000000000000000000028152600401808573ffffffffffffffffffffffffffffffffffffffff1681526020018481526020018373ffffffffffffffffffffffffffffffffffffffff1681526020018280519060200190808383829060006004602084601f0104600302600f01f150905090810190601f168015610ca75780820380516001836020036101000a031916815260200191505b509450505050506000604051808303816000876161da5a03f1925050501515610ccf57610002565b60019050610cd8565b9392505050565b6000600160005060008473ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005060008373ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050549050610d42565b9291505056"

				if !compareResult(result, expected) {
					logger.Info(as)

					fmt.Printf("Result does not match expected value. Result: %s, Expected: %s\n", result, expected)
					// Fail the test here
					listErrorAddress = append(listErrorAddress, as.Address()) // Thêm địa chỉ vào listErrorAddress nếu có lỗi

					count = count + 1
					fmt.Printf("Resul match expected value. Result: %s, Expected: %s\n", result, expected)

				} else {
					assert.Equal(t, result, expected)
					countCheck = countCheck + 1

				}
			}
		}
	}
	///[0xebE01CdEd59b64eB34584E2E14F5510D23Bf6De2 0xF5ADB85f874B224D7a80aB09a257230b804f8479 0x5eFf52A7dbd70aF4901Da5966123e699AFe4B8Bc 0xda73035a3Db008764DC5C7734045aD403E2CC9be]
	logger.Info("listErrorAddress: ", listErrorAddress)
	logger.Info("count: ", count)
	logger.Info("countCheck done: ", countCheck)

}

func TestCallTxs(t *testing.T) {

	clientConfigJson := `
	{
		"version": "0.0.1.0",
		"private_key": "2b3aa0f620d2d73c046cd93eb64f2eb687a95b22e278500aa251c8c9dda1203b",
		"parent_address": "9065e2caaa4bf6533ded0b2a763e4ee81cd31bd2",
		"parent_connection_address": "0.0.0.0:4200",
		"parent_connection_type": "node",
		"dns_link": "http://127.0.0.1:7080/api/dns/connection-address/"
	}
	`
	config := &c_config.ClientConfig{}
	json.Unmarshal([]byte(clientConfigJson), config)
	client, _ := NewClient(
		config,
	)
	txBytesList, err := LoadTransactionsFromFile("txs")
	if err != nil {
		logger.Error("error loading transactions from file: %w", err)
	}
	txBytesList1, err := LoadTransactionsFromFile("txs1")

	if err != nil {
		logger.Error("error loading transactions from file: %w", err)
	}

	txs2 := make([]types.Transaction, 0, len(txBytesList))
	txs := make([]types.Transaction, 0, len(txBytesList1))
	for _, txBytes := range txBytesList {
		tx := &transaction.Transaction{}
		err := tx.Unmarshal(txBytes)
		if err != nil {
			logger.Error("error unmarshaling transaction: %w", err)
		}
		txs2 = append(txs2, tx)

	}

	for _, txBytes1s := range txBytesList1 {
		tx := &transaction.Transaction{}
		err := tx.Unmarshal(txBytes1s)
		if err != nil {
			logger.Error("error unmarshaling transaction: %w", err)
		}
		txs = append(txs, tx)

	}

	allTransactions := append(txs, txs2...)

	txBytesExample, err := LoadTransactionsFromFile("send_tx_example")
	if err != nil {
		logger.Error("error loading transactions from file: %w", err)
	}
	txE := &transaction.Transaction{}
	err = txE.Unmarshal(txBytesExample[0])
	if err != nil {
		logger.Error("error unmarshaling transaction: %w", err)
	}
	txsCall := make([]types.Transaction, 0, len(allTransactions))
	for _, tx := range allTransactions {
		txNew := txE.CopyTransaction()
		txNew.SetFromAddress(tx.FromAddress())
		txNew.SetToAddress(tx.ToAddress())
		listRelatedAddress := []common.Address{tx.FromAddress(), tx.ToAddress()}
		bRelatedAddresses := make([][]byte, len(listRelatedAddress))
		for i, v := range listRelatedAddress {
			bRelatedAddresses[i] = v.Bytes()
		}
		txNew.UpdateRelatedAddresses(bRelatedAddresses)
		txNew.SetNonce(2)
		if err != nil {
			logger.Error("error unmarshaling transaction: %w", err)
		}
		txsCall = append(txsCall, txNew)

	}
	logger.Info(len(txsCall))

	client.transactionController.SendTransactions(txsCall[:1000])

}

func TestCallTxsOne(t *testing.T) {

	clientConfigJson := `
	{
		"version": "0.0.1.0",
		"private_key": "2b3aa0f620d2d73c046cd93eb64f2eb687a95b22e278500aa251c8c9dda1203b",
		"parent_address": "9065e2caaa4bf6533ded0b2a763e4ee81cd31bd2",
		"parent_connection_address": "0.0.0.0:4200",
		"parent_connection_type": "node",
		"dns_link": "http://127.0.0.1:7080/api/dns/connection-address/"
	}
	`
	config := &c_config.ClientConfig{}
	json.Unmarshal([]byte(clientConfigJson), config)
	client, _ := NewClient(
		config,
	)
	txBytesList, err := LoadTransactionsFromFile("txs")
	if err != nil {
		logger.Error("error loading transactions from file: %w", err)
	}
	txBytesList1, err := LoadTransactionsFromFile("txs1")

	if err != nil {
		logger.Error("error loading transactions from file: %w", err)
	}

	txs2 := make([]types.Transaction, 0, len(txBytesList))
	txs := make([]types.Transaction, 0, len(txBytesList1))
	for _, txBytes := range txBytesList {
		tx := &transaction.Transaction{}
		err := tx.Unmarshal(txBytes)
		if err != nil {
			logger.Error("error unmarshaling transaction: %w", err)
		}
		txs2 = append(txs2, tx)

	}

	for _, txBytes1s := range txBytesList1 {
		tx := &transaction.Transaction{}
		err := tx.Unmarshal(txBytes1s)
		if err != nil {
			logger.Error("error unmarshaling transaction: %w", err)
		}
		txs = append(txs, tx)

	}

	allTransactions := append(txs, txs2...)

	txBytesExample, err := LoadTransactionsFromFile("send_tx_example")
	if err != nil {
		logger.Error("error loading transactions from file: %w", err)
	}
	txE := &transaction.Transaction{}
	err = txE.Unmarshal(txBytesExample[0])
	if err != nil {
		logger.Error("error unmarshaling transaction: %w", err)
	}
	txsCall := make([]types.Transaction, 0, len(allTransactions))
	for _, tx := range allTransactions {
		txNew := txE.CopyTransaction()
		txNew.SetFromAddress(tx.FromAddress())
		txNew.SetToAddress(tx.ToAddress())
		listRelatedAddress := []common.Address{tx.FromAddress(), tx.ToAddress()}
		bRelatedAddresses := make([][]byte, len(listRelatedAddress))
		for i, v := range listRelatedAddress {
			bRelatedAddresses[i] = v.Bytes()
		}
		txNew.UpdateRelatedAddresses(bRelatedAddresses)
		txNew.SetNonce(1)
		if err != nil {
			logger.Error("error unmarshaling transaction: %w", err)
		}
		txsCall = append(txsCall, txNew)

	}
	logger.Info(len(txsCall))
	logger.Info(txsCall[0].CallData())
	txOne := txsCall[0].CopyTransaction()

	client.transactionController.SendTransactions([]types.Transaction{txOne})

}

func TestCheckCallTxsOne(t *testing.T) {
	// Đường dẫn đến file log
	filePath := "txs.og"

	// Đọc nội dung file
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Lỗi khi đọc file: %v", err)
	}

	// Tách các dòng trong file
	lines := bytes.Split(data, []byte("\n"))

	// Mảng để lưu trữ các địa chỉ
	addresses := []string{}

	// Duyệt qua từng dòng
	for _, line := range lines {
		// Kiểm tra xem dòng có chứa dấu hai chấm không
		if bytes.Contains(line, []byte(":")) {
			// Tìm vị trí của dấu hai chấm cuối cùng
			lastColonIndex := bytes.LastIndex(line, []byte(":"))

			// Nếu tìm thấy dấu hai chấm, trích xuất địa chỉ
			if lastColonIndex != -1 {
				address := string(bytes.TrimSpace(line[lastColonIndex+1:]))
				addresses = append(addresses, address)
			}
		}
	}

	// Kiểm tra kết quả
	assert.Greater(t, len(addresses), 0, "Không tìm thấy địa chỉ nào trong file")
	fmt.Println("Số địa chỉ được tìm thấy:", len(addresses))
	count := 0
	countCheckDone := 0

	listErrorAddress := []common.Address{} // Khởi tạo listErrorAddress
	for _, address := range addresses {
		as, err := client.AccountState(
			common.HexToAddress(address),
		)
		if err != nil {
			logger.Info(err)
		} else {

			response, err := makeRequestCheckErc20(as.Address().String())
			if err != nil {
				logger.Error("Get error: ", as.Address().String(), err)
			} else {
				result := strings.TrimSpace(response.Result)                    // Loại bỏ khoảng trắng thừa
				result = regexp.MustCompile(`\s+`).ReplaceAllString(result, "") // Loại bỏ khoảng trắng thừa
				result = strings.ReplaceAll(result, "\n", "")                   // Loại bỏ ký tự xuống dòng
				//result = sanitizeString(result) // Hàm để làm sạch chuỗi (nếu cần)
				expected := "0x00000000000000000000000000000000000000000000000000000000000003e8"

				if !compareResult(result, expected) {
					logger.Info(as)

					fmt.Printf("Result does not match expected value. Result: %s, Expected: %s\n", result, expected)
					// Fail the test here
					// assert.Fail(t, "Result does not match expected value")
					listErrorAddress = append(listErrorAddress, as.Address()) // Thêm địa chỉ vào listErrorAddress nếu có lỗi

					count = count + 1
					// fmt.Printf("Resul match expected value. Result: %s, Expected: %s\n", result, expected)

				} else {
					// assert.Equal(t, result, expected)
					countCheckDone = countCheckDone + 1
				}
			}
		}
	}

	///[0xebE01CdEd59b64eB34584E2E14F5510D23Bf6De2 0xF5ADB85f874B224D7a80aB09a257230b804f8479 0x5eFf52A7dbd70aF4901Da5966123e699AFe4B8Bc 0xda73035a3Db008764DC5C7734045aD403E2CC9be]
	logger.Info("listErrorAddress: ", listErrorAddress)
	logger.Info("count: ", count)
	logger.Info("countCheckDone: ", countCheckDone)
}

func TestCheckCallTxsTow(t *testing.T) {
	// Đường dẫn đến file log
	filePath := "txs.og"

	// Đọc nội dung file
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Lỗi khi đọc file: %v", err)
	}

	// Tách các dòng trong file
	lines := bytes.Split(data, []byte("\n"))

	// Mảng để lưu trữ các địa chỉ
	addresses := []string{}

	// Duyệt qua từng dòng
	for _, line := range lines {
		// Kiểm tra xem dòng có chứa dấu hai chấm không
		if bytes.Contains(line, []byte(":")) {
			// Tìm vị trí của dấu hai chấm cuối cùng
			lastColonIndex := bytes.LastIndex(line, []byte(":"))

			// Nếu tìm thấy dấu hai chấm, trích xuất địa chỉ
			if lastColonIndex != -1 {
				address := string(bytes.TrimSpace(line[lastColonIndex+1:]))
				addresses = append(addresses, address)
			}
		}
	}

	// Kiểm tra kết quả
	assert.Greater(t, len(addresses), 0, "Không tìm thấy địa chỉ nào trong file")
	fmt.Println("Số địa chỉ được tìm thấy:", len(addresses))
	count := 0
	countCheckDone := 0

	listErrorAddress := []common.Address{} // Khởi tạo listErrorAddress
	for _, address := range addresses {

		as, err := client.AccountState(
			common.HexToAddress(address),
		)
		if err != nil {
			logger.Info(err)
		} else {

			response, err := makeRequestCheckErc20(as.Address().String())
			if err != nil {
				logger.Error("Get error: ", as.Address().String(), err)
			} else {
				result := strings.TrimSpace(response.Result)                    // Loại bỏ khoảng trắng thừa
				result = regexp.MustCompile(`\s+`).ReplaceAllString(result, "") // Loại bỏ khoảng trắng thừa
				result = strings.ReplaceAll(result, "\n", "")                   // Loại bỏ ký tự xuống dòng
				//result = sanitizeString(result) // Hàm để làm sạch chuỗi (nếu cần)
				expected := "0x00000000000000000000000000000000000000000000000000000000000007d0"

				if !compareResult(result, expected) {
					logger.Info(as)

					fmt.Printf("Result does not match expected value. Result: %s, Expected: %s\n", result, expected)
					// Fail the test here
					// assert.Fail(t, "Result does not match expected value")
					listErrorAddress = append(listErrorAddress, as.Address()) // Thêm địa chỉ vào listErrorAddress nếu có lỗi

					count = count + 1
					// fmt.Printf("Resul match expected value. Result: %s, Expected: %s\n", result, expected)

				} else {
					// assert.Equal(t, result, expected)
					countCheckDone = countCheckDone + 1
				}
			}
		}
	}

	///[0xebE01CdEd59b64eB34584E2E14F5510D23Bf6De2 0xF5ADB85f874B224D7a80aB09a257230b804f8479 0x5eFf52A7dbd70aF4901Da5966123e699AFe4B8Bc 0xda73035a3Db008764DC5C7734045aD403E2CC9be]
	logger.Info("listErrorAddress: ", listErrorAddress)
	logger.Info("count: ", count)
	logger.Info("countCheckDone: ", countCheckDone)
}

func compareResult(result, expected string) bool {
	return result == expected
}

type Request struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type Response struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  string           `json:"result"`
	Error   *json.RawMessage `json:"error"`
	ID      int              `json:"id"`
}

func makeRequest(address string) (*Response, error) {
	req := Request{
		JSONRPC: "2.0",
		Method:  "eth_getCode",
		Params:  []interface{}{address, "latest"},
		ID:      1,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post("http://localhost:8646", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if response.Error != nil {
		return &response, fmt.Errorf("error from server: %s", *response.Error)
	}

	return &response, nil
}

func makeRequestCheckErc20(address string) (*Response, error) {

	params := []interface{}{
		map[string]interface{}{
			"to":   address,
			"data": "0x70a08231000000000000000000000000f41e7ec0cbf22f8224328656a990a8548b70976a",
		},
		"latest",
	}
	req := Request{
		JSONRPC: "2.0",
		Method:  "eth_call",
		Params:  params,
		ID:      1,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post("http://localhost:8545", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if response.Error != nil {
		return &response, fmt.Errorf("error from server: %s", *response.Error)
	}

	return &response, nil
}

func LoadTransactionsFromFile(filename string) ([][]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return [][]byte{}, nil
		}
		return nil, fmt.Errorf("error reading transactions file: %w", err)
	}

	if len(data) == 0 {
		return [][]byte{}, nil
	}

	var txBytesList [][]byte
	err = gob.NewDecoder(bytes.NewReader(data)).Decode(&txBytesList)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling transactions: %w", err)
	}
	return txBytesList, nil
}

func getDataForDeploySmartContract() ([]byte, error) {
	deployData := transaction.NewDeployData(
		common.FromHex(string("0x6060604052604060405190810160405280600581526020017f312e312e3300000000000000000000000000000000000000000000000000000081526020015060066000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061008c57805160ff19168380011785556100bd565b828001600101855582156100bd579182015b828111156100bc57825182600050559160200191906001019061009e565b5b5090506100e891906100ca565b808211156100e457600081815060009055506001016100ca565b5090565b50505b6a52b7d2dcc80cd2e4000000600060005060003373ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050819055506a52b7d2dcc80cd2e4000000600260005081905550604060405190810160405280600581526020017f526f626f7400000000000000000000000000000000000000000000000000000081526020015060036000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106101c757805160ff19168380011785556101f8565b828001600101855582156101f8579182015b828111156101f75782518260005055916020019190600101906101d9565b5b5090506102239190610205565b8082111561021f5760008181506000905550600101610205565b5090565b50506012600460006101000a81548160ff02191690837f0100000000000000000000000000000000000000000000000000000000000000908102040217905550604060405190810160405280600381526020017f524254000000000000000000000000000000000000000000000000000000000081526020015060056000509080519060200190828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106102eb57805160ff191683800117855561031c565b8280016001018555821561031c579182015b8281111561031b5782518260005055916020019190600101906102fd565b5b5090506103479190610329565b808211156103435760008181506000905550600101610329565b5090565b50505b610d48806103586000396000f3606060405236156100ab576000357c01000000000000000000000000000000000000000000000000000000009004806306fdde03146100bd578063095ea7b31461013d57806318160ddd1461017957806323b872dd146101a1578063313ce567146101e657806354fd4d501461021157806370a082311461029157806395d89b41146102c2578063a9059cbb14610342578063cae9ca511461037e578063dd62ed3e14610401576100ab565b34610002576100bb5b610002565b565b005b34610002576100cf600480505061043b565b60405180806020018281038252838181518152602001915080519060200190808383829060006004602084601f0104600302600f01f150905090810190601f16801561012f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b346100025761016160048080359060200190919080359060200190919050506104dc565b60405180821515815260200191505060405180910390f35b346100025761018b60048050506105b0565b6040518082815260200191505060405180910390f35b34610002576101ce60048080359060200190919080359060200190919080359060200190919050506105b9565b60405180821515815260200191505060405180910390f35b34610002576101f860048050506107c5565b604051808260ff16815260200191505060405180910390f35b346100025761022360048050506107d8565b60405180806020018281038252838181518152602001915080519060200190808383829060006004602084601f0104600302600f01f150905090810190601f1680156102835780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34610002576102ac6004808035906020019091905050610879565b6040518082815260200191505060405180910390f35b34610002576102d460048050506108b7565b60405180806020018281038252838181518152602001915080519060200190808383829060006004602084601f0104600302600f01f150905090810190601f1680156103345780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34610002576103666004808035906020019091908035906020019091905050610958565b60405180821515815260200191505060405180910390f35b34610002576103e96004808035906020019091908035906020019091908035906020019082018035906020019191908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050909091905050610a98565b60405180821515815260200191505060405180910390f35b34610002576104256004808035906020019091908035906020019091905050610cdf565b6040518082815260200191505060405180910390f35b60036000508054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156104d45780601f106104a9576101008083540402835291602001916104d4565b820191906000526020600020905b8154815290600101906020018083116104b757829003601f168201915b505050505081565b600081600160005060003373ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005060008573ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040518082815260200191505060405180910390a3600190506105aa565b92915050565b60026000505481565b600081600060005060008673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000505410158015610653575081600160005060008673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005060003373ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000505410155b801561065f5750600082115b156107b45781600060005060008573ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282825054019250508190555081600060005060008673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282825054039250508190555081600160005060008673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005060003373ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828282505403925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a3600190506107be566107bd565b600090506107be565b5b9392505050565b600460009054906101000a900460ff1681565b60066000508054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108715780601f1061084657610100808354040283529160200191610871565b820191906000526020600020905b81548152906001019060200180831161085457829003601f168201915b505050505081565b6000600060005060008373ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000505490506108b2565b919050565b60056000508054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156109505780601f1061092557610100808354040283529160200191610950565b820191906000526020600020905b81548152906001019060200180831161093357829003601f168201915b505050505081565b600081600060005060003373ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005054101580156109995750600082115b15610a885781600060005060003373ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282825054039250508190555081600060005060008573ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828282505401925050819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a360019050610a9256610a91565b60009050610a92565b5b92915050565b600082600160005060003373ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005060008673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050819055508373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925856040518082815260200191505060405180910390a38373ffffffffffffffffffffffffffffffffffffffff1660405180807f72656365697665417070726f76616c28616464726573732c75696e743235362c81526020017f616464726573732c627974657329000000000000000000000000000000000000815260200150602e01905060405180910390207c0100000000000000000000000000000000000000000000000000000000900433853086604051857c0100000000000000000000000000000000000000000000000000000000028152600401808573ffffffffffffffffffffffffffffffffffffffff1681526020018481526020018373ffffffffffffffffffffffffffffffffffffffff1681526020018280519060200190808383829060006004602084601f0104600302600f01f150905090810190601f168015610ca75780820380516001836020036101000a031916815260200191505b509450505050506000604051808303816000876161da5a03f1925050501515610ccf57610002565b60019050610cd8565b9392505050565b6000600160005060008473ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060005060008373ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600050549050610d42565b9291505056")),
		common.HexToAddress("da7284fac5e804f8b9d71aa39310f0f86776b51d"),
	)
	return deployData.Marshal()
}

type KeyPair struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Address    string `json:"address"`
	Key        string `json:"key"` // Có vẻ như đây là một trường không cần thiết
}

func AddressToIndex(input string) int {
	if len(input) < 2 {
		return 0 // Xử lý trường hợp chuỗi quá ngắn
	}

	// Lấy hai ký tự đầu tiên
	twoChars := input[:2]

	// Chuyển đổi thành số nguyên
	num, err := strconv.ParseInt(twoChars, 16, 64) // Sử dụng cơ số 16 (hexadecimal)
	if err != nil {
		return 0 // Xử lý lỗi chuyển đổi
	}

	// Trả về số dư khi chia cho 16
	return int(num % 16)
}

func TestConcurrentSendTransaction2Check(t *testing.T) {

	data, err := os.ReadFile("kp1.json")
	if err != nil {
		logger.Error(err)
	}

	// Khởi tạo mảng để lưu trữ các KeyPair
	var keyPairs []KeyPair

	// Giải mã JSON thành mảng KeyPair
	err = json.Unmarshal(data, &keyPairs)
	if err != nil {
		logger.Error(err)
	}

	var wg sync.WaitGroup
	// results := make(chan types.AccountState, len(keyPairs)) // Sử dụng channel để thu thập kết quả
	for _, kp := range keyPairs {
		fmt.Printf("Private Key: %s\n", kp.PrivateKey)
		fmt.Printf("Public Key: %s\n", kp.PublicKey)
		fmt.Printf("Address: %s\n", kp.Address)
		// Bỏ qua trường "key" vì nó không cần thiết
		fmt.Println("------------------")
		wg.Add(1)
		go func(kp KeyPair) {
			defer wg.Done()
			logger.Info("SendTransaction: ", kp.Address)
			// Sử dụng một channel để lưu trữ kết quả của mỗi giao dịch
			// results <- as
		}(kp)
	}
	// Đảm bảo rằng tất cả goroutines hoàn thành trước khi tiếp tục
	wg.Wait()
	// close(results)

}

func TestConcurrentAccountState(t *testing.T) {
	// fileLogger, errLog := logger.NewFileLogger("application.log")
	// if errLog != nil {
	// 	log.Error("Could not create log file: %v", errLog)
	// }
	// defer fileLogger.Close()
	// Khởi tạo cấu hình client (giống như trong ví dụ ban đầu)
	// clientConfigJson := `
	// {
	// 	"version": "0.0.1.0",
	// 	"private_key": "2b3aa0f620d2d73c046cd93eb64f2eb687a95b22e278500aa251c8c9dda1203b",
	// 	"parent_address": "9065e2caaa4bf6533ded0b2a763e4ee81cd31bd2",
	// 	"parent_connection_address": "0.0.0.0:4200",
	// 	"parent_connection_type": "node",
	// 	"dns_link": "http://127.0.0.1:7777/api/dns/connection-address/"
	// }
	// `
	// config := &c_config.ClientConfig{}
	// json.Unmarshal([]byte(clientConfigJson), config)
	// client, _ = NewClient(config)

	// Địa chỉ tài khoản cần kiểm tra (thay thế bằng địa chỉ cần thiết)
	targetAddress := common.HexToAddress("0x97126b71376f7e55fba904fdaa9df0dbd396612f")

	// Số lần chạy
	numIterations := 1000000

	// Sử dụng goroutines và wait group để xử lý đồng thời
	var wg sync.WaitGroup
	results := make(chan AccountStateResult, numIterations) // Sử dụng channel để thu thập kết quả

	for i := 0; i < numIterations; i++ {
		wg.Add(1)
		go func(iteration int) {
			defer wg.Done()
			as, err := client.AccountState(targetAddress)
			results <- AccountStateResult{Iteration: iteration, AccountState: as, Error: err}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Xử lý kết quả từ channel
	for result := range results {
		if result.Error != nil {
			t.Errorf("Lỗi trong lần chạy %d: %v", result.Iteration, result.Error)
		}
		fmt.Printf("Lần chạy %d: %+v\n", result.Iteration, result.AccountState)
	}
}

// Cấu trúc để lưu trữ kết quả của mỗi lần chạy
type AccountStateResult struct {
	Iteration    int
	AccountState types.AccountState
	Error        error
}

type SendTransactionResult struct {
	Iteration int
	Receipt   types.Receipt
	Error     error
}

type SendTransactionResult2 struct {
	Address string
	Receipt types.Receipt
	Error   error
}

func TestMultipleClient(t *testing.T) {
	clientConfigJson := `
	{
		"version": "0.0.1.0",
		"private_key": "1c923d7764cb712f2f007a53f5c14f898fc3fcc3f00c609c55ec5cb4d7443211",
		"parent_address": "b4970c8ff037fcbef0e7aba0cdc3aedc332820ba",
		"parent_connection_address": "127.0.0.1:3011",
		"parent_connection_type": "node"
	}
	`
	config := &c_config.ClientConfig{}
	json.Unmarshal([]byte(clientConfigJson), config)
	client, _ := NewClient(
		config,
	)

	clientConfigJson2 := `
	{
		"version": "0.0.1.0",
		"private_key": "16f497a07ff4df0dc12488c06864de8b5d8566572c9604f99e2a394f0991e3b3",
		"parent_address": "b4970c8ff037fcbef0e7aba0cdc3aedc332820ba",
		"parent_connection_address": "127.0.0.1:3011",
		"parent_connection_type": "node"
	}
	`
	config2 := &c_config.ClientConfig{}
	json.Unmarshal([]byte(clientConfigJson2), config2)
	client2, _ := NewClient(
		config2,
	)

	as, err := client.AccountState(
		common.HexToAddress("ae357fc27436ed8aeeb2df11cbd745a12b2e2093"),
	)
	assert.Nil(t, err)
	logger.Info(as)

	as2, err := client2.AccountState(
		common.HexToAddress("5c65574e19415b72b6f458f9510d65d84c034c4c"),
	)
	assert.Nil(t, err)
	logger.Info(as2)
}

func TestSubscribe(t *testing.T) {
	eventChan, err := client.Subcribe(
		common.HexToAddress("da7284fac5e804f8b9d71aa39310f0f86776b51d"),
		common.HexToAddress("0xb3E65f6e1f1Cccb759Be53c743B0DDC9C2ecaf62"),
	)
	assert.Nil(t, err)
	eventlog := <-eventChan
	logger.Info(eventlog)
}

func TestClose(t *testing.T) {
	client.Close()
}

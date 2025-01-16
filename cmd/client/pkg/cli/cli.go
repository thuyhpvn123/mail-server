package cli

import (
	"bufio"
	// "encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	// "time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"gomail/cmd/client/command"
	"gomail/cmd/client/pkg/client_context"
	client_types "gomail/cmd/client/types"
	"gomail/mtn/bls"
	p_common "gomail/mtn/common"
	"gomail/mtn/logger"
	"gomail/mtn/network"
	pb "gomail/mtn/proto"
	"gomail/mtn/transaction"
	"gomail/mtn/types"
	"gomail/cmd/client"
)

var (
	ErrorGetAccountStateTimedOut = errors.New("get account state timed out")
	ErrorInvalidAction           = errors.New("invalid action")
)

type Cli struct {
	client *client.Client
	clientContext *client_context.ClientContext

	//
	stop     bool
	commands map[int]string
	reader   *bufio.Reader

	transactionController client_types.TransactionController
	// accountStateChan      chan types.AccountState
	defaultRelatedAddress map[common.Address][][]byte
}

func NewCli(
	client *client.Client,
	clientContext *client_context.ClientContext,
	transactionController client_types.TransactionController,
	// accountStateChan chan types.AccountState,
) client_types.Cli {

	commands := map[int]string{
		0: "Exit",
		1: "Send transaction",
		2: "Change account",
		3: "Create account",
		4: "Get account state",
		5: "Subscribe",
		6: "Get stake state",
		// 7: "Get smart contract data",
		8: "Get stats",
		9: "Change log level",
	}
	return &Cli{
		client: client,
		clientContext:         clientContext,
		stop:                  false,
		commands:              commands,
		transactionController: transactionController,
		// accountStateChan:      accountStateChan,
		defaultRelatedAddress: make(map[common.Address][][]byte),
	}
}

func (cli *Cli) Start() {
	cli.reader = bufio.NewReader(os.Stdin)
	for {
		if cli.stop {
			return
		}
		cli.PrintCommands()

		command := cli.ReadInput()
		switch command {
		case "0":
		// TODO
		case "1":
			err := cli.SendTransaction()
			if err != nil {
				logger.Warn("err", err)
			}
		case "2":
			cli.ChangeAccount()
		case "3":
			cli.CreateAccount()
		case "4":
			cli.PrintMessage("Enter address: ", "")
			cli.client.AccountState(cli.ReadInputAddress())
			// TODO4
		case "5":
			cli.Subscribe()
		case "6":
			cli.PrintMessage("Enter address: ", "")
			cli.StakeState(cli.ReadInputAddress())
		case "7":
		case "8":
			cli.GetStats()
		case "9":
			cli.ChangeLogLevel()
		}
	}
}

func (cli *Cli) Subscribe() {
	cli.PrintMessage("Enter smart contract storage address: ", "")
	storageAddress := cli.ReadInput()
	cli.PrintMessage("Enter smart contract address: ", "")
	contractAddress := cli.ReadInput()

	storageConnection := network.NewConnection(
		common.HexToAddress(storageAddress),
		p_common.STORAGE_CONNECTION_TYPE,
		cli.clientContext.Config.DnsLink(),
	)

	err := storageConnection.Connect()
	if err != nil {
		logger.Error("Subscribe fail", err)
		return
	}
	go cli.clientContext.SocketServer.HandleConnection(storageConnection)

	err = cli.clientContext.MessageSender.SendBytes(
		storageConnection,
		command.SubscribeToAddress,
		common.HexToAddress(contractAddress).Bytes(),
	)
	if err != nil {
		logger.Error("Subscribe fail", err)
	}
	logger.Debug("Subscribe address: ", contractAddress)
}

// TODE Cli stop
func (cli *Cli) Stop() {
}

func (cli *Cli) PrintCommands() {
	str := p_common.Cyan + "======= Commands =======\n" + p_common.Purple
	for i := 0; i < len(cli.commands); i++ {
		str += fmt.Sprintf("%v: %v\n", i, cli.commands[i])
	}
	str += p_common.Reset
	fmt.Print(str)
}

func (cli *Cli) SendTransaction() error {
	cli.PrintMessage("Enter to address: ", "")
	toAddress := cli.ReadInputAddress()
	cli.PrintMessage("Enter to amount (default 10*10^18): ", "")
	amount := cli.ReadBigInt()
	cli.PrintMessage(`Enter action (default 0):
	0: None
	1: Stake
	2: Unstake
	3: Deploy smart contract
	4: Call smart contract
	8: Open state channel
	9: Join state channel
	10: Commit state channel account state
	11: Commit state channel
	12: update storage host 
	`, "")

	actionStr := cli.ReadInput()
	var action pb.ACTION
	if actionStr == "" {
		action = 0
	} else {
		actionI, _ := strconv.Atoi(actionStr)
		action = pb.ACTION(int32(actionI))
	}
	if action < 0 || action > 20 {
		return ErrorInvalidAction
	}

	var data []byte
	if action == pb.ACTION_UPDATE_STORAGE_HOST {
		cli.PrintMessage("Enter new storage host: ", "")
		storageHost := cli.ReadInput()
		cli.PrintMessage("Enter new storage address: ", "")
		storageAddress := cli.ReadInputAddress()
		updateData := transaction.NewUpdateStorageHostData(storageHost, storageAddress)
		data, _ = updateData.Marshal()
	}

	var err error
	as, err := cli.client.AccountState(cli.clientContext.KeyPair.Address())
	if err != nil {
		return err
	}

	if action == pb.ACTION_DEPLOY_SMART_CONTRACT {
		data, err = cli.getDataForDeploySmartContract()
		if err != nil {
			panic(err)
		}
		toAddress = common.BytesToAddress(
			crypto.Keccak256(
				append(
					as.Address().Bytes(),
					as.LastHash().Bytes()...),
			)[12:],
		)
	}
	// var commissionPrivateKey []byte
	if action == pb.ACTION_CALL_SMART_CONTRACT {
		data, err = cli.getDataForCallSmartContract()
		if err != nil {
			panic(err)
		}

		cli.PrintMessage("Enter to private key for commission sign: ", "")
		// hexCommissionPrivateKey := cli.ReadInput()
		// commissionPrivateKey = common.FromHex(hexCommissionPrivateKey)
	}

	var relatedAddresses [][]byte
	if action == pb.ACTION_CALL_SMART_CONTRACT || action == pb.ACTION_DEPLOY_SMART_CONTRACT {
		relatedAddresses = cli.ReadRelatedAddress(toAddress)
	}

	if action == pb.ACTION_OPEN_CHANNEL {
		validatorAddresses := cli.ReadValidatorAddresses()
		data, _ = transaction.NewOpenStateChannelData(validatorAddresses).Marshal()
		toAddress = common.BytesToAddress(
			crypto.Keccak256(
				append(
					as.Address().Bytes(),
					as.LastHash().Bytes()...),
			)[12:],
		)
	}

	cli.PrintMessage("Enter max gas (default 500000): ", "")
	maxGas, err := strconv.ParseUint(cli.ReadInput(), 10, 64)
	if err != nil {
		if action == pb.ACTION_OPEN_CHANNEL {
			maxGas = p_common.OPEN_CHANNEL_GAS_COST
		} else {
			// maxGas = 500000
			maxGas = 100000

		}
	}

	cli.PrintMessage("Enter max gas price in gwei (default 10 gwei): ", "")
	maxGasPriceGwei, err := strconv.ParseUint(cli.ReadInput(), 10, 64)
	if err != nil {
		maxGasPriceGwei = 10
	}
	// maxGasPrice := 1000000000 * maxGasPriceGwei
	maxGasPrice := 1000000 * maxGasPriceGwei

	cli.PrintMessage("Enter max time use in milli second (default 1000 milli second): ", "")
	maxTimeUse, err := strconv.ParseUint(cli.ReadInput(), 10, 64)
	if err != nil {
		maxTimeUse = 1000
	}

	var transaction types.Receipt


	relatedAddress := make([]common.Address, 0)
	for i,v := range relatedAddresses {
		relatedAddress[i] = common.HexToAddress(string(v))
	}
	transaction, err = cli.client.SendTransactionWithDeviceKey(
		toAddress,
		amount,
		action,
		data,
		relatedAddress,
		maxGas,
		maxGasPrice,
		maxTimeUse,
	)

	logger.Debug("Sending transaction", transaction)
	if err != nil {
		logger.Warn(err)
	}

	return err
}

func (cli *Cli) ChangeAccount() {
	cli.PrintMessage("Enter to private key(hex): ", "")
	hexPrivateKey := cli.ReadInput()
	keyPair := bls.NewKeyPair(common.FromHex(hexPrivateKey))
	cli.clientContext.KeyPair = keyPair
	cli.clientContext.SocketServer.SetKeyPair(keyPair)
	// disconnect parent connection and init new conneciton

	parentConn := cli.clientContext.ConnectionsManager.ParentConnection()
	if parentConn != nil {
		newParentConn := parentConn.Clone()
		cli.clientContext.ConnectionsManager.AddParentConnection(newParentConn)
		parentConn.Disconnect()
		err := newParentConn.Connect()
		if err != nil {
			logger.Warn("error when connect to parent", err)
		} else {
			cli.clientContext.SocketServer.OnConnect(newParentConn)
			go cli.clientContext.SocketServer.HandleConnection(newParentConn)
		}
	}
	logger.Info("Running with key pair:", keyPair)
}

func (cli *Cli) CreateAccount() {
	keyPair := bls.GenerateKeyPair()
	logger.Info(fmt.Sprintf("Key pair:\n%v", keyPair))
}

func (cli *Cli) ReadInput() string {
	input, err := cli.reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	input = strings.Replace(input, "\n", "", -1)
	return input
}

func (cli *Cli) ReadInputAddress() common.Address {
	input := cli.ReadInput()
	address := common.HexToAddress(input)
	return address
}

func (cli *Cli) ReadBigInt() *big.Int {
	input := cli.ReadInput()
	if input == "" {
		input = "10000000000000000000"
	}
	bigInt := big.NewInt(0)
	bigInt.SetString(input, 10)
	return big.NewInt(0).SetBytes(bigInt.Bytes())
}

func (cli *Cli) PrintMessage(message string, color string) {
	if color == "" {
		color = p_common.Purple
	}
	fmt.Printf(color+"%v\n"+p_common.Reset, message)
}

func (cli *Cli) AccountState(address common.Address) (types.AccountState, error) {
	parentConn := cli.clientContext.ConnectionsManager.ParentConnection()
	cli.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetAccountState,
		address.Bytes(),
	)
	// select {
	// case accountState := <-cli.accountStateChan:
	// 	return accountState, nil
	// case <-time.After(2 * time.Second):
		return nil, ErrorGetAccountStateTimedOut
	// }
}

func (cli *Cli) StakeState(address common.Address) {
	parentConn := cli.clientContext.ConnectionsManager.ParentConnection()
	cli.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetStakeState,
		address.Bytes(),
	)
}

func (cli *Cli) getDataForDeploySmartContract() ([]byte, error) {
	cli.PrintMessage("Enter to smart contract file name (in contracts folder): ", "")
	contractFileName := cli.ReadInput()
	b, _ := os.ReadFile("./contracts/" + contractFileName)
	cli.PrintMessage("Enter smart contract storage host: ", "")
	contractStorageHost := cli.ReadInput()
	if contractStorageHost == "" {
		contractStorageHost = "127.0.0.1:3051"
	}
	cli.PrintMessage("Enter smart contract storage address: ", "")
	contractStorageAddress := cli.ReadInput()
	if contractStorageAddress == "" {
		contractStorageAddress = "da7284fac5e804f8b9d71aa39310f0f86776b51d"
	}
	deployData := transaction.NewDeployData(
		common.FromHex(string(b)),
		common.HexToAddress(contractStorageAddress),
	)
	return deployData.Marshal()
}

func (cli *Cli) getDataForCallSmartContract() ([]byte, error) {
	cli.PrintMessage("Enter to input for call smart contract (hex): ", "")
	input := cli.ReadInput()
	callData := transaction.NewCallData(common.FromHex(input))
	return callData.Marshal()
}

func (cli *Cli) ReadRelatedAddress(smartcontractAddress common.Address) [][]byte {
	cli.PrintMessage("Enter Related Address: ", "")
	stringRelatedAddresses := cli.ReadInput()
	if stringRelatedAddresses == "" {
		if cli.defaultRelatedAddress[smartcontractAddress] == nil {
			return [][]byte{}
		}
		return cli.defaultRelatedAddress[smartcontractAddress]
	}
	hexRelatedAddresses := strings.Split(stringRelatedAddresses, ",")
	relatedAddresses := make([][]byte, len(hexRelatedAddresses))
	logger.Debug("Temp Related Address")
	for idx, hexAddress := range hexRelatedAddresses {
		address := common.HexToAddress(hexAddress)
		logger.Debug(address)
		relatedAddresses[idx] = address.Bytes()
	}
	cli.defaultRelatedAddress[smartcontractAddress] = append(
		cli.defaultRelatedAddress[smartcontractAddress],
		relatedAddresses...)
	return relatedAddresses
}

func (cli *Cli) ReadValidatorAddresses() []common.Address {
	cli.PrintMessage("Enter Validator Addresses: ", "")
	stringvalidatorAddresses := cli.ReadInput()
	if stringvalidatorAddresses == "" {
		return []common.Address{}
	}
	hexValidatorAddresses := strings.Split(stringvalidatorAddresses, ",")
	validatorAddresses := make([]common.Address, len(hexValidatorAddresses))
	for idx, hexAddress := range hexValidatorAddresses {
		address := common.HexToAddress(hexAddress)
		logger.Debug(address)
		validatorAddresses[idx] = address
	}
	return validatorAddresses
}

func (cli *Cli) GetStats() {
	parentConn := cli.clientContext.ConnectionsManager.ParentConnection()
	cli.clientContext.MessageSender.SendBytes(
		parentConn,
		command.GetStats,
		[]byte{},
	)
}

func (cli *Cli) ChangeLogLevel() {
	parentConn := cli.clientContext.ConnectionsManager.ParentConnection()
	str := p_common.Cyan + "======= Log level =======\n" + p_common.Purple
	loglevel := map[int]string{
		0: "DEBUGP",
		1: "ERROR",
		2: "WARN",
		3: "INFO",
		4: "DEBUG",
		5: "TRACE",
	}
	for i := 0; i < len(loglevel); i++ {
		str += fmt.Sprintf("%v: %v\n", i, loglevel[i])
	}
	fmt.Print(str)
	level, err := strconv.Atoi(cli.ReadInput())
	if err != nil {
		logger.Error(err)
		return
	}
	cli.clientContext.MessageSender.SendBytes(
		parentConn,
		command.ChangeLogLevel,
		uint256.NewInt(
			uint64(level)).Bytes(),
	)
}

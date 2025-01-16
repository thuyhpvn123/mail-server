package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	// "github.com/holiman/uint256"

	"gomail/cmd/client"
	c_config "gomail/cmd/client/pkg/config"
	"gomail/mtn/bls"
	p_common "gomail/mtn/common"
	"gomail/mtn/logger"
	pb "gomail/mtn/proto"
	"gomail/mtn/transaction"
)

const (
	defaultLogLevel        = logger.FLAG_DEBUG
	defaultConfigPath      = "config.json"
	defaultDataFile        = "data.json"
	defaultSkipKeyPress    = false
	REPLACE_ADDRESS        = "1510151015101510151015101510151015101510"
	REPLACE_VAR            = "$"
	OUTPUT_RELATED_ADDRESS = "related_address"
)

var (
	CONFIG_FILE_PATH string
	DATA_FILE_PATH   string
	LOG_LEVEL        int
	SKIP_KEY_PRESS   bool
	//
)

type SCData struct {
	Action         string   `json:"action"`
	Input          string   `json:"input"`
	Amount         string   `json:"amount"`
	Address        string   `json:"address"`
	RelatedAddress []string `json:"related_address"`
	StorageHost    string   `json:"storage_host"`
	StorageAddress string   `json:"storage_address"`
	Export         string   `json:"export"`
	ReplaceAddress []int    `json:"replace_address"`
	Type           string   `json:"type"`
	InputPath      string   `json:"input_path"`
	ParentAddress  string   `json:"parent_address"`
	Name           string   `json:"name"`
	Output         struct {
		Index int    `json:"index"`
		Type  string `json:"type"`
	} `json:"output"`
}

func main() {
	flag.IntVar(&LOG_LEVEL, "log-level", defaultLogLevel, "Log level")
	flag.IntVar(&LOG_LEVEL, "ll", defaultLogLevel, "Log level (shorthand)")

	flag.StringVar(&CONFIG_FILE_PATH, "config", defaultConfigPath, "Config path")
	flag.StringVar(&CONFIG_FILE_PATH, "c", defaultConfigPath, "Config path (shorthand)")

	flag.StringVar(&DATA_FILE_PATH, "data", defaultDataFile, "Data file path")
	flag.StringVar(&DATA_FILE_PATH, "d", defaultDataFile, "Data file path (shorthand)")

	flag.BoolVar(&SKIP_KEY_PRESS, "skip", defaultSkipKeyPress, "Skip press to run new transaction")
	flag.BoolVar(
		&SKIP_KEY_PRESS,
		"s",
		defaultSkipKeyPress,
		"Skip press to run new transaction (shorthand)",
	)

	flag.Parse()
	// set logger config
	loggerConfig := &logger.LoggerConfig{
		Flag:    LOG_LEVEL,
		Outputs: []*os.File{os.Stdout},
	}
	logger.SetConfig(loggerConfig)

	config, err := c_config.LoadConfig(CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.ClientConfig)

	sendTransactionsLoop(
		cConfig,
	)

	logger.Debug("Done. Press any key to exist")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

func sendTransactionsLoop(
	config *c_config.ClientConfig,
) {
	output := []string{}
	addressList := []string{}
	c, err := client.NewClient(
		config,
	)
	if err != nil {
		logger.Error("Error when init client", err)
		panic(err)
	}

	datas := getDatas()

	logger.Info("Total request", len(datas))

	for i := 0; i < len(datas); i++ {
		maxGas := uint64(10000000)
		maxGasPrice := uint64(p_common.MINIMUM_BASE_FEE)
		amount, success := big.NewInt(0).SetString(datas[i].Amount, 10)
		if !success {
			logger.Error("Unable to parse amount")
			panic("Unable to parse amount")
		}
		var action pb.ACTION
		var toAddress common.Address
		var bData []byte

		for j := 0; j < len(datas[i].ReplaceAddress); j++ {
			datas[i].Input = strings.Replace(
				datas[i].Input,
				REPLACE_ADDRESS,
				addressList[datas[i].ReplaceAddress[j]],
				1,
			)
		}

		if datas[i].Action == "deploy" {
			as, err := c.AccountState(config.Address())
			if err != nil {
				logger.Error("Error when aget account state", err)
				panic(err)
			}
			action = pb.ACTION_DEPLOY_SMART_CONTRACT
			toAddress = common.BytesToAddress(
				crypto.Keccak256(
					append(
						as.Address().Bytes(),
						as.LastHash().Bytes()...),
				)[12:],
			)
			deployData := transaction.NewDeployData(
				common.FromHex(datas[i].Input),
				// datas[i].StorageHost,
				common.HexToAddress(datas[i].StorageAddress),
			)
			bData, err = deployData.Marshal()
			if err != nil {
				panic(err)
			}
		} else if datas[i].Action == "call" {
			action = pb.ACTION_CALL_SMART_CONTRACT
			if len(datas[i].Address) < 40 {
				index, err := strconv.Atoi(datas[i].Address)
				if err != nil {
					panic(err)
				}
				datas[i].Address = addressList[index]
			}
			toAddress = common.HexToAddress(datas[i].Address)
			callData := transaction.NewCallData(common.FromHex(datas[i].Input))
			bData, err = callData.Marshal()
			if err != nil {
				panic(err)
			}
		} else if datas[i].Action == "pos" {
			action = pb.ACTION_CALL_SMART_CONTRACT
			if len(datas[i].Address) < 40 {
				index, err := strconv.Atoi(datas[i].Address)
				if err != nil {
					panic(err)
				}
				datas[i].Address = addressList[index]
			}
			logger.Error(datas[i].Address)
			dat, _ := os.ReadFile(datas[i].InputPath)
			lBranch := [][]struct {
				Address string `json:"address"`
				Parent  string `json:"parent_address"`
				Type    string `json:"type"`
			}{}
			err := json.Unmarshal(dat, &lBranch)
			if err != nil {
				panic(err)
			}
			for _, nodes := range lBranch {
				for _, node := range nodes {
					inputCallData := datas[i].Input
					nType := "0"
					switch node.Type {
					case "node":
						nType = "1"
					case "vminer":
						nType = "2"
					case "eminer":
						nType = "3"
					}
					inputCallData = strings.Replace(inputCallData, REPLACE_VAR, nType, 1)
					inputCallData = strings.Replace(inputCallData, REPLACE_ADDRESS, hex.EncodeToString(common.HexToAddress(node.Address).Bytes()), 1)
					inputCallData = strings.Replace(inputCallData, REPLACE_ADDRESS, hex.EncodeToString(common.HexToAddress(node.Parent).Bytes()), 1)
					println(inputCallData)
					callData := transaction.NewCallData(common.FromHex(inputCallData))
					bData, err = callData.Marshal()
					if err != nil {
						panic(err)
					}
					toAddress = common.HexToAddress(datas[i].Address)
					relatedAddress := make([]common.Address, 0)

					receipt, err := c.SendTransactionWithDeviceKey(
						toAddress,
						big.NewInt(0),
						action,
						bData,
						relatedAddress,
						maxGas,
						maxGasPrice,
						0,
					)
					if err != nil {
						logger.Error("error when send transaction", err)
						panic("error when send transaction")
					}
					logger.Info("Receive receipt", receipt)

				}
			}

			return
		}
		lenRelatedAddress := len(datas[i].RelatedAddress)
		relatedAddress := make([]common.Address, lenRelatedAddress+1)
		for i, v := range datas[i].RelatedAddress {
			if len(v) < 40 {
				index, err := strconv.Atoi(v)
				if err != nil {
					panic(err)
				}
				v = addressList[index]
			}
			relatedAddress[i] = common.HexToAddress(v)
		}
		relatedAddress[lenRelatedAddress] = bls.NewKeyPair(config.PrivateKey()).Address()
		logger.Info(relatedAddress)
		logger.Info(toAddress)

		receipt, err := c.SendTransactionWithDeviceKey(
			toAddress,
			big.NewInt(0).SetBytes(amount.Bytes()),
			action,
			bData,
			relatedAddress,
			maxGas,
			maxGasPrice,
			0,
		)
		if err != nil {
			logger.Error("error when send transaction", err)
			panic("error when send transaction")
		}
		logger.Info("Receive receipt", receipt)

		switch datas[i].Output.Type {
		case OUTPUT_RELATED_ADDRESS:
			outputAddress := common.BytesToAddress(receipt.Return()).String()
			if datas[i].Output.Index == 0 {
				addressList = append(addressList, outputAddress)
			} else {
				addressList[datas[i].Output.Index] = outputAddress
			}
		}
		if datas[i].Action == "deploy" {
			addressList = append(addressList, hex.EncodeToString(receipt.ToAddress().Bytes()))

		}
		if receipt.Status() != pb.RECEIPT_STATUS_RETURNED &&
			receipt.Status() != pb.RECEIPT_STATUS_HALTED {
			panic("Fail transaction")
		}
		fPath := datas[i].Export
		if fPath != "" {
			output := make([]SCData, 1)
			lsRelatedAddress := make([]string, len(relatedAddress))
			for k, v := range relatedAddress {
				lsRelatedAddress[k] = v.String()
			}

			datas[i].RelatedAddress = lsRelatedAddress
			datas[i].Export = ""
			datas[i].ReplaceAddress = make([]int, 0)
			output[0] = datas[i]
			jsonData, err := json.MarshalIndent(output, "", "    ")
			if err != nil {
				panic("Error marshaling JSON")
			}

			err = os.WriteFile(fPath, jsonData, 0644)
			if err != nil {
				panic("Error writing to file:")
			}
			logger.Debug("Data has been successfully written ", datas[i].Export)
		}
		output = append(output, datas[i].Name)
		output = append(output, common.BytesToAddress(receipt.Return()).String())
		jsonData, err := json.MarshalIndent(output, "", "    ")
		if err != nil {
			panic("Error marshaling JSON")
		}

		err = os.WriteFile("output.json", jsonData, 0644)
		if err != nil {
			panic("Error writing to file:")
		}

		if !SKIP_KEY_PRESS {
			logger.Debug("Press any key to continue")
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
		}
	}
}

func getDatas() []SCData {
	dat, _ := os.ReadFile(DATA_FILE_PATH)
	scDatas := []SCData{}
	err := json.Unmarshal(dat, &scDatas)
	if err != nil {
		panic(err)
	}
	return scDatas
}
package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"

	// "github.com/holiman/uint256"

	"gomail/cmd/client"
	c_config "gomail/cmd/client/pkg/config"
	"gomail/pkg/logger"
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
	PrivateKeySecp string `json:"private_key_secp"`
	ChainId        string `json:"chain_id"`
	Name           string `json:"name"`
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

	var nonce uint64
	logger.Info("Total request", len(datas))

	for i := 0; i < len(datas); i++ {

		nonce++

		receipt, err := c.AddAccountForClient(
			datas[i].PrivateKeySecp,
			datas[i].ChainId,
		)
		if err != nil {
			logger.Error("error when send transaction", err)
			panic("error when send transaction")
		}

		switch datas[i].Output.Type {
		case OUTPUT_RELATED_ADDRESS:
			outputAddress := common.BytesToAddress(receipt.Return()).String()
			if datas[i].Output.Index == 0 {
				addressList = append(addressList, outputAddress)
			} else {
				addressList[datas[i].Output.Index] = outputAddress
			}
		}

		output = append(output, datas[i].Name)

		output = append(output, hex.EncodeToString(receipt.Return()))

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

package client

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"

	c_config "gomail/cmd/client/pkg/config"
	"gomail/mtn/logger"
	pb "gomail/mtn/proto"
	p_transaction "gomail/mtn/transaction"
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
		"parent_address": "b4970c8ff037fcbef0e7aba0cdc3aedc332820ba",
		"parent_connection_address": "35.243.233.132:3061",
		"parent_connection_type": "node"
	}
	`
	config := &c_config.ClientConfig{}
	json.Unmarshal([]byte(clientConfigJson), config)
	client, _ = NewClient(
		config,
	)
}

func TestSendTransaction(t *testing.T) {
	receipt, err := client.SendTransaction(
		common.HexToAddress("9c7b6a12d1ea977a95acfed604079e01220d5cb3"),
		big.NewInt(10),
		pb.ACTION_EMPTY,
		[]byte{},
		nil,
		100000,
		1000000000,
		0,
	)
	assert.Nil(t, err)
	logger.Info(receipt)
}

func TestGetAccountState(t *testing.T) {
	as, err := client.AccountState(
		common.HexToAddress("0x97126B71376F7e55fBA904FdaA9dF0dBd396612f"),
	)
	assert.Nil(t, err)
	logger.Info(as)
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

func TestMy(t *testing.T) {
	clientConfigJson := `
	{
		"version": "0.0.1.0",
		"private_key": "2b3aa0f620d2d73c046cd93eb64f2eb687a95b22e278500aa251c8c9dda1203b",
		"parent_address": "b4970c8ff037fcbef0e7aba0cdc3aedc332820ba",
		"parent_connection_address": "35.243.233.132:6011",
		"parent_connection_type": "node"
	}
	`
	config := &c_config.ClientConfig{}
	json.Unmarshal([]byte(clientConfigJson), config)
	myClient, _ := NewClient(
		config,
	)
	callData := p_transaction.NewCallData(
		common.FromHex(
			"c653c9a10000000000000000000000000000000000000000000000000000000000000001000000000000000000000000fdd11471417109d88c48030e579f3523e485f6fa",
		),
	)
	bCallData, _ := callData.Marshal()
	_, err := myClient.SendTransaction(
		common.HexToAddress("E9D6dedC8f914f9CeC39D91176A10321f7992fC8"),
		big.NewInt(0),
		pb.ACTION_CALL_SMART_CONTRACT,
		bCallData,
		nil,
		200000,
		1000000000,
		0,
	)

	if err != nil {
		log.Warn("error:", err)
	} else {
		logger.Info("Done send transaction from ")
	}
}

func TestClose(t *testing.T) {
	client.Close()
}

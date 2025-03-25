package config

import (
	"encoding/json"
	"os"

	"gomail/pkg/bls"
	p_common "gomail/pkg/common"
	"gomail/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

type ClientConfig struct {
	PrivateKey_ string `json:"private_key"`

	ConnectionAddress_       string `json:"connection_address"`
	PublicConnectionAddress_ string `json:"public_connection_address"`
	DnsLink_                 string `json:"dns_link"`

	Version_          string       `json:"version"`
	TransactionFeeHex string       `json:"transaction_fee"`
	TransactionFee    *uint256.Int `json:"-"`

	ParentAddress           string `json:"parent_address"`
	ParentConnectionAddress string `json:"parent_connection_address"`
	ParentConnectionType    string `json:"parent_connection_type"`
	ChainId                 uint64 `json:"chain_id"`
}

func (c *ClientConfig) ConnectionAddress() string {
	return c.ConnectionAddress_
}

func (c *ClientConfig) PublicConnectionAddress() string {
	return c.PublicConnectionAddress_
}

func (c *ClientConfig) Version() string {
	return c.Version_
}

func (c *ClientConfig) PrivateKey() []byte {
	return common.FromHex(c.PrivateKey_)
}

func (c *ClientConfig) Address() common.Address {
	_, _, address := bls.GenerateKeyPairFromSecretKey(c.PrivateKey_)
	return address
}

func (c *ClientConfig) NodeType() string {
	return p_common.CLIENT_CONNECTION_TYPE
}

func (c *ClientConfig) DnsLink() string {
	return c.DnsLink_
}

func LoadConfig(configPath string) (types.Config, error) {
	// general config
	config := &ClientConfig{}
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(raw, config)
	if err != nil {
		return nil, err
	}
	config.TransactionFee = uint256.NewInt(0).SetBytes(common.FromHex(config.TransactionFeeHex))
	return config, nil
}

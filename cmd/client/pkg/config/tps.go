package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type TpsConfig struct {
	PrivateKey_   string
	ParentAddress string
	DnsLink_      string
	ScAddresses   []string
}

var Config *TpsConfig

func LoadTpsConfig(configFilePath string) (*TpsConfig, error) {
	viper.SetConfigFile(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config TpsConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

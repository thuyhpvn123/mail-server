package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type LastHashData struct {
	LastHash string `json:"lastHash"`
}

func SaveLastHash(filePath string, lastHash common.Hash) error {
	data := LastHashData{LastHash: lastHash.Hex()}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func ReadLastHash(filePath string) (common.Hash, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return common.Hash{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to read file: %w", err)
	}

	var transactionData LastHashData
	err = json.Unmarshal(data, &transactionData)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	hash := common.HexToHash(transactionData.LastHash)
	return hash, nil
}

// Uint64ToBytes converts a uint64 to a byte array.
func Uint64ToBytes(value uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, value)
	return bytes
}

// BytesToUint64 converts a byte array to a uint64.  Handles potential errors gracefully.
func BytesToUint64(bytes []byte) (uint64, error) {
	if len(bytes) != 8 {
		return 0, fmt.Errorf("byte array must be 8 bytes long")
	}
	return binary.BigEndian.Uint64(bytes), nil
}

// Uint64ToBigInt converts a uint64 to a big.Int.
func Uint64ToBigInt(value uint64) *big.Int {
	return new(big.Int).SetUint64(value)
}

func ParseBlockNumber(blockNumberStr string) (uint64, error) {
	// Kiểm tra xem blockNumberStr đã có tiền tố "0x" chưa
	blockNumberStr = strings.TrimPrefix(blockNumberStr, "0x")

	blockNumber, err := strconv.ParseUint(blockNumberStr, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("lỗi chuyển đổi block number: %w", err)
	}

	return blockNumber, nil
}

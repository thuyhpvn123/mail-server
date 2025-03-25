package txs_eth

import (
	"fmt"
	"os"
	"path/filepath"

	"gomail/pkg/logger"

	"github.com/ethereum/go-ethereum/common"
	eth_types "github.com/ethereum/go-ethereum/core/types"
)

func SaveTransactionToFile(tx *eth_types.Transaction, path string) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logger.Error("Failed to create directory", err)
		return err
	}

	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func LoadTransactionFromFile(path string) (*eth_types.Transaction, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tx := new(eth_types.Transaction)
	err = tx.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func LoadTransactionByHash(hash common.Hash, dataDir string) (*eth_types.Transaction, error) {
	filePath := filepath.Join(dataDir, fmt.Sprintf("%v.dat", hash))
	tx, err := LoadTransactionFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return tx, nil
}

func SaveTransactionByHash(tx *eth_types.Transaction, hash common.Hash, dataDir string) error {
	filePath := filepath.Join(dataDir, fmt.Sprintf("%v.dat", hash))
	// Lưu transaction sử dụng hàm SaveTransactionToFile
	err := SaveTransactionToFile(tx, filePath)
	if err != nil {
		return err
	}

	// Lưu trữ text hash vào file dataDir/map/tx.Hash
	mapFilePath := filepath.Join(dataDir, "map", fmt.Sprintf("%v.txt", tx.Hash()))
	err = os.MkdirAll(filepath.Dir(mapFilePath), os.ModePerm) // Tạo thư mục nếu chưa tồn tại
	if err != nil {
		return fmt.Errorf("failed to create directory for map file: %w", err)
	}
	err = os.WriteFile(mapFilePath, []byte(hash.Hex()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write hash to map file: %w", err)
	}

	return nil
}

func SaveTransactionMap(txHash common.Hash, hash common.Hash, dataDir string) error {

	// Lưu trữ text hash vào file dataDir/map/tx.Hash
	mapFilePath := filepath.Join(dataDir, "map", fmt.Sprintf("%v.txt", txHash))
	err := os.MkdirAll(filepath.Dir(mapFilePath), os.ModePerm) // Tạo thư mục nếu chưa tồn tại
	if err != nil {
		return fmt.Errorf("failed to create directory for map file: %w", err)
	}
	err = os.WriteFile(mapFilePath, []byte(hash.Hex()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write hash to map file: %w", err)
	}

	return nil
}

func GetHashFromFile(dataDir, mapDir string, hash common.Hash) (common.Hash, error) {
	// Xây dựng đường dẫn file
	mapFilePath := filepath.Join(dataDir, mapDir, fmt.Sprintf("%v.txt", hash))

	// Kiểm tra xem file có tồn tại hay không
	if _, err := os.Stat(mapFilePath); os.IsNotExist(err) {
		return common.Hash{}, fmt.Errorf("file not found: %w", err)
	}

	data, err := os.ReadFile(mapFilePath)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to read hash from file: %w", err)
	}

	hashResult := common.HexToHash(string(data))

	return hashResult, nil
}

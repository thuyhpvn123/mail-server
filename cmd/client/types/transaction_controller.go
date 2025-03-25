package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	pb "gomail/pkg/proto"
	"gomail/types"
)

type TransactionController interface {
	SendTransaction(
		lastHash common.Hash,
		fromAddress common.Address,
		toAddress common.Address,
		pendingUse *big.Int,
		amount *big.Int,
		maxGas uint64,
		maxGasFee uint64,
		maxTimeUse uint64,
		action pb.ACTION,
		data []byte,
		relatedAddress [][]byte,
		lastDeviceKey common.Hash,
		newDeviceKey common.Hash,
		commissionPrivateKey []byte,
		nonce uint64,
		chainId uint64,
	) (types.Transaction, error)

	ReadTransaction(
		lastHash common.Hash,
		fromAddress common.Address,
		toAddress common.Address,
		pendingUse *big.Int,
		amount *big.Int,
		maxGas uint64,
		maxGasFee uint64,
		maxTimeUse uint64,
		action pb.ACTION,
		data []byte,
		relatedAddress [][]byte,
		lastDeviceKey common.Hash,
		newDeviceKey common.Hash,
		commissionPrivateKey []byte,
		nonce uint64,
		chainId uint64,
	) (types.Transaction, error)

	SendTransactionWithDeviceKey(
		lastHash common.Hash,
		fromAddress common.Address,
		toAddress common.Address,
		pendingUse *big.Int,
		amount *big.Int,
		maxGas uint64,
		maxGasFee uint64,
		maxTimeUse uint64,
		action pb.ACTION,
		data []byte,
		relatedAddress [][]byte,
		lastDeviceKey common.Hash,
		newDeviceKey common.Hash,
		commissionPrivateKey []byte,
		nonce uint64,
		deviceKey []byte,
		chainId uint64,
	) (types.Transaction, error)

	SendTransactions(transactions []types.Transaction) error
}

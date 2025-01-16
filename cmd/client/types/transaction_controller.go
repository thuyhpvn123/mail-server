package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	pb "gomail/mtn/proto"
	"gomail/mtn/types"
)

type TransactionController interface {
	SendTransaction(
		lastHash common.Hash,
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
	) (types.Transaction, error)
	SendTransactionWithDeviceKey(
		lastHash common.Hash,
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
		deviceKey []byte,
	) (types.Transaction, error)
}

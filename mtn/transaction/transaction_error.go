package transaction

import (
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"

	pb "gomail/mtn/proto"
)

type TransactionError struct {
	Code        int64
	Description string
}

var (
	InvalidTransactionHash              = &TransactionError{1, "invalid transaction hash"}
	InvalidNewDeviceKey                 = &TransactionError{2, "invalid new device key"}
	NotMatchLastHash                    = &TransactionError{3, "not match last hash"}
	InvalidLastDeviceKey                = &TransactionError{4, "invalid last device key"}
	InvalidAmount                       = &TransactionError{5, "invalid amount"}
	InvalidPendingUse                   = &TransactionError{6, "invalid pending use"}
	InvalidDeploySmartContractToAccount = &TransactionError{
		7,
		"invalid deploy smart contract to account",
	}
	InvalidCallSmartContractToAccount = &TransactionError{
		8,
		"invalid call smart contract to account",
	}
	InvalidCallSmartContractData     = &TransactionError{9, "invalid call smart contract data"}
	InvalidStakeAddress              = &TransactionError{10, "invalid stake address"}
	InvalidUnstakeAddress            = &TransactionError{11, "invalid unstake address"}
	InvalidUnstakeAmount             = &TransactionError{12, "invalid unstake amount"}
	InvalidMaxGas                    = &TransactionError{13, "invalid max gas"}
	InvalidMaxGasPrice               = &TransactionError{14, "invalid max gas price"}
	InvalidCommissionSign            = &TransactionError{15, "invalid commission sign"}
	NotEnoughBalanceForCommissionFee = &TransactionError{
		16,
		"smart contract not enough balance for commission fee",
	}
	InvalidOpenChannelToAccount = &TransactionError{17, "invalid open channel to account"}
	InvalidSign                 = &TransactionError{18, "invalid sign"}
	InvalidCommitAddress        = &TransactionError{19, "invalid commit address"}
	InvalidOpenAccountAmount    = &TransactionError{20, "invalid open account amount"}
	IsNotCreator                = &TransactionError{21, "is not creator"}

	//
	NotEnoughVerifyMinerToVerifyTransation = &TransactionError{
		22,
		"not enough verify miner to verify transaction",
	}
	VerifyTransactionSignTimedOut = &TransactionError{
		23,
		"verify transaction sign timed out",
	}
	InvalidCodeStorage        = &TransactionError{24, "invalid code storage"}
	InvalidCallNonExitAccount = &TransactionError{25, "invalid call smart contract to non-existent account"}
)

var CodeToError = map[int64]*TransactionError{
	1:  InvalidTransactionHash,
	2:  InvalidNewDeviceKey,
	3:  NotMatchLastHash,
	4:  InvalidLastDeviceKey,
	5:  InvalidAmount,
	6:  InvalidPendingUse,
	7:  InvalidDeploySmartContractToAccount,
	8:  InvalidCallSmartContractToAccount,
	9:  InvalidCallSmartContractData,
	10: InvalidStakeAddress,
	11: InvalidUnstakeAddress,
	12: InvalidUnstakeAmount,
	13: InvalidMaxGas,
	14: InvalidMaxGasPrice,
	15: InvalidCommissionSign,
	16: NotEnoughBalanceForCommissionFee,
	17: InvalidOpenChannelToAccount,
	18: InvalidSign,
	19: InvalidCommitAddress,
	20: InvalidOpenAccountAmount,
	21: IsNotCreator,
	22: NotEnoughVerifyMinerToVerifyTransation,
	23: VerifyTransactionSignTimedOut,
	24: InvalidCodeStorage,
	25: InvalidCallNonExitAccount,
}

func (te *TransactionError) Proto() *pb.TransactionError {
	return &pb.TransactionError{
		Code:        te.Code,
		Description: te.Description,
	}
}

func (te *TransactionError) FromProto(pbData *pb.TransactionError) {
	te.Code = pbData.Code
	te.Description = pbData.Description
}

func (transactionErr *TransactionError) Marshal() ([]byte, error) {
	return proto.Marshal(transactionErr.Proto())
}

func (transactionErr *TransactionError) Unmarshal(data []byte) error {
	pbData := &pb.TransactionError{}
	if err := proto.Unmarshal(data, pbData); err != nil {
		return err
	}
	transactionErr.FromProto(pbData)
	return nil
}

type TransactionHashWithErrorCode struct {
	transactionHash common.Hash
	errorCode       int64
}

func NewTransactionHashWithErrorCode(
	transactionHash common.Hash,
	errorCode int64,
) *TransactionHashWithErrorCode {
	return &TransactionHashWithErrorCode{
		transactionHash: transactionHash,
		errorCode:       errorCode,
	}
}

func (te *TransactionHashWithErrorCode) Proto() *pb.TransactionHashWithErrorCode {
	return &pb.TransactionHashWithErrorCode{
		TransactionHash: te.transactionHash[:],
		Code:            te.errorCode,
	}
}

func (te *TransactionHashWithErrorCode) FromProto(
	pbData *pb.TransactionHashWithErrorCode,
) {
	te.transactionHash = common.BytesToHash(pbData.TransactionHash)
	te.errorCode = pbData.Code
}

func (transactionErr *TransactionHashWithErrorCode) Marshal() ([]byte, error) {
	return proto.Marshal(transactionErr.Proto())
}

func (transactionErr *TransactionHashWithErrorCode) Unmarshal(data []byte) error {
	pbData := &pb.TransactionHashWithErrorCode{}
	if err := proto.Unmarshal(data, pbData); err != nil {
		return err
	}

	transactionErr.FromProto(pbData)
	return nil
}

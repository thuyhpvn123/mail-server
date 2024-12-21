package stake

import (
	"errors"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	p_common "gomail/mtn/common"
	"gomail/mtn/logger"
	"gomail/mvm"
	pb "gomail/mtn/proto"
)

type StakeGetter struct {
	accountStatesDB mvm.AccountStateDB
	smartContractDB mvm.SmartContractDB
}

func NewStakeGetter(
	accountStatesDB mvm.AccountStateDB,
	smartContractDB mvm.SmartContractDB,
) *StakeGetter {
	return &StakeGetter{
		accountStatesDB: accountStatesDB,
		smartContractDB: smartContractDB,
	}
}

func (s *StakeGetter) checkMvm() (*mvm.MVMApi, error) {
	if s.accountStatesDB == nil || s.smartContractDB == nil {
		logger.Error("accountStatesDB or smartContractDB is nil")
		return nil, errors.New("accountStatesDB or smartContractDB is nil")
	}
	// reinit mvm api
	mvmApi := mvm.MVMApiInstance()
	if mvmApi == nil {
		mvm.InitMVMApi(s.smartContractDB, s.accountStatesDB)
		mvmApi = mvm.MVMApiInstance()
	}
	mvmApi.SetAccountStateDb(s.accountStatesDB)
	mvmApi.SetSmartContractDb(s.smartContractDB)
	mvm.MVMApiInstance().SetRelatedAddresses(
		[]common.Address{
			p_common.NATIVE_SMART_CONTRACT_STAKE_ADDRESS,
		},
	)
	return mvmApi, nil
}

func (s *StakeGetter) GetValidatorsWithStakeAmount() (map[common.Address]*big.Int, error) {
	mvmApi, err := s.checkMvm()
	if err != nil {
		return nil, err
	}
	// generate inputs
	getValidatorsWithStakeAmount, err := StakeABI().Pack(
		"getValidatorsWithStakeAmount",
	)
	if err != nil {
		logger.Error("Failed to pack getValidatorsWithStakeAmount input", err)
		return nil, err
	}

	// execute smart contract to distribute reward for leader, node and miners
	callRs := mvmApi.Call(
		// transaction data
		common.Address{}.Bytes(),
		p_common.NATIVE_SMART_CONTRACT_STAKE_ADDRESS.Bytes(),
		getValidatorsWithStakeAmount,
		big.NewInt(0),
		math.MaxUint64,
		math.MaxUint64,
		// block context data, skip theses fields
		0,
		0,
		0,
		0,
		0,
		common.Address{},
	)

	if callRs.Status != pb.RECEIPT_STATUS_RETURNED {
		logger.Error("Failed to call getValidatorsWithStakeAmount", callRs)
		return nil, errors.New("Failed to call getValidatorsWithStakeAmount")
	}

	// parse result
	mapRs := make(map[string]interface{})
	err = StakeABI().UnpackIntoMap(
		mapRs,
		"getValidatorsWithStakeAmount",
		callRs.Return,
	)
	if err != nil {
		logger.Error("error when unpack get child nodes return data")
		return nil, err
	}

	validatorsWithStakeAmount := make(map[common.Address]*big.Int)
	for i, v := range mapRs["addresses"].([]common.Address) {
		validatorsWithStakeAmount[v] = mapRs["amounts"].([]*big.Int)[i]
	}

	return validatorsWithStakeAmount, nil
}

func (s *StakeGetter) GetStakeInfo(
	nodeAddress common.Address,
) (*StakeInfo, error) {
	mvmApi, err := s.checkMvm()
	if err != nil {
		return nil, err
	}
	// generate inputs
	getStakeInfoInput, err := StakeABI().Pack(
		"getStakeInfo",
		nodeAddress,
	)
	if err != nil {
		logger.Error("Failed to pack getStakeInfo input", err)
		return nil, err
	}

	// execute smart contract to distribute reward for leader, node and miners
	callRs := mvmApi.Call(
		// transaction data
		common.Address{}.Bytes(),
		p_common.NATIVE_SMART_CONTRACT_STAKE_ADDRESS.Bytes(),
		getStakeInfoInput,
		big.NewInt(0),
		math.MaxUint64,
		math.MaxUint64,
		// block context data, skip theses fields
		0,
		0,
		0,
		0,
		0,
		common.Address{},
	)

	if callRs.Status != pb.RECEIPT_STATUS_RETURNED {
		logger.Error("Failed to call getValidatorsWithStakeAmount", callRs)
		return nil, errors.New("Failed to call getValidatorsWithStakeAmount")
	}

	// parse result
	mapRs := make(map[string]interface{})
	err = StakeABI().UnpackIntoMap(
		mapRs,
		"getStakeInfo",
		callRs.Return,
	)
	if err != nil {
		logger.Error("error when unpack get child nodes return data")
		return nil, err
	}
	stakeInfo := NewStakeInfo(
		mapRs["owner"].(common.Address),
		mapRs["amount"].(*big.Int),
		mapRs["childNodes"].([]common.Address),
		mapRs["childExecuteMiners"].([]common.Address),
		mapRs["childVerifyMiners"].([]common.Address),
	)

	return stakeInfo, nil
}

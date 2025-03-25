package stake

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type StakeInfo struct {
	owner              common.Address
	amount             *big.Int
	childNodes         []common.Address
	childExecuteMiners []common.Address
	childVerifyMiners  []common.Address
}

func NewStakeInfo(
	owner common.Address,
	amount *big.Int,
	childNodes []common.Address,
	childExecuteMiners []common.Address,
	childVerifyMiners []common.Address,
) *StakeInfo {
	return &StakeInfo{
		owner:              owner,
		amount:             amount,
		childNodes:         childNodes,
		childExecuteMiners: childExecuteMiners,
		childVerifyMiners:  childVerifyMiners,
	}
}

func (s *StakeInfo) Owner() common.Address {
	return s.owner
}

func (s *StakeInfo) Amount() *big.Int {
	return s.amount
}

func (s *StakeInfo) ChildNodes() []common.Address {
	return s.childNodes
}

func (s *StakeInfo) ChildExecuteMiners() []common.Address {
	return s.childExecuteMiners
}

func (s *StakeInfo) ChildVerifyMiners() []common.Address {
	return s.childVerifyMiners
}

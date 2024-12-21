package execute_sc_grouper

import (
	"github.com/ethereum/go-ethereum/common"

	"gomail/mtn/types"
)

type ExecuteSmartContractsGrouper struct {
	groupCount                  uint64
	mapAddressGroup             map[common.Address]uint64
	mapGroupExecuteTransactions map[uint64][]types.Transaction
}

func NewExecuteSmartContractsGrouper() *ExecuteSmartContractsGrouper {
	return &ExecuteSmartContractsGrouper{
		groupCount:                  0,
		mapAddressGroup:             make(map[common.Address]uint64),
		mapGroupExecuteTransactions: make(map[uint64][]types.Transaction),
	}
}

func (e *ExecuteSmartContractsGrouper) AddTransactions(
	transactions []types.Transaction,
) {
	// extract addresses
	for _, v := range transactions {
		addresses := v.RelatedAddresses()
		addresses = append(addresses, v.FromAddress())
		addresses = append(addresses, v.ToAddress())
		e.assignGroup(addresses)
	}
	// add transactions to group
	for _, v := range transactions {
		groupId := e.mapAddressGroup[v.ToAddress()]
		e.mapGroupExecuteTransactions[groupId] = append(e.mapGroupExecuteTransactions[groupId], v)
	}
}

func (e *ExecuteSmartContractsGrouper) GetGroupTransactions() map[uint64][]types.Transaction {
	return e.mapGroupExecuteTransactions
}

func (e *ExecuteSmartContractsGrouper) CountGroupWithTransactions() int {
	return len(e.mapGroupExecuteTransactions)
}

func (e *ExecuteSmartContractsGrouper) Clear() {
	e.groupCount = 0
	e.mapAddressGroup = make(map[common.Address]uint64)
	e.mapGroupExecuteTransactions = make(map[uint64][]types.Transaction)
}

func (e *ExecuteSmartContractsGrouper) assignGroup(addresses []common.Address) uint64 {
	groupId := uint64(0)
	for _, address := range addresses {
		rGroup := e.mapAddressGroup[address]
		// If the group is 0, skip it
		if rGroup == 0 {
			continue
		}
		if groupId == 0 {
			// If the group is 0, assign it to the address group
			groupId = rGroup
			continue
		}
		if rGroup < groupId {
			// If the group is less than the current group, assign it to the address group
			groupId = rGroup
		}
	}
	// If the group is 0, assign a new group
	if groupId == 0 {
		e.groupCount++
		groupId = e.groupCount
	}
	// assign all address to group id, remove old group id if exists
	for _, address := range addresses {
		e.mapAddressGroup[address] = groupId
	}

	return groupId
}

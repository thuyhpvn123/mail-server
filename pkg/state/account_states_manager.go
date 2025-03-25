package state

//
// import (
// 	"sync"
//
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/holiman/uint256"
//
// 	"gomail/pkg/logger"
// 	"gomail/pkg/storage"
// 	trie_package "gomail/pkg/trie"
// 	"gomail/types"
// )
//
// type AccountStatesManager struct {
// 	sync.RWMutex
// 	lastStateRoot common.Hash
//
// 	trie         *trie_package.MerklePatriciaTrie
// 	db           storage.Storage
// 	originStates map[common.Address]types.AccountState
//
// 	commitedTransactionState map[common.Address]types.AccountState
// 	liveStates               map[common.Address]types.AccountState
// 	dirtyStates              map[common.Address]struct{}
// }
//
// func NewAccountStatesManager(
// 	trie *trie_package.MerklePatriciaTrie,
// 	db storage.Storage,
// ) *AccountStatesManager {
// 	return &AccountStatesManager{
// 		trie:                     trie,
// 		db:                       db,
// 		originStates:             make(map[common.Address]types.AccountState),
// 		commitedTransactionState: make(map[common.Address]types.AccountState),
// 		liveStates:               make(map[common.Address]types.AccountState),
// 		dirtyStates:              make(map[common.Address]struct{}),
// 	}
// }
//
// // getter
// func (am *AccountStatesManager) AccountState(address common.Address) types.AccountState {
// 	am.Lock()
// 	defer am.Unlock()
// 	return am.accountState(address)
// }
//
// func (am *AccountStatesManager) OriginAccountState(address common.Address) types.AccountState {
// 	am.Lock()
// 	defer am.Unlock()
// 	if originState, ok := am.originStates[address]; ok {
// 		return originState.Copy()
// 	}
//
// 	bData, _ := am.trie.Get(address.Bytes())
// 	var accountState types.AccountState
// 	if len(bData) == 0 {
// 		accountState = NewAccountState(address)
// 	} else {
// 		accountState = &AccountState{}
// 		accountState.Unmarshal(bData)
// 	}
// 	am.originStates[address] = accountState
// 	return am.originStates[address]
// }
//
// func (am *AccountStatesManager) Exist(address common.Address) bool {
// 	am.RLock()
// 	defer am.RUnlock()
// 	if _, ok := am.liveStates[address]; ok {
// 		return true
// 	}
//
// 	if _, ok := am.originStates[address]; ok {
// 		return true
// 	}
//
// 	_, err := am.trie.Get(address.Bytes())
// 	return err == nil
// }
//
// func (am *AccountStatesManager) accountState(address common.Address) (as types.AccountState) {
// 	defer func() {
// 		logger.Trace("get account state", as)
// 	}()
// 	if liveState, ok := am.liveStates[address]; ok {
// 		as = liveState.Copy()
// 		return
// 	}
//
// 	if originState, ok := am.originStates[address]; ok {
// 		as = originState.Copy()
// 		return
// 	}
//
// 	bData, _ := am.trie.Get(address.Bytes())
// 	var accountState types.AccountState
// 	if len(bData) == 0 {
// 		accountState = NewAccountState(address)
// 	} else {
// 		accountState = &AccountState{}
// 		err := accountState.Unmarshal(bData)
// 		if err != nil {
// 			logger.Error("error when unmarshal account state", err)
// 		}
// 	}
//
// 	am.originStates[address] = accountState
// 	as = accountState.Copy()
// 	return
// }
//
// func (am *AccountStatesManager) setState(newState types.AccountState) {
// 	am.dirtyStates[newState.Address()] = struct{}{}
// 	am.liveStates[newState.Address()] = newState
// }
//
// func (am *AccountStatesManager) SetState(newState types.AccountState) {
// 	am.Lock()
// 	defer am.Unlock()
// 	am.setState(newState)
// 	logger.Trace("set new state", newState)
// }
//
// func (am *AccountStatesManager) GetStorageSnapShot() storage.SnapShot {
// 	am.RLock()
// 	defer am.RUnlock()
// 	return am.db.GetSnapShot()
// }
//
// func (am *AccountStatesManager) Storage() storage.Storage {
// 	return am.db
// }
//
// func (am *AccountStatesManager) Copy() *AccountStatesManager {
// 	am.RLock()
// 	defer am.RUnlock()
// 	cp := &AccountStatesManager{
// 		trie:         am.trie.Copy(),
// 		originStates: make(map[common.Address]types.AccountState, len(am.originStates)),
// 		liveStates:   make(map[common.Address]types.AccountState, len(am.liveStates)),
// 		dirtyStates:  make(map[common.Address]struct{}, len(am.dirtyStates)),
// 	}
// 	//
// 	for i, v := range am.originStates {
// 		cp.originStates[i] = v
// 	}
//
// 	for i, v := range am.liveStates {
// 		cp.liveStates[i] = v
// 	}
//
// 	for i, v := range am.dirtyStates {
// 		cp.dirtyStates[i] = v
// 	}
// 	return cp
// }
//
// // setter
// func (am *AccountStatesManager) SetSmartContractState(
// 	address common.Address,
// 	smState types.SmartContractState,
// ) {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	as.SetSmartContractState(smState)
// 	am.setState(as)
// }
//
// func (am *AccountStatesManager) SetNewDeviceKey(address common.Address, newDeviceKey common.Hash) {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	as.SetNewDeviceKey(newDeviceKey)
// 	am.setState(as)
// }
//
// func (am *AccountStatesManager) SetLastHash(address common.Address, newLastHash common.Hash) {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	as.SetLastHash(newLastHash)
// 	am.setState(as)
// }
//
// func (am *AccountStatesManager) AddPendingBalance(address common.Address, amount *uint256.Int) {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	as.AddPendingBalance(amount)
// 	am.setState(as)
// }
//
// func (am *AccountStatesManager) SubPendingBalance(
// 	address common.Address,
// 	amount *uint256.Int,
// ) error {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	err := as.SubPendingBalance(amount)
// 	if err != nil {
// 		return err
// 	}
// 	am.setState(as)
// 	return nil
// }
//
// func (am *AccountStatesManager) SubBalance(address common.Address, amount *uint256.Int) error {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	err := as.SubBalance(amount)
// 	if err != nil {
// 		return err
// 	}
// 	am.setState(as)
// 	return nil
// }
//
// func (am *AccountStatesManager) AddBalance(address common.Address, amount *uint256.Int) {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	as.AddBalance(amount)
// 	am.setState(as)
// }
//
// func (am *AccountStatesManager) SubTotalBalance(address common.Address, amount *uint256.Int) error {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	err := as.SubTotalBalance(amount)
// 	if err != nil {
// 		logger.Error("error when sub total balance", err)
// 		return err
// 	}
// 	am.setState(as)
// 	return nil
// }
//
// func (am *AccountStatesManager) SetCodeHash(address common.Address, hash common.Hash) {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	as.SetCodeHash(hash)
// 	am.setState(as)
// }
//
// func (am *AccountStatesManager) SetStorageHost(address common.Address, storageHost string) {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	as.SetStorageHost(storageHost)
// 	am.setState(as)
// }
//
// func (am *AccountStatesManager) SetStorageAddress(
// 	address common.Address,
// 	storageAddress common.Address,
// ) {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	as.SetStorageAddress(storageAddress)
// 	am.setState(as)
// }
//
// func (am *AccountStatesManager) SetStorageRoot(address common.Address, hash common.Hash) {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	as.SetStorageRoot(hash)
// 	am.setState(as)
// }
//
// func (am *AccountStatesManager) SetStateChannelState(
// 	address common.Address,
// 	scState types.StateChannelState,
// ) {
// 	am.Lock()
// 	defer am.Unlock()
// 	as := am.accountState(address)
// 	as.SetStateChannelState(scState)
// 	am.setState(as)
// }
//
// func (am *AccountStatesManager) Hash() common.Hash {
// 	return am.trie.Hash()
// }
//
// func (am *AccountStatesManager) AccountStateChanges() map[common.Address]types.AccountState {
// 	return am.commitedTransactionState
// }
//
// func (am *AccountStatesManager) Commit() (common.Hash, error) {
// 	am.Lock()
// 	defer am.Unlock()
//
// 	hash, nodeSet, oldKeys, err := am.trie.Commit(true)
// 	if err != nil {
// 		return common.Hash{}, err
// 	}
// 	// delete before insert to avoid conflic
// 	for i := 0; i < len(oldKeys); i++ {
// 		am.db.Delete(oldKeys[i])
// 	}
// 	if nodeSet != nil {
// 		// save nodeSet to db
// 		batch := [][2][]byte{}
// 		for _, node := range nodeSet.Nodes {
// 			batch = append(batch, [2][]byte{node.Hash.Bytes(), node.Blob})
// 		}
// 		err := am.db.BatchPut(batch)
// 		if err != nil {
// 			return common.Hash{}, err
// 		}
// 	}
// 	am.originStates = make(map[common.Address]types.AccountState)
// 	am.commitedTransactionState = make(map[common.Address]types.AccountState)
// 	am.liveStates = make(map[common.Address]types.AccountState)
// 	am.dirtyStates = make(map[common.Address]struct{})
//
// 	am.trie, err = trie_package.New(hash, am.db)
// 	if err != nil {
// 		return common.Hash{}, err
// 	}
// 	am.lastStateRoot = hash
// 	return hash, err
// }
//
// func (am *AccountStatesManager) CommitTransaction() (common.Hash, error) {
// 	am.Lock()
// 	defer am.Unlock()
// 	for addr := range am.dirtyStates {
// 		logger.Trace("committing transaction for address", addr.String())
// 		state := am.liveStates[addr]
// 		logger.Trace("state", state)
// 		am.commitedTransactionState[addr] = state.Copy()
// 		bData, err := state.Marshal()
// 		if err != nil {
// 			logger.Error("error when marshal account state", err)
// 			return common.Hash{}, err
// 		}
// 		err = am.trie.Update(state.Address().Bytes(), bData)
// 		logger.Trace("commiting address", addr, "state", state, "err", err)
// 		if err != nil {
// 			logger.Error("error when update account state", err)
// 			return common.Hash{}, err
// 		}
// 	}
// 	// clear dirty states
// 	am.dirtyStates = make(map[common.Address]struct{})
// 	// calculate hash
// 	return am.trie.Hash(), nil
// }
//
// func (am *AccountStatesManager) RollbackToLastTransaction() error { // rollback to last transaction commit
// 	am.Lock()
// 	defer am.Unlock()
// 	am.originStates = map[common.Address]types.AccountState{}
// 	am.liveStates = map[common.Address]types.AccountState{}
// 	am.dirtyStates = map[common.Address]struct{}{}
// 	return nil
// }
//
// func (am *AccountStatesManager) RollbackToLastBlock() (err error) { // rollback to last block state
// 	am.Lock()
// 	defer am.Unlock()
// 	am.trie, err = trie_package.New(am.lastStateRoot, am.db)
// 	am.originStates = map[common.Address]types.AccountState{}
// 	am.commitedTransactionState = map[common.Address]types.AccountState{}
// 	am.liveStates = map[common.Address]types.AccountState{}
// 	am.dirtyStates = map[common.Address]struct{}{}
// 	return err
// }
//
// func (am *AccountStatesManager) CopyFrom(origin types.AccountStatesManager) {
// 	am.RLock()
// 	defer am.RUnlock()
// 	originAccountStatesManager := origin.(*AccountStatesManager)
// 	am.trie = originAccountStatesManager.trie.Copy()
// 	am.db = originAccountStatesManager.db
// 	am.lastStateRoot = originAccountStatesManager.lastStateRoot
// 	am.originStates = make(
// 		map[common.Address]types.AccountState,
// 		len(originAccountStatesManager.originStates),
// 	)
// 	am.liveStates = make(
// 		map[common.Address]types.AccountState,
// 		len(originAccountStatesManager.liveStates),
// 	)
// 	am.dirtyStates = make(map[common.Address]struct{}, len(originAccountStatesManager.dirtyStates))
// 	for i, v := range originAccountStatesManager.originStates {
// 		am.originStates[i] = v
// 	}
// 	for i, v := range originAccountStatesManager.liveStates {
// 		am.liveStates[i] = v
// 	}
// 	for i, v := range originAccountStatesManager.dirtyStates {
// 		am.dirtyStates[i] = v
// 	}
// }

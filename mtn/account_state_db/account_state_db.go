package account_state_db

import (
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	p_common "gomail/mtn/common"
	"gomail/mtn/logger"
	"gomail/mtn/state"
	"gomail/mtn/storage"
	p_trie "gomail/mtn/trie"
	"gomail/mtn/types"
)

type AccountStateDB struct {
	trie           *p_trie.MerklePatriciaTrie
	originRootHash common.Hash
	db             storage.Storage
	dirtyAccounts  map[common.Address]types.AccountState
	mu             sync.Mutex
}

func NewAccountStateDB(
	trie *p_trie.MerklePatriciaTrie,
	db storage.Storage,
) *AccountStateDB {
	return &AccountStateDB{
		trie:           trie,
		db:             db,
		originRootHash: trie.Hash(),
		dirtyAccounts:  make(map[common.Address]types.AccountState),
	}
}

func (db *AccountStateDB) AccountState(address common.Address) (types.AccountState, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.getOrCreateAccountState(address)
}

func (db *AccountStateDB) SubPendingBalance(address common.Address, amount *big.Int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	as, err := db.getOrCreateAccountState(address)
	if err != nil {
		return err
	}
	err = as.SubPendingBalance(amount)
	if err != nil {
		return err
	}
	db.setDirtyAccountState(as)
	return nil
}

func (db *AccountStateDB) AddPendingBalance(address common.Address, amount *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	as, _ := db.getOrCreateAccountState(address)
	logger.Info("Adding pending balance ", as, amount)
	as.AddPendingBalance(amount)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) AddBalance(address common.Address, amount *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	as, _ := db.getOrCreateAccountState(address)
	as.AddBalance(amount)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) SubBalance(address common.Address, amount *big.Int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	as, _ := db.getOrCreateAccountState(address)
	err := as.SubBalance(amount)
	if err != nil {
		return err
	}
	db.setDirtyAccountState(as)
	return nil
}

func (db *AccountStateDB) SubTotalBalance(address common.Address, amount *big.Int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	as, _ := db.getOrCreateAccountState(address)
	err := as.SubTotalBalance(amount)
	if err != nil {
		return err
	}
	db.setDirtyAccountState(as)
	return nil
}

func (db *AccountStateDB) SetLastHash(address common.Address, hash common.Hash) {
	db.mu.Lock()
	defer db.mu.Unlock()

	as, _ := db.getOrCreateAccountState(address)
	as.SetLastHash(hash)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) SetNewDeviceKey(address common.Address, newDeviceKey common.Hash) {
	db.mu.Lock()
	defer db.mu.Unlock()

	as, _ := db.getOrCreateAccountState(address)
	as.SetNewDeviceKey(newDeviceKey)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) SetState(as types.AccountState) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.setDirtyAccountState(as)
}

// smart contract state
func (db *AccountStateDB) SetCreatorPublicKey(
	address common.Address,
	creatorPublicKey p_common.PublicKey,
) {
	db.mu.Lock()
	defer db.mu.Unlock()

	as, _ := db.getOrCreateAccountState(address)
	as.SetCreatorPublicKey(creatorPublicKey)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) SetCodeHash(address common.Address, codeHash common.Hash) {
	db.mu.Lock()
	defer db.mu.Unlock()

	as, _ := db.getOrCreateAccountState(address)
	as.SetCodeHash(codeHash)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) SetStorageRoot(address common.Address, storageRoot common.Hash) {
	db.mu.Lock()
	defer db.mu.Unlock()

	as, _ := db.getOrCreateAccountState(address)
	as.SetStorageRoot(storageRoot)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) SetStorageAddress(
	address common.Address,
	storageAddress common.Address,
) {
	db.mu.Lock()
	defer db.mu.Unlock()

	as, _ := db.getOrCreateAccountState(address)
	as.SetStorageAddress(storageAddress)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) AddLogHash(address common.Address, logsHash common.Hash) {
	db.mu.Lock()
	defer db.mu.Unlock()

	as, _ := db.getOrCreateAccountState(address)
	as.AddLogHash(logsHash)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) Discard() (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	logger.Info("Discard")
	// Discard dirty accounts
	db.dirtyAccounts = make(map[common.Address]types.AccountState)
	// Load new trie from db
	db.trie, err = p_trie.New(db.originRootHash, db.db)
	return err
}

func (db *AccountStateDB) Commit() (common.Hash, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	logger.Info("Commit Account State DB")
	// call intermidiate root to update dirty accounts to trie
	rootHash, err := db.IntermediateRoot()
	if err != nil {
		return common.Hash{}, err
	}
	// commit trie
	hash, nodeSet, oldKeys, err := db.trie.Commit(true)
	if err != nil {
		return common.Hash{}, err
	}
	if rootHash != hash {
		return common.Hash{}, errors.New("In")
	}
	for i := 0; i < len(oldKeys); i++ {
		// should save oldKeys to archive db for backup
		db.db.Delete(oldKeys[i])
	}
	if nodeSet != nil {
		// save nodeSet to db
		batch := [][2][]byte{}
		for _, node := range nodeSet.Nodes {
			batch = append(batch, [2][]byte{node.Hash.Bytes(), node.Blob})
		}
		err := db.db.BatchPut(batch)
		if err != nil {
			return common.Hash{}, err
		}
	}
	db.dirtyAccounts = make(map[common.Address]types.AccountState)

	db.trie, err = p_trie.New(hash, db.db)
	if err != nil {
		return common.Hash{}, err
	}
	db.originRootHash = hash
	return hash, err
}

func (db *AccountStateDB) IntermediateRoot() (common.Hash, error) {
	// update dirty accounts to trie
	for address, as := range db.dirtyAccounts {
		b, err := as.Marshal()
		if err != nil {
			return common.Hash{}, err
		}
		err = db.trie.Update(address.Bytes(), b)
		if err != nil {
			return common.Hash{}, err
		}
	}
	// get root hash
	rootHash := db.trie.Hash()

	return rootHash, nil
}

func (db *AccountStateDB) setDirtyAccountState(as types.AccountState) {
	db.dirtyAccounts[as.Address()] = as
}

func (db *AccountStateDB) getOrCreateAccountState(
	address common.Address,
) (types.AccountState, error) {
	// get from dirty
	as := db.dirtyAccounts[address]
	if as != nil {
		return as, nil
	}
	// if not exist in dirty then get from trie
	bData, _ := db.trie.Get(address.Bytes())
	if len(bData) == 0 {
		// if not exist in trie create new once
		return state.NewAccountState(address), nil
	}
	// exist in trie, unmarshal
	as = &state.AccountState{}
	err := as.Unmarshal(bData)
	if err != nil {
		return nil, err
	}
	return as, nil
}

func (db *AccountStateDB) Storage() storage.Storage {
	return db.db
}

func (db *AccountStateDB) CopyFrom(as types.AccountStateDB) error {
	asDB := as.(*AccountStateDB)
	db.trie = asDB.trie
	db.originRootHash = asDB.originRootHash
	db.db = asDB.db
	db.dirtyAccounts = asDB.dirtyAccounts
	return nil
}

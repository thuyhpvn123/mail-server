package account_state_db

import (
	"bytes"
	"errors"
	"math/big"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	p_common "gomail/pkg/common"
	"gomail/pkg/logger"
	pb "gomail/pkg/proto"
	"gomail/pkg/state"
	"gomail/pkg/storage"
	p_trie "gomail/pkg/trie"
	"gomail/types"
)

type AccountStateDB struct {
	trie *p_trie.MerklePatriciaTrie

	originRootHash common.Hash
	db             storage.Storage
	dirtyAccounts  sync.Map
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
		dirtyAccounts:  sync.Map{},
	}
}

func (db *AccountStateDB) AccountState(address common.Address) (types.AccountState, error) {
	return db.getOrCreateAccountState(address)
}

func (db *AccountStateDB) SubPendingBalance(address common.Address, amount *big.Int) error {
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

func (db *AccountStateDB) RefreshPendingBalance(address common.Address) error {
	as, err := db.getOrCreateAccountState(address)
	if err != nil {
		return err
	}
	pendingBalance := as.PendingBalance()
	err = as.SubPendingBalance(pendingBalance)
	if err != nil {
		return err
	}
	as.AddBalance(pendingBalance)
	db.setDirtyAccountState(as)
	return nil
}

func (db *AccountStateDB) AddPendingBalance(address common.Address, amount *big.Int) error {
	as, err := db.getOrCreateAccountState(address)
	if err != nil {
		return err
	}
	as.AddPendingBalance(amount)
	db.setDirtyAccountState(as)
	return nil
}

func (db *AccountStateDB) PlusOneNonce(address common.Address) error {
	as, err := db.getOrCreateAccountState(address)
	if err != nil {
		return err
	}
	as.PlusOneNonce()
	db.setDirtyAccountState(as)
	return nil
}

func (db *AccountStateDB) SetAccountType(address common.Address, accountTypeNew pb.ACCOUNT_TYPE) error {
	as, err := db.getOrCreateAccountState(address)
	if err != nil {
		return err
	}
	err = as.SetAccountType(accountTypeNew)
	if err != nil {
		return err
	}
	db.setDirtyAccountState(as)
	return nil
}

func (db *AccountStateDB) SetPublicKeyBls(address common.Address, publicKeyBls []byte) error {
	as, err := db.getOrCreateAccountState(address)
	if err != nil {
		return err
	}
	as.SetPublicKeyBls(publicKeyBls)
	db.setDirtyAccountState(as)
	return nil
}

func (db *AccountStateDB) AddBalance(address common.Address, amount *big.Int) error {
	as, err := db.getOrCreateAccountState(address)
	if err != nil {
		return err
	}
	as.AddBalance(amount)
	db.setDirtyAccountState(as)
	return nil
}

func (db *AccountStateDB) SubBalance(address common.Address, amount *big.Int) error {
	as, _ := db.getOrCreateAccountState(address)
	err := as.SubBalance(amount)
	if err != nil {
		return err
	}
	db.setDirtyAccountState(as)
	return nil
}

func (db *AccountStateDB) SubTotalBalance(address common.Address, amount *big.Int) error {
	as, _ := db.getOrCreateAccountState(address)
	err := as.SubTotalBalance(amount)
	if err != nil {
		return err
	}
	db.setDirtyAccountState(as)
	return nil
}

func (db *AccountStateDB) SetLastHash(address common.Address, hash common.Hash) {
	as, _ := db.getOrCreateAccountState(address)
	as.SetLastHash(hash)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) SetNewDeviceKey(address common.Address, newDeviceKey common.Hash) {
	as, _ := db.getOrCreateAccountState(address)
	as.SetNewDeviceKey(newDeviceKey)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) SetState(as types.AccountState) {
	db.setDirtyAccountState(as)
}

// smart contract state
func (db *AccountStateDB) SetCreatorPublicKey(
	address common.Address,
	creatorPublicKey p_common.PublicKey,
) {
	as, _ := db.getOrCreateAccountState(address)
	as.SetCreatorPublicKey(creatorPublicKey)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) SetCodeHash(address common.Address, codeHash common.Hash) {
	as, _ := db.getOrCreateAccountState(address)
	as.SetCodeHash(codeHash)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) SetStorageRoot(address common.Address, storageRoot common.Hash) {
	as, _ := db.getOrCreateAccountState(address)
	as.SetStorageRoot(storageRoot)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) SetStorageAddress(
	address common.Address,
	storageAddress common.Address,
) {
	as, _ := db.getOrCreateAccountState(address)
	as.SetStorageAddress(storageAddress)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) AddLogHash(address common.Address, logsHash common.Hash) {
	as, _ := db.getOrCreateAccountState(address)
	as.AddLogHash(logsHash)
	db.setDirtyAccountState(as)
}

func (db *AccountStateDB) Discard() (err error) {
	db.dirtyAccounts = sync.Map{}
	db.trie, err = p_trie.New(db.originRootHash, db.db, true)
	return err
}

func (db *AccountStateDB) Commit() (common.Hash, error) {

	db.mu.Lock()
	defer db.mu.Unlock()
	// call intermidiate root to update dirty accounts to trie
	rootHash, err := db.IntermediateRoot()
	if err != nil {
		return common.Hash{}, err
	}
	trieCommit := db.trie.Copy()
	// commit trie
	hash, nodeSet, oldKeys, err := trieCommit.Commit(true)
	if err != nil {
		return common.Hash{}, err
	}
	if rootHash != hash {
		return common.Hash{}, errors.New("In")
	}
	for i := 0; i < len(oldKeys); i++ {
		// should save oldKeys to archive db for backup
		// db.db.Delete(oldKeys[i])
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

	db.trie, err = p_trie.New(hash, db.db, true)
	if err != nil {
		return common.Hash{}, err
	}
	db.originRootHash = hash

	return hash, err
}

func (db *AccountStateDB) IntermediateRoot() (common.Hash, error) {
	var sortedAddresses []common.Address
	db.dirtyAccounts.Range(func(key, value interface{}) bool {
		address := key.(common.Address)

		sortedAddresses = append(sortedAddresses, address)
		return true
	})
	sort.Slice(sortedAddresses, func(i, j int) bool {
		return bytes.Compare(sortedAddresses[i].Bytes(), sortedAddresses[j].Bytes()) < 0
	})

	for _, address := range sortedAddresses {
		value, _ := db.dirtyAccounts.Load(address)
		as := value.(types.AccountState)
		b, err := as.Marshal()
		if err != nil {
			logger.Error("Error marshaling account state:", err)
			return common.Hash{}, err
		}
		if err := db.trie.Update(address.Bytes(), b); err != nil {
			logger.Error("Error updating trie:", err)
			return common.Hash{}, err
		}
	}

	rootHash := db.trie.Hash()
	db.dirtyAccounts = sync.Map{}

	return rootHash, nil
}

func (db *AccountStateDB) setDirtyAccountState(as types.AccountState) {
	db.dirtyAccounts.Store(as.Address(), as)
}

func (db *AccountStateDB) PublicSetDirtyAccountState(as types.AccountState) {
	db.dirtyAccounts.Store(as.Address(), as)
}

func (db *AccountStateDB) getOrCreateAccountState(
	address common.Address,
) (types.AccountState, error) {

	value, ok := db.dirtyAccounts.Load(address)
	if ok {
		return value.(types.AccountState), nil
	}
	bData, err := db.trie.Get(address.Bytes())

	if err != nil {
		logger.Error("getOrCreateAccountState 1", err)
		return nil, err

		// bData, err = db.trieRead.Get(address.Bytes())

		// if err != nil {
		// 	logger.Error("getOrCreateAccountState 1.1", err)

		// 	return nil, err
		// }
	}
	if len(bData) == 0 {
		return state.NewAccountState(address), nil
	}
	as := &state.AccountState{}
	err = as.Unmarshal(bData)
	if err != nil {
		logger.Error("getOrCreateAccountState 2", err)

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
	asDB.dirtyAccounts.Range(func(key, value interface{}) bool {
		db.dirtyAccounts.Store(key, value)
		return true
	})
	return nil
}

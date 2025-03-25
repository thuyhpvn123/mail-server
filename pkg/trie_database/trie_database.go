package trie_database

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/trie/trienode"
	"github.com/ethereum/go-ethereum/triedb"

	"gomail/pkg/account_state_db"
	"gomail/pkg/logger"
	"gomail/pkg/storage"
)

type TrieDatabase struct {
	trieR          *trie.Trie
	trieDB         *triedb.Database
	originRootHash common.Hash
	db             *storage.ShardelDB
	dirtyAccounts  sync.Map
	mu             sync.Mutex
	address        common.Address
	dbName         string
	accountStateDB *account_state_db.AccountStateDB
}

func NewTrieDatabase(
	hash common.Hash,
	db *storage.ShardelDB,
	address common.Address,
	dbName string,
	accountStateDB *account_state_db.AccountStateDB,
) *TrieDatabase {

	trieDB := triedb.NewDatabase(rawdb.NewDatabase(db), &triedb.Config{})
	// Tạo một đối tượng Trie mới
	var trieR *trie.Trie
	var err error

	if (hash == common.Hash{}) {
		trieR, err = trie.New(trie.TrieID(common.Hash{}), trieDB)

	} else {
		// Thử tối đa 3 lần với độ trễ giữa các lần thử
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			trieR, err = trie.New(trie.TrieID(hash), trieDB)
			if err == nil {
				break
			}
			// Nếu không phải lần thử cuối cùng, đợi trước khi thử lại
			if i < maxRetries-1 {
				time.Sleep(100 * time.Millisecond)
			}
		}

	}
	if err != nil {

		logger.Info("Error creating trie: %v", err)
		return nil
	}

	return &TrieDatabase{
		trieR:          trieR,
		trieDB:         trieDB,
		db:             db,
		originRootHash: trieR.Hash(),
		dirtyAccounts:  sync.Map{},
		address:        address,
		dbName:         dbName,
		accountStateDB: accountStateDB,
	}
}

func (trieDatabae *TrieDatabase) Commit() (common.Hash, error) {
	trieDatabae.IntermediateRoot()
	trieCopy := trieDatabae.trieR.Copy()
	root, nodes := trieCopy.Commit(false)

	if nodes == nil {
		return root, nil
	}

	nodeSet := trienode.NewWithNodeSet(nodes)

	if err := trieDatabae.trieDB.Update(root, trieDatabae.originRootHash, 0, nodeSet, nil); err != nil {
		log.Fatalf("lỗi khi cập nhật trie database: %v", err)
	}

	if err := trieDatabae.trieDB.Commit(root, false); err != nil {
		log.Fatalf("lỗi khi commit trie database: %v", err)
	}

	// Create a new trie based on the new root hash
	newTrie, err := trie.New(trie.TrieID(root), trieDatabae.trieDB)
	if err != nil {
		logger.Error("Error creating new trie after commit: %v", err)
		return common.Hash{}, err
	}

	trieDatabae.trieR = newTrie
	trieDatabae.originRootHash = root

	return root, nil
}

func (trieDatabae *TrieDatabase) RestoreTrieFromRootHash(rootHash common.Hash) (*trie.Trie, error) {
	// Thử tối đa 3 lần với độ trễ giữa các lần thử
	maxRetries := 3
	var err error
	var tr *trie.Trie
	for i := 0; i < maxRetries; i++ {
		tr, err = trie.New(trie.TrieID(rootHash), trieDatabae.trieDB)
		if err == nil {
			return tr, nil
		}

		// Nếu không phải lần thử cuối cùng, đợi trước khi thử lại
		if i < maxRetries-1 {
			time.Sleep(100 * time.Millisecond)
		}
		logger.Error("Error creating trie after restore, retrying: %v", err)
	}

	// Nếu đến đây, tất cả các lần thử đều thất bại
	logger.Error("Error creating trie after multiple retries")
	return nil, err
}

func (trieDatabae *TrieDatabase) IntermediateRoot() (common.Hash, error) {
	var sortedKeys []string // Thay đổi kiểu thành string
	trieDatabae.dirtyAccounts.Range(func(key, value interface{}) bool {
		address := key.(string) // Thay đổi kiểu thành string
		sortedKeys = append(sortedKeys, address)
		return true
	})
	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i] < sortedKeys[j] // So sánh chuỗi trực tiếp
	})

	for _, key := range sortedKeys {
		value, _ := trieDatabae.dirtyAccounts.Load(key)
		valStr := value.(string)
		if err := trieDatabae.trieR.Update([]byte(key), []byte(valStr)); err != nil { // Chuyển đổi cả key và value thành []byte
			return common.Hash{}, err
		}
	}

	rootHash := trieDatabae.trieR.Hash()
	trieDatabae.dirtyAccounts = sync.Map{}

	return rootHash, nil
}

func (trieDatabae *TrieDatabase) Storage() storage.Storage {
	return trieDatabae.db
}

func (trieDatabae *TrieDatabase) setDirty(key string, value string) {
	trieDatabae.dirtyAccounts.Store(key, value)
}

func (trieDatabae *TrieDatabase) Get(
	key string,
) (string, error) {

	value, ok := trieDatabae.dirtyAccounts.Load(key)
	if ok {
		return value.(string), nil
	}
	bData, err := trieDatabae.trieR.Get([]byte(key))
	if err != nil {
		logger.Error("TrieDatabase Get", err)
		return "", err
	}
	return string(bData), nil
}

func (trieDatabae *TrieDatabase) Put(
	key string,
	value string,
) error {
	trieDatabae.setDirty(key, value)
	return nil
}

// GetAllKeyValues retrieves all key-value pairs from both dirtyAccounts and the trie.
// It returns a map[string]string containing all the data.  If a key exists in both
// dirtyAccounts and the trie, the value from dirtyAccounts takes precedence.
func (trieDatabae *TrieDatabase) GetAllKeyValues() (map[string]string, error) {
	allKeyValues := make(map[string]string)

	// Iterate over dirtyAccounts and add/update key-value pairs in the map.
	trieDatabae.dirtyAccounts.Range(func(key, value interface{}) bool {
		allKeyValues[key.(string)] = value.(string)
		return true
	})

	// Get key-value pairs from the trie using NodeIterator, only for keys not in dirtyAccounts
	iter, err := trieDatabae.trieR.NodeIterator(nil)
	if err != nil {
		return nil, err
	}
	it := trie.NewIterator(iter)
	for it.Next() {
		key := string(it.Key)
		if _, ok := allKeyValues[key]; !ok {
			allKeyValues[key] = string(it.Value)
		}
	}

	return allKeyValues, nil
}

// Discard abandons all changes made since the last Commit.
func (trieDatabae *TrieDatabase) Discard() error {
	trieDatabae.mu.Lock()
	defer trieDatabae.mu.Unlock()

	newTrie, err := trieDatabae.RestoreTrieFromRootHash(trieDatabae.originRootHash)
	if err != nil {
		return err
	}

	trieDatabae.trieR = newTrie
	trieDatabae.dirtyAccounts = sync.Map{}
	return nil
}

// ... existing code ...

// SearchKeyValuesByValue searches for key-value pairs with the given value.
// It returns a map[string]string containing all matching key-value pairs.
func (trieDatabae *TrieDatabase) SearchByValue(searchValue string) (map[string]string, error) {
	matchingKeyValues := make(map[string]string)
	// Search in dirtyAccounts
	trieDatabae.dirtyAccounts.Range(func(key, value interface{}) bool {
		if value.(string) == searchValue {
			matchingKeyValues[key.(string)] = value.(string)
		}
		return true
	})

	// Search in trieR
	iter, err := trieDatabae.trieR.NodeIterator(nil)
	if err != nil {
		return nil, err
	}
	it := trie.NewIterator(iter)
	for it.Next() {
		if string(it.Value) == searchValue {
			key := string(it.Key)
			if _, ok := matchingKeyValues[key]; !ok {
				matchingKeyValues[key] = string(it.Value)
			}
		}
	}

	return matchingKeyValues, nil
}

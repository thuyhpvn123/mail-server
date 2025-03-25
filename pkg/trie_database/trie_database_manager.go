package trie_database

import (
	"sync"

	"gomail/pkg/account_state_db"
	"gomail/pkg/logger"
	"gomail/pkg/storage"

	"github.com/ethereum/go-ethereum/common"
)

// TrieDatabaseManager quản lý nhiều TrieDatabase
type TrieDatabaseManager struct {
	trieDatabases  map[common.Hash]*TrieDatabase
	database       *storage.ShardelDB
	accountStateDB *account_state_db.AccountStateDB
}

var (
	instance *TrieDatabaseManager
	once     sync.Once
)

func CreateTrieDatabaseManager(db *storage.ShardelDB, accountStateDB *account_state_db.AccountStateDB) *TrieDatabaseManager {
	logger.Info("CreateTrieDatabaseManager")
	once.Do(func() {
		instance = &TrieDatabaseManager{
			trieDatabases:  make(map[common.Hash]*TrieDatabase),
			database:       db,
			accountStateDB: accountStateDB,
		}
	})
	return instance
}
func GetTrieDatabaseManager() *TrieDatabaseManager {
	logger.Error("GetTrieDatabaseManager")
	return instance
}

// CommitAllTrieDatabases duyệt qua tất cả các TrieDatabase và commit chúng.
func (manager *TrieDatabaseManager) CommitAllTrieDatabases() error {
	for id, trieDB := range manager.trieDatabases {
		root, err := trieDB.Commit()
		if err != nil {
			logger.Error("Failed to commit TrieDatabase", "id", id, "error", err)
			return err // Trả về lỗi nếu bất kỳ commit nào không thành công
		}
		as, err := manager.accountStateDB.AccountState(trieDB.address)
		if err != nil {
			logger.Error("Failed to commit TrieDatabase get AccountState", "id", id, "error", err)
			return err // Trả về lỗi nếu bất kỳ commit nào không thành công
		}
		as.SmartContractState().SetTrieDatabaseMapValue(trieDB.dbName, root.Bytes())
		manager.accountStateDB.PublicSetDirtyAccountState(as)
		logger.Info("Committed TrieDatabase", "id", id, "root", root)
	}

	return nil
}

// DiscardAllTrieDatabases loại bỏ tất cả các thay đổi đang chờ xử lý trong tất cả các TrieDatabase.
func (manager *TrieDatabaseManager) DiscardAllTrieDatabases() {
	for id, trieDB := range manager.trieDatabases {
		trieDB.Discard()
		logger.Info("Discarded TrieDatabase", "id", id)
	}
}

// GetTrieDatabase lấy một TrieDatabase theo ID của nó.
// func (manager *TrieDatabaseManager) GetTrieDatabase(id common.Hash, hash common.Hash) (*TrieDatabase, bool) {
// 	trieDB, exists := manager.trieDatabases[id]
// 	if !exists {
// 		return nil, false // trả về true nếu nó đã tồn tại, false nếu nó vừa được tạo

// 	}
// 	logger.Info("return GetTrieDatabase trieDB", trieDB)
// 	return trieDB, true // trả về true nếu nó đã tồn tại, false nếu nó vừa được tạo
// }

// GetTrieDatabase lấy một TrieDatabase theo ID của nó.
func (manager *TrieDatabaseManager) GetOrCrateTrieDatabase(id common.Hash, hash common.Hash, address common.Address, dbName string) (*TrieDatabase, bool) {
	trieDB, exists := manager.trieDatabases[id]
	if !exists {
		trieDB = NewTrieDatabase(hash, manager.database, address, dbName, manager.accountStateDB)
		if trieDB == nil {
			return nil, false
		}
		manager.trieDatabases[id] = trieDB
	}
	return trieDB, true // trả về true nếu nó đã tồn tại, false nếu nó vừa được tạo
}

// RemoveTrieDatabase xóa một TrieDatabase khỏi danh sách quản lý
func (manager *TrieDatabaseManager) RemoveTrieDatabase(id common.Hash) {
	delete(manager.trieDatabases, id)
}

// ListAllIDs lấy danh sách tất cả các ID của TrieDatabase
func (manager *TrieDatabaseManager) ListAllIDs() []common.Hash {
	ids := make([]common.Hash, 0, len(manager.trieDatabases))
	for id := range manager.trieDatabases {
		ids = append(ids, id)
	}
	return ids
}

package smart_contract_db

import (
	"encoding/hex"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/maypok86/otter"

	p_common "gomail/mtn/common"
	"gomail/mtn/logger"
	"gomail/mtn/network"
	remote_db "gomail/remote_storage_db"
	"gomail/mtn/smart_contract"
	"gomail/mtn/storage"
	"gomail/mtn/trie"
	"gomail/mtn/types"
	t_network "gomail/mtn/types/network"
)

type RemoteStorageDB interface {
	Get(
		key []byte,
	) ([]byte, error)
	GetCode(
		address common.Address,
	) ([]byte, error)
	SetBlockNumber(blockNumber uint64)
	Close()
}

type AccountStateDB interface {
	AccountState(address common.Address) (types.AccountState, error)
}

type SmartContractDB struct {
	cacheRemoteDBs   otter.Cache[common.Address, RemoteStorageDB]
	cacheCode        otter.Cache[common.Hash, []byte]
	cacheStorageTrie otter.Cache[common.Address, *trie.MerklePatriciaTrie]

	messageSender      t_network.MessageSender
	dnsLink            string
	accountStateDB     AccountStateDB
	currentBlockNumber uint64

	sync.RWMutex

	updateDatas map[common.Address]types.SmartContractUpdateData
}

func NewSmartContractDB(
	messageSender t_network.MessageSender,
	dnsLink string,
	accountStateDB AccountStateDB,
	currentBlockNumber uint64,
) *SmartContractDB {
	cacheRemoteDBs, err := otter.MustBuilder[common.Address, RemoteStorageDB](20).
		CollectStats().
		DeletionListener(func(key common.Address, value RemoteStorageDB, cause otter.DeletionCause) {
			// close db to free connection and close read request
			value.Close()
		}).
		Build()
	if err != nil {
		logger.Error("error creating cacheRemoteDBs", "error", err)
		return nil
	}

	cacheCode, err := otter.MustBuilder[common.Hash, []byte](1_000).
		CollectStats().
		// Cost(func(key common.Hash, value []byte) int {
		//   return len(value)
		// }).
		Build()
	if err != nil {
		logger.Error("error creating cacheCode", "error", err)
		return nil
	}

	cacheStorageTrie, err := otter.MustBuilder[common.Address, *trie.MerklePatriciaTrie](1_000).
		CollectStats().
		Build()
	if err != nil {
		logger.Error("error creating cacheCode", "error", err)
		return nil
	}

	return &SmartContractDB{
		cacheRemoteDBs:     cacheRemoteDBs,
		cacheCode:          cacheCode,
		cacheStorageTrie:   cacheStorageTrie,
		dnsLink:            dnsLink,
		accountStateDB:     accountStateDB,
		messageSender:      messageSender,
		currentBlockNumber: currentBlockNumber,
		updateDatas:        make(map[common.Address]types.SmartContractUpdateData),
	}
}

func (scdb *SmartContractDB) SetAccountStateDB(asdb types.AccountStateDB) {
	scdb.accountStateDB = asdb
}

func (scdb *SmartContractDB) CreateRemoteStorageDB(as types.AccountState) (RemoteStorageDB, error) {
	// create connection to storage
	connection := network.NewConnection(
		as.SmartContractState().StorageAddress(),
		p_common.STORAGE_CONNECTION_TYPE,
		scdb.dnsLink,
	)
	err := connection.Connect()
	if err != nil {
		logger.Error("error creating remote db connection", "error", err)
		return nil, err
	}
	// read requests
	go connection.ReadRequest()
	remoteDB := remote_db.NewRemoteStorageDB(
		connection,
		scdb.messageSender,
		as.Address(),
	)
	scdb.cacheRemoteDBs.Set(as.Address(), remoteDB)
	return remoteDB, nil
}

func (scdb *SmartContractDB) Code(address common.Address) []byte {
	// get code hash
	as, _ := scdb.accountStateDB.AccountState(address)
	if as == nil ||
		as.SmartContractState() == nil { // check if account state exists and is a smart contract
		logger.Error("account state does not exist or is not a smart contract")
		return nil
	}
	// check cache
	if code, ok := scdb.cacheCode.Get(as.SmartContractState().CodeHash()); ok {
		return code
	}
	// get code from remote db
	remoteDB, ok := scdb.cacheRemoteDBs.Get(address)
	if !ok {
		var err error
		remoteDB, err = scdb.CreateRemoteStorageDB(as)
		if err != nil {
			logger.Error("error creating remote db", "error", err)
			return nil
		}
	}
	// set block number
	remoteDB.SetBlockNumber(scdb.currentBlockNumber)
	code, err := remoteDB.GetCode(
		as.Address(),
	)
	// verify code and code hash ??

	// cache code
	if err == nil {
		scdb.cacheCode.Set(as.SmartContractState().CodeHash(), code)
	}

	return code
}

func (scdb *SmartContractDB) StorageValue(address common.Address, key []byte) ([]byte, bool) {
	// get code hash
	as, _ := scdb.accountStateDB.AccountState(address)
	if as == nil ||
		as.SmartContractState() == nil { // check if account state exists and is a smart contract
		return common.Hash{}.Bytes(), true
	}

	scdb.Lock()
	defer scdb.Unlock()
	// set block number for remote db
	remoteDB, ok := scdb.cacheRemoteDBs.Get(address)
	if !ok {
		var err error
		remoteDB, err = scdb.CreateRemoteStorageDB(as)
		if err != nil {
			logger.Error("error creating remote db", "error", err)
			return common.Hash{}.Bytes(), false
		}
	}
	remoteDB.SetBlockNumber(scdb.currentBlockNumber)

	// get storage from trie
	storageTrie, ok := scdb.cacheStorageTrie.Get(as.Address())
	if !ok {
		// get storage root, and create trie
		var err error
		storageTrie, err = trie.New(as.SmartContractState().StorageRoot(), remoteDB)
		if err != nil {
			logger.Error("error getting storage value", "error", err)
			return common.Hash{}.Bytes(), false
		}
		scdb.cacheStorageTrie.Set(as.Address(), storageTrie)
	} else {
		logger.Info("storage trie from cache")
	}
	value, err := storageTrie.Get(
		key,
	)
	if err != nil {
		return value, false
	}

	// todo verify value and key
	rootHash := storageTrie.Hash()
	if rootHash != as.SmartContractState().StorageRoot() {
		logger.Error(
			"storage root does not match root hash " +
				rootHash.String() +
				"as.SmartContractState().StorageRoot() " +
				as.SmartContractState().
					StorageRoot().
					String(),
		)
	} else {
		logger.Info("storage root matches")
	}
	if value == nil {
		return common.Hash{}.Bytes(), true
	}
	return value, true
}

func (scdb *SmartContractDB) SetBlockNumber(blockNumber uint64) {
	scdb.currentBlockNumber = blockNumber
}

func (scdb *SmartContractDB) SetCode(
	address common.Address,
	codeHash common.Hash,
	code []byte,
) {
	scdb.Lock()
	defer scdb.Unlock()
	scdb.cacheCode.Set(codeHash, code)
	if _, ok := scdb.updateDatas[address]; !ok {
		updateData := smart_contract.NewSmartContractUpdateData(
			[]byte{},
			map[string][]byte{},
			[]types.EventLog{},
		)
		scdb.updateDatas[address] = updateData
	}
	updateData := scdb.updateDatas[address]
	updateData.SetCode(code)
}

func (scdb *SmartContractDB) SetStorageValue(
	address common.Address,
	key []byte,
	value []byte,
) error {
	scdb.Lock()
	defer scdb.Unlock()
	storageTrie, ok := scdb.cacheStorageTrie.Get(address)
	if !ok {
		// create memory storage trie
		memoryDB := storage.NewMemoryDb()
		storageTrie, _ = trie.New(common.Hash{}, memoryDB)
		// set to cache
		scdb.cacheStorageTrie.Set(address, storageTrie)
	}
	storageTrie.Update(key, value)

	if _, ok := scdb.updateDatas[address]; !ok {
		updateData := smart_contract.NewSmartContractUpdateData(
			[]byte{},
			map[string][]byte{},
			[]types.EventLog{},
		)
		scdb.updateDatas[address] = updateData
	}
	updateData := scdb.updateDatas[address]
	updateData.UpdateStorage(map[string][]byte{
		hex.EncodeToString(key): value,
	})
	return nil
}

func (scdb *SmartContractDB) AddEventLogs(eventLogs []types.EventLog) {
	for _, eventLog := range eventLogs {
		address := eventLog.Address()
		if _, ok := scdb.updateDatas[address]; !ok {
			updateData := smart_contract.NewSmartContractUpdateData(
				[]byte{},
				map[string][]byte{},
				[]types.EventLog{},
			)
			scdb.updateDatas[address] = updateData
		}
		updateData := scdb.updateDatas[address]
		updateData.AddEventLog(eventLog)
	}
}

func (scdb *SmartContractDB) NewTrieStorage(
	address common.Address,
) common.Hash {
	scdb.Lock()
	defer scdb.Unlock()
	// get code hash
	as, _ := scdb.accountStateDB.AccountState(address)
	if as == nil ||
		as.SmartContractState() == nil { // check if account state exists and is a smart contract
		return common.Hash{}
	}
	// set block number for remote db
	remoteDB, ok := scdb.cacheRemoteDBs.Get(address)
	if !ok {
		var err error
		remoteDB, err = scdb.CreateRemoteStorageDB(as)
		if err != nil {
			logger.Error("error creating remote db", "error", err)
			return common.Hash{}
		}
	}
	remoteDB.SetBlockNumber(scdb.currentBlockNumber)

	// get storage from trie
	storageTrie, ok := scdb.cacheStorageTrie.Get(as.Address())
	if !ok {
		// get storage root, and create trie
		var err error
		storageTrie, err = trie.New(as.SmartContractState().StorageRoot(), remoteDB)
		if err != nil {
			logger.Error("error getting storage value", "error", err)
			return common.Hash{}
		}
		scdb.cacheStorageTrie.Set(as.Address(), storageTrie)
	} else {
		logger.Info("storage trie from cache")
	}
	return storageTrie.Hash()
}

func (scdb *SmartContractDB) StorageRoot(
	address common.Address,
) common.Hash {
	scdb.RLock()
	defer scdb.RUnlock()
	if storageTrie, ok := scdb.cacheStorageTrie.Get(address); ok {
		return storageTrie.Hash()
	}
	return common.Hash{}
}

func (scdb *SmartContractDB) DeleteAddress(
	address common.Address,
) {
	scdb.Lock()
	defer scdb.Unlock()
	scdb.cacheRemoteDBs.Delete(address)
	scdb.cacheStorageTrie.Delete(address)
}

func (scdb *SmartContractDB) GetSmartContractUpdateDatas() map[common.Address]types.SmartContractUpdateData {
	return scdb.updateDatas
}

func (scdb *SmartContractDB) ClearSmartContractUpdateDatas() {
	scdb.updateDatas = make(map[common.Address]types.SmartContractUpdateData)
}

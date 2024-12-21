package stake

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"

	"gomail/mtn/logger"
	"gomail/mtn/storage"
	"gomail/mtn/trie"
)

type StakeSmartContractDb struct {
	codePath        string
	storageRootPath string
	storageDBPath   string
	checkPointPath  string

	code        []byte
	storageTrie *trie.MerklePatriciaTrie
	storageDB   storage.Storage
	storageRoot common.Hash
	hasDirty    bool
}

func NewStakeSmartContractDb(
	codePath string,
	storageRootPath string,
	storageDBPath string,
	checkPointPath string,
) (*StakeSmartContractDb, error) {
	// load latest stake storage root
	bStorageRoot, err := os.ReadFile(storageRootPath)
	if err != nil {
		panic(err)
	}
	// load latest stake storage root
	hexCode, err := os.ReadFile(codePath)
	if err != nil {
		panic(err)
	}
	code := common.FromHex(string(hexCode))

	levelDB, err := storage.NewLevelDB(storageDBPath)
	if err != nil {
		return nil, err
	}
	storageRoot := common.HexToHash(string(bStorageRoot))
	logger.Info("load stake storage root", storageRoot)
	t, err := trie.New(storageRoot, levelDB)
	if err != nil {
		return nil, err
	}

	return &StakeSmartContractDb{
		codePath:        codePath,
		storageRootPath: storageRootPath,
		storageDBPath:   storageDBPath,
		checkPointPath:  checkPointPath,

		code:        code,
		storageTrie: t,
		storageRoot: storageRoot,
		storageDB:   levelDB,
	}, nil
}

// get code
func (s *StakeSmartContractDb) Code(address common.Address) []byte {
	return s.code
}

// get storage
func (s *StakeSmartContractDb) StorageValue(address common.Address, key []byte) ([]byte, bool) {
	v, err := s.storageTrie.Get(key)
	if err != nil {
		return v, false
	}
	return v, true
}

// set storage
func (s *StakeSmartContractDb) UpdateStorageValue(key []byte, value []byte) {
	s.storageTrie.Update(key, value)
}

// Commit changes to storage
func (s *StakeSmartContractDb) Commit() (common.Hash, error) {
	hash, nodeSet, oldKeys, err := s.storageTrie.Commit(true)
	if err != nil {
		return common.Hash{}, err
	}
	// save to db
	for i := 0; i < len(oldKeys); i++ {
		// should save oldKeys to archive db for backup
		s.storageDB.Delete(oldKeys[i])
	}
	if nodeSet != nil {
		// save nodeSet to db
		batch := [][2][]byte{}
		for _, node := range nodeSet.Nodes {
			batch = append(batch, [2][]byte{node.Hash.Bytes(), node.Blob})
		}
		err := s.storageDB.BatchPut(batch)
		if err != nil {
			return common.Hash{}, err
		}
	}
	s.hasDirty = false
	s.storageRoot = hash
	s.storageTrie, _ = trie.New(hash, s.storageDB)

	// save latest stake storage root
	err = os.WriteFile(s.storageRootPath, []byte(hex.EncodeToString(hash.Bytes())), 0644)
	return hash, nil
}

// set storage
func (s *StakeSmartContractDb) Discard() {
	s.storageTrie, _ = trie.New(s.storageRoot, s.storageDB)
	s.hasDirty = false
}

// Copy from
func (s *StakeSmartContractDb) CopyFrom(from *StakeSmartContractDb) {
	s.code = from.code
	s.storageTrie = from.storageTrie
}

func (s *StakeSmartContractDb) HasDirty() bool {
	return s.hasDirty
}

func (s *StakeSmartContractDb) StakeStorageDB() storage.Storage {
	return s.storageDB
}

func (s *StakeSmartContractDb) StakeStorageRoot() common.Hash {
	return s.storageRoot
}

func (s *StakeSmartContractDb) SaveCheckPoint(blockNumber uint64) error {
	logger.Info("save stake smart contract db check point", blockNumber)
	// create db
	dbSnapShot := s.storageDB.GetSnapShot()
	checkPointDB, err := storage.NewLevelDB(
		fmt.Sprintf("%v%v", s.checkPointPath, blockNumber),
	)
	//
	if err != nil {
		return err
	}
	//
	go func(snapShot storage.SnapShot, checkPointDB storage.Storage) {
		// save to db
		iter := snapShot.GetIterator()
		defer iter.Release()
		defer checkPointDB.Close()
		batch := [][2][]byte{}
		for iter.Next() {
			cKey := make([]byte, len(iter.Key()))
			cValue := make([]byte, len(iter.Value()))
			copy(cKey, iter.Key())
			copy(cValue, iter.Value())
			batch = append(batch, [2][]byte{cKey, cValue})
		}
		iter.Release()
		err := checkPointDB.BatchPut(batch)
		if err != nil {
			logger.Error("save stake smart contract db check point failed", err)
		}
		// close db
	}(dbSnapShot, checkPointDB)

	return nil
}

func (s *StakeSmartContractDb) UpdateFromCheckPointData(
	storageRoot common.Hash,
	storageData [][2][]byte,
) error {
	s.storageRoot = storageRoot
	err := s.storageDB.BatchPut(storageData)
	if err != nil {
		return err
	}
	s.storageTrie, err = trie.New(storageRoot, s.storageDB)
	if err != nil {
		return err
	}
	return nil
}

func (s *StakeSmartContractDb) CheckPointPath() string {
	return s.checkPointPath
}

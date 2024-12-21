package storage

import (
	"time"

	"gomail/mtn/logger"
	pb "gomail/mtn/proto"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDB struct {
	db        *leveldb.DB
	closed    bool
	path      string
	closeChan chan bool
}

type LevelDBSnapShot struct {
	snapShot *leveldb.Snapshot
}

func NewLevelDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	lvDb := &LevelDB{db, false, path, make(chan bool)}
	// Run garbage collection every hour.
	go func() {
		for {
			select {
			case <-time.After(15 * time.Minute):
				lvDb.Compact()
			case <-lvDb.closeChan:
				logger.Debug("LevelDB closeChan")
				return
			}
		}
	}()
	return lvDb, nil
}

func (ldb *LevelDB) Get(key []byte) ([]byte, error) {
	return ldb.db.Get(key, nil)
}

func (ldb *LevelDB) Put(key, value []byte) error {
	return ldb.db.Put(key, value, nil)
}

func (ldb *LevelDB) Compact() error {
	return ldb.db.CompactRange(util.Range{
		Start: nil,
		Limit: nil,
	})
}

func (ldb *LevelDB) Has(key []byte) bool {
	has, _ := ldb.db.Has(key, nil)
	return has
}

func (ldb *LevelDB) Delete(key []byte) error {
	return ldb.db.Delete(key, nil)
}

func (ldb *LevelDB) BatchPut(kvs [][2][]byte) error {
	batch := new(leveldb.Batch)
	for i := range kvs {
		batch.Put(kvs[i][0], kvs[i][1])
	}
	return ldb.db.Write(batch, nil)
}

func (ldb *LevelDB) Open() error {
	var err error
	if ldb.closed {
		ldb.db, err = leveldb.OpenFile(ldb.path, nil)
		if err != nil {
			return err
		}
		ldb.closed = false
	}
	return nil
}

func (ldb *LevelDB) Close() error {
	if !ldb.closed {
		ldb.closeChan <- true
		err := ldb.db.Close()
		if err != nil {
			return err
		}
		ldb.closed = true
	}
	return nil
}

func (ldb *LevelDB) GetIterator() IIterator {
	return ldb.db.NewIterator(nil, nil)
}

func (ldb *LevelDB) GetSnapShot() SnapShot {
	snapShot, _ := ldb.db.GetSnapshot()
	return NewLevelDBSnapShot(snapShot)
}

func (ldb *LevelDB) Stats() *pb.LevelDBStats {
	stats := &leveldb.DBStats{}
	ldb.db.Stats(stats)
	levelSizes := make([]uint64, len(stats.LevelSizes))
	for i, v := range stats.LevelSizes {
		levelSizes[i] = uint64(v)
	}
	levelRead := make([]uint64, len(stats.LevelRead))
	for i, v := range stats.LevelRead {
		levelRead[i] = uint64(v)
	}
	levelTablesCounts := make([]uint64, len(stats.LevelTablesCounts))
	for i, v := range stats.LevelTablesCounts {
		levelTablesCounts[i] = uint64(v)
	}
	levelWrite := make([]uint64, len(stats.LevelWrite))
	for i, v := range stats.LevelWrite {
		levelWrite[i] = uint64(v)
	}
	levelDurations := make([]uint64, len(stats.LevelDurations))
	for i, v := range stats.LevelDurations {
		levelDurations[i] = uint64(v)
	}
	pbStat := &pb.LevelDBStats{
		LevelSizes:        levelSizes,
		LevelTablesCounts: levelTablesCounts,
		LevelRead:         levelRead,
		LevelWrite:        levelWrite,
		LevelDurations:    levelDurations,
		MemComp:           stats.MemComp,
		Level0Comp:        stats.Level0Comp,
		NonLevel0Comp:     stats.NonLevel0Comp,
		SeekComp:          stats.SeekComp,
		AliveSnapshots:    stats.AliveSnapshots,
		AliveIterators:    stats.AliveIterators,
		IOWrite:           stats.IOWrite,
		IORead:            stats.IORead,
		BlockCacheSize:    int32(stats.BlockCacheSize),
		OpenedTablesCount: int32(stats.OpenedTablesCount),
		Path:              ldb.path,
	}
	return pbStat
}

func NewLevelDBSnapShot(snapShot *leveldb.Snapshot) *LevelDBSnapShot {
	return &LevelDBSnapShot{
		snapShot: snapShot,
	}
}

func (snapShot *LevelDBSnapShot) GetIterator() IIterator {
	return snapShot.snapShot.NewIterator(nil, nil)
}

func (snapShot *LevelDBSnapShot) Release() {
	snapShot.snapShot.Release()
}

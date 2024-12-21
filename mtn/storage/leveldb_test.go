package storage

import (
	"encoding/hex"
	fmt "fmt"
	"runtime"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/maypok86/otter"
	"github.com/stretchr/testify/assert"

	"gomail/mtn/logger"
)

var ldb *LevelDB

func initLevelDB() *LevelDB {
	path := "./ldb_test"
	ldb, _ = NewLevelDB(path)
	return ldb
}

func TestInitLevelDB(t *testing.T) {
	ldb = initLevelDB()
	assert.NotEmpty(t, ldb)
}

func TestPut(t *testing.T) {
	err1 := ldb.Put(common.FromHex("0x001"), common.FromHex("0xf1f1"))
	err2 := ldb.Put(common.FromHex("0x002"), common.FromHex("0xh2"))
	err3 := ldb.Put(common.FromHex("0x003"), common.FromHex("0xh3h3"))
	assert.Empty(t, err1)
	assert.Empty(t, err2)
	assert.Empty(t, err3)
}

func TestGet(t *testing.T) {
	ldb = initLevelDB()
	db1, err := ldb.Get(
		common.FromHex("0x85d91f19f9fc0c4c3985e3ff3c384e867da222a144d918b6f3cf8daecf232162"),
	)
	logger.Info("db1", hex.EncodeToString(db1))
	logger.Info("err", err)
}

func TestHas(t *testing.T) {
	chk1 := ldb.Has(common.FromHex("0x001"))
	chk2 := ldb.Has(common.FromHex("0x002"))
	chk3 := ldb.Has(common.FromHex("0x003"))
	chk4 := ldb.Has(common.FromHex("0x004"))
	assert.True(t, chk1)
	assert.True(t, chk2)
	assert.True(t, chk3)
	assert.False(t, chk4)
}

func TestDelete(t *testing.T) {
	ldb.Put(common.FromHex("0x004"), common.FromHex("0x00f0ff"))
	chk1 := ldb.Has(common.FromHex("0x004"))
	assert.True(t, chk1)
	ldb.Delete(common.FromHex("0x004"))
	chk2 := ldb.Has(common.FromHex("0x004"))
	assert.False(t, chk2)
}

func TestBatchPut(t *testing.T) {
	arr := [][2][]byte{
		{
			common.FromHex("0x001"),
			common.FromHex("0xh1111"),
		},
		{
			common.FromHex("0x002"),
			common.FromHex("0xf1f2"),
		},
		{
			common.FromHex("0x003"),
			common.FromHex("0x0000"),
		},
	}
	assert.NoError(t, ldb.BatchPut(arr))
	db1, _ := ldb.Get(common.FromHex("0x001"))
	assert.Equal(t, db1, common.FromHex("0xh1111"))
	db2, _ := ldb.Get(common.FromHex("0x002"))
	assert.Equal(t, db2, common.FromHex("0xf1f2"))
	db3, _ := ldb.Get(common.FromHex("0x003"))
	assert.Equal(t, db3, common.FromHex("0x0000"))
}

func TestGetIterator(t *testing.T) {
	ldb = initLevelDB()
	iterator := ldb.GetIterator()
	assert.NotEmpty(t, iterator)
}

func TestCloseDB(t *testing.T) {
	ldb = initLevelDB()
	err := ldb.Open()
	fmt.Printf("ERR1 %v", err)
	// _, err = NewLevelDB(path)
	// fmt.Printf("ERR2 %v", err)
}

func TestCacheDB(t *testing.T) {
	path1 := "./ldb_test/1"
	ldb1, _ := NewLevelDB(path1)
	//
	path2 := "./ldb_test/2"
	ldb2, _ := NewLevelDB(path2)
	//
	path3 := "./ldb_test/3"
	ldb3, _ := NewLevelDB(path3)
	cache, err := otter.MustBuilder[string, *LevelDB](100_000).
		CollectStats().
		DeletionListener(func(key string, value *LevelDB, cause otter.DeletionCause) {
			value.Close()
		}).
		Build()
	if err != nil {
		panic(err)
	}
	cache.Set(path1, ldb1)
	cache.Set(path2, ldb2)
	cache.Set(path3, ldb3)
	ldb1, _ = cache.Get(path1)
	ldb2, _ = cache.Get(path2)
	ldb3, _ = cache.Get(path3)

	var o runtime.MemStats
	runtime.ReadMemStats(&o)
	fmt.Printf("%v", o.Alloc/1024/1024)
}

package storage

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gomail/pkg/logger"
	"gomail/pkg/loggerfile"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Cấu trúc quản lý một LevelDB với snapshot
type ReplicatedLevelDB struct {
	primaryDB *leveldb.DB
	snapshot  *leveldb.Snapshot
	path      string
	mu        sync.RWMutex // Thêm mutex để kiểm soát truy cập đồng thời

}

// Tạo mới `ReplicatedLevelDB`
func NewReplicatedLevelDB(path string) *ReplicatedLevelDB {
	return &ReplicatedLevelDB{path: path}
}

// Mở LevelDB và tạo snapshot ban đầu
func (r *ReplicatedLevelDB) Open(parallelism int) error {
	if err := createDirIfNotExists(r.path); err != nil {
		logger.Error("Không thể tạo thư mục DB:", err)
		return err
	}

	// Mở LevelDB với các options đã được cấu hình
	var err error

	// Open the database in read-only mode
	r.primaryDB, err = leveldb.OpenFile(r.path, nil)

	if err != nil {
		logger.Error("Lỗi mở LevelDB:", err, r.path)
		return err
	}

	// Tạo snapshot ban đầu
	r.snapshot, err = r.primaryDB.GetSnapshot()
	if err != nil {
		logger.Error("Lỗi tạo snapshot ban đầu:", err)
		return err
	}

	// logger.Info("Database đã mở thành công:", r.path)
	return nil
}

// Ghi dữ liệu vào Primary và cập nhật snapshot
func (r *ReplicatedLevelDB) Put(key, value []byte) error {
	// logger.Error("Put call")
	err := r.primaryDB.Put(key, value, nil)
	if err != nil {
		return err
	}

	// Cập nhật snapshot sau khi ghi
	return r.updateSnapshot()
}

func (r *ReplicatedLevelDB) BatchPut(kvs [][2][]byte) error {
	batch := new(leveldb.Batch)
	for _, kv := range kvs {

		batch.Put(kv[0], kv[1])
	}
	writeOptions := &opt.WriteOptions{
		Sync: false, // Bắt buộc ghi vào đĩa ngay lập tức
	}
	// Ghi batch vào primary database
	err := r.primaryDB.Write(batch, writeOptions)
	if err != nil {
		return err
	}

	// Cập nhật snapshot sau khi batch write
	return r.updateSnapshot()
}

// Ưu tiên đọc dữ liệu từ snapshot nếu lỗi đọc từ db
func (r *ReplicatedLevelDB) Get(key []byte) ([]byte, error) {
	// r.mu.RLock() // Chờ đến khi không có updateSnapshot() đang chạy
	// defer r.mu.RUnlock()
	if r.snapshot == nil {
		logger.Info("Snapshot chưa được khởi tạo")

		return nil, fmt.Errorf("snapshot chưa được khởi tạo")
	}
	value, err := r.snapshot.Get(key, nil)
	if err != nil {
		// Debug
		logger.Info("Get from snapshot err", err)
		value, err = r.primaryDB.Get(key, nil)
		if err != nil {
			logger.Error("Get from primaryDB err", err)
			fileLogger, _ := loggerfile.NewFileLogger("error/getLevelDb.log")
			fileLogger.Info(fmt.Sprintf("getting key from primaryDB: %s value: %s", hex.EncodeToString(key), value))
		}
		// panic(fmt.Sprintf("Dừng chương trình do Get db thất bại: key=%s", hex.EncodeToString(key)))
	}
	// Thêm lệnh debug ở đây
	return value, err
}

// Xóa key khỏi Primary và cập nhật snapshot
func (r *ReplicatedLevelDB) Delete(key []byte) error {

	err := r.primaryDB.Delete(key, nil)
	if err != nil {
		return err
	}

	// Cập nhật snapshot sau khi xóa
	return r.updateSnapshot()
}

// Kiểm tra key có tồn tại không
func (r *ReplicatedLevelDB) Has(key []byte) bool {
	if r.snapshot == nil {
		return false
	}
	exists, _ := r.snapshot.Has(key, nil)
	return exists
}

// Lấy tất cả key trong database (chỉ dùng để debug)
func (r *ReplicatedLevelDB) GetAllKeys() ([]string, error) {

	var keys []string
	iter := r.primaryDB.NewIterator(nil, nil)
	for iter.Next() {
		keys = append(keys, string(iter.Key()))
	}
	iter.Release()

	if err := iter.Error(); err != nil {
		return nil, err
	}
	return keys, nil
}

// Đóng LevelDB
func (r *ReplicatedLevelDB) Close() error {

	var err error
	if r.snapshot != nil {
		r.snapshot.Release()
	}
	if r.primaryDB != nil {
		err = r.primaryDB.Close()
	}
	return err
}

func (r *ReplicatedLevelDB) updateSnapshot() error {
	// start := time.Now() // Ghi nhận thời điểm bắt đầu
	// r.mu.Lock() // Khóa ghi để đảm bảo không có Get() nào đang chạy
	// defer r.mu.Unlock()
	if r.snapshot != nil {
		r.snapshot.Release()
	}

	var err error
	r.snapshot, err = r.primaryDB.GetSnapshot()
	if err != nil {
		logger.Error("Lỗi cập nhật snapshot:", err)
		return err
	}

	// elapsed := time.Since(start) // Tính toán thời gian đã trôi qua
	// logger.Info("Snapshot đã được cập nhật. Thời gian:", elapsed) // Ghi log thời gian
	return nil
}

// Tạo thư mục nếu chưa tồn tại
func createDirIfNotExists(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

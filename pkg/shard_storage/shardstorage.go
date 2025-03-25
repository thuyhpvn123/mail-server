package shard_storage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ShardStorage là cấu trúc lưu trữ dữ liệu phân mảnh
type ShardStorage struct {
	maxBlocksPerShard int
	shardDir          string
	lineByte          int
}

// NewShardStorage tạo một đối tượng ShardStorage mới
func NewShardStorage(maxBlocksPerShard int, shardDir string, lineByte int) (*ShardStorage, error) {
	// Kiểm tra tính hợp lệ của tham số
	if maxBlocksPerShard <= 0 || lineByte <= 0 {
		return nil, fmt.Errorf("invalid parameters: maxBlocksPerShard and lineByte must be positive")
	}
	// Tạo thư mục shards nếu chưa tồn tại
	if err := os.MkdirAll(shardDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create shard directory: %w", err)
	}

	return &ShardStorage{
		maxBlocksPerShard: maxBlocksPerShard,
		shardDir:          shardDir,
		lineByte:          lineByte,
	}, nil
}

// getShardFileName lấy tên file shard dựa trên số thứ tự shard
func (ss *ShardStorage) getShardFileName(shardIndex int) string {
	return fmt.Sprintf("%s/shard_%d.txt", ss.shardDir, shardIndex)
}

// createShardFileIfNeeded tạo file shard nếu chưa tồn tại, xử lý lỗi tốt hơn
func (ss *ShardStorage) createShardFileIfNeeded(shardFile string) error {
	if _, err := os.Stat(shardFile); os.IsNotExist(err) {
		file, err := os.Create(shardFile)
		if err != nil {
			return fmt.Errorf("failed to create shard file: %w", err)
		}
		defer file.Close() // Đóng file sau khi sử dụng
	}
	return nil
}

func (ss *ShardStorage) ensureLinesExist(file *os.File, lineIndex int) error {
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	fileSize := fileInfo.Size()
	// Thay đổi ở đây: Sử dụng phép chia int64
	neededLines := lineIndex - int(fileSize/(int64(ss.lineByte+1))) + 1

	if neededLines > 0 {
		_, err = file.Seek(fileSize, 0)
		if err != nil {
			return fmt.Errorf("failed to seek to end of file: %w", err)
		}
		for i := 0; i < neededLines; i++ {
			_, err = file.WriteString(strings.Repeat(" ", ss.lineByte) + "\n")
			if err != nil {
				return fmt.Errorf("failed to write empty line: %w", err)
			}
		}
	}
	return nil
}

// updateLine cập nhật một dòng trong file shard
func (ss *ShardStorage) updateLine(shardFile string, lineIndex int, newLine string) error {
	file, err := os.OpenFile(shardFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open shard file: %w", err)
	}
	defer file.Close()

	// Kiểm tra và tạo các dòng trống nếu cần thiết
	if err := ss.ensureLinesExist(file, lineIndex); err != nil {
		return fmt.Errorf("failed to ensure lines exist: %w", err)
	}

	// Viết dòng mới vào file
	_, err = file.Seek(int64(lineIndex*(ss.lineByte+1)), 0)
	if err != nil {
		return fmt.Errorf("failed to seek to line: %w", err)
	}
	_, err = file.WriteString(newLine + "\n")
	if err != nil {
		return fmt.Errorf("failed to write line: %w", err)
	}

	return nil
}

// getLine đọc dòng thứ n của file
func (ss *ShardStorage) getLine(shardFile string, lineIndex int) (string, error) {
	file, err := os.Open(shardFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to open shard file: %w", err)
	}
	defer file.Close()

	// Kiểm tra file rỗng
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}
	if fileInfo.Size() == 0 {
		return "", nil
	}
	lineByteOffset := int64(lineIndex) * int64(ss.lineByte+1) // Sửa đổi dòng này
	if lineByteOffset >= fileInfo.Size() {
		return "", fmt.Errorf("line number out of range")
	}

	_, err = file.Seek(lineByteOffset, 0)
	if err != nil {
		return "", fmt.Errorf("failed to seek to line: %w", err)
	}

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading line: %w", err)
	}
	return "", nil // Dòng không tồn tại
}

// saveBlock lưu block vào file shard, đảm bảo dòng cuối cùng khớp. Sử dụng updateLine để cập nhật.
func (ss *ShardStorage) SetIndexValue(blockNumber int, blockHash string) error {
	shardIndex := (blockNumber - 1) / ss.maxBlocksPerShard
	lineIndex := (blockNumber - 1) % ss.maxBlocksPerShard
	shardFile := ss.getShardFileName(shardIndex)

	// Kiểm tra nếu file shard chưa tồn tại và tạo nếu cần
	if err := ss.createShardFileIfNeeded(shardFile); err != nil {
		return fmt.Errorf("failed to create or access shard file: %w", err)
	}

	// Sử dụng updateLine để ghi block vào file
	err := ss.updateLine(shardFile, lineIndex, blockHash)
	if err != nil {
		return fmt.Errorf("failed to update line: %w", err)
	}

	return nil
}

// findBlockHashByBlockNumber tìm blockHash dựa trên blockNumber
func (ss *ShardStorage) FindValueByIndex(blockNumber int) (string, error) {
	shardIndex := (blockNumber - 1) / ss.maxBlocksPerShard
	lineIndex := (blockNumber - 1) % ss.maxBlocksPerShard
	shardFile := ss.getShardFileName(shardIndex)

	blockHash, err := ss.getLine(shardFile, lineIndex)
	if err != nil {
		return "", fmt.Errorf("error finding block hash: %w", err)
	}

	return blockHash, nil
}

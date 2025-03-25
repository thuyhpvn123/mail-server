package main

import (
	"encoding/hex"
	"fmt"

	"gomail/pkg/logger"
	pb "gomail/pkg/proto"
	"gomail/pkg/transaction"

	"google.golang.org/protobuf/proto"
)

func main() {
	// Ví dụ sử dụng
	hexString := "0a208dbce4440e64fef64a0a907b550da8e6079b81394cb226b4bb918712b1d0b07a121426d209379611be4829eede2d20232d9cfc7ef7f41a07470de4df8200002220000000000000000000000000000000000000000000000000002386f26fc1000028a09c0130c0843d38f02e5a20290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e5636220bd9d91783ae8bbc7a3e69008772e266d2ff3eded1b09265942e9d43d889c31486a08000000000000000172142f4cb880116850929d8b44fac82e907bc21f19d0"

	transactionProto, err := transactionInfoFromHex(hexString)
	if err != nil {
		logger.Error("Lỗi khi lấy thông tin giao dịch:", err)
		return
	}

	tx := transaction.TransactionFromProto(transactionProto)
	fmt.Println("Transaction:", tx)

	fmt.Println("Hash:", tx.Hash())

}

func transactionInfoFromHex(hexData string) (*pb.Transaction, error) {
	// Chuyển đổi chuỗi hex thành byte slice
	data, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, fmt.Errorf("lỗi chuyển đổi hex sang byte: %w", err)
	}

	// Khởi tạo cấu trúc Transaction
	transaction := &pb.Transaction{}

	// Giải mã dữ liệu vào cấu trúc Transaction
	err = proto.Unmarshal(data, transaction)
	if err != nil {
		return nil, fmt.Errorf("lỗi giải mã dữ liệu: %w", err)
	}

	return transaction, nil
}

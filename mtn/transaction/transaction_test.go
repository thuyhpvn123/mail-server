package transaction

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/proto"

	p_common "gomail/mtn/common"
	"gomail/mtn/logger"
	pb "gomail/mtn/proto"
)

func TestMarshal(t *testing.T) {
	transaction := NewTransaction(
		common.HexToHash(
			"0x0000000000000000000000000000000000000000000000000000000000000000",
		),
		p_common.PubkeyFromBytes(
			common.FromHex(
				"b993b5e8daf8c84f86ed25456a635bf894a7c82d641bddbf404432fd459b0d04b5db1968d0b671c0d0211784c3c90e36",
			),
		),
		common.HexToAddress("ba00cfe5a5697c1f7e6d7243a6608f69e7a0a99c"),
		big.NewInt(1000000),
		big.NewInt(1000000),
		1000000,
		1000000,
		1000000,
		pb.ACTION_EMPTY,
		common.FromHex("3"),
		[][]byte{common.FromHex("4")},
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000005"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000006"),
	)
	txWd, _ := transaction.Marshal()
	msg := &pb.Message{
		Body: txWd,
	}
	bData, err := proto.Marshal(msg)
	logger.Info("cc %v", hex.EncodeToString(bData))
	if err != nil {
		panic(err)
		logger.Error("Error when marshal ", err)
		return
	}
	// unmarshal
	msg = &pb.Message{}
	err = proto.Unmarshal(bData, msg)
	pbTransaction := &pb.Transaction{}
	logger.Info("tx: %v", err)
	err = proto.Unmarshal(msg.Body, pbTransaction)
	if err != nil {
		panic(err)
	}
	tx := &Transaction{}
	tx.FromProto(pbTransaction)
	logger.Info("tx: ", tx.String())
}

func TestUnmarshal(t *testing.T) {
	// bData := common.FromHex(
	// 	"0ad1010a1c53656e645472616e73616374696f6e576974684465766963654b657912308f4f64095a726ec28a2854ed287c4e1a91104881f10438bb44e6add15228621ac5b7be6a1ca2fd6cc0591ce1c3c8e0bb1a1400000000000000000000000000000000000000002260a5bf81d74098014676646863bbdb88ca0e15ae4a07a2260c309686a9c98cbfe53bce7dbe10a69f9dedfa5950f0962c901183e6421e0e12379e3257a95898cd039bc0c1e0e5a9db89965800093f64beddb23ca981df772c6185c0acebb9a17bfc2a07302e302e302e3112cc040aa7040a2041198f43f854c36d1a64370a3cc84bd018f6bc72cabb1269e1b13a88626278b0122000000000000000000000000000000000000000000000000000000000000000001a308f4f64095a726ec28a2854ed287c4e1a91104881f10438bb44e6add15228621ac5b7be6a1ca2fd6cc0591ce1c3c8e0bb2214ba00cfe5a5697c1f7e6d7243a6608f69e7a0a99c2a2000000000000000000000000000000000000000000000000000000000000000003220000000000000000000000000000000000000000000000000000000000000000038c096b1024080c8afa02548e80750045a260a245e6c6c010000000000000000000000000000000000000000000000000000000000000000621477c51d415d49186325daf00ec8b20fe1456919086a2000000000000000000000000000000000000000000000000000000000000000007220510e4e770828ddbf7f7b00ab00a9f6adaf81c0dc9cc85f1f8249c256942d61d97a6095300b51fed0c36873f7447d1a75bd44bfbb5a6b18a6d0408fdbbfeb1dc529efb44cc79f035c2df4b477b6ffede84cda0eddce8fd01c833ce70ebcbc5ec6e0147371089de21dfb6a17f31dfc959f2edfc5399dec112e482f896b5ab6a629d7d2820160b119416dc3ef1e4be8706f3f2b074e5461d9f5ad03bb62354460a1c5230c6b8ea211467175081faae485faead2735e440ed66ca2ba21bceba8beeb08cd2126b5989962699f5b3ee50fcad41502f41096ba89f4018dc6781fca90c2d33d6f73091220290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563",
	// )
	// 12ec010a20fe79f45c0a994b5f5bb95a42f03852e95153107a15f04e5a5c22b75ae21ef5f5122000000000000000000000000000000000000000000000000000000000000000001a30b993b5e8daf8c84f86ed25456a635bf894a7c82d641bddbf404432fd459b0d04b5db1968d0b671c0d0211784c3c90e362214ba00cfe5a5697c1f7e6d7243a6608f69e7a0a99c2a030f424032030f424038c0843d40c0843d48c0843d5a01036201046a20000000000000000000000000000000000000000000000000000000000000000572200000000000000000000000000000000000000000000000000000000000000006
	// 0a2000000000000000000000000000000000000000000000000000000000000000001230b993b5e8daf8c84f86ed25456a635bf894a7c82d641bddbf404432fd459b0d04b5db1968d0b671c0d0211784c3c90e361a14ba00cfe5a5697c1f7e6d7243a6608f69e7a0a99c22030f42402a030f424030c0843d38c0843d40c0843d5201035a0104622000000000000000000000000000000000000000000000000000000000000000056a200000000000000000000000000000000000000000000000000000000000000006
	bData := common.FromHex(
		"12ec010a20fe79f45c0a994b5f5bb95a42f03852e95153107a15f04e5a5c22b75ae21ef5f5122000000000000000000000000000000000000000000000000000000000000000001a30b993b5e8daf8c84f86ed25456a635bf894a7c82d641bddbf404432fd459b0d04b5db1968d0b671c0d0211784c3c90e362214ba00cfe5a5697c1f7e6d7243a6608f69e7a0a99c2a030f424032030f424038c0843d40c0843d48c0843d5a01036201046a20000000000000000000000000000000000000000000000000000000000000000572200000000000000000000000000000000000000000000000000000000000000006",
	)
	msg := &pb.Message{}
	proto.Unmarshal(bData, msg)
	// txWd := &pb.TransactionWithDeviceKey{}
	txWd := &pb.Transaction{}
	proto.Unmarshal(msg.Body, txWd)
	transaction := &Transaction{}
	// transaction.FromProto(txWd.Transaction)
	transaction.FromProto(txWd)
	fmt.Printf("Transaction %v", transaction)
	fmt.Printf("Transaction %v", transaction.RelatedAddresses())
	fmt.Printf("Transaction %v", transaction.LastDeviceKey())
	fmt.Printf("Transaction %v", transaction.NewDeviceKey())
}

func TestErrorUnmarshal(t *testing.T) {
	bData := common.FromHex(
		"0805120e696e76616c696420616d6f756e74",
	)

	transactionErr := &TransactionError{}
	err := transactionErr.Unmarshal(bData)
	if err != nil {
		logger.Error("Error when unmarshal ", err)
	}
	fmt.Printf("Transaction err %v", transactionErr)
}

func TestUnmarshalTransactionHashData(t *testing.T) {
	bData := common.FromHex(
		"0a20000000000000000000000000000000000000000000000000000000000000000012308f4f64095a726ec28a2854ed287c4e1a91104881f10438bb44e6add15228621ac5b7be6a1ca2fd6cc0591ce1c3c8e0bb1a14ba00cfe5a5697c1f7e6d7243a6608f69e7a0a99c222000000000000000000000000000000000000000000000000000000000000000002a20000000000000000000000000000000000000000000000000000000000000000030c096b1023880c8afa02540e807480452260a245e6c6c0100000000000000000000000000000000000000000000000000000000000000005a1477c51d415d49186325daf00ec8b20fe1456919086230d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d34d346a30e75d04e04efbd3cdbc0c3045ec5ec1d34001d3403d17a003005f350b40c2f420bce45d45f36e3d0b6e7af78d83eb50fd",
	)
	txHashData := &pb.TransactionHashData{}
	err := proto.Unmarshal(bData, txHashData)
	if err != nil {
		logger.Error("Error when unmarshal ", err)
		return
	}
	logger.Info("txHashData XX",
		hex.EncodeToString(txHashData.LastHash),
		hex.EncodeToString(txHashData.PublicKey),
		hex.EncodeToString(txHashData.ToAddress),
		hex.EncodeToString(txHashData.PendingUse),
		hex.EncodeToString(txHashData.Amount),
		txHashData.MaxGas,
		txHashData.MaxGasPrice,
		txHashData.MaxTimeUse,
		txHashData.Action,
		hex.EncodeToString(txHashData.Data),
		hex.EncodeToString(txHashData.RelatedAddresses[0]),
		hex.EncodeToString(txHashData.LastDeviceKey),
		hex.EncodeToString(txHashData.NewDeviceKey),
	)
	// LastHash         []byte   `protobuf:"bytes,1,opt,name=LastHash,proto3" json:"LastHash,omitempty"`
	// PublicKey        []byte   `protobuf:"bytes,2,opt,name=PublicKey,proto3" json:"PublicKey,omitempty"`
	// ToAddress        []byte   `protobuf:"bytes,3,opt,name=ToAddress,proto3" json:"ToAddress,omitempty"`
	// PendingUse       []byte   `protobuf:"bytes,4,opt,name=PendingUse,proto3" json:"PendingUse,omitempty"`
	// Amount           []byte   `protobuf:"bytes,5,opt,name=Amount,proto3" json:"Amount,omitempty"`
	// MaxGas           uint64   `protobuf:"varint,6,opt,name=MaxGas,proto3" json:"MaxGas,omitempty"`
	// MaxGasPrice      uint64   `protobuf:"varint,7,opt,name=MaxGasPrice,proto3" json:"MaxGasPrice,omitempty"`
	// MaxTimeUse       uint64   `protobuf:"varint,8,opt,name=MaxTimeUse,proto3" json:"MaxTimeUse,omitempty"`
	// Action           ACTION   `protobuf:"varint,9,opt,name=Action,proto3,enum=transaction.ACTION" json:"Action,omitempty"`
	// Data             []byte   `protobuf:"bytes,10,opt,name=Data,proto3" json:"Data,omitempty"`
	// RelatedAddresses [][]byte `protobuf:"bytes,11,rep,name=RelatedAddresses,proto3" json:"RelatedAddresses,omitempty"`
	// LastDeviceKey    []byte   `protobuf:"bytes,12,opt,name=LastDeviceKey,proto3" json:"LastDeviceKey,omitempty"` // hash last transaction deviceKey
	// NewDeviceKey     []byte   `protobuf:"bytes,13,opt,name=NewDeviceKey,proto3" json:"NewDeviceKey,omitempty"`   // hash of hash for new deviceKey
	logger.Info("txHashData", crypto.Keccak256Hash(bData))
}

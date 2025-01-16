package services

import (
	"encoding/hex"
	"fmt"
	"math/big"

	pb "gomail/mtn/proto"
	"gomail/mtn/transaction"
	"gomail/mtn/logger"

	"gomail/cmd/client"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	e_common "github.com/ethereum/go-ethereum/common"
	// "gomail/emailstorage"
	rc "gomail/mtn/receipt"
	"gomail/mtn/types"
	"gomail/emailstorage"
)

type SendTransactionService interface {
	CreateEmail(
		mailStorageAdd common.Address,
		sender string,
		subject string,
		body string,
		fileKeys [][32]byte,
		createdAt uint64,
		discription string,
	) (interface{}, error)
	GetEmailStorage(add string) (interface{}, error)
	PushFileInfos(
		infos []emailstorage.Info ,
	) (interface{}, error) 
	UploadChunk(
		fileKey [32]byte,
		chunkData []byte,
		chunkHash [32]byte,
	) (interface{}, error) 
}

type sendTransactionService struct {
	chainClient        *client.Client
	mailFactoryAbi     *abi.ABI
	mailFactoryAddress e_common.Address
	mailStorageAbi     *abi.ABI
	mailStorageAddress e_common.Address
	notiAddress			e_common.Address
	adminAddress e_common.Address
	fileAbi *abi.ABI
	fileAddress e_common.Address
}

func NewSendTransactionService(
	chainClient *client.Client,
	mailFactoryAbi *abi.ABI,
	mailFactoryAddress e_common.Address,
	mailStorageAbi *abi.ABI,
	mailStorageAddress e_common.Address,
	notiAddress		e_common.Address,
	adminAddress e_common.Address,
	fileAbi *abi.ABI,
	fileAddress e_common.Address,
) SendTransactionService {
	return &sendTransactionService{
		chainClient:        chainClient,
		mailFactoryAbi:     mailFactoryAbi,
		mailFactoryAddress: mailFactoryAddress,
		mailStorageAbi:     mailStorageAbi,
		mailStorageAddress: mailStorageAddress,
		notiAddress:        notiAddress,
		adminAddress:adminAddress,
		fileAbi: fileAbi,
		fileAddress: fileAddress,
	}
}

func (h *sendTransactionService) CreateEmail(
	mailStorageAdd common.Address,
	sender string,
	subject string,
	body string,
	fileKeys [][32]byte,
	createdAt uint64,
	discription string,
) (interface{}, error) {
	var result interface{}
	input, err := h.mailStorageAbi.Pack(
		"createEmail",
		sender,
		subject,
		body,
		fileKeys,
		createdAt,
		discription,
	)
	if err != nil {
		logger.Error("error when pack call data createEmail", err)
		return nil, err
	}
	callData := transaction.NewCallData(input)

	bData, err := callData.Marshal()
	if err != nil {
		logger.Error("error when marshal call data createEmail", err)
		return nil, err
	}
	fmt.Println("input: ", hex.EncodeToString(bData))
	relatedAddress := []e_common.Address{
		h.notiAddress,
		h.mailStorageAddress,
		h.adminAddress,
		common.HexToAddress("0xC09459d3f3B58597A5f48E60849744006E60bB10"),
		common.HexToAddress("0xF9d620dae6B224578A4eae25b94B7C3c85ad487B"),
	}
	maxGas := uint64(5_000_000)
	maxGasPrice := uint64(1_000_000_000)
	timeUse := uint64(0)
	receipt, err := h.chainClient.SendTransactionWithDeviceKey(
		mailStorageAdd,
		big.NewInt(0),
		4,
		bData,
		relatedAddress,
		maxGas,
		maxGasPrice,
		timeUse,
	)
	fmt.Println("rc createEmail:", receipt)
	if receipt.Status() == pb.RECEIPT_STATUS_RETURNED {
		kq := make(map[string]interface{})
		err = h.mailStorageAbi.UnpackIntoMap(kq, "createEmail", receipt.Return())
		if err != nil {
			logger.Error("UnpackIntoMap")
			return nil, err
		}
		var receipts []types.Receipt
		receipts = append(receipts,receipt )
		receiptsProto := rc.ReceiptsToProto(receipts)
		result = common.BytesToHash(receiptsProto[0].TransactionHash)
		logger.Info("CreateEmail - Result - ", kq)
	} else {
		result = hex.EncodeToString(receipt.Return())
		logger.Info("CreateEmail - Result - ", result)

	}
	return result, nil
}
func (h *sendTransactionService) GetEmailStorage(add string) (interface{}, error) {
	var result interface{}

	input, err := h.mailFactoryAbi.Pack(
		"getEmailStorageBySender",
		common.HexToAddress(add),
	)
	if err != nil {
		logger.Error("error when pack call data GetEmailStorage", err)
		return nil, err
	}
	callData := transaction.NewCallData(input)

	bData, err := callData.Marshal()
	if err != nil {
		logger.Error("error when marshal call data GetEmailStorage", err)
		return nil, err
	}
	fmt.Println("input getEmailStorageBySender: ", hex.EncodeToString(bData))
	relatedAddress := []e_common.Address{}
	maxGas := uint64(5_000_000)
	maxGasPrice := uint64(1_000_000_000)
	timeUse := uint64(0)
	receipt, err := h.chainClient.SendTransactionWithDeviceKey(
		h.mailFactoryAddress,
		big.NewInt(0),
		4,
		bData,
		relatedAddress,
		maxGas,
		maxGasPrice,
		timeUse,
	)
	fmt.Println("rc:", receipt)
	if receipt.Status() == pb.RECEIPT_STATUS_RETURNED {
		var kq common.Address
		err = h.mailFactoryAbi.UnpackIntoInterface(&kq, "getEmailStorageBySender", receipt.Return())
		if err != nil {
			logger.Error("UnpackIntoMap")
			return nil, err
		}
		result = kq
		logger.Info("GetEmailStorage - Result - ", result)
	} else {
		result = hex.EncodeToString(receipt.Return())
		logger.Info("GetEmailStorage - Result - ", result)

	}
	return result, nil
}
func (h *sendTransactionService) PushFileInfos(
	infos []emailstorage.Info ,
) (interface{}, error) {
	var result interface{}
	input, err := h.fileAbi.Pack(
		"pushFileInfos",
		infos,
	)
	if err != nil {
		logger.Error("error when pack call data PushFileInfos", err)
		return nil, err
	}
	callData := transaction.NewCallData(input)

	bData, err := callData.Marshal()
	if err != nil {
		logger.Error("error when marshal call data PushFileInfos", err)
		return nil, err
	}
	fmt.Println("input: ", hex.EncodeToString(bData))
	relatedAddress := []e_common.Address{
	}
	maxGas := uint64(5_000_000)
	maxGasPrice := uint64(1_000_000_000)
	timeUse := uint64(0)
	receipt, err := h.chainClient.SendTransactionWithDeviceKey(
		h.fileAddress,
		big.NewInt(0),
		4,
		bData,
		relatedAddress,
		maxGas,
		maxGasPrice,
		timeUse,
	)
	fmt.Println("rc pushFileInfo:", receipt)
	if receipt.Status() == pb.RECEIPT_STATUS_RETURNED {
		var kq [][32]byte
		err = h.fileAbi.UnpackIntoInterface(&kq, "pushFileInfos", receipt.Return())
		if err != nil {
			logger.Error("UnpackIntoInterface pushFileInfo")
			return nil, err
		}
		result = kq
		logger.Info("PushFileInfo - Result - ", kq)
	} else {
		result = hex.EncodeToString(receipt.Return())
		logger.Info("PushFileInfo - Result - ", result)

	}
	return result, nil
}
func (h *sendTransactionService) UploadChunk(
	fileKey [32]byte,
	chunkData []byte,
	chunkHash [32]byte,
) (interface{}, error) {
	var result interface{}
	input, err := h.fileAbi.Pack(
		"uploadChunk",
		fileKey,
		chunkData,
		chunkHash,
	)
	if err != nil {
		logger.Error("error when pack call data uploadChunk", err)
		return nil, err
	}
	callData := transaction.NewCallData(input)

	bData, err := callData.Marshal()
	if err != nil {
		logger.Error("error when marshal call data uploadChunk", err)
		return nil, err
	}
	fmt.Println("input: ", hex.EncodeToString(bData))
	relatedAddress := []e_common.Address{
	}
	maxGas := uint64(5_000_000)
	maxGasPrice := uint64(1_000_000_000)
	timeUse := uint64(0)
	receipt, err := h.chainClient.SendTransactionWithDeviceKey(
		h.fileAddress,
		big.NewInt(0),
		4,
		bData,
		relatedAddress,
		maxGas,
		maxGasPrice,
		timeUse,
	)
	fmt.Println("rc uploadChunk:", receipt)
	if receipt.Status() == pb.RECEIPT_STATUS_RETURNED {
		var kq bool
		err = h.fileAbi.UnpackIntoInterface(&kq, "uploadChunk", receipt.Return())
		if err != nil {
			logger.Error("UnpackIntoInterface uploadChunk")
			return nil, err
		}
		result = kq
		logger.Info("UploadChunk - Result - ", kq)
	} else {
		result = hex.EncodeToString(receipt.Return())
		logger.Info("UploadChunk - Result - ", result)

	}
	return result, nil
}

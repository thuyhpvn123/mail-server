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
	"gomail/emailstorage"
	rc "gomail/mtn/receipt"
	"gomail/mtn/types"
)

type SendTransactionService interface {
	CreateEmail(
		mailStorageAdd common.Address,
		sender string,
		subject string,
		fromHeader string,
		replyTo string,
		messageId string,
		body string,
		html string,
		fileData []emailstorage.File,
		createdAt *big.Int,
	) (interface{}, error)
	GetEmailStorage(add string) (interface{}, error)
}

type sendTransactionService struct {
	chainClient        *client.Client
	mailFactoryAbi     *abi.ABI
	mailFactoryAddress e_common.Address
	mailStorageAbi     *abi.ABI
	notiAddress			e_common.Address
}

func NewSendTransactionService(
	chainClient *client.Client,
	mailFactoryAbi *abi.ABI,
	mailFactoryAddress e_common.Address,
	mailStorageAbi *abi.ABI,
	notiAddress		e_common.Address,
) SendTransactionService {
	return &sendTransactionService{
		chainClient:        chainClient,
		mailFactoryAbi:     mailFactoryAbi,
		mailFactoryAddress: mailFactoryAddress,
		mailStorageAbi:     mailStorageAbi,
		notiAddress:        notiAddress,
	}
}

func (h *sendTransactionService) CreateEmail(
	mailStorageAdd common.Address,
	sender string,
	subject string,
	fromHeader string,
	replyTo string,
	messageId string,
	body string,
	html string,
	fileData []emailstorage.File,
	createdAt *big.Int,
) (interface{}, error) {
	var result interface{}
	input, err := h.mailStorageAbi.Pack(
		"createEmail",
		sender,
		subject,
		fromHeader,
		replyTo,
		messageId,
		body,
		html,
		fileData,
		createdAt,
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
	relatedAddress := []e_common.Address{h.notiAddress}
	maxGas := uint64(5_000_000)
	maxGasPrice := uint64(1_000_000_000)
	timeUse := uint64(0)
	receipt, err := h.chainClient.SendTransaction(
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
	receipt, err := h.chainClient.SendTransaction(
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

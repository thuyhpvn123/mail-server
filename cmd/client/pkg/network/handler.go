package network

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"gomail/cmd/client/command"
	"gomail/mtn/logger"
	pb "gomail/mtn/proto"
	"gomail/mtn/receipt"
	"gomail/mtn/smart_contract"
	"gomail/mtn/state"
	"gomail/mtn/stats"
	"gomail/mtn/transaction"
	"gomail/mtn/types"
	"gomail/mtn/types/network"
)

var ErrorCommandNotFound = errors.New("command not found")

type Handler struct {
	accountStateChan     chan types.AccountState
	receiptChan          chan types.Receipt
	eventLogChan         chan types.EventLogs
	transactionErrorChan chan types.TransactionError
}

func NewHandler(
	accountStateChan chan types.AccountState,
	receiptChan chan types.Receipt,
	transactionErrorChan chan types.TransactionError,
) *Handler {
	return &Handler{
		accountStateChan:     accountStateChan,
		receiptChan:          receiptChan,
		transactionErrorChan: transactionErrorChan,
	}
}

func (h *Handler) HandleRequest(request network.Request) (err error) {
	cmd := request.Message().Command()
	logger.Debug("handling command: " + cmd)
	switch cmd {
	case command.InitConnection:
		return h.handleInitConnection(request)
	case command.AccountState:
		return h.handleAccountState(request)
	case command.TransactionError:
		transactionError := &transaction.TransactionHashWithErrorCode{}
		err = transactionError.Unmarshal(request.Message().Body())
		logger.Debug("Receive Transaction error: ", transactionError)
		return nil
	case command.Receipt:
		return h.handleReceipt(request)
	case command.EventLogs:
		return h.handleEventLogs(request)
	case command.QueryLogs:
		return h.handleEventLogs(request)
	case command.Stats:
		return h.handleStats(request)
	}
	return ErrorCommandNotFound
}

func (h *Handler) SetEventLogsChan(ch chan types.EventLogs) {
	h.eventLogChan = ch
}

func (h *Handler) GetEventLogsChan() chan types.EventLogs {
	return h.eventLogChan
}

/*
handleInitConnection will receive request from connection
then init that connection with data in request then
add it to connection manager
*/
func (h *Handler) handleInitConnection(request network.Request) (err error) {
	conn := request.Connection()
	initData := &pb.InitConnection{}
	err = request.Message().Unmarshal(initData)
	if err != nil {
		return err
	}
	address := common.BytesToAddress(initData.Address)
	logger.Debug(fmt.Sprintf(
		"init connection from %v type %v", address, initData.Type,
	))
	conn.Init(address, initData.Type)
	return nil
}

/*
handleAccountState will receive account state from connection
then push it to account state chan
*/
func (h *Handler) handleAccountState(request network.Request) (err error) {
	accountState := &state.AccountState{}
	err = accountState.Unmarshal(request.Message().Body())
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("Receive Account state: \n%v", accountState))
	h.accountStateChan <- accountState
	return nil
}

/*
handleAccountState will receive receipt from connection
then print it out
*/
func (h *Handler) handleReceipt(request network.Request) (err error) {
	receipt := &receipt.Receipt{}
	err = receipt.Unmarshal(request.Message().Body())
	if err != nil {
		return err
	}
	if h.receiptChan != nil {
		h.receiptChan <- receipt
	} else {
		logger.Debug(fmt.Sprintf("Receive receipt: %v", receipt))
		logger.Debug(fmt.Sprintf("Receive To address: %v", request.Message().ToAddress()))
		if receipt.Status() == pb.RECEIPT_STATUS_TRANSACTION_ERROR {
			transactionErr := &transaction.TransactionError{}
			transactionErr.Unmarshal(receipt.Return())
			logger.Debug("Receive Transaction error: ", transactionErr)
		}
	}
	return nil
}

/*
handleTransactionError will receive transaction error from parent node connection
then print it out
*/
func (h *Handler) handleEventLogs(request network.Request) error {
	eventLogs := &smart_contract.EventLogs{}
	err := eventLogs.Unmarshal(request.Message().Body())
	if err != nil {
		logger.Error("Handle Event Logs Error", err)
		return err
	}
	eventLogList := eventLogs.EventLogList()
	for _, eventLog := range eventLogList {
		logger.Debug("EventLogs: ", eventLog.String())
	}
	if h.eventLogChan != nil {
		h.eventLogChan <- eventLogs
	}
	return nil
}

/*
handleStats will receive stats from connection
then print it our
*/
func (h *Handler) handleStats(request network.Request) (err error) {
	stats := &stats.Stats{}
	err = stats.Unmarshal(request.Message().Body())
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Receive Stats: \n%v", stats))
	return nil
}

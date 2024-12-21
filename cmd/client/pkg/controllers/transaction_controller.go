package controllers

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"gomail/cmd/client/command"
	"gomail/cmd/client/pkg/client_context"
	client_types "gomail/cmd/client/types"
	p_common "gomail/mtn/common"
	pb "gomail/mtn/proto"
	"gomail/mtn/transaction"
	t "gomail/mtn/transaction"
	"gomail/mtn/types"
)

type TransactionController struct {
	clientContext *client_context.ClientContext
}

func NewTransactionController(
	clientContext *client_context.ClientContext,
) client_types.TransactionController {
	return &TransactionController{
		clientContext: clientContext,
	}
}

func (tc *TransactionController) SendTransaction(
	lastHash common.Hash,
	toAddress common.Address,
	pendingUse *big.Int,
	amount *big.Int,
	maxGas uint64,
	maxGasFee uint64,
	maxTimeUse uint64,
	action pb.ACTION,
	data []byte,
	relatedAddress [][]byte,
	lastDeviceKey common.Hash,
	newDeviceKey common.Hash,
	commissionPrivateKey []byte,
) (types.Transaction, error) {
	transaction := t.NewTransaction(
		lastHash,
		tc.clientContext.KeyPair.PublicKey(),
		toAddress,
		pendingUse,
		amount,
		maxGas,
		maxGasFee,
		maxTimeUse,
		action,
		data,
		relatedAddress,
		lastDeviceKey,
		newDeviceKey,
	)
	transaction.SetSign(tc.clientContext.KeyPair.PrivateKey())
	if commissionPrivateKey != nil {
		transaction.SetCommissionSign(p_common.PrivateKeyFromBytes(commissionPrivateKey))
	}
	bTransaction, err := transaction.Marshal()
	if err != nil {
		return nil, err
	}
	parentConnection := tc.clientContext.ConnectionsManager.ParentConnection()
	err = tc.clientContext.MessageSender.SendBytes(
		parentConnection,
		command.SendTransaction,
		bTransaction,
	)
	return transaction, err
}

func (tc *TransactionController) SendTransactions(
	transactions []types.Transaction,
) error {

	bTransaction, err := transaction.MarshalTransactions(transactions)
	if err != nil {
		return err
	}
	parentConnection := tc.clientContext.ConnectionsManager.ParentConnection()
	err = tc.clientContext.MessageSender.SendBytes(
		parentConnection,
		command.SendTransactions,
		bTransaction,
	)
	return err
}

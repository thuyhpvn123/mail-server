package controllers

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"gomail/cmd/client/command"
	"gomail/cmd/client/pkg/client_context"
	client_types "gomail/cmd/client/types"
	p_common "gomail/pkg/common"
	pb "gomail/pkg/proto"
	"gomail/pkg/transaction"
	"gomail/types"

	"fmt" // For formatted error messages

	"google.golang.org/protobuf/proto" // For marshaling protobuf messages
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
	fromAddress common.Address,
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
	nonce uint64,
	chainId uint64,
) (types.Transaction, error) {
	transaction := transaction.NewTransaction(
		lastHash,
		fromAddress,
		toAddress,
		pendingUse,
		amount,
		maxGas,
		maxGasFee,
		maxTimeUse,
		data,
		relatedAddress,
		lastDeviceKey,
		newDeviceKey,
		nonce,
		chainId,
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

func (tc *TransactionController) ReadTransaction(
	lastHash common.Hash,
	fromAddress common.Address,
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
	nonce uint64,
	chainId uint64,
) (types.Transaction, error) {
	transaction := transaction.NewTransaction(
		lastHash,
		fromAddress,
		toAddress,
		pendingUse,
		amount,
		maxGas,
		maxGasFee,
		maxTimeUse,
		data,
		relatedAddress,
		lastDeviceKey,
		newDeviceKey,
		nonce,
		chainId,
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
		command.ReadTransaction,
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

func (tc *TransactionController) SendTransactionWithDeviceKey(
	lastHash common.Hash,
	fromAddress common.Address,
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
	nonce uint64,
	deviceKey []byte,
	chainId uint64,
) (types.Transaction, error) {
	transaction := transaction.NewTransaction(
		lastHash,
		fromAddress,
		toAddress,
		pendingUse,
		amount,
		maxGas,
		maxGasFee,
		maxTimeUse,
		data,
		relatedAddress,
		lastDeviceKey,
		newDeviceKey,
		nonce,
		chainId,
	)
	transaction.SetSign(tc.clientContext.KeyPair.PrivateKey())
	if commissionPrivateKey != nil {
		transaction.SetCommissionSign(p_common.PrivateKeyFromBytes(commissionPrivateKey))
	}

	// Create TransactionWithDeviceKey
	transactionWithDeviceKey := &pb.TransactionWithDeviceKey{
		Transaction: transaction.Proto().(*pb.Transaction),
		DeviceKey:   deviceKey,
	}

	// Serialize to bytes
	bTransactionWithDeviceKey, err := proto.Marshal(transactionWithDeviceKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TransactionWithDeviceKey: %w", err)
	}
	// dataB := common.FromHex("0a208dbce4440e64fef64a0a907b550da8e6079b81394cb226b4bb918712b1d0b07a121426d209379611be4829eede2d20232d9cfc7ef7f41a07470de4df8200002220000000000000000000000000000000000000000000000000002386f26fc1000028a09c0130c0843d38f02e5a20290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e5636220bd9d91783ae8bbc7a3e69008772e266d2ff3eded1b09265942e9d43d889c31486a08000000000000000172142f4cb880116850929d8b44fac82e907bc21f19d0")
	// unmarshalledTransactionWithDeviceKey := &pb.Transaction{}
	// err = proto.Unmarshal(dataB, unmarshalledTransactionWithDeviceKey)

	// logger.Error(hex.EncodeToString(unmarshalledTransactionWithDeviceKey.Hash))
	// logger.Error(hex.EncodeToString(unmarshalledTransactionWithDeviceKey.LastHash))

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal TransactionWithDeviceKey: %w", err)
	}
	parentConnection := tc.clientContext.ConnectionsManager.ParentConnection()
	err = tc.clientContext.MessageSender.SendBytes(
		parentConnection,
		command.SendTransactionWithDeviceKey,
		bTransactionWithDeviceKey,
	)
	return transaction, err
}

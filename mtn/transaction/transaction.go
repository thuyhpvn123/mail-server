package transaction

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"gomail/mtn/bls"
	p_common "gomail/mtn/common"
	"gomail/mtn/logger"
	pb "gomail/mtn/proto"
	"gomail/mtn/types"
)

type Transaction struct {
	proto *pb.Transaction
}

func NewTransaction(
	lastHash common.Hash,
	publicKey p_common.PublicKey,
	toAddress common.Address,
	pendingUse *big.Int,
	amount *big.Int,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
	action pb.ACTION,
	data []byte,
	relatedAddresses [][]byte,
	lastDeviceKey common.Hash,
	newDeviceKey common.Hash,
) types.Transaction {
	proto := &pb.Transaction{
		LastHash:         lastHash.Bytes(),
		PublicKey:        publicKey.Bytes(),
		ToAddress:        toAddress.Bytes(),
		PendingUse:       pendingUse.Bytes(),
		Amount:           amount.Bytes(),
		MaxGas:           maxGas,
		MaxGasPrice:      maxGasPrice,
		MaxTimeUse:       maxTimeUse,
		Action:           action,
		Data:             data,
		RelatedAddresses: relatedAddresses,
		LastDeviceKey:    lastDeviceKey.Bytes(),
		NewDeviceKey:     newDeviceKey.Bytes(),
	}
	tx := &Transaction{
		proto: proto,
	}
	tx.SetHash(tx.CalculateHash())
	return tx
}

func TransactionsToProto(transactions []types.Transaction) []*pb.Transaction {
	rs := make([]*pb.Transaction, len(transactions))
	for i, v := range transactions {
		rs[i] = v.Proto().(*pb.Transaction)
	}
	return rs
}

func TransactionFromProto(txPb *pb.Transaction) types.Transaction {
	return &Transaction{
		proto: txPb,
	}
}

func TransactionsFromProto(pbTxs []*pb.Transaction) []types.Transaction {
	rs := make([]types.Transaction, len(pbTxs))
	for i, v := range pbTxs {
		rs[i] = TransactionFromProto(v)
	}
	return rs
}

// general

func (t *Transaction) Unmarshal(b []byte) error {
	pbTransaction := &pb.Transaction{}
	err := proto.Unmarshal(b, pbTransaction)
	if err != nil {
		return err
	}
	t.FromProto(pbTransaction)
	return nil
}

func (t *Transaction) Marshal() ([]byte, error) {
	return proto.Marshal(t.proto)
}

func (t *Transaction) Proto() protoreflect.ProtoMessage {
	return t.proto
}

func (t *Transaction) FromProto(pbMessage protoreflect.ProtoMessage) {
	pbTransaction := pbMessage.(*pb.Transaction)
	t.proto = pbTransaction
}

func (t *Transaction) String() string {
	str := fmt.Sprintf(`
	Hash: %v
	From: %v
	To: %v
	Amount: %v
	Action: %v
	Data: %v
	Last Hash: %v
	Max Gas: %v
	Max Gas Price: %v
	Max Time Use: %v
	Sign: %v
	Commission Sign:  %v
	Pending Use: %v
`,
		hex.EncodeToString(t.proto.Hash),
		hex.EncodeToString(t.FromAddress().Bytes()),
		hex.EncodeToString(t.proto.ToAddress),
		big.NewInt(0).SetBytes(t.proto.Amount),
		t.proto.Action,
		hex.EncodeToString(t.proto.Data),
		t.LastHash(),
		t.MaxGas(),
		t.MaxGasPrice(),
		t.MaxTimeUse(),
		hex.EncodeToString(t.proto.Sign),
		hex.EncodeToString(t.proto.CommissionSign),
		hex.EncodeToString(t.proto.PendingUse),
	)
	return str
}

// getter
func (t *Transaction) Hash() common.Hash {
	return common.BytesToHash(t.proto.Hash)
}

func (t *Transaction) NewDeviceKey() common.Hash {
	return common.BytesToHash(t.proto.NewDeviceKey)
}

func (t *Transaction) LastDeviceKey() common.Hash {
	return common.BytesToHash(t.proto.LastDeviceKey)
}

func (t *Transaction) FromAddress() common.Address {
	return common.BytesToAddress(
		crypto.Keccak256(t.proto.PublicKey),
	)
}

func (t *Transaction) ToAddress() common.Address {
	return common.BytesToAddress(t.proto.ToAddress)
}

func (t *Transaction) Pubkey() p_common.PublicKey {
	return p_common.PubkeyFromBytes(t.proto.PublicKey)
}

func (t *Transaction) LastHash() common.Hash {
	return common.BytesToHash(t.proto.LastHash)
}

func (t *Transaction) Sign() p_common.Sign {
	return p_common.SignFromBytes(t.proto.Sign)
}

func (t *Transaction) Amount() *big.Int {
	return big.NewInt(0).SetBytes(t.proto.Amount)
}

func (t *Transaction) PendingUse() *big.Int {
	return big.NewInt(0).SetBytes(t.proto.PendingUse)
}

func (t *Transaction) Action() pb.ACTION {
	return t.proto.Action
}

func (t *Transaction) BRelatedAddresses() [][]byte {
	return t.proto.RelatedAddresses
}

func (t *Transaction) RelatedAddresses() []common.Address {
	relatedAddresses := make([]common.Address, len(t.proto.RelatedAddresses)+1)
	for i, v := range t.proto.RelatedAddresses {
		relatedAddresses[i] = common.BytesToAddress(v)
	}
	// append to address
	relatedAddresses[len(t.proto.RelatedAddresses)] = t.ToAddress()
	return relatedAddresses
}

func (t *Transaction) Fee(currentGasPrice uint64) *big.Int {
	fee := big.NewInt(int64(t.proto.MaxGas))
	fee = fee.Mul(fee, big.NewInt(int64(currentGasPrice)))
	fee = fee.Mul(fee, big.NewInt(int64((t.proto.MaxTimeUse/1000)+1.0)))
	return fee
}

func (t *Transaction) Data() []byte {
	return t.proto.Data
}

func (t *Transaction) DeployData() types.DeployData {
	deployData := &DeployData{}
	deployData.Unmarshal(t.Data())
	return deployData
}

func (t *Transaction) CallData() types.CallData {
	callData := &CallData{}
	callData.Unmarshal(t.Data())
	return callData
}

func (t *Transaction) OpenStateChannelData() types.OpenStateChannelData {
	openData := &OpenStateChannelData{}
	openData.Unmarshal(t.Data())
	return openData
}

func (t *Transaction) UpdateStorageHostData() types.UpdateStorageHostData {
	data := &UpdateStorageHostData{}
	data.Unmarshal(t.Data())
	return data
}

func (t *Transaction) CommissionSign() p_common.Sign {
	return p_common.SignFromBytes(t.proto.CommissionSign)
}

func (t *Transaction) MaxGas() uint64 {
	return t.proto.MaxGas
}

func (t *Transaction) MaxGasPrice() uint64 {
	return t.proto.MaxGasPrice
}

func (t *Transaction) MaxFee() *big.Int {
	return big.NewInt(0).Mul(
		big.NewInt(int64(t.MaxGasPrice())),
		big.NewInt(int64(t.MaxGas())),
	)
}

func (t *Transaction) MaxTimeUse() uint64 {
	return t.proto.MaxTimeUse
}

// setter

func (t *Transaction) CalculateHash() common.Hash {
	hashPb := &pb.TransactionHashData{
		LastHash:         t.proto.LastHash,
		PublicKey:        t.proto.PublicKey,
		ToAddress:        t.proto.ToAddress,
		PendingUse:       t.proto.PendingUse,
		Amount:           t.proto.Amount,
		MaxGas:           t.proto.MaxGas,
		MaxGasPrice:      t.proto.MaxGasPrice,
		MaxTimeUse:       t.proto.MaxTimeUse,
		Action:           t.proto.Action,
		Data:             t.proto.Data,
		RelatedAddresses: t.proto.RelatedAddresses,
		LastDeviceKey:    t.proto.LastDeviceKey,
		NewDeviceKey:     t.proto.NewDeviceKey,
	}
	bHashPb, _ := proto.Marshal(hashPb)
	// logger.Info("CalculateHash", hex.EncodeToString(bHashPb))
	return crypto.Keccak256Hash(bHashPb)
}

func (t *Transaction) SetHash(hash common.Hash) {
	t.proto.Hash = hash.Bytes()
}

func (t *Transaction) SetSign(privateKey p_common.PrivateKey) {
	t.proto.Sign = bls.Sign(privateKey, t.proto.Hash).Bytes()
}

func (t *Transaction) SetCommissionSign(privateKey p_common.PrivateKey) {
	t.proto.CommissionSign = bls.Sign(privateKey, t.proto.Hash).Bytes()
}

// validate

func (t *Transaction) ValidSign() bool {
	return bls.VerifySign(
		t.Pubkey(),
		t.Sign(),
		t.Hash().Bytes(),
	)
}

func (t *Transaction) ValidTransactionHash() bool {
	return t.CalculateHash() == t.Hash()
}

func (t *Transaction) ValidLastHash(fromAccountState types.AccountState) bool {
	return t.LastHash() == fromAccountState.LastHash()
}

func (t *Transaction) ValidDeviceKey(fromAccountState types.AccountState) bool {
	return fromAccountState.DeviceKey() == common.Hash{} || // skip check device key if account state doesn't have device key
		crypto.Keccak256Hash(t.LastDeviceKey().Bytes()) == fromAccountState.DeviceKey()
}

func (t *Transaction) ValidMaxGas() bool {
	if t.Action() == pb.ACTION_OPEN_CHANNEL {
		return t.MaxGas() >= p_common.OPEN_CHANNEL_GAS_COST
	}
	return t.MaxGas() >= p_common.TRANSFER_GAS_COST
}

func (t *Transaction) ValidMaxGasPrice(currentGasPrice uint64) bool {
	if t.ToAddress() == p_common.NATIVE_SMART_CONTRACT_REWARD_ADDRESS &&
		t.Action() == pb.ACTION_CALL_SMART_CONTRACT {
		// skip check gas price for native smart contract
		return true
	}
	return currentGasPrice <= t.MaxGasPrice()
}

func (t *Transaction) ValidAmountSpend(
	fromAccountState types.AccountState,
	spendAmount *big.Int,
) bool {
	totalBalance := big.NewInt(0).Add(fromAccountState.Balance(), t.PendingUse())
	totalSpend := big.NewInt(0).Add(spendAmount, t.Amount())
	return totalBalance.Cmp(totalSpend) >= 0
}

func (t *Transaction) ValidAmount(
	fromAccountState types.AccountState,
	currentGasPrice uint64,
) bool {
	fee := t.Fee(currentGasPrice)
	return t.ValidAmountSpend(fromAccountState, fee)
}

func (t *Transaction) ValidPendingUse(fromAccountState types.AccountState) bool {
	pendingBalance := fromAccountState.PendingBalance()
	pendingUse := t.PendingUse()
	return pendingUse.Cmp(pendingBalance) <= 0
}

func (t *Transaction) ValidDeploySmartContractToAccount(fromAccountState types.AccountState) bool {
	validToAddress := common.BytesToAddress(
		crypto.Keccak256(
			append(
				fromAccountState.Address().Bytes(),
				fromAccountState.LastHash().Bytes()...),
		)[12:],
	)
	if validToAddress != t.ToAddress() {
		logger.Warn("Not match deploy address", validToAddress, t.ToAddress())
	}
	return validToAddress == t.ToAddress()
}

func (t *Transaction) ValidOpenChannelToAccount(fromAccountState types.AccountState) bool {
	validToAddress := common.BytesToAddress(
		crypto.Keccak256(
			append(
				fromAccountState.Address().Bytes(),
				fromAccountState.LastHash().Bytes()...),
		)[12:],
	)
	if validToAddress != t.ToAddress() {
		logger.Warn("Not match open channel address", validToAddress, t.ToAddress())
	}
	return validToAddress == t.ToAddress()
}

func (t *Transaction) ValidCallSmartContractToAccount(toAccountState types.AccountState) bool {
	scState := toAccountState.SmartContractState()
	return scState != nil
}

func MarshalTransactions(txs []types.Transaction) ([]byte, error) {
	return proto.Marshal(&pb.Transactions{Transactions: TransactionsToProto(txs)})
}

func UnmarshalTransactions(b []byte) ([]types.Transaction, error) {
	pbTxs := &pb.Transactions{}
	err := proto.Unmarshal(b, pbTxs)
	if err != nil {
		return nil, err
	}
	return TransactionsFromProto(pbTxs.Transactions), nil
}

func MarshalTransactionsWithBlockNumber(
	txs []types.Transaction,
	blockNumber uint64,
) ([]byte, error) {
	pbTxs := make([]*pb.Transaction, len(txs))
	for i, v := range txs {
		pbTxs[i] = v.Proto().(*pb.Transaction)
	}
	return proto.Marshal(&pb.TransactionsWithBlockNumber{
		Transactions: pbTxs,
		BlockNumber:  blockNumber,
	})
}

func UnmarshalTransactionsWithBlockNumber(b []byte) ([]types.Transaction, uint64, error) {
	pbTxs := &pb.TransactionsWithBlockNumber{}
	err := proto.Unmarshal(b, pbTxs)
	if err != nil {
		return nil, 0, err
	}
	rs := make([]types.Transaction, len(pbTxs.Transactions))
	for i, v := range pbTxs.Transactions {
		rs[i] = &Transaction{proto: v}
	}
	return rs, pbTxs.BlockNumber, nil
}

package transaction

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"gomail/pkg/bls"
	p_common "gomail/pkg/common"
	"gomail/pkg/logger"
	pb "gomail/pkg/proto"
	"gomail/pkg/utils"
	"gomail/types"

	e_types "github.com/ethereum/go-ethereum/core/types"
)

type Transaction struct {
	proto      *pb.Transaction
	cachedHash atomic.Pointer[common.Hash]
}

func NewTransaction(
	lastHash common.Hash,
	fromAddress common.Address,
	toAddress common.Address,
	pendingUse *big.Int,
	amount *big.Int,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
	data []byte,
	relatedAddresses [][]byte,
	lastDeviceKey common.Hash,
	newDeviceKey common.Hash,
	nonce uint64,
	chainId uint64,
) types.Transaction {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	proto := &pb.Transaction{
		LastHash:         lastHash.Bytes(),
		FromAddress:      fromAddress.Bytes(),
		ToAddress:        toAddress.Bytes(),
		PendingUse:       pendingUse.Bytes(),
		Amount:           amount.Bytes(),
		MaxGas:           maxGas,
		MaxGasPrice:      maxGasPrice,
		MaxTimeUse:       maxTimeUse,
		Data:             data,
		RelatedAddresses: relatedAddresses,
		LastDeviceKey:    lastDeviceKey.Bytes(),
		NewDeviceKey:     newDeviceKey.Bytes(),
		Nonce:            nonceBytes,
		ChainID:          chainId,
	}
	tx := &Transaction{
		proto: proto,
	}
	return tx // Return the *Transaction directly
}

func NewTransactionOffChain(
	lastHash common.Hash,
	fromAddress common.Address,
	toAddress common.Address,
	pendingUse *big.Int,
	amount *big.Int,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
	data []byte,
	relatedAddresses [][]byte,
	lastDeviceKey common.Hash,
	newDeviceKey common.Hash,
	nonce uint64,
) types.Transaction {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	proto := &pb.Transaction{
		LastHash:         lastHash.Bytes(),
		FromAddress:      fromAddress.Bytes(),
		ToAddress:        toAddress.Bytes(),
		PendingUse:       pendingUse.Bytes(),
		Amount:           amount.Bytes(),
		MaxGas:           maxGas,
		MaxGasPrice:      maxGasPrice,
		MaxTimeUse:       maxTimeUse,
		Data:             data,
		RelatedAddresses: relatedAddresses,
		LastDeviceKey:    lastDeviceKey.Bytes(),
		NewDeviceKey:     newDeviceKey.Bytes(),
		Nonce:            nonceBytes,
	}
	tx := &Transaction{
		proto: proto,
	}
	return tx // Return the *Transaction directly
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

func (t *Transaction) GetEthTx() *e_types.Transaction {
	input := "0x"
	if len(t.GetData()) != 0 {
		input = (*hexutil.Big)(big.NewInt(0).SetBytes(t.GetData())).String()
	}
	v, s, r := t.RawSignatureValues()
	baseJson := map[string]interface{}{
		// type is filled in by the test
		"chainId":  (*hexutil.Big)(big.NewInt(0).SetBytes(utils.Uint64ToBytes(t.proto.ChainID))).String(),
		"nonce":    (*hexutil.Big)(big.NewInt(0).SetBytes(t.proto.Nonce)).String(),
		"gas":      (*hexutil.Big)(big.NewInt(0).SetBytes(utils.Uint64ToBytes(t.proto.MaxGas))).String(),
		"gasPrice": (*hexutil.Big)(big.NewInt(0).SetBytes(utils.Uint64ToBytes(t.proto.MaxGasPrice))).String(),
		"value":    (*hexutil.Big)(big.NewInt(0).SetBytes(t.proto.Amount)).String(),
		"input":    input,
		"v":        (*hexutil.Big)(big.NewInt(0).SetBytes(v.Bytes())).String(),
		"r":        (*hexutil.Big)(big.NewInt(0).SetBytes(s.Bytes())).String(),
		"s":        (*hexutil.Big)(big.NewInt(0).SetBytes(r.Bytes())).String(),
	}
	if t.IsDeployContract() {
		baseJson["to"] = t.ToAddress().Hex()
	}
	// Marshal the JSON
	jsonBytes, err := json.Marshal(baseJson)
	if err != nil {
		return nil
	}
	// Unmarshal the tx
	tx := new(e_types.Transaction)

	err = tx.UnmarshalJSON(jsonBytes)
	if err != nil {
		return nil
	}
	return tx
}

func (t *Transaction) Unmarshal(b []byte) error {
	pbTransaction := &pb.Transaction{}
	err := proto.Unmarshal(b, pbTransaction)
	if err != nil {
		return err
	}
	t.FromProto(pbTransaction)
	return nil
}

// Kiểm tra giao dịch có phải là Deploy Contract không
func (t *Transaction) IsDeployContract() bool {
	return t.ToAddress() == (common.Address{})
}

// Kiểm tra giao dịch có phải là Call Contract không
func (t *Transaction) IsCallContract() bool {
	return t.ToAddress() != (common.Address{}) && len(t.Data()) > 0
}

// Kiểm tra giao dịch có phải là chuyển tiền thông thường không
func (t *Transaction) IsRegularTransaction() bool {
	return t.ToAddress() != (common.Address{}) && len(t.Data()) == 0
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

func (t *Transaction) CopyTransaction() types.Transaction {
	// Tạo một bản sao của protobuf message
	newProto := proto.Clone(t.proto).(*pb.Transaction)

	// Tạo một Transaction mới từ bản sao protobuf message
	newTx := &Transaction{
		proto: newProto,
	}

	return newTx
}

func (t *Transaction) String() string {
	str := fmt.Sprintf(`
	Hash: %v
	From: %v
	To: %v
	Amount: %v
	Data: %v
	Last Hash: %v
	Max Gas: %v
	Max Gas Price: %v
	Max Time Use: %v
	Sign: %v
	Commission Sign:  %v
	Pending Use: %v
`, t.Hash(),
		hex.EncodeToString(t.proto.FromAddress),
		hex.EncodeToString(t.proto.ToAddress),
		big.NewInt(0).SetBytes(t.proto.Amount),
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
	// Kiểm tra cache có giá trị hay chưa
	if cached := t.cachedHash.Load(); cached != nil {
		return *cached // Trả về giá trị đã cache nếu có
	}

	// Tạo dữ liệu để băm
	hashPb := &pb.TransactionHashData{
		LastHash:      t.proto.LastHash,
		FromAddress:   t.proto.FromAddress,
		ToAddress:     t.proto.ToAddress,
		PendingUse:    t.proto.PendingUse,
		Amount:        t.proto.Amount,
		MaxGas:        t.proto.MaxGas,
		MaxGasPrice:   t.proto.MaxGasPrice,
		MaxTimeUse:    t.proto.MaxTimeUse,
		Data:          t.proto.Data,
		LastDeviceKey: t.proto.LastDeviceKey,
		NewDeviceKey:  t.proto.NewDeviceKey,
		Nonce:         t.proto.Nonce,
	}

	bHashPb, _ := proto.Marshal(hashPb)

	// Tính giá trị băm
	hash := crypto.Keccak256Hash(bHashPb)

	// Lưu vào cache (atomic)
	t.cachedHash.Store(&hash)

	return hash
}

func (t *Transaction) GetNonce() uint64 {
	if len(t.proto.Nonce) == 8 { // Check for valid length
		return binary.BigEndian.Uint64(t.proto.Nonce)
	} else {
		return 0 // Or any default value
	}
}

func (t *Transaction) GetChainID() uint64 {
	return t.proto.ChainID
}

func (t *Transaction) NewDeviceKey() common.Hash {
	return common.BytesToHash(t.proto.NewDeviceKey)
}

func (t *Transaction) LastDeviceKey() common.Hash {
	return common.BytesToHash(t.proto.LastDeviceKey)
}

func (t *Transaction) FromAddress() common.Address {
	return common.BytesToAddress(t.proto.FromAddress)
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

func (t *Transaction) BRelatedAddresses() [][]byte {
	return t.proto.RelatedAddresses
}

func (t *Transaction) UpdateRelatedAddresses(relatedAddresses [][]byte) {
	t.proto.RelatedAddresses = relatedAddresses
}

func (t *Transaction) SetReadOnly(readOnly bool) {
	t.proto.ReadOnly = readOnly
}

func (t *Transaction) GetReadOnly() bool {
	return t.proto.ReadOnly
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

func (t *Transaction) GetData() []byte {
	if t.IsCallContract() {
		return t.CallData().Input()
	}
	if t.IsDeployContract() {
		return t.DeployData().Code()
	}
	return t.Data()
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

func (t *Transaction) RawSignatureValues() (v, r, s *big.Int) {
	if t == nil {
		return nil, nil, nil
	}
	v = new(big.Int).SetBytes(t.proto.V)
	r = new(big.Int).SetBytes(t.proto.R)
	s = new(big.Int).SetBytes(t.proto.S)
	return v, r, s
}

// setSignatureValues sets the signature values (v, r, s) from *big.Int.
func (t *Transaction) SetSignatureValues(chainID, v, r, s *big.Int) {
	if t == nil {
		return
	}
	t.proto.V = v.Bytes()
	t.proto.R = r.Bytes()
	t.proto.S = s.Bytes()
	// Cập nhật ChainID nếu cần thiết
	t.proto.ChainID = uint64(chainID.Int64())
}

func (t *Transaction) SetSign(privateKey p_common.PrivateKey) {
	t.proto.Sign = bls.Sign(privateKey, t.Hash().Bytes()).Bytes()
}
func (t *Transaction) SetSignBytes(bytes []byte) {
	t.proto.Sign = bytes
}
func (t *Transaction) SetCommissionSign(privateKey p_common.PrivateKey) {
	t.proto.CommissionSign = bls.Sign(privateKey, t.Hash().Bytes()).Bytes()
}
func (t *Transaction) SetNonce(nonce uint64) {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)
	t.proto.Nonce = nonceBytes
}

func (t *Transaction) SetFromAddress(address common.Address) {
	t.proto.FromAddress = address.Bytes()
}

func (t *Transaction) SetToAddress(address common.Address) {
	t.proto.ToAddress = address.Bytes()
}

// validate
func (t *Transaction) ValidEthSign() bool {

	ethTx := t.GetEthTx()
	if ethTx == nil {
		return false
	}
	signer := e_types.NewCancunSigner(utils.Uint64ToBigInt(t.proto.ChainID))

	from, err := e_types.Sender(signer, ethTx)
	if err != nil {
		return false
	}
	return from == t.FromAddress()
}

func (t *Transaction) ValidSign() bool {
	return bls.VerifySign(
		t.Pubkey(),
		t.Sign(),
		t.Hash().Bytes(),
	)
}

func (t *Transaction) ValidChainID(chainId uint64) bool {
	return t.GetChainID() == chainId
}

func (t *Transaction) ValidLastHash(fromAccountState types.AccountState) bool {
	return t.LastHash() == fromAccountState.LastHash()
}

func (t *Transaction) ValidNonce(fromAccountState types.AccountState) bool {
	return t.GetNonce() == fromAccountState.Nonce()
}

func (t *Transaction) ValidMinTimeUse() bool {
	return t.MaxTimeUse() >= p_common.MIN_TX_TIME && t.MaxTimeUse() <= p_common.MAX_GROUP_TIME
}

func (t *Transaction) ValidTx0(fromAccountState types.AccountState, chainId string) (bool, int64) {
	if t.GetNonce() == 0 && len(fromAccountState.PublicKeyBls()) == 0 {
		dataInput := t.CallData().Input()

		if len(dataInput) != 113 {
			return false, InvalidDataInputLengthForTx0.Code
		}

		message := append(append([]byte{}, dataInput[:48]...), []byte(chainId)...)
		hash := crypto.Keccak256(message)

		valid := bls.VerifySign(
			p_common.PubkeyFromBytes(dataInput[:48]),
			t.Sign(),
			t.Hash().Bytes(),
		)
		if !valid {
			return false, InvalidBLSSignatureForTx0.Code
		}

		pb, err := secp256k1.RecoverPubkey(hash, dataInput[48:])
		if err != nil {
			return false, FailedToRecoverPubkeyForTx0.Code
		}

		var addr common.Address
		copy(addr[:], crypto.Keccak256(pb[1:])[12:])

		if t.FromAddress() == t.ToAddress() && addr == t.FromAddress() {
			return true, 0
		} else {
			return false, InvalidAddressMatchForTx0.Code
		}
	} else if t.GetNonce() == 0 && len(fromAccountState.PublicKeyBls()) != 0 {
		return false, NonceZeroButPublicKeyNotEmpty.Code
	}

	return true, 0
}

func (t *Transaction) ValidDeviceKey(fromAccountState types.AccountState) bool {
	return fromAccountState.DeviceKey() == common.Hash{} || // skip check device key if account state doesn't have device key
		crypto.Keccak256Hash(t.LastDeviceKey().Bytes()) == fromAccountState.DeviceKey()
}

func (t *Transaction) ValidMaxGas() bool {
	return t.MaxGas() >= p_common.TRANSFER_GAS_COST
}

func (t *Transaction) ValidMaxGasPrice(currentGasPrice uint64) bool {
	if t.ToAddress() == p_common.NATIVE_SMART_CONTRACT_REWARD_ADDRESS &&
		t.IsCallContract() {
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

	validToAddress := crypto.CreateAddress(fromAccountState.Address(), fromAccountState.Nonce())

	if validToAddress != t.ToAddress() {
		logger.Warn("Not match deploy address", validToAddress, t.ToAddress())
	}
	return validToAddress == t.ToAddress()
}

func (t *Transaction) ValidDeployData() bool {
	if t.DeployData() == nil || len(t.DeployData().Code()) == 0 {
		logger.Warn("Deploy data is nil")
		return false
	}
	return true
}

func (t *Transaction) ValidCallData() bool {
	if t.CallData() == nil || len(t.CallData().Input()) == 0 {
		logger.Warn("Deploy data is nil")
		return false
	}
	return true
}

func (t *Transaction) ValidOpenChannelToAccount(fromAccountState types.AccountState) bool {
	// validToAddress := common.BytesToAddress(
	// 	crypto.Keccak256(
	// 		append(
	// 			fromAccountState.Address().Bytes(),
	// 			fromAccountState.LastHash().Bytes()...),
	// 	)[12:],
	// )
	validToAddress := crypto.CreateAddress(fromAccountState.Address(), fromAccountState.Nonce())

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

func (t *Transaction) MarshalJSON() ([]byte, error) {
	// Tạo một map để lưu trữ dữ liệu của transaction
	data := map[string]interface{}{
		"Hash":             hex.EncodeToString(t.Hash().Bytes()),
		"FromAddress":      hex.EncodeToString(t.proto.FromAddress),
		"ToAddress":        hex.EncodeToString(t.proto.ToAddress),
		"PendingUse":       hex.EncodeToString(t.proto.PendingUse),
		"Amount":           hex.EncodeToString(t.proto.Amount),
		"MaxGas":           t.proto.MaxGas,
		"MaxGasPrice":      t.proto.MaxGasPrice,
		"MaxTimeUse":       t.proto.MaxTimeUse,
		"Data":             hex.EncodeToString(t.proto.Data),
		"RelatedAddresses": t.proto.RelatedAddresses,
		"LastDeviceKey":    hex.EncodeToString(t.proto.LastDeviceKey),
		"NewDeviceKey":     hex.EncodeToString(t.proto.NewDeviceKey),
		"Nonce":            t.GetNonce(),
		"Sign":             hex.EncodeToString(t.proto.Sign),
		"CommissionSign":   hex.EncodeToString(t.proto.CommissionSign),
		"ReadOnly":         t.GetReadOnly(),
	}

	// Chuyển đổi map thành JSON
	return json.Marshal(data)
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	// Tạo một map để lưu trữ dữ liệu JSON
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}

	// Lấy dữ liệu từ map và gán cho các trường của transaction
	// ... (Cần thêm logic để chuyển đổi dữ liệu từ JSON về các kiểu dữ liệu tương ứng của Go) ...

	return nil
}

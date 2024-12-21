package state

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"

	p_common "gomail/mtn/common"
	pb "gomail/mtn/proto"
	"gomail/mtn/types"
)

type SmartContractState struct {
	// proto *pb.SmartContractState
	createPublicKey p_common.PublicKey
	storageAddress  common.Address
	codeHash        common.Hash
	storageRoot     common.Hash
	logsHash        common.Hash
}

func NewSmartContractState(
	createPublicKey p_common.PublicKey,
	storageAddress common.Address,
	codeHash common.Hash,
	storageRoot common.Hash,
	logsHash common.Hash,
) types.SmartContractState {
	return &SmartContractState{
		createPublicKey: createPublicKey,
		storageAddress:  storageAddress,
		codeHash:        codeHash,
		storageRoot:     storageRoot,
		logsHash:        logsHash,
	}
}

func NewEmptySmartContractState() types.SmartContractState {
	return &SmartContractState{}
}

// general
func (ss *SmartContractState) Proto() *pb.SmartContractState {
	return &pb.SmartContractState{
		CreatorPublicKey: ss.createPublicKey.Bytes(),
		StorageAddress:   ss.storageAddress.Bytes(),
		CodeHash:         ss.codeHash.Bytes(),
		StorageRoot:      ss.storageRoot.Bytes(),
		LogsHash:         ss.logsHash.Bytes(),
	}
}

func (ss *SmartContractState) Marshal() ([]byte, error) {
	return proto.Marshal(ss.Proto())
}

func (ss *SmartContractState) FromProto(pbData *pb.SmartContractState) {
	ss.createPublicKey = p_common.PubkeyFromBytes(pbData.CreatorPublicKey)
	ss.storageAddress = common.BytesToAddress(pbData.StorageAddress)
	ss.codeHash = common.BytesToHash(pbData.CodeHash)
	ss.storageRoot = common.BytesToHash(pbData.StorageRoot)
	ss.logsHash = common.BytesToHash(pbData.LogsHash)
}

func (ss *SmartContractState) Unmarshal(b []byte) error {
	ssProto := &pb.SmartContractState{}
	err := proto.Unmarshal(b, ssProto)
	if err != nil {
		return err
	}
	ss.FromProto(ssProto)
	return nil
}

func (ss *SmartContractState) String() string {
	jsonSmartContractState := &JsonSmartContractState{}
	jsonSmartContractState.FromSmartContractState(ss)
	b, _ := json.MarshalIndent(jsonSmartContractState, "", " ")
	return string(b)
}

func (ss *SmartContractState) CreatorPublicKey() p_common.PublicKey {
	return ss.createPublicKey
}

func (ss *SmartContractState) CreatorAddress() common.Address {
	return p_common.AddressFromPubkey(ss.CreatorPublicKey())
}

func (ss *SmartContractState) StorageAddress() common.Address {
	return ss.storageAddress
}

func (ss *SmartContractState) CodeHash() common.Hash {
	return ss.codeHash
}

func (ss *SmartContractState) StorageRoot() common.Hash {
	return ss.storageRoot
}

func (ss *SmartContractState) LogsHash() common.Hash {
	return ss.logsHash
}

// setter
func (ss *SmartContractState) SetCreatorPublicKey(pk p_common.PublicKey) {
	ss.createPublicKey = pk
}

func (ss *SmartContractState) SetStorageAddress(address common.Address) {
	ss.storageAddress = address
}

func (ss *SmartContractState) SetCodeHash(hash common.Hash) {
	ss.codeHash = hash
}

func (ss *SmartContractState) SetStorageRoot(hash common.Hash) {
	ss.storageRoot = hash
}

func (ss *SmartContractState) SetLogsHash(hash common.Hash) {
	ss.logsHash = hash
}

func (ss *SmartContractState) Copy() types.SmartContractState {
	cpSs := &SmartContractState{}
	copy(cpSs.createPublicKey[:], ss.createPublicKey[:])
	copy(cpSs.storageAddress[:], ss.storageAddress[:])
	copy(cpSs.codeHash[:], ss.codeHash[:])
	copy(cpSs.storageRoot[:], ss.storageRoot[:])
	copy(cpSs.logsHash[:], ss.logsHash[:])
	return cpSs
}

type JsonSmartContractState struct {
	CreatorPublicKey string `json:"creator_public_key"`
	StorageAddress   string `json:"storage_address"`
	CodeHash         string `json:"code_hash"`
	StorageRoot      string `json:"storage_root"`
	LogsHash         string `json:"logs_hash"`
}

func (jss *JsonSmartContractState) FromSmartContractState(ss types.SmartContractState) {
	jss.CreatorPublicKey = ss.CreatorPublicKey().String()
	jss.StorageAddress = ss.StorageAddress().String()
	jss.CodeHash = ss.CodeHash().String()
	jss.StorageRoot = ss.StorageRoot().String()
	jss.LogsHash = ss.LogsHash().String()
}

func (jss *JsonSmartContractState) ToSmartContractState() types.SmartContractState {
	createPublicKey := p_common.PubkeyFromBytes(common.FromHex(jss.CreatorPublicKey))
	storageAddress := common.HexToAddress(jss.StorageAddress)
	codeHash := common.HexToHash(jss.CodeHash)
	storageRoot := common.HexToHash(jss.StorageRoot)
	logsHash := common.HexToHash(jss.LogsHash)
	return NewSmartContractState(createPublicKey, storageAddress, codeHash, storageRoot, logsHash)
}

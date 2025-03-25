package smart_contract

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"gomail/pkg/logger"
	pb "gomail/pkg/proto"
	"gomail/types"
)

type ExecuteSCResult struct {
	transactionHash common.Hash
	status          pb.RECEIPT_STATUS
	exception       pb.EXCEPTION
	returnData      []byte
	gasUsed         uint64
	logsHash        common.Hash

	mapAddBalance     map[string][]byte
	mapSubBalance     map[string][]byte
	mapStorageRoot    map[string][]byte
	mapCodeHash       map[string][]byte
	mapStorageAddress map[string]common.Address
	mapCreatorPubkey  map[string][]byte

	mapStorageAddressTouchedAddresses map[common.Address][]common.Address

	mapNativeSmartContractUpdateStorage map[common.Address][][2][]byte
	eventLogs                           []types.EventLog
}

func NewExecuteSCResult(
	transactionHash common.Hash,
	status pb.RECEIPT_STATUS,
	exception pb.EXCEPTION,
	rt []byte,
	gasUsed uint64,
	logsHash common.Hash,

	mapAddBalance map[string][]byte,
	mapSubBalance map[string][]byte,

	mapCodeHash map[string][]byte,
	mapStorageRoot map[string][]byte,

	mapStorageAddress map[string]common.Address,
	mapCreatorPubkey map[string][]byte,

	mapStorageAddressTouchedAddresses map[common.Address][]common.Address,

	mapNativeSmartContractUpdateStorage map[common.Address][][2][]byte,

	eventLogs []types.EventLog,
) *ExecuteSCResult {
	rs := &ExecuteSCResult{
		transactionHash: transactionHash,
		status:          status,
		exception:       exception,
		returnData:      rt,
		gasUsed:         gasUsed,
		logsHash:        logsHash,

		mapAddBalance:  mapAddBalance,
		mapSubBalance:  mapSubBalance,
		mapCodeHash:    mapCodeHash,
		mapStorageRoot: mapStorageRoot,

		mapStorageAddress: mapStorageAddress,
		mapCreatorPubkey:  mapCreatorPubkey,

		mapStorageAddressTouchedAddresses: mapStorageAddressTouchedAddresses,

		mapNativeSmartContractUpdateStorage: mapNativeSmartContractUpdateStorage,

		eventLogs: eventLogs,
	}
	logger.DebugP("len mapStorageAddressTouchedAddresses", len(mapStorageAddressTouchedAddresses))

	return rs
}

func NewErrorExecuteSCResult(
	transactionHash common.Hash,
	status pb.RECEIPT_STATUS,
	exception pb.EXCEPTION,
	rt []byte,
) *ExecuteSCResult {
	rs := &ExecuteSCResult{
		transactionHash: transactionHash,
		status:          status,
		exception:       exception,
		returnData:      rt,
		gasUsed:         0,
	}

	return rs
}

// general
func (r *ExecuteSCResult) Unmarshal(b []byte) error {
	pbRequest := &pb.ExecuteSCResult{}
	err := proto.Unmarshal(b, pbRequest)
	if err != nil {
		return err
	}
	r.FromProto(pbRequest)
	return nil
}

func (r *ExecuteSCResult) Marshal() ([]byte, error) {
	return proto.Marshal(r.Proto())
}

func (ex *ExecuteSCResult) String() string {
	str := fmt.Sprintf(`
	Transaction Hash: %v
	Add Balance Change:
	`,
		common.Bytes2Hex(ex.transactionHash[:]),
	)
	for i, v := range ex.mapAddBalance {
		str += fmt.Sprintf("%v: %v \n", i, uint256.NewInt(0).SetBytes(v))
	}
	str += fmt.Sprintln("Sub Balance Change: ")
	for i, v := range ex.mapSubBalance {
		str += fmt.Sprintf("%v: %v \n", i, uint256.NewInt(0).SetBytes(v))
	}
	str += fmt.Sprintln("Code Hash: ")
	for i, v := range ex.mapCodeHash {
		str += fmt.Sprintf("%v: %v \n", i, hex.EncodeToString(v))
	}
	str += fmt.Sprintln("Storage Root: ")
	for i, v := range ex.mapStorageRoot {
		str += fmt.Sprintf("%v: %v \n", i, hex.EncodeToString(v))
	}
	str += fmt.Sprintf(`
	Status: %v
	Exception: %v
	Return: %v
	GasUsed: %v
	`,
		ex.status,
		ex.exception,
		hex.EncodeToString(ex.returnData),
		ex.gasUsed,
	)
	return str
}

// getter
func (r *ExecuteSCResult) Proto() protoreflect.ProtoMessage {
	mapStorageAddressTouchedAddresses := make(map[string]*pb.TouchedAddresses)
	for k, v := range r.mapStorageAddressTouchedAddresses {
		touchedAddresses := &pb.TouchedAddresses{
			Addresses: make([][]byte, len(v)),
		}
		for i, a := range v {
			touchedAddresses.Addresses[i] = a.Bytes()
		}
		mapStorageAddressTouchedAddresses[hex.EncodeToString(k.Bytes())] = touchedAddresses
	}
	//
	mapStorageAddress := make(map[string][]byte, len(r.mapStorageAddress))
	for k, v := range r.mapStorageAddress {
		mapStorageAddress[k] = v.Bytes()
	}
	//
	mapNativeSmartContractUpdateStorage := make(
		map[string]*pb.StorageDatas,
		len(r.mapNativeSmartContractUpdateStorage),
	)
	for k, v := range r.mapNativeSmartContractUpdateStorage {
		datas := make([]*pb.StorageData, len(v))
		for i, vv := range v {
			datas[i] = &pb.StorageData{
				Key:   vv[0],
				Value: vv[1],
			}
		}
		mapNativeSmartContractUpdateStorage[k.Hex()] = &pb.StorageDatas{
			Datas: datas,
		}
	}
	eventLogs := make([]*pb.EventLog, len(r.eventLogs))
	for k, v := range r.eventLogs {
		eventLogs[k] = v.Proto()
	}

	protoData := &pb.ExecuteSCResult{
		TransactionHash: r.transactionHash.Bytes(),
		MapAddBalance:   r.mapAddBalance,
		MapSubBalance:   r.mapSubBalance,
		MapCodeHash:     r.mapCodeHash,
		MapStorageRoot:  r.mapStorageRoot,
		Status:          r.status,
		Exception:       r.exception,
		Return:          r.returnData,
		GasUsed:         r.gasUsed,

		MapStorageAddress: mapStorageAddress,
		MapCreatorPubkey:  r.mapCreatorPubkey,

		MapStorageAddressTouchedAddresses: mapStorageAddressTouchedAddresses,

		MapNativeSmartContractUpdateStorage: mapNativeSmartContractUpdateStorage,

		EventLogs: eventLogs,
	}
	return protoData
}

func (r *ExecuteSCResult) FromProto(pbData *pb.ExecuteSCResult) {
	r.transactionHash = common.BytesToHash(pbData.TransactionHash)
	r.mapAddBalance = pbData.MapAddBalance
	r.mapSubBalance = pbData.MapSubBalance
	r.mapCodeHash = pbData.MapCodeHash
	r.mapStorageRoot = pbData.MapStorageRoot
	r.status = pbData.Status
	r.exception = pbData.Exception
	r.returnData = pbData.Return
	r.gasUsed = pbData.GasUsed
	r.mapStorageAddressTouchedAddresses = make(map[common.Address][]common.Address)
	for k, v := range pbData.MapStorageAddressTouchedAddresses {
		address := common.HexToAddress(k)
		touchedAddresses := make([]common.Address, len(v.Addresses))
		for i, a := range v.Addresses {
			touchedAddresses[i] = common.BytesToAddress(a)
		}
		r.mapStorageAddressTouchedAddresses[address] = touchedAddresses
	}
	if len(pbData.MapCreatorPubkey) > 0 {
		r.mapStorageAddress = make(map[string]common.Address)
		for k, v := range pbData.MapStorageAddress {
			r.mapStorageAddress[k] = common.BytesToAddress(v)
		}
		r.mapCreatorPubkey = pbData.MapCreatorPubkey
	}
	r.mapNativeSmartContractUpdateStorage = make(
		map[common.Address][][2][]byte,
		len(pbData.MapNativeSmartContractUpdateStorage),
	)
	for k, v := range pbData.MapNativeSmartContractUpdateStorage {
		address := common.HexToAddress(k)
		r.mapNativeSmartContractUpdateStorage[address] = make([][2][]byte, len(v.Datas))
		for i, vv := range v.Datas {
			r.mapNativeSmartContractUpdateStorage[address][i] = [2][]byte{vv.Key, vv.Value}
		}
	}

	r.eventLogs = make([]types.EventLog, len(pbData.EventLogs))
	for idx, eventLog := range pbData.EventLogs {
		r.eventLogs[idx] = &EventLog{}
		r.eventLogs[idx].FromProto(eventLog)
	}

	logger.DebugP("len mapStorageAddressTouchedAddresses", len(r.mapStorageAddressTouchedAddresses))
}

func (r *ExecuteSCResult) TransactionHash() common.Hash {
	return r.transactionHash
}

func (r *ExecuteSCResult) MapAddBalance() map[string][]byte {
	return r.mapAddBalance
}

func (r *ExecuteSCResult) MapSubBalance() map[string][]byte {
	return r.mapSubBalance
}

func (r *ExecuteSCResult) MapStorageRoot() map[string][]byte {
	return r.mapStorageRoot
}

func (r *ExecuteSCResult) MapCodeHash() map[string][]byte {
	return r.mapCodeHash
}

func (r *ExecuteSCResult) MapStorageAddress() map[string]common.Address {
	return r.mapStorageAddress
}

func (r *ExecuteSCResult) MapCreatorPubkey() map[string][]byte {
	return r.mapCreatorPubkey
}

func (r *ExecuteSCResult) GasUsed() uint64 {
	return r.gasUsed
}

func (r *ExecuteSCResult) ReceiptStatus() pb.RECEIPT_STATUS {
	return r.status
}

func (r *ExecuteSCResult) Exception() pb.EXCEPTION {
	return r.exception
}

func (r *ExecuteSCResult) Return() []byte {
	return r.returnData
}

func (r *ExecuteSCResult) LogsHash() common.Hash {
	return r.logsHash
}

func (r *ExecuteSCResult) EventLogs() []types.EventLog {
	return r.eventLogs
}

func (r *ExecuteSCResult) MapStorageAddressTouchedAddresses() map[common.Address][]common.Address {
	return r.mapStorageAddressTouchedAddresses
}

func (r *ExecuteSCResult) MapNativeSmartContractUpdateStorage() map[common.Address][][2][]byte {
	return r.mapNativeSmartContractUpdateStorage
}

func ExecuteSCResultsFromProto(pbData []*pb.ExecuteSCResult) []types.ExecuteSCResult {
	results := make([]types.ExecuteSCResult, len(pbData))
	for i, v := range pbData {
		rs := &ExecuteSCResult{}
		rs.FromProto(v)
		results[i] = rs
	}
	return results
}

func ExecuteSCResultsToProto(results []types.ExecuteSCResult) []*pb.ExecuteSCResult {
	pbResults := make([]*pb.ExecuteSCResult, len(results))
	for i, v := range results {
		pbResults[i] = v.Proto().(*pb.ExecuteSCResult)
	}
	return pbResults
}

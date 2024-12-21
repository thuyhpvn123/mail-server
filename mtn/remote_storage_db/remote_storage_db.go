package remote_storage_db

import (
	"bytes"
	"errors"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"

	p_common "gomail/mtn/common"
	"gomail/mtn/logger"
	pb "gomail/mtn/proto"
	"gomail/mtn/types/network"
)

// Remote db is used to get data from storage connection

type RemoteStorageDB struct {
	remoteConnection   network.Connection
	messageSender      network.MessageSender
	address            common.Address
	currentBlockNumber uint64
	sync.Mutex
}

func NewRemoteStorageDB(
	remoteConnection network.Connection,
	messageSender network.MessageSender,
	address common.Address,
) *RemoteStorageDB {
	return &RemoteStorageDB{
		remoteConnection: remoteConnection,
		messageSender:    messageSender,
		address:          address,
	}
}

func (remoteDB *RemoteStorageDB) checkConnection() error {
	if !remoteDB.remoteConnection.IsConnect() {
		err := remoteDB.remoteConnection.Connect()
		if err != nil {
			return err
		}
		go remoteDB.remoteConnection.ReadRequest()
	}
	return nil
}

func (remoteDB *RemoteStorageDB) Get(
	key []byte,
) ([]byte, error) {
	remoteDB.Lock()
	defer remoteDB.Unlock()
	err := remoteDB.checkConnection()
	if err != nil {
		logger.Error("RemoteStorageDB.Get() checkConnection error: %v", err)
		return nil, err
	}
	// send get request to remote connection
	remoteDB.messageSender.SendMessage(
		remoteDB.remoteConnection,
		p_common.GetSmartContractStorage,
		&pb.GetSmartContractStorage{
			Address:     remoteDB.address.Bytes(),
			Key:         key,
			BlockNumber: remoteDB.currentBlockNumber,
		},
	)
	// wait for response
	rqChan, errChan := remoteDB.remoteConnection.RequestChan()
	timeOutChan := time.After(10 * time.Second)
	for {
		select {
		case response := <-rqChan:
			logger.Info("RemoteStorageDB.Get() command: ", response.Message().Command())
			switch response.Message().Command() {
			case p_common.GetSmartContractStorageResponse:
				logger.Info("case GetSmartContractStorageResponse")
				data := &pb.GetSmartContractStorageResponse{}
				err := proto.Unmarshal(response.Message().Body(), data)
				if err != nil {
					logger.Error("RemoteStorageDB.Get() proto.Unmarshal error: %v", err)
					return nil, err
				}
				if bytes.Equal(remoteDB.address.Bytes(), data.Address) &&
					bytes.Equal(key, data.Key) &&
					data.BlockNumber == remoteDB.currentBlockNumber {
					return data.Value, nil
				}
				if data.BlockNumber == 0 {
					logger.Error("RemoteStorageDB.Get()", data)
					return nil, errors.New("RemoteStorageDB.Get() data.BlockNumber == 0")
				}
			}
		case err := <-errChan:
			logger.Error("RemoteStorageDB.Get() error: %v", err)
			return nil, err
		case <-timeOutChan:
			logger.Error("RemoteStorageDB.Get() timeout")
			return nil, errors.New("timeout")
		}
	}
}

func (remoteDB *RemoteStorageDB) GetCode(
	address common.Address,
) ([]byte, error) {
	remoteDB.Lock()
	defer remoteDB.Unlock()
	err := remoteDB.checkConnection()
	if err != nil {
		logger.Error("RemoteStorageDB.Get() checkConnection error: %v", err)
		return nil, err
	}
	// send get request to remote connection
	remoteDB.messageSender.SendMessage(
		remoteDB.remoteConnection,
		p_common.GetSmartContractCode,
		&pb.GetSmartContractCode{
			Address:     address.Bytes(),
			BlockNumber: remoteDB.currentBlockNumber,
		},
	)
	// wait for response
	rqChan, errChan := remoteDB.remoteConnection.RequestChan()
	timeOutChan := time.After(10 * time.Second)
	for {
		select {
		case response := <-rqChan:
			logger.Info("RemoteStorageDB.Get() command: ", response.Message().Command())
			switch response.Message().Command() {
			case p_common.GetSmartContractCodeResponse:
				logger.Info("case GetSmartContractCodeResponse")
				data := &pb.GetSmartContractCodeResponse{}
				err := proto.Unmarshal(response.Message().Body(), data)
				if err != nil {
					logger.Error("RemoteStorageDB.Get() proto.Unmarshal error: %v", err)
					return nil, err
				}
				if address == common.Address(data.Address) &&
					data.BlockNumber == remoteDB.currentBlockNumber {
					return data.Code, nil
				}
			}
		case err := <-errChan:
			logger.Error("RemoteStorageDB.Get() error: %v", err)
			return nil, err
		case <-timeOutChan:
			logger.Error("RemoteStorageDB.Get() timeout")
			return nil, errors.New("timeout")
		}
	}
}

func (remoteDB *RemoteStorageDB) SetBlockNumber(blockNumber uint64) {
	remoteDB.currentBlockNumber = blockNumber
}

func (remoteDB *RemoteStorageDB) Close() {
	remoteDB.remoteConnection.Disconnect()
}

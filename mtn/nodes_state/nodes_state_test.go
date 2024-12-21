package nodes_state

import (
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"gomail/mtn/bls"
	p_common "gomail/mtn/common"
	"gomail/mtn/logger"
	p_network "gomail/mtn/network"
	"gomail/test/mock"
	"gomail/mtn/types/network"
)

func TestNodesState_GetStateRoot(t *testing.T) {
	logger.SetFlag(logger.FLAG_DEBUG)
	type fields struct {
		messageSender       network.MessageSender
		connectionsManager  network.ConnectionsManager
		childNodes          []common.Address
		childNodeStateRoots [TOTAL_NODES]common.Hash
		receivedChan        chan bool
		receivedNodesState  int
	}
	childNodes := []common.Address{}
	for i := 0; i < TOTAL_NODES; i++ {
		childNodes = append(childNodes, common.Address{byte(i)})
	}
	tests := []struct {
		name    string
		fields  fields
		want    common.Hash
		wantErr bool
	}{
		{
			name: "Test success",
			fields: fields{
				messageSender:       p_network.NewMessageSender(bls.GenerateKeyPair(), ""),
				connectionsManager:  createTestConnectionsManager(),
				childNodes:          childNodes,
				childNodeStateRoots: [TOTAL_NODES]common.Hash{},
				receivedChan:        make(chan bool),
				receivedNodesState:  0,
			},
			want: common.HexToHash(
				"0x7d47e3a1e68215898a2d1a6c8e5560e874e34fc8eaa90a9624c5f065c31c249c",
			),
			wantErr: false,
		},
		{
			name: "Test timeout",
			fields: fields{
				messageSender:       p_network.NewMessageSender(bls.GenerateKeyPair(), ""),
				connectionsManager:  createTestConnectionsManager(),
				childNodes:          childNodes,
				childNodeStateRoots: [TOTAL_NODES]common.Hash{},
				receivedChan:        make(chan bool),
				receivedNodesState:  0,
			},
			want:    common.Hash{},
			wantErr: true,
		},
		{
			name: "Test connection not found",
			fields: fields{
				messageSender:       p_network.NewMessageSender(bls.GenerateKeyPair(), ""),
				connectionsManager:  p_network.NewConnectionsManager(),
				childNodes:          childNodes,
				childNodeStateRoots: [TOTAL_NODES]common.Hash{},
				receivedChan:        make(chan bool),
				receivedNodesState:  0,
			},
			want:    common.Hash{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NodesState{
				messageSender:       tt.fields.messageSender,
				connectionsManager:  tt.fields.connectionsManager,
				childNodes:          tt.fields.childNodes,
				childNodeStateRoots: tt.fields.childNodeStateRoots,
				receivedChan:        tt.fields.receivedChan,
				receivedNodesState:  tt.fields.receivedNodesState,
			}
			if tt.name == "Test success" {
				go func() {
					time.Sleep(1 * time.Second)
					for i := 0; i < TOTAL_NODES; i++ {
						go n.ReceiveNodeState(common.Address{byte(i)}, common.Hash{byte(i)})
					}
				}()
			}
			got, err := n.GetStateRoot()
			if (err != nil) != tt.wantErr {
				t.Errorf("NodesState.GetStateRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodesState.GetStateRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createMockConnection() []network.Connection {
	var testConnection []network.Connection
	for i := 0; i < TOTAL_NODES; i++ {
		testConnection = append(testConnection, mock.NewTestConnection(
			common.Address{byte(i)},
			"testConnection",
			p_common.CHILD_NODE_CONNECTION_TYPE,
			func(message network.Message) (network.Request, error) {
				return nil, nil
			},
		))
	}
	return testConnection
}

func createTestConnectionsManager() network.ConnectionsManager {
	connections := createMockConnection()
	connectionsManager := p_network.NewConnectionsManager()
	for _, connection := range connections {
		connectionsManager.AddConnection(connection, true)
	}
	return connectionsManager
}

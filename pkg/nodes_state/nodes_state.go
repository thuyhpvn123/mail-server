package nodes_state

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	p_common "gomail/pkg/common"
	"gomail/pkg/logger"
	"gomail/pkg/state"
	p_sync "gomail/pkg/sync"
	"gomail/pkg/trie/node"
	"gomail/types/network"
)

const TOTAL_NODES = 16

type NodesState struct {
	sync.RWMutex

	prefix []byte

	socketServer       network.SocketServer
	messageSender      network.MessageSender
	connectionsManager network.ConnectionsManager

	childNodes          []common.Address
	childNodeStateRoots [TOTAL_NODES]common.Hash

	receivedChan        chan bool
	receivedNodesState  int
	getStateConnections []network.Connection

	currentSession string
}

func NewNodesState(
	childNodes []common.Address,

	messageSender network.MessageSender,
	connectionsManager network.ConnectionsManager,
) *NodesState {
	return &NodesState{
		childNodes:          childNodes,
		messageSender:       messageSender,
		connectionsManager:  connectionsManager,
		getStateConnections: make([]network.Connection, TOTAL_NODES),
	}
}

func (n *NodesState) SetSocketServer(s network.SocketServer) {
	n.socketServer = s
}

func (n *NodesState) GetStateRoot(blockNumber uint64) (common.Hash, error) {
	// reset counter
	n.receivedNodesState = 0
	n.receivedChan = make(chan bool)
	n.currentSession = uuid.New().String()
	// send request to get node state
	n.sendGetNodeStateRequest(blockNumber)
	// wait for all nodes state
	select {
	case <-time.After(30 * time.Second):
		return common.Hash{}, errors.New("Timeout to get all nodes state")
	case <-n.receivedChan:
	}

	fNode := &node.FullNode{
		Children: [17]node.Node{},
	}
	for i, v := range n.childNodeStateRoots {
		h := node.HashNode(v.Bytes())
		fNode.Children[i] = h
		logger.Debug("GetStateRoot: ", v)
	}
	b, _ := fNode.Marshal()
	logger.Debug("GetStateRoot: ", hex.EncodeToString(b))
	return crypto.Keccak256Hash(b), nil
}

func (n *NodesState) GetChildNode(i int) common.Address {
	return n.childNodes[i]
}

func (n *NodesState) GetChildNodeIdx(nodeAddress common.Address) int {
	i := slices.Index(n.childNodes[:], nodeAddress)
	return i
}

func (n *NodesState) GetChildNodeStateRoot(address common.Address) common.Hash {
	i := slices.Index(n.childNodes[:], address)
	return n.childNodeStateRoots[i]
}

func (n *NodesState) SetChildNode(i int, childNode common.Address) {
	n.childNodes[i] = childNode
}

func (n *NodesState) ReceiveNodeState(
	nodeAddress common.Address,
	uuid string,
	stateRoot common.Hash,
) {
	n.Lock()
	defer n.Unlock()
	if n.currentSession != uuid {
		return
	}

	i := slices.Index(n.childNodes[:], nodeAddress)
	err := n.setChildNodeStateRoot(i, stateRoot)
	logger.Debug("ReceiveNodeState", "nodeAddress", nodeAddress, "stateRoot", stateRoot)
	if err != nil {
		return
	}
	n.receivedNodesState++
	logger.Debug("ReceiveNodeState", n.receivedNodesState, TOTAL_NODES)
	if n.receivedNodesState == TOTAL_NODES {
		n.receivedChan <- true
	}
}

func (n *NodesState) setChildNodeStateRoot(i int, childNodeStateRoot common.Hash) error {
	if i >= TOTAL_NODES {
		return errors.New("index out of range")
	}
	n.childNodeStateRoots[i] = childNodeStateRoot
	return nil
}

func (n *NodesState) sendGetNodeStateRequest(blockNumber uint64) error {
	var bBlockNumber [8]byte
	binary.LittleEndian.PutUint64(bBlockNumber[:], blockNumber)
	bCurrentSession := []byte(n.currentSession)
	b := make([]byte, 0, 8+len(bCurrentSession))
	b = append(b, bBlockNumber[:]...)
	b = append(b, bCurrentSession...)

	// Send request to get node state
	for _, childNode := range n.childNodes {
		connection := n.connectionsManager.ConnectionByTypeAndAddress(
			p_common.CHILD_NODE_CONNECTION_IDX,
			childNode,
		)
		if connection == nil {
			return errors.New("child node connection not found for address: " + childNode.String())
		}
		go func(conn network.Connection) {
			n.messageSender.SendBytes(
				conn,
				p_common.GetNodeStateRoot,
				b,
			)
		}(connection)
	}
	return nil
}

func (n *NodesState) SendCancelPendingStates() {
	wg := sync.WaitGroup{}
	for _, childNode := range n.childNodes {
		connection := n.connectionsManager.ConnectionByTypeAndAddress(
			p_common.CHILD_NODE_CONNECTION_IDX,
			childNode,
		)

		if connection == nil {
			continue
		}
		wg.Add(1)
		go func(conn network.Connection) {
			n.messageSender.SendBytes(
				conn,
				p_common.CancelNodePendingState,
				[]byte{},
			)
			wg.Done()
		}(connection)
	}
	wg.Wait()
}

func (n *NodesState) SendGetAccountState(address common.Address, id string) error {
	nibbles := p_common.KeybytesToHex(crypto.Keccak256(address.Bytes()))
	// remove prefix
	nibbles = nibbles[len(n.prefix):]
	// add to group
	idx := nibbles[0]
	nodeAddress := n.childNodes[idx]
	if (nodeAddress == common.Address{}) {
		logger.Error("node address not found")
		return errors.New("node address not found")
	}
	connection := n.connectionsManager.ConnectionByTypeAndAddress(
		p_common.CHILD_NODE_CONNECTION_IDX,
		nodeAddress,
	)
	if connection == nil {
		logger.Error("connection not found")
		return errors.New("connection not found")
	}
	if n.getStateConnections[idx] == nil {
		// need to clone connection because main connection is processing request, so it wont be able to receive response
		n.getStateConnections[idx] = connection.Clone()
		err := n.getStateConnections[idx].Connect()
		if err != nil {
			return errors.New("new connection not found")
		}
		go n.socketServer.HandleConnection(n.getStateConnections[idx])
	}

	bData, err := state.MarshalGetAccountStateWithIdRequest(address, id)
	if err != nil {
		return err
	}

	return n.messageSender.SendBytes(
		n.getStateConnections[idx],
		p_common.GetAccountStateWithIdRequest,
		bData,
	)
}

func (n *NodesState) SendGetNodeSyncData(
	latestCheckPointBlockNumber uint64,
	validatorAddress common.Address,
) {
	// nodeStatesIndex int,
	wg := sync.WaitGroup{}
	for i, childNode := range n.childNodes {
		connection := n.connectionsManager.ConnectionByTypeAndAddress(
			p_common.CHILD_NODE_CONNECTION_IDX,
			childNode,
		)

		if connection == nil {
			continue
		}
		getNodeSyncData := p_sync.NewGetNodeSyncData(
			latestCheckPointBlockNumber,
			validatorAddress,
			i,
		)
		wg.Add(1)
		go func(conn network.Connection, data *p_sync.GetNodeSyncData) {
			n.messageSender.SendMessage(
				conn,
				p_common.GetNodeSyncData,
				p_sync.GetNodeSyncDataToProto(data),
			)
			wg.Done()
		}(connection, getNodeSyncData)
	}
	wg.Wait()
}

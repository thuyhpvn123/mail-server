package network

import (
	"sync"

	p_common "gomail/pkg/common"
	"gomail/pkg/logger"
	"gomail/types/network"

	"github.com/holiman/uint256"

	pb "gomail/pkg/proto"

	"github.com/ethereum/go-ethereum/common"
)

type ConnectionsManager struct {
	mu                          sync.RWMutex
	parentConnection            network.Connection
	typeToMapAddressConnections []map[common.Address]network.Connection
}

func NewConnectionsManager() network.ConnectionsManager {
	cm := &ConnectionsManager{}
	cm.typeToMapAddressConnections = make([]map[common.Address]network.Connection, 20)
	for i := range cm.typeToMapAddressConnections {
		cm.typeToMapAddressConnections[i] = make(map[common.Address]network.Connection)
	}
	return cm
}

// getter
func (cm *ConnectionsManager) ConnectionsByType(cType int) map[common.Address]network.Connection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.typeToMapAddressConnections[cType]
}

func (cm *ConnectionsManager) ConnectionByTypeAndAddress(cType int, address common.Address) network.Connection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.typeToMapAddressConnections[cType][address]
}

func (cm *ConnectionsManager) ConnectionsByTypeAndAddresses(cType int, addresses []common.Address) map[common.Address]network.Connection {
	rs := make(map[common.Address]network.Connection, len(addresses))
	for _, v := range addresses {
		rs[v] = cm.ConnectionByTypeAndAddress(cType, v)
	}

	return rs
}

func (cm *ConnectionsManager) FilterAddressAvailable(cType int, addresses map[common.Address]*uint256.Int) map[common.Address]*uint256.Int {
	availableAddresses := make(map[common.Address]*uint256.Int)
	for address := range addresses {
		if cm.ConnectionByTypeAndAddress(cType, address) != nil {
			availableAddresses[address] = addresses[address]
		}
	}
	return availableAddresses
}

func (cm *ConnectionsManager) ParentConnection() network.Connection {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.parentConnection
}

func (cm *ConnectionsManager) Stats() *pb.NetworkStats {
	pbNetworkStats := &pb.NetworkStats{
		TotalConnectionByType: make(map[string]int32, len(cm.typeToMapAddressConnections)),
	}
	for i, v := range cm.typeToMapAddressConnections {
		pbNetworkStats.TotalConnectionByType[p_common.MapIndexToConnectionType(i)] = int32(len(v))
	}
	return pbNetworkStats
}

// setter
func (cm *ConnectionsManager) AddParentConnection(conn network.Connection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.parentConnection = conn
}

func (cm *ConnectionsManager) RemoveConnection(conn network.Connection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cType := p_common.MapConnectionTypeToIndex(conn.Type())
	if cm.typeToMapAddressConnections[cType][conn.Address()] == conn {
		delete(cm.typeToMapAddressConnections[cType], conn.Address())
		logger.Debug("Removing connection", conn.Address())
	}
}

func (cm *ConnectionsManager) AddConnection(conn network.Connection, replace bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	address := conn.Address()
	cType := p_common.MapConnectionTypeToIndex(conn.Type())

	if (address != common.Address{} &&
		cm.typeToMapAddressConnections[cType][address] == nil) ||
		replace {
		cm.typeToMapAddressConnections[cType][address] = conn
	}
}

func MapAddressConnectionToInterface(data map[common.Address]network.Connection) map[common.Address]interface{} {
	rs := make(map[common.Address]interface{})
	for i, v := range data {
		rs[i] = v
	}
	return rs
}

package network

import (
	"context"
	"fmt"
	"net"
	"time"

	"gomail/pkg/bls"
	p_common "gomail/pkg/common"
	"gomail/pkg/logger"
	pb "gomail/pkg/proto"
	"gomail/types/network"
)

type SocketServer struct {
	connectionsManager network.ConnectionsManager
	listener           net.Listener
	handler            network.Handler

	nodeType   string
	version    string
	dnsLink    string
	keyPair    *bls.KeyPair
	ctx        context.Context
	cancelFunc context.CancelFunc

	onConnectedCallBack    []func(network.Connection)
	onDisconnectedCallBack []func(network.Connection)
}

func NewSocketServer(
	keyPair *bls.KeyPair,
	connectionsManager network.ConnectionsManager,
	handler network.Handler,
	nodeType string,
	version string,
	dnsLink string,
) network.SocketServer {
	s := &SocketServer{
		keyPair:            keyPair,
		connectionsManager: connectionsManager,
		handler:            handler,
		nodeType:           nodeType,
		version:            version,
		dnsLink:            dnsLink,
	}
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())
	return s
}

func (s *SocketServer) SetContext(ctx context.Context, cancelFunc context.CancelFunc) {
	s.cancelFunc = cancelFunc
	s.ctx = ctx
}

func (s *SocketServer) AddOnConnectedCallBack(callBack func(network.Connection)) {
	s.onConnectedCallBack = append(s.onConnectedCallBack, callBack)
}

func (s *SocketServer) AddOnDisconnectedCallBack(callBack func(network.Connection)) {
	s.onDisconnectedCallBack = append(s.onDisconnectedCallBack, callBack)
}

func (s *SocketServer) Listen(listenAddress string) error {
	var err error
	s.listener, err = net.Listen("tcp", listenAddress)
	if err != nil {
		return err
	}
	defer func() {
		if s.listener != nil {
			s.listener.Close()
			s.listener = nil
		}
	}()
	logger.Debug(fmt.Sprintf("Listening at %v", listenAddress))
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			tcpConn, err := s.listener.Accept()
			if err != nil {
				logger.Warn(fmt.Sprintf("Error when accept connection %v\n", err))
				continue
			}
			conn, err := ConnectionFromTcpConnection(tcpConn, s.dnsLink)
			if err != nil {
				logger.Warn(
					fmt.Sprintf("error when create connection from tcp connection: %v", err),
				)
				continue
			}
			go s.OnConnect(conn)
			go s.HandleConnection(conn)
		}
	}
}

func (s *SocketServer) Stop() {
	s.cancelFunc()
}

func (s *SocketServer) OnConnect(conn network.Connection) {
	logger.Info(fmt.Sprintf("On Connect with %s", conn.RemoteAddr()))
	SendMessage(conn, p_common.InitConnection, &pb.InitConnection{
		Address: s.keyPair.Address().Bytes(),
		Type:    s.nodeType,
	}, s.version)

	for _, v := range s.onConnectedCallBack {
		v(conn)
	}
}

func (s *SocketServer) OnDisconnect(conn network.Connection) {
	logger.Warn(
		fmt.Sprintf(
			"On Disconnect with %s",
			conn,
		),
	)
	s.connectionsManager.RemoveConnection(conn)

	for _, v := range s.onDisconnectedCallBack {
		v(conn)
	}
}

func (s *SocketServer) HandleConnection(conn network.Connection) error {
	logger.Debug(fmt.Sprintf("handle connection %v", conn.Address()))
	go conn.ReadRequest()
	defer func() {
		conn.Disconnect()
		s.OnDisconnect(conn)
	}()
	requestChan, errorChan := conn.RequestChan()
	for {
		select {
		case <-s.ctx.Done():
			return nil
		case request := <-requestChan:
			if request == nil {
				return nil
			}
			err := s.handler.HandleRequest(request)
			if err != nil {
				logger.Warn(fmt.Sprintf("error when process request %v", err))
				continue
			}
		case err := <-errorChan:
			if err != ErrDisconnected {
				logger.Warn(fmt.Sprintf("error when read request %v", err))
			}
			return err
		}
	}
}

func (s *SocketServer) SetKeyPair(newKeyPair *bls.KeyPair) {
	s.keyPair = newKeyPair
}

func (s *SocketServer) StopAndRetryConnectToParent(conn network.Connection) {
	if conn == s.connectionsManager.ParentConnection() {
		logger.Warn("Disconnected with parent")
		// stop running if disconnected with parent
		s.Stop()
		// if connection is parent connection then retry connect
		go func(_conn network.Connection) {
			for {
				<-time.After(5 * time.Second)
				err := _conn.Connect()
				if err != nil {
					logger.Warn(fmt.Sprintf("error when retry connect to parent %v", err))
				} else {
					s.ctx, s.cancelFunc = context.WithCancel(context.Background())
					s.connectionsManager.AddParentConnection(conn)
					s.OnConnect(conn)
					go s.HandleConnection(conn)
					return
				}
			}
		}(conn)
	}
}

func (s *SocketServer) RetryConnectToParent(conn network.Connection) {
	if conn == s.connectionsManager.ParentConnection() {
		logger.Warn("Disconnected with parent")
		// if connection is parent connection then retry connect
		go func(_conn network.Connection) {
			for {
				<-time.After(5 * time.Second)
				err := _conn.Connect()
				if err != nil {
					logger.Warn(fmt.Sprintf("error when retry connect to parent %v", err))
				} else {
					s.ctx, s.cancelFunc = context.WithCancel(context.Background())
					s.connectionsManager.AddParentConnection(conn)
					s.OnConnect(conn)
					go s.HandleConnection(conn)
					return
				}
			}
		}(conn)
	}
}

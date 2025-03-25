package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"

	p_common "gomail/pkg/common"
	"gomail/pkg/logger"
	pb "gomail/pkg/proto"
	"gomail/types/network"
)

var (
	ErrDisconnected         = errors.New("error connection disconnected")
	ErrInvalidMessageLength = errors.New("invalid message length")
	ErrExceedMessageLength  = errors.New("exceed message length")
	ErrNilConnection        = errors.New("nil connection")
)

func ConnectionFromTcpConnection(tcpConn net.Conn, dnsLink string) (network.Connection, error) {
	return &Connection{
		dnsLink:      dnsLink,
		address:      common.Address{},
		cType:        "",
		tcpConn:      tcpConn,
		requestChan:  make(chan network.Request, 1000000),
		errorChan:    make(chan error, 1),
		realConnAddr: tcpConn.RemoteAddr().String(),
	}, nil
}

func NewConnection(
	address common.Address,
	cType string,
	dnsLink string,
) network.Connection {
	return &Connection{
		address:     address,
		dnsLink:     dnsLink,
		cType:       cType,
		requestChan: make(chan network.Request, 1000000),
		errorChan:   make(chan error, 1),
		connect:     false,
	}
}

type Connection struct {
	mu      sync.Mutex
	address common.Address
	cType   string

	requestChan chan network.Request
	errorChan   chan error
	tcpConn     net.Conn
	connect     bool

	dnsLink      string
	realConnAddr string
}

// getter
func (c *Connection) Address() common.Address {
	return c.address
}

func (c *Connection) ConnectionAddress() (string, error) {
	var err error
	if c.realConnAddr == "" {
		c.realConnAddr, err = p_common.GetRealConnectionAddress(
			c.dnsLink,
			c.address,
		)
	}

	return c.realConnAddr, err
}

func (c *Connection) RequestChan() (chan network.Request, chan error) {
	return c.requestChan, c.errorChan
}

func (c *Connection) Type() string {
	return c.cType
}

func (c *Connection) String() string {
	connectionAddress, _ := c.ConnectionAddress()
	return fmt.Sprintf(
		`Address: %v 
		Type %v
		Connection Address %v`,
		c.address,
		c.cType,
		connectionAddress,
	)
}

// setter
func (c *Connection) Init(
	address common.Address,
	cType string,
) {
	c.address = address
	c.cType = cType
}

func (c *Connection) SetRealConnAddr(realConnAddr string) {
	c.realConnAddr = realConnAddr
}

// other
func (c *Connection) SendMessage(message network.Message) error {
	if c == nil {
		return ErrNilConnection
	}
	b, err := message.Marshal()
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(b)))
	if c.tcpConn == nil {
		return errors.New("tcpConn is closed")
	}
	_, err = c.tcpConn.Write(length)
	if err != nil {
		return err
	}
	// Write length and message in a loop to ensure complete transmission
	totalSent := 0
	for totalSent < len(b) { // Ensure both length and message are sent
		n, err := c.tcpConn.Write(b[totalSent:])
		if err != nil {
			return err
		}
		totalSent += n
	}
	return nil
}

func (c *Connection) Connect() (err error) {
	realConnectionAddress, err := c.ConnectionAddress()
	if err != nil {
		return err
	}

	c.tcpConn, err = net.Dial("tcp", realConnectionAddress)
	if err == nil {
		c.requestChan = make(chan network.Request, 100)
		c.errorChan = make(chan error, 1)
		c.connect = true
	}
	return err
}

func (c *Connection) Disconnect() error {
	return c.tcpConn.Close()
}

func (c *Connection) IsConnect() bool {
	return c.connect
}

func (c *Connection) ReadRequest() {
	defer func() {
		logger.Info("Connection closed")
		c.connect = false
		close(c.errorChan)
		close(c.requestChan)
	}()

	for {
		bLength := make([]byte, 8)
		_, err := io.ReadFull(c.tcpConn, bLength)
		if err != nil {
			switch err {
			case io.EOF:
				c.errorChan <- ErrDisconnected
			default:
				c.errorChan <- err
			}
			return
		}
		messageLength := binary.LittleEndian.Uint64(bLength)
		start := time.Now()
		maxMsgLength := uint64(1073741824)
		if messageLength > maxMsgLength {
			c.errorChan <- ErrExceedMessageLength
			return
		}

		data := make([]byte, messageLength)
		byteRead, err := io.ReadFull(c.tcpConn, data)
		if err != nil {
			switch err {
			case io.EOF:
				c.errorChan <- ErrDisconnected
			default:
				c.errorChan <- err
			}
			return

		}

		if uint64(byteRead) != messageLength {
			c.errorChan <- ErrExceedMessageLength
			return
		}

		msg := &pb.Message{}
		err = proto.Unmarshal(data[:messageLength], msg)
		clear(data)
		if err != nil {
			c.errorChan <- err
			return
		}

		c.requestChan <- NewRequest(c, NewMessage(msg))
		logger.Trace(
			"Process time for read request: "+time.Since(start).String(),
			c.Address(),
			c.tcpConn.RemoteAddr(),
		)
	}
}

func (c *Connection) Clone() network.Connection {
	newConn := NewConnection(
		c.address,
		c.cType,
		c.dnsLink,
	)
	return newConn
}

func (c *Connection) RemoteAddr() string {
	return c.tcpConn.RemoteAddr().String()
}

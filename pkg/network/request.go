package network

import "gomail/types/network"

type Request struct {
	connection network.Connection
	message    network.Message
}

func NewRequest(
	connection network.Connection,
	message network.Message,
) network.Request {
	return &Request{
		connection: connection,
		message:    message,
	}
}

func (r *Request) Message() network.Message {
	return r.message
}

func (r *Request) Connection() network.Connection {
	return r.connection
}

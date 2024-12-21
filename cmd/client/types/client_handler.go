package types

import "gomail/mtn/types/network"

type ClientHandler interface {
	HandleRequest(request network.Request) (err error)
}

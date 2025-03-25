package types

import "gomail/types/network"

type ClientHandler interface {
	HandleRequest(request network.Request) (err error)
}

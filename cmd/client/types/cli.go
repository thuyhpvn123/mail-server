package types

import (
	"github.com/ethereum/go-ethereum/common"
	"gomail/mtn/types"
)

type Cli interface {
	Start()
	Stop()
	PrintCommands()
	PrintMessage(string, string)
	SendTransaction() error
	CreateAccount()
	ChangeAccount()
	AccountState(common.Address) (types.AccountState, error)
	ReadInput() string
	ReadInputAddress() common.Address
}

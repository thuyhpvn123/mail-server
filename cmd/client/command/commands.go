package command

const (
	//General
	InitConnection = "InitConnection"

	GetStats       = "GetStats"
	Stats          = "Stats"
	ChangeLogLevel = "ChangeLogLevel"

	// Send messages
	SendTransaction      = "SendTransaction"
	SendTransactions     = "SendTransactions"
	GetAccountState      = "GetAccountState"
	SubscribeToAddress   = "SubscribeToAddress"
	GetStakeState        = "GetStakeState"
	GetSmartContractData = "GetSmartContractData"

	// Receive message
	AccountState      = "AccountState"
	StakeState        = "StakeState"
	Receipt           = "Receipt"
	TransactionError  = "TransactionError"
	EventLogs         = "EventLogs"
	QueryLogs         = "QueryLogs"
	SmartContractData = "SmartContractData"
)

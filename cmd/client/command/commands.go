package command

const (
	//General
	InitConnection = "InitConnection"

	GetStats       = "GetStats"
	Stats          = "Stats"
	ChangeLogLevel = "ChangeLogLevel"

	// Send messages
	ReadTransaction              = "ReadTransaction"
	SendTransaction              = "SendTransaction"
	SendTransactionWithDeviceKey = "SendTransactionWithDeviceKey"
	SendTransactions             = "SendTransactions"
	GetAccountState              = "GetAccountState"
	SubscribeToAddress           = "SubscribeToAddress"
	GetStakeState                = "GetStakeState"
	GetSmartContractData         = "GetSmartContractData"

	GetDeviceKey = "GetDeviceKey"

	// Receive message
	AccountState      = "AccountState"
	StakeState        = "StakeState"
	Receipt           = "Receipt"
	TransactionError  = "TransactionError"
	EventLogs         = "EventLogs"
	QueryLogs         = "QueryLogs"
	SmartContractData = "SmartContractData"
	DeviceKey         = "DeviceKey"
)

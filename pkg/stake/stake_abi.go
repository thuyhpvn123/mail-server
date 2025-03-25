package stake

import (
	"bytes"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func StakeABI() abi.ABI {
	a, _ := abi.JSON(
		bytes.NewReader([]byte(`[
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "addr",
				"type": "address"
			}
		],
		"name": "getStakeInfo",
		"outputs": [
			{
				"internalType": "address",
				"name": "owner",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			},
			{
				"internalType": "address[]",
				"name": "childNodes",
				"type": "address[]"
			},
			{
				"internalType": "address[]",
				"name": "childExecuteMiners",
				"type": "address[]"
			},
			{
				"internalType": "address[]",
				"name": "childVerifyMiners",
				"type": "address[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getValidatorsWithStakeAmount",
		"outputs": [
			{
				"internalType": "address[]",
				"name": "addresses",
				"type": "address[]"
			},
			{
				"internalType": "uint256[]",
				"name": "amounts",
				"type": "uint256[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`)),
	)
	return a
}

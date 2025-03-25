package execute_sc_grouper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"

	p_common "gomail/pkg/common"
	"gomail/pkg/transaction"
	"gomail/types"
)

func TestExecuteSmartContractsGrouper_AddTransactions(t *testing.T) {
	type fields struct {
		groupCount                  uint64
		mapAddressGroup             map[common.Address]uint64
		mapGroupExecuteTransactions map[uint64][]types.Transaction
	}
	type args struct {
		transactions []types.Transaction
	}

	tx1 := transaction.NewTransaction(
		common.Hash{},
		p_common.PubkeyFromBytes([]byte{}),
		common.Address{0x01},
		uint256.NewInt(0),
		uint256.NewInt(0),
		0,
		0,
		0,
		0,
		nil,
		nil,
		common.Hash{},
		common.Hash{},
	)

	tx2 := transaction.NewTransaction(
		common.Hash{},
		p_common.PubkeyFromBytes([]byte{}),
		common.Address{0x02},
		uint256.NewInt(0),
		uint256.NewInt(0),
		0,
		0,
		0,
		0,
		nil,
		nil,
		common.Hash{},
		common.Hash{},
	)

	tx3 := transaction.NewTransaction(
		common.Hash{},
		p_common.PubkeyFromBytes([]byte{0x1}),
		common.Address{0x02},
		uint256.NewInt(0),
		uint256.NewInt(0),
		0,
		0,
		0,
		0,
		nil,
		nil,
		common.Hash{},
		common.Hash{},
	)

	tx4 := transaction.NewTransaction(
		common.Hash{},
		p_common.PubkeyFromBytes([]byte{0x02}),
		common.Address{0x03},
		uint256.NewInt(0),
		uint256.NewInt(0),
		0,
		0,
		0,
		0,
		nil,
		nil,
		common.Hash{},
		common.Hash{},
	)

	tests := []struct {
		name                              string
		fields                            fields
		args                              args
		wantedMapGroupExecuteTransactions map[uint64][]types.Transaction
	}{
		{
			name: "Test 1",
			fields: fields{
				groupCount:                  0,
				mapAddressGroup:             make(map[common.Address]uint64),
				mapGroupExecuteTransactions: make(map[uint64][]types.Transaction),
			},
			args: args{
				transactions: []types.Transaction{
					tx1,
					tx2,
					tx3,
					tx4,
				},
			},
			wantedMapGroupExecuteTransactions: map[uint64][]types.Transaction{
				1: {tx1, tx2, tx3},
				2: {tx4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ExecuteSmartContractsGrouper{
				groupCount:                  tt.fields.groupCount,
				mapAddressGroup:             tt.fields.mapAddressGroup,
				mapGroupExecuteTransactions: tt.fields.mapGroupExecuteTransactions,
			}
			e.AddTransactions(tt.args.transactions)
			assert.EqualValues(
				t,
				tt.wantedMapGroupExecuteTransactions,
				e.mapGroupExecuteTransactions,
				"The two maps should be the same.",
			)
		})
	}
}

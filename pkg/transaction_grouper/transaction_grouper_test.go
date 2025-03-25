package transaction_grouper

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"

	p_common "gomail/pkg/common"
	"gomail/pkg/transaction"
	"gomail/types"
)

func TestTransactionGrouper_AddFromTransactions(t *testing.T) {
	type fields struct {
		groups [16][]types.Transaction
		prefix []byte
	}
	type args struct {
		transactions []types.Transaction
	}
	tx1 := transaction.NewTransaction(
		common.Hash{},
		p_common.PubkeyFromBytes([]byte{}),
		common.Address{},
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
		name       string
		fields     fields
		args       args
		wantGroups [16][]types.Transaction
	}{
		{
			name: "Test without prefix",
			fields: fields{
				groups: [16][]types.Transaction{},
				prefix: []byte{},
			},
			args: args{
				transactions: []types.Transaction{
					tx1,
				},
			},
			wantGroups: [16][]types.Transaction{
				0x0f: {tx1},
			},
		},
		{
			name: "Test with prefix",
			fields: fields{
				groups: [16][]types.Transaction{},
				prefix: []byte{0x0f},
			},
			args: args{
				transactions: []types.Transaction{
					tx1,
				},
			},
			wantGroups: [16][]types.Transaction{
				0x04: {tx1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TransactionGrouper{
				groups: tt.fields.groups,
				prefix: tt.fields.prefix,
			}
			tr.AddFromTransactions(tt.args.transactions)
			assert.EqualValues(t, tt.wantGroups, tr.groups)
		})
	}
}

func TestTransactionGrouper_AddToTransactions(t *testing.T) {
	type fields struct {
		groups [16][]types.Transaction
		prefix []byte
	}
	type args struct {
		transactions []types.Transaction
	}

	tx1 := transaction.NewTransaction(
		common.Hash{},
		p_common.PubkeyFromBytes([]byte{}),
		common.Address{0x12},
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
		name       string
		fields     fields
		args       args
		wantGroups [16][]types.Transaction
	}{
		{
			name: "Test without prefix",
			fields: fields{
				groups: [16][]types.Transaction{},
				prefix: []byte{},
			},
			args: args{
				transactions: []types.Transaction{
					tx1,
				},
			},
			wantGroups: [16][]types.Transaction{
				0x01: {tx1},
			},
		},
		{
			name: "Test with prefix",
			fields: fields{
				groups: [16][]types.Transaction{},
				prefix: []byte{0x01},
			},
			args: args{
				transactions: []types.Transaction{
					tx1,
				},
			},
			wantGroups: [16][]types.Transaction{
				0x02: {tx1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TransactionGrouper{
				groups: tt.fields.groups,
				prefix: tt.fields.prefix,
			}
			tr.AddToTransactions(tt.args.transactions)
			assert.EqualValues(t, tt.wantGroups, tr.groups)
		})
	}
}

func TestTransactionGrouper_GetTransactionsGroups(t *testing.T) {
	type fields struct {
		groups [16][]types.Transaction
		prefix []byte
	}
	tx1 := transaction.NewTransaction(
		common.Hash{},
		p_common.PubkeyFromBytes([]byte{}),
		common.Address{0x12},
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
		name   string
		fields fields
		want   [16][]types.Transaction
	}{
		{
			name: "Test get transaction groups",
			fields: fields{
				groups: [16][]types.Transaction{
					0x01: {tx1},
				},
				prefix: []byte{0x01},
			},
			want: [16][]types.Transaction{
				0x01: {tx1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TransactionGrouper{
				groups: tt.fields.groups,
				prefix: tt.fields.prefix,
			}
			if got := tr.GetTransactionsGroups(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransactionGrouper.GetTransactionsGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionGrouper_Clear(t *testing.T) {
	type fields struct {
		groups [16][]types.Transaction
		prefix []byte
	}
	tx1 := transaction.NewTransaction(
		common.Hash{},
		p_common.PubkeyFromBytes([]byte{}),
		common.Address{0x12},
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
		name   string
		fields fields
	}{
		{
			name: "Test clear",
			fields: fields{
				groups: [16][]types.Transaction{
					0x01: {tx1},
				},
				prefix: []byte{0x01},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TransactionGrouper{
				groups: tt.fields.groups,
				prefix: tt.fields.prefix,
			}
			tr.Clear()
			assert.Empty(t, tr.groups)
		})
	}
}

func TestNewTransactionGrouper(t *testing.T) {
	type args struct {
		prefix []byte
	}
	tests := []struct {
		name string
		args args
		want *TransactionGrouper
	}{
		{
			name: "Test new transaction grouper",
			args: args{
				prefix: []byte{0x01},
			},
			want: &TransactionGrouper{
				groups: [16][]types.Transaction{},
				prefix: []byte{0x01},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransactionGrouper(tt.args.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionGrouper() = %v, want %v", got, tt.want)
			}
		})
	}
}

package account_state_db

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"gomail/pkg/state"
	"gomail/pkg/storage"
	p_trie "gomail/pkg/trie"
	"gomail/types"
)

func TestAccountStateDB_AddPendingBalance(t *testing.T) {
	type fields struct {
		trie           *p_trie.MerklePatriciaTrie
		originRootHash common.Hash
		db             storage.Storage
		dirtyAccounts  map[common.Address]types.AccountState
	}
	type args struct {
		address common.Address
		amount  *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test 1",
			fields: fields{
				trie:           nil,
				originRootHash: common.Hash{},
				db:             nil,
				dirtyAccounts: map[common.Address]types.AccountState{
					{0x01}: state.NewAccountState(common.Address{0x01}),
				},
			},
			args: args{
				address: common.Address{0x01},
				amount:  big.NewInt(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &AccountStateDB{
				trie:           tt.fields.trie,
				originRootHash: tt.fields.originRootHash,
				db:             tt.fields.db,
				dirtyAccounts:  tt.fields.dirtyAccounts,
			}
			db.AddPendingBalance(tt.args.address, tt.args.amount)
			if tt.name == "Test 1" {
				db.AddPendingBalance(tt.args.address, tt.args.amount)
				if db.dirtyAccounts[tt.args.address].PendingBalance().Cmp(big.NewInt(2)) != 0 {
					t.Errorf(
						"AddPendingBalance() = %v, want %v",
						db.dirtyAccounts[tt.args.address].PendingBalance(),
						big.NewInt(2),
					)
				}
			}
		})
	}
}

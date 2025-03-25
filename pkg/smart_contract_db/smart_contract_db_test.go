package smart_contract_db

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/maypok86/otter"

	"gomail/pkg/bls"
	"gomail/pkg/logger"
	"gomail/pkg/network"
	pb "gomail/pkg/proto"
	"gomail/pkg/state"
	"gomail/pkg/storage"
	"gomail/pkg/trie"
	t_network "gomail/types/network"
)

func TestSmartContractDB_Code(t *testing.T) {
	type fields struct {
		cacheRemoteDBs     otter.Cache[common.Address, RemoteStorageDB]
		cacheCode          otter.Cache[common.Hash, []byte]
		cacheStorageTrie   otter.Cache[common.Address, *trie.MerklePatriciaTrie]
		messageSender      t_network.MessageSender
		dnsLink            string
		accountStateDB     AccountStateDB
		currentBlockNumber *uint256.Int
	}

	cacheRemoteDBs, err := otter.MustBuilder[common.Address, RemoteStorageDB](20).
		CollectStats().
		DeletionListener(func(key common.Address, value RemoteStorageDB, cause otter.DeletionCause) {
			value.Close()
		}).
		Build()
	if err != nil {
		panic(err)
	}

	cacheCode, err := otter.MustBuilder[common.Hash, []byte](1_000).
		CollectStats().
		Build()

	cacheStorageTrie, err := otter.MustBuilder[common.Address, *trie.MerklePatriciaTrie](1_000).
		CollectStats().
		Build()
	if err != nil {
		panic(err)
	}

	messageSender := network.NewMessageSender(bls.GenerateKeyPair(), "127.0.0.1")
	accountStatesDb := storage.NewMemoryDb()
	asTrie, err := trie.New(common.Hash{}, accountStatesDb)
	accountStatesManager := state.NewAccountStatesManager(asTrie, accountStatesDb)
	if err != nil {
		panic(err)
	}
	// put account
	as := state.AccountStateFromProto(
		&pb.AccountState{
			Address: common.FromHex("0x02"),
			SmartContractState: &pb.SmartContractState{
				CreatorPublicKey: common.FromHex(
					"86d5de6f7c9c13cc0d959a553cc0e4853ba5faae45a28da9bddc8ef8e104eb5d3dece8dfaa24f11b4243ec27537e3184",
				),
				StorageHost:    "storage",
				StorageAddress: common.FromHex("da7284fac5e804f8b9d71aa39310f0f86776b51d"),
				CodeHash: common.FromHex(
					"0x4760c37bbf051e02eb7ae63c68e49a9caee114b39b84f9e4528ab5466f326403",
				),
				StorageRoot: common.FromHex(
					"0x1b6ee20ffa6a16e16685b3f6a946e471193a63cb7392d0bb3f2d2b15c216a798",
				),
			},
		},
	)
	accountStatesManager.SetState(as)
	accountStatesManager.CommitTransaction()
	_, err = accountStatesManager.Commit()
	if err != nil {
		panic(err)
	}

	type args struct {
		address common.Address
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "Test 1",
			fields: fields{
				cacheRemoteDBs:     cacheRemoteDBs,
				cacheCode:          cacheCode,
				cacheStorageTrie:   cacheStorageTrie,
				messageSender:      messageSender,
				dnsLink:            "http://127.0.0.1:7080/api/dns/connection-address/",
				accountStateDB:     accountStatesManager,
				currentBlockNumber: uint256.NewInt(0),
			},
			args: args{
				address: common.HexToAddress("0x02"),
			},
			want: common.FromHex(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scdb := &SmartContractDB{
				cacheRemoteDBs:     tt.fields.cacheRemoteDBs,
				cacheCode:          tt.fields.cacheCode,
				cacheStorageTrie:   tt.fields.cacheStorageTrie,
				messageSender:      tt.fields.messageSender,
				dnsLink:            tt.fields.dnsLink,
				accountStateDB:     tt.fields.accountStateDB,
				currentBlockNumber: tt.fields.currentBlockNumber,
			}
			got := scdb.Code(tt.args.address)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("SmartContractDB.Code() = %v, want %v", got, tt.want)
			// }
			logger.Info("SmartContractDB.Code()", hex.EncodeToString(got))
			scdb.Code(tt.args.address)
		})
	}
}

func TestSmartContractDB_StorageValue(t *testing.T) {
	type fields struct {
		cacheRemoteDBs     otter.Cache[common.Address, RemoteStorageDB]
		cacheCode          otter.Cache[common.Hash, []byte]
		cacheStorageTrie   otter.Cache[common.Address, *trie.MerklePatriciaTrie]
		messageSender      t_network.MessageSender
		dnsLink            string
		accountStateDB     AccountStateDB
		currentBlockNumber *uint256.Int
	}

	cacheRemoteDBs, err := otter.MustBuilder[common.Address, RemoteStorageDB](1_000).
		CollectStats().
		DeletionListener(func(key common.Address, value RemoteStorageDB, cause otter.DeletionCause) {
			value.Close()
		}).
		Build()
	if err != nil {
		panic(err)
	}

	cacheCode, err := otter.MustBuilder[common.Hash, []byte](1_000).
		CollectStats().
		Build()

	cacheStorageTrie, err := otter.MustBuilder[common.Address, *trie.MerklePatriciaTrie](1_000).
		CollectStats().
		Build()

	messageSender := network.NewMessageSender(bls.GenerateKeyPair(), "127.0.0.1")
	accountStatesDb := storage.NewMemoryDb()
	asTrie, err := trie.New(common.Hash{}, accountStatesDb)
	accountStatesManager := state.NewAccountStatesManager(asTrie, accountStatesDb)
	if err != nil {
		panic(err)
	}

	// put account
	as := state.AccountStateFromProto(
		&pb.AccountState{
			Address: common.FromHex("0x02"),
			SmartContractState: &pb.SmartContractState{
				CreatorPublicKey: common.FromHex(
					"86d5de6f7c9c13cc0d959a553cc0e4853ba5faae45a28da9bddc8ef8e104eb5d3dece8dfaa24f11b4243ec27537e3184",
				),
				StorageHost:    "storage",
				StorageAddress: common.FromHex("da7284fac5e804f8b9d71aa39310f0f86776b51d"),
				CodeHash: common.FromHex(
					"0x4760c37bbf051e02eb7ae63c68e49a9caee114b39b84f9e4528ab5466f326403",
				),
				StorageRoot: common.FromHex(
					"0x1b6ee20ffa6a16e16685b3f6a946e471193a63cb7392d0bb3f2d2b15c216a798",
				),
			},
		},
	)
	accountStatesManager.SetState(as)
	accountStatesManager.CommitTransaction()
	_, err = accountStatesManager.Commit()
	if err != nil {
		panic(err)
	}

	type args struct {
		address common.Address
		key     []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "Test 1",
			fields: fields{
				cacheRemoteDBs:     cacheRemoteDBs,
				cacheCode:          cacheCode,
				cacheStorageTrie:   cacheStorageTrie,
				messageSender:      messageSender,
				dnsLink:            "http://127.0.0.1:7080/api/dns/connection-address/",
				accountStateDB:     accountStatesManager,
				currentBlockNumber: uint256.NewInt(0),
			},
			args: args{
				address: common.HexToAddress("0x02"),
				key: common.FromHex(
					"0000000000000000000000000000000000000000000000000000000000000004",
				),
			},
			want: common.FromHex(
				"00000000000000000000000000000000000000000000000000000002540be400",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scdb := &SmartContractDB{
				cacheRemoteDBs:     tt.fields.cacheRemoteDBs,
				cacheCode:          tt.fields.cacheCode,
				cacheStorageTrie:   tt.fields.cacheStorageTrie,
				messageSender:      tt.fields.messageSender,
				dnsLink:            tt.fields.dnsLink,
				accountStateDB:     tt.fields.accountStateDB,
				currentBlockNumber: tt.fields.currentBlockNumber,
			}
			if got := scdb.StorageValue(tt.args.address, tt.args.key); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("SmartContractDB.StorageValue() = %v, want %v", got, tt.want)
			}
			// test second time for cache
			if got := scdb.StorageValue(tt.args.address, tt.args.key); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("SmartContractDB.StorageValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

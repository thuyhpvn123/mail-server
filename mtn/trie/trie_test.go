package trie

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	e_common "github.com/ethereum/go-ethereum/common"

	"gomail/mtn/logger"
	"gomail/mtn/network"
	"gomail/remote_storage_db"
	"gomail/mtn/storage"
)

func TestEmptyTrie(t *testing.T) {
	trie, err := New(
		e_common.Hash{},
		storage.NewMemoryDb(),
	)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	res := trie.Hash()
	exp := EmptyRootHash
	if res != exp {
		t.Errorf("expected %x got %x", exp, res)
	}
}

func TestInsert(t *testing.T) {
	trie, err := New(
		e_common.Hash{},
		storage.NewMemoryDb(),
	)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	updateString(trie, "doe", "reindeer")
	gDoe, err := trie.Get([]byte("doe"))
	fmt.Print("xxx" + string(gDoe))
	updateString(trie, "dog", "puppy")
	updateString(trie, "dogglesworth", "cat")

	exp := e_common.HexToHash("8aad789dff2f538bca5d8ea56e8abe10f4c7ba3a5dea95fea4cd6e7c3a1168d3")
	root := trie.Hash()
	if root != exp {
		t.Errorf("case 1: exp %x got %x", exp, root)
	}
	dbx := storage.NewMemoryDb()
	trie, err = New(
		e_common.Hash{},
		dbx,
	)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	updateString(trie, "A", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dox", "dox")
	updateString(trie, "dogglesworth", "cat")
	updateString(trie, "dogy", "ya")
	updateString(trie, "dogyy", "ya1")

	gA, err := trie.Get([]byte("A"))
	logger.DebugP("xxx2", string(gA))

	gDoe, err = trie.Get([]byte("dogglesworth"))
	logger.DebugP("xxx3" + string(gDoe))

	gdogy, err := trie.Get([]byte("dogy"))
	logger.DebugP("xxx4" + string(gdogy))

	gdogyy, err := trie.Get([]byte("dogyy"))
	logger.DebugP("xxx4" + string(gdogyy))

	exp = e_common.HexToHash("d23786fb4a010da3ce639d66d5e904a11dbc02746d1ce25029e53290cabf28ab")
	root, nodeSet, _, err := trie.Commit(true)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	if root != exp {
		t.Errorf("case 2: exp %x got %x", exp, root)
	}
	for i, v := range nodeSet.Nodes {
		logger.DebugP("Node set", i, v)
	}
	gA, err = trie.Get([]byte("A"))
	logger.DebugP("XXXXXXXXXXXX", string(gA))

	// save nodeSet to db
	for _, node := range nodeSet.Nodes {
		logger.DebugP("Commit", "node", node)
		err := dbx.Put(node.Hash.Bytes(), node.Blob)
		if err != nil {
			logger.Error("dbx put error", err)
		}
	}
	logger.DebugP("xxx root" + hex.EncodeToString(root.Bytes()))
	mpt, err := New(root, dbx)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	gA, err = mpt.Get([]byte("A"))
	if err != nil {
		logger.Error("mpt get error", err)
	}
	logger.DebugP("xxx4", string(gA))
	gXX, err := mpt.Get([]byte("dogglesworth"))
	if err != nil {
		logger.Error("mpt get error", err)
	}
	logger.DebugP("MMMM", string(gXX))
	gXXX, err := mpt.Get([]byte("dogy"))
	if err != nil {
		logger.Error("mpt get error", err)
	}
	logger.DebugP("MMMM2", string(gXXX))

	gXXXXX, err := mpt.Get([]byte("dogyy"))
	if err != nil {
		logger.Error("mpt get error", err)
	}
	logger.DebugP("MMMM3", string(gXXXXX))

	gX, err := mpt.Get([]byte("dox"))
	if err != nil {
		logger.Error("mpt get error", err)
	}
	logger.DebugP("MMMM3", string(gX))
}

func TestInsertAndRemove(t *testing.T) {
	dbx := storage.NewMemoryDb()

	trie, err := New(
		e_common.Hash{},
		dbx,
	)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	updateString(trie, "do", "reindeer")
	updateString(trie, "do1", "puppy1")
	updateString(trie, "bu", "puppy2")
	root, nodeSet, oldKeys, err := trie.Commit(true)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}

	for _, node := range nodeSet.Nodes {
		logger.DebugP("Commit", "node", node.Hash)
		err := dbx.Put(node.Hash.Bytes(), node.Blob)
		if err != nil {
			logger.Error("dbx put error", err)
		}
	}

	trie2, err := New(root, dbx)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	txt, err := trie2.Get([]byte("do1"))
	if err != nil {
		logger.Error("trie 2 get error", err)
	}
	if string(txt) != "puppy1" {
		t.Errorf("case 2: exp puppy1 got %x", string(txt))
	}
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	updateString(trie2, "do1", "puppy3")

	root2, nodeSet, oldKeys, err := trie2.Commit(true)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}

	for _, node := range nodeSet.Nodes {
		logger.DebugP("Commit", "node", node.Hash)
		err := dbx.Put(node.Hash.Bytes(), node.Blob)
		if err != nil {
			logger.Error("dbx put error", err)
		}
	}

	for i := 0; i < len(oldKeys); i++ {
		logger.Error(hex.EncodeToString(oldKeys[i]))
		dbx.Delete(oldKeys[i])
	}

	_, err = New(root, dbx)
	if err == nil {
		t.Errorf("expected nil error got %x", err)
	}

	trie3, err := New(root2, dbx)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	txt, err = trie3.Get([]byte("do1"))
	if err != nil {
		logger.Error("trie 2 get error", err)
	}

	logger.Error(string(txt))
}

func updateString(trie *MerklePatriciaTrie, k, v string) {
	trie.Update([]byte(k), []byte(v))
}

func TestMerklePatriciaTrie_String(t *testing.T) {
	// 0xf8edce0de9a848231c444b71dd2676aaf7909c426306ae26c3fa75462762788e
	db, err := storage.NewLevelDB("./test/db/account_states/")
	if err != nil {
		panic(err)
	}
	newRoot := e_common.HexToHash(
		"0x230546c294eb221d99cf6cc5d7d0464c7020821f5efe7069874170528648141b",
	)
	// oldRoot := e_common.HexToHash("f8edce0de9a848231c444b71dd2676aaf7909c426306ae26c3fa75462762788e")
	trie, err := New(
		newRoot,
		db,
	)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
		return
	}

	err = trie.Update(e_common.Hex2Bytes("c5b109d06fcedde202fae3774f9c6a3bc2684f49"), []byte("hmm"))
	if err != nil {
		panic(err)
	}
	root, nodeSet, _, err := trie.Commit(true)
	logger.Info("new root", root)
	if err != nil {
		panic(err)
	}

	// save nodeSet to db
	if nodeSet != nil {
		for _, node := range nodeSet.Nodes {
			logger.DebugP("Commit", "node", node)
			err := db.Put(node.Hash.Bytes(), node.Blob)
			if err != nil {
				logger.Error("dbx put error", err)
			}
		}
	}
	trie2, err := New(
		root,
		db,
	)
	v, err := trie2.Get(e_common.Hex2Bytes("c5b109d06fcedde202fae3774f9c6a3bc2684f49"))
	if err != nil {
		panic(err)
	}
	logger.Info(
		hex.EncodeToString(
			v,
		),
	)
	// logger.Info("trie", trie2.String())
}

func TestMerklePatriciaTrie_GetStorageKeys(t *testing.T) {
	trie, err := New(
		e_common.Hash{},
		storage.NewMemoryDb(),
	)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	updateString(trie, "doe", "reindeer")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dogglesworth", "cat")
	updateString(trie, "A", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dox", "dox")
	updateString(trie, "dogglesworth", "cat")
	updateString(trie, "dogy", "ya")
	updateString(trie, "dogyy", "ya1")
	root, nodeSet, _, err := trie.Commit(true)
	logger.Info("new root", root)
	if err != nil {
		panic(err)
	}

	db := storage.NewMemoryDb()
	// save nodeSet to db
	for _, node := range nodeSet.Nodes {
		logger.DebugP("Commit", "node", node)
		err := db.Put(node.Hash.Bytes(), node.Blob)
		if err != nil {
			logger.Error("dbx put error", err)
		}
	}
	trie2, err := New(
		root,
		db,
	)
	logger.Info("trie ", trie2)
	logger.Info("storage keys", trie2.GetStorageKeys())
}

func TestNewTrieWithRemoteDB(t *testing.T) {
	connection := network.NewConnection(
		common.HexToAddress("da7284fac5e804f8b9d71aa39310f0f86776b51d"),
		"storage",
		"http://127.0.0.1:7080/api/dns/connection-address/",
	)
	connection.Connect()
	messageSender := network.NewMessageSender(
		"1.0.0",
	)
	remoteDB := remote_storage_db.NewRemoteStorageDB(
		connection.(*network.Connection),
		messageSender.(*network.MessageSender),
		common.HexToAddress("02"),
	)

	trie, err := New(
		e_common.HexToHash("0x1b6ee20ffa6a16e16685b3f6a946e471193a63cb7392d0bb3f2d2b15c216a798"),
		remoteDB,
	)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	if err != nil {
		panic(err)
	}
	v, err := trie.Get(
		common.FromHex("0000000000000000000000000000000000000000000000000000000000000003"),
	)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	if err != nil {
		panic(err)
	}
	logger.Info("value", hex.EncodeToString(v))
}

func TestStorageTrie(t *testing.T) {
	trie, err := New(
		e_common.Hash{},
		storage.NewMemoryDb(),
	)
	if err != nil {
		t.Errorf("expected nil error got %x", err)
	}
	trie.Update(
		common.FromHex("0000000000000000000000000000000000000000000000000000000000000000"),
		common.FromHex("00000000000000000000000097126b71376f7e55fba904fdaa9df0dbd396612f"),
	)
	trie.Update(
		common.FromHex("0000000000000000000000000000000000000000000000000000000000000003"),
		common.FromHex("0000000000000000000000000000000000000000000000000000000000000001"),
	)
	trie.Update(
		common.FromHex("0000000000000000000000000000000000000000000000000000000000000004"),
		common.FromHex("00000000000000000000000000000000000000000000000000000002540be400"),
	)
	trie.Update(
		common.FromHex("0000000000000000000000000000000000000000000000000000000000000005"),
		common.FromHex("4d65746120446f6c6c6172205265776172640000000000000000000000000024"),
	)
	trie.Update(
		common.FromHex("0000000000000000000000000000000000000000000000000000000000000006"),
		common.FromHex("5553444d5200000000000000000000000000000000000000000000000000000a"),
	)
	trie.Update(
		common.FromHex("0000000000000000000000000000000000000000000000000000000000000007"),
		common.FromHex("0000000000000000000000000000000000000000000000000000000000000000"),
	)
	trie.Update(
		common.FromHex("000000000000000000000000000000000000000000000000000000000000000d"),
		common.FromHex("0000000000000000000000000000000000000000000000000000000000000012"),
	)
	trie.Update(
		common.FromHex("965cc3758cf4d4860015aab1d7336972ce4a797ab1225b752484364fb6cd824a"),
		common.FromHex("00000000000000000000000000000000000000000000000000000002540be400"),
	)

	res := trie.Hash()
	exp := common.HexToHash("0x1b6ee20ffa6a16e16685b3f6a946e471193a63cb7392d0bb3f2d2b15c216a798")
	if res != exp {
		t.Errorf("expected %x got %x", exp, res)
	}
}

package trie

import (
	e_common "github.com/ethereum/go-ethereum/common"
)

type Tracer struct {
	inserts    map[string]struct{}
	deletes    map[string]struct{}
	oldKeys    [][]byte
	accessList map[string][]byte
}

// newTracer initializes the tracer for capturing trie changes.
func newTracer() *Tracer {
	return &Tracer{
		inserts:    make(map[string]struct{}),
		deletes:    make(map[string]struct{}),
		oldKeys:    [][]byte{},
		accessList: make(map[string][]byte),
	}
}

// onRead tracks the newly loaded trie node and caches the rlp-encoded
// blob internally. Don't change the value outside of function since
// it's not deep-copied.
func (t *Tracer) onRead(path []byte, val []byte) {
	t.accessList[string(path)] = val
}

// onInsert tracks the newly inserted trie node. If it's already
// in the deletion set (resurrected node), then just wipe it from
// the deletion set as it's "untouched".
func (t *Tracer) onInsert(path []byte) {
	if _, present := t.deletes[string(path)]; present {
		delete(t.deletes, string(path))
		return
	}
	t.inserts[string(path)] = struct{}{}
}

// onDelete tracks the newly deleted trie node. If it's already
// in the addition set, then just wipe it from the addition set
// as it's untouched.
func (t *Tracer) onDelete(path []byte) {
	if _, present := t.inserts[string(path)]; present {
		delete(t.inserts, string(path))
		return
	}
	t.deletes[string(path)] = struct{}{}
}

// reset clears the content tracked by Tracer.
func (t *Tracer) reset() {
	t.inserts = make(map[string]struct{})
	t.deletes = make(map[string]struct{})
	t.accessList = make(map[string][]byte)
	t.oldKeys = [][]byte{}
}

// copy returns a deep copied Tracer instance.
func (t *Tracer) copy() *Tracer {
	var (
		inserts    = make(map[string]struct{})
		deletes    = make(map[string]struct{})
		accessList = make(map[string][]byte)
	)
	for path := range t.inserts {
		inserts[path] = struct{}{}
	}
	for path := range t.deletes {
		deletes[path] = struct{}{}
	}
	for path, blob := range t.accessList {
		accessList[path] = e_common.CopyBytes(blob)
	}
	return &Tracer{
		inserts:    inserts,
		deletes:    deletes,
		accessList: accessList,
	}
}

// deletedNodes returns a list of node paths which are deleted from the trie.
func (t *Tracer) deletedNodes() []string {
	var paths []string
	for path := range t.deletes {
		// It's possible a few deleted nodes were embedded
		// in their parent before, the deletions can be no
		// effect by deleting nothing, filter them out.
		_, ok := t.accessList[path]
		if !ok {
			continue
		}
		paths = append(paths, path)
	}
	return paths
}

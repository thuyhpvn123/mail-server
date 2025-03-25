package pack_pool

import (
	"sync"

	"gomail/types"
)

type PackPool struct {
	packs []types.Pack
	mutex sync.Mutex
}

func NewPackPool() *PackPool {
	return &PackPool{}
}

func (pp *PackPool) AddPack(pack types.Pack) {
	pp.mutex.Lock()
	defer pp.mutex.Unlock()

	pp.packs = append(pp.packs, pack)
}

func (pp *PackPool) Addpacks(packs []types.Pack) {
	pp.mutex.Lock()
	defer pp.mutex.Unlock()

	pp.packs = append(pp.packs, packs...)
}

func (pp *PackPool) Getpacks() []types.Pack {
	pp.mutex.Lock()
	defer pp.mutex.Unlock()
	packs := pp.packs
	pp.packs = make([]types.Pack, 0)
	return packs
}

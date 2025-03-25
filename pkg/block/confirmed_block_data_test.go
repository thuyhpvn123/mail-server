package block

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"gomail/types"
)

func TestSaveConfirmedBlockDataToFile(t *testing.T) {
	// Test code here
	initConfirmedBlockData := NewConfirmedBlockData(
		&BlockHeader{},
		[]types.Receipt{},
		common.Hash{},
		map[common.Address][]byte{},
	)
	err := SaveConfirmedBlockDataToFile(initConfirmedBlockData, "./")
	assert.Nil(t, err)
}

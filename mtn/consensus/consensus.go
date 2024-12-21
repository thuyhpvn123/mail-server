package consensus

import (
	"gomail/mtn/types"
)

type ConsensusEngine interface {
	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine.
	VerifyRootBlockHeader(header types.He) error
	VerifyShardBlockHeader(header types.He) error
}

package iep

import (
	"github.com/artheranet/arthera-node/inter"
	"github.com/artheranet/arthera-node/inter/ier"
)

type LlrEpochPack struct {
	Votes  []inter.LlrSignedEpochVote
	Record ier.LlrIdxFullEpochRecord
}

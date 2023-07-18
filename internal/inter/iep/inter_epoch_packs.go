package iep

import (
	"github.com/artheranet/arthera-node/internal/inter"
	"github.com/artheranet/arthera-node/internal/inter/ier"
)

type LlrEpochPack struct {
	Votes  []inter.LlrSignedEpochVote
	Record ier.LlrIdxFullEpochRecord
}

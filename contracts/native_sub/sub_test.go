package native_sub

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAbiPacking(t *testing.T) {
	state, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)

	sp := SubscriptionPlan{
		PlanId:       1,
		Name:         "plan1 is the best for eoa",
		Description:  "long description of plan1 is lorem ipsum dolor sit amet consectetuer adipiscing elit",
		Duration:     30,
		Units:        100000,
		Price:        100,
		CapFrequency: 1,
		CapUnits:     10,
		ForContract:  true,
		Active:       true,
	}

	StoreSubscriptionPlan(state, sp)
	sp2 := RetrieveSubscriptionPlan(state, sp.PlanId)
	require.Equal(t, sp, sp2)

	//hash := common.BytesToHash(data)
	//sp2, err := UnpackSubscriptionPlan(hash.Bytes())
	//require.Nil(t, err)
	//require.Equal(t, sp, sp2)
}

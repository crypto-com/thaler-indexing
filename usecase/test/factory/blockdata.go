package usecasefactory

import (
	"github.com/crypto-com/chainindex"
	. "github.com/crypto-com/chainindex/test/factory"
	_ "github.com/crypto-com/chainindex/test/factory"
	"github.com/crypto-com/chainindex/usecase"
)

func RandomBlockData() usecase.BlockData {
	block := RandomBlock()
	reward := RandomBlockReward()
	return usecase.BlockData{
		Block:              block,
		Signatures:         RandomBlockSignaturesOfSize(3, block.Height),
		Activities:         RandomAllActivities(),
		Reward:             &reward,
		CouncilNodeUpdates: RandomCouncilNodeUpdatesOfSize(3),
	}
}

func RandomBlockSignaturesOfSize(n int, blockHeight uint64) []chainindex.BlockSignature {
	signatures := make([]chainindex.BlockSignature, 0, n)
	for i := 0; i < n; i += 1 {
		blockSignature := RandomBlockSignature()
		blockSignature.BlockHeight = blockHeight
		signatures = append(signatures, blockSignature)
	}

	return signatures
}

func RandomCouncilNodeUpdatesOfSize(n int) []chainindex.CouncilNodeUpdate {
	updates := make([]chainindex.CouncilNodeUpdate, 0, n)
	for i := 0; i < n; i += 1 {
		updates = append(updates, RandomCouncilNodeLeftUpdate())
	}

	return updates
}

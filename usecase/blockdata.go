package usecase

import (
	"github.com/luci/go-render/render"

	"github.com/crypto-com/chainindex"
)

// BlockData is an aggregate of data inside a block. It is not a major concern
// in the domain, but more a usecase of how block should be persisted.
type BlockData struct {
	Block              chainindex.Block
	Signatures         []chainindex.BlockSignature
	Activities         []chainindex.Activity
	Reward             *chainindex.BlockReward
	CouncilNodeUpdates []chainindex.CouncilNodeUpdate
}

func (data *BlockData) String() string {
	return render.Render(data)
}

type BlockDataRepository interface {
	Store(blockData *BlockData) error
}

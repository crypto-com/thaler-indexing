package chainindex

import (
	"math/big"
	"time"

	"github.com/luci/go-render/render"
)

type Block struct {
	Height  uint64
	Hash    string
	Time    time.Time
	AppHash string
}

func (block *Block) String() string {
	return render.Render(block)
}

type BlockReward struct {
	BlockHeight uint64
	Minted      *big.Int
}

func (reward *BlockReward) String() string {
	return render.Render(reward)
}

type BlockSignature struct {
	BlockHeight        uint64
	CouncilNodeAddress string
	Signature          string
	IsProposer         bool
}

func (signature *BlockSignature) String() string {
	return render.Render(signature)
}

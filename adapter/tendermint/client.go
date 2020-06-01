package tendermint

import "github.com/crypto-com/chainindex/adapter/tendermint/types"

type Client interface {
	Genesis() (*types.Genesis, error)
	LatestBlockHeight() (uint64, error)
	BlockResults(height uint64) (*types.BlockResults, error)
	Block(height uint64) (*types.Block, error)
}

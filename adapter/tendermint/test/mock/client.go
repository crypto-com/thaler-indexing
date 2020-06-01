package adaptertendermintmock

import (
	"github.com/stretchr/testify/mock"

	"github.com/crypto-com/chainindex/adapter/tendermint/types"
)

type MockTendermintClient struct {
	mock.Mock
}

func (client *MockTendermintClient) Genesis() (*types.Genesis, error) {
	args := client.Called()
	return args.Get(0).(*types.Genesis), args.Error(1)
}

func (client *MockTendermintClient) LatestBlockHeight() (uint64, error) {
	args := client.Called()
	return args.Get(0).(uint64), args.Error(1)
}

func (client *MockTendermintClient) BlockResults(height uint64) (*types.BlockResults, error) {
	args := client.Called(height)
	return args.Get(0).(*types.BlockResults), args.Error(1)
}

func (client *MockTendermintClient) Block(height uint64) (*types.Block, error) {
	args := client.Called(height)
	return args.Get(0).(*types.Block), args.Error(1)
}

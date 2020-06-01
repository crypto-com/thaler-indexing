package syncservice

import (
	"time"

	"github.com/crypto-com/chainindex/adapter/tendermint"
	"github.com/crypto-com/chainindex/usecase"
)

// Subscribe to a blockheight news feed
type BlocksFeedSubscriber interface {
	Run(BlocksFeedSubscriberParams)
}

type BlocksFeedSubscriberParams struct {
	TendermintHeight RWSerialUint64

	OnExitCh chan<- bool
}

type DefaultBlocksFeedSubscriber struct {
	logger     usecase.Logger
	blocksFeed tendermint.BlocksFeed
}

func NewDefaultBlocksFeedSubscriber(
	logger usecase.Logger,
	blocksFeed tendermint.BlocksFeed,
) *DefaultBlocksFeedSubscriber {
	return &DefaultBlocksFeedSubscriber{
		logger: logger.WithFields(usecase.LogFields{
			"module": "DefaultBlocksFeedSubscriber",
		}),
		blocksFeed: blocksFeed,
	}
}

func (subscriber *DefaultBlocksFeedSubscriber) Run(params BlocksFeedSubscriberParams) {
	defer func() {
		if r := recover(); r != nil {
			subscriber.logger.Panicf("panic when running: %v", r)
		}
		subscriber.logger.Info("shutting down")
		params.OnExitCh <- true
	}()

	blocksFeedCh := make(chan uint64)
	if err := subscriber.blocksFeed.Subscribe(blocksFeedCh); err != nil {
		subscriber.logger.Panicf("error subscribing to blocks feed: %v", err)
	}
	for {
		select {
		case latestBlockHeight := <-blocksFeedCh:
			subscriber.logger.WithFields(usecase.LogFields{
				"latestBlockHeight": latestBlockHeight,
			}).Infof("received latest block height")

			params.TendermintHeight.SetIfLarger(latestBlockHeight)
		case <-time.After(1 * time.Minute):
			subscriber.logger.Debug("awaiting new block from blocks feed for a while")
		}
	}
}

package tendermint

import (
	"errors"
	"fmt"
	"time"

	tendermintadapter "github.com/crypto-com/chainindex/adapter/tendermint"
	"github.com/crypto-com/chainindex/usecase"
)

type PollingBlocksFeed struct {
	logger  usecase.Logger
	client  tendermintadapter.Client
	options PollingBlocksFeedOptions

	lastBroadcastedHeight uint64
	subscribed            bool
}

func NewPollingBlocksFeed(
	logger usecase.Logger,
	client tendermintadapter.Client,
	options PollingBlocksFeedOptions,
) *PollingBlocksFeed {
	return &PollingBlocksFeed{
		logger: logger.WithFields(usecase.LogFields{
			"module": "PollingBlocksFeed",
		}),
		client:  client,
		options: options,
	}
}

func (feed *PollingBlocksFeed) Subscribe(subscriberCh chan<- uint64) error {
	if feed.subscribed {
		return errors.New("error subscribing to blocks feed: already subscrbed")
	}

	feed.subscribed = true
	go func() {
		_ = feed.pollInBackground(subscriberCh)
	}()

	return nil
}

func (feed *PollingBlocksFeed) pollInBackground(subscriberCh chan<- uint64) error {
	for {
		err := feed.Poll(subscriberCh)
		if err != nil {
			feed.logger.Errorf("error polling latest block height: %v", err)
		}

		<-time.After(feed.options.PollingInterval)
	}
}

func (feed *PollingBlocksFeed) Poll(subscriberCh chan<- uint64) error {
	var err error

	var latestBlockHeight uint64
	latestBlockHeight, err = feed.client.LatestBlockHeight()
	if err != nil {
		return fmt.Errorf("error getting latest block height: %v", err)
	}

	logger := feed.logger.WithFields(usecase.LogFields{
		"latestBlockHeight":     latestBlockHeight,
		"lastBroadcastedHeight": feed.lastBroadcastedHeight,
	})
	logger.Debug("latest block height polled")

	if latestBlockHeight <= feed.lastBroadcastedHeight {
		return nil
	}

	logger.Info("broadcasting latest block height to subscriber")
	select {
	case subscriberCh <- latestBlockHeight:
		feed.lastBroadcastedHeight = latestBlockHeight
	default:
		feed.logger.Info("unable to send latest block height to subscribe: channel is blocked")
	}

	return nil
}

type PollingBlocksFeedOptions struct {
	PollingInterval time.Duration
}

package tendermint_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	tendermintadapter "github.com/crypto-com/chainindex/adapter/tendermint"
	. "github.com/crypto-com/chainindex/adapter/tendermint/test/mock"
	"github.com/crypto-com/chainindex/infrastructure/tendermint"
	. "github.com/crypto-com/chainindex/usecase/test/fake"
)

var _ = Describe("PollingBlocksFeed", func() {
	It("should implement BlocksFeed", func() {
		anyLogger := new(FakeLogger)
		anyTendermintClient := new(MockTendermintClient)
		anyOptions := tendermint.PollingBlocksFeedOptions{
			PollingInterval: 0 * time.Second,
		}

		var _ tendermintadapter.BlocksFeed = tendermint.NewPollingBlocksFeed(
			anyLogger,
			anyTendermintClient,
			anyOptions,
		)
	})

	Describe("Poll", func() {
		It("should send latest block height to the subscriber channel", func() {
			mockLogger := new(FakeLogger)
			mockTendermintClient := new(MockTendermintClient)
			options := tendermint.PollingBlocksFeedOptions{
				PollingInterval: 0 * time.Second,
			}

			feed := tendermint.NewPollingBlocksFeed(mockLogger, mockTendermintClient, options)

			anyBlockHeight := uint64(1000)
			mockTendermintClient.On("LatestBlockHeight").Return(anyBlockHeight, nil)

			subscriberCh := make(chan uint64, 1)
			err := feed.Poll(subscriberCh)

			Expect(err).To(BeNil())
			Eventually(subscriberCh).Should(Receive(Equal(anyBlockHeight)))
		})

		It("should not send duplicate block height to the subscriber channel", func() {
			var err error

			mockLogger := new(FakeLogger)
			mockTendermintClient := new(MockTendermintClient)
			options := tendermint.PollingBlocksFeedOptions{
				PollingInterval: 0 * time.Second,
			}

			feed := tendermint.NewPollingBlocksFeed(mockLogger, mockTendermintClient, options)

			anyBlockHeight := uint64(1000)
			mockTendermintClient.On("LatestBlockHeight").Return(anyBlockHeight, nil)

			subscriberCh := make(chan uint64, 1)
			err = feed.Poll(subscriberCh)
			Expect(err).To(BeNil())
			Eventually(subscriberCh).Should(Receive(Equal(anyBlockHeight)))

			err = feed.Poll(subscriberCh)
			Expect(err).To(BeNil())
			Consistently(subscriberCh).ShouldNot(Receive())
		})

		It("should not send block height to the subscriber channel when there is update", func() {
			var err error

			mockLogger := new(FakeLogger)
			mockTendermintClient := new(MockTendermintClient)
			options := tendermint.PollingBlocksFeedOptions{
				PollingInterval: 0 * time.Second,
			}

			feed := tendermint.NewPollingBlocksFeed(mockLogger, mockTendermintClient, options)

			anyBlockHeight := uint64(1000)
			mockTendermintClient.On("LatestBlockHeight").Once().Return(anyBlockHeight, nil)

			subscriberCh := make(chan uint64, 1)
			err = feed.Poll(subscriberCh)
			Expect(err).To(BeNil())
			Eventually(subscriberCh).Should(Receive(Equal(anyBlockHeight)))

			updatedBlockHeight := uint64(1001)
			mockTendermintClient.On("LatestBlockHeight").Return(updatedBlockHeight, nil)

			err = feed.Poll(subscriberCh)
			Expect(err).To(BeNil())
			Eventually(subscriberCh).Should(Receive(Equal(updatedBlockHeight)))
		})

		It("should not block when subscribe channel is blocked", func() {
			var err error

			mockLogger := new(FakeLogger)
			mockTendermintClient := new(MockTendermintClient)
			options := tendermint.PollingBlocksFeedOptions{
				PollingInterval: 0 * time.Second,
			}

			feed := tendermint.NewPollingBlocksFeed(mockLogger, mockTendermintClient, options)

			anyBlockHeight := uint64(1000)
			mockTendermintClient.On("LatestBlockHeight").Once().Return(anyBlockHeight, nil)

			subscriberCh := make(chan uint64, 1)
			err = feed.Poll(subscriberCh)
			Expect(err).To(BeNil())

			updatedBlockHeight := uint64(1001)
			mockTendermintClient.On("LatestBlockHeight").Return(updatedBlockHeight, nil)

			err = feed.Poll(subscriberCh)
			Expect(err).To(BeNil())
		})
	})
})

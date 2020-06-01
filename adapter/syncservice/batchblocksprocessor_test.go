package syncservice_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chainindex/adapter/syncservice"
	"github.com/crypto-com/chainindex/usecase"
	. "github.com/crypto-com/chainindex/usecase/test/factory"
	. "github.com/crypto-com/chainindex/usecase/test/fake"
)

var _ = Describe("Batchblocksworker", func() {
	Describe("BatchBlocksProcessorAggregator", func() {
		It("should pass nothing to the block data channel when the new block height data is not the consecutive one", func() {
			anyLogger := new(FakeLogger)
			anyInputBlockDataCh := make(chan *usecase.BlockData, 1)
			anyMaxSize := 3
			anyInitHeight := uint64(0)
			anyOutputBlockDataCh := make(chan *usecase.BlockData, 1)
			anyOnWorkerAvailableCh := make(chan int, 1)

			syncservice.NewBatchBlocksAggregator(syncservice.BatchBlocksAggregatorParams{
				Logger:           anyLogger,
				BlockDataInputCh: anyInputBlockDataCh,
				MaxSize:          anyMaxSize,
				LastSyncHeight:   anyInitHeight,

				BlockDataOutputCh: anyOutputBlockDataCh,

				OnWorkerAvailableCh: anyOnWorkerAvailableCh,
			})

			anyBlockData := RandomBlockData()
			anyInputBlockDataCh <- &anyBlockData

			Eventually(anyOutputBlockDataCh).ShouldNot(Receive())
		})
	})

	Describe("BatchBlocksProcessorSlidingWindow", func() {
		Describe("Insert", func() {
			It("should panic when map is full", func() {
				sizeOf3 := 3
				anyInitHeight := uint64(0)
				window := syncservice.NewBatchBlocksSlidingWindow(sizeOf3, anyInitHeight)

				InsertRandomBlockData(window, uint64(0))
				InsertRandomBlockData(window, uint64(1))
				InsertRandomBlockData(window, uint64(2))

				anyBlockData := RandomBlockData()
				Expect(func() {
					window.Insert(uint64(3), &anyBlockData)
				}).To(Panic())
			})

			It("should insert height to block data pointer mapping into the sliding window", func() {
				anySize := 5
				anyInitHeight := uint64(0)
				window := syncservice.NewBatchBlocksSlidingWindow(anySize, anyInitHeight)

				anyBlockData := RandomBlockData()
				window.Insert(uint64(0), &anyBlockData)

				actualBlockHeight, found := window.Get(uint64(0))
				Expect(found).To(BeTrue())
				Expect(actualBlockHeight).To(Equal(&anyBlockData))
			})
		})

		Describe("GetSuccessiveBlockData", func() {
			It("should return empty slice when the first block height data is not ready", func() {
				anySize := 3
				initHeight := uint64(0)
				window := syncservice.NewBatchBlocksSlidingWindow(anySize, initHeight)

				Expect(window.PopSuccessiveBlockData()).To(BeEmpty())
			})

			It("should return empty slice even when the second block height data is ready", func() {
				anySize := 3
				initHeight := uint64(0)
				window := syncservice.NewBatchBlocksSlidingWindow(anySize, initHeight)

				InsertRandomBlockData(window, uint64(1))

				Expect(window.PopSuccessiveBlockData()).To(BeEmpty())
			})

			It("should return first block height data when it is ready but second is not", func() {
				anySize := 3
				initHeight := uint64(0)
				window := syncservice.NewBatchBlocksSlidingWindow(anySize, initHeight)

				anyBlockData := RandomBlockData()
				window.Insert(uint64(0), &anyBlockData)

				actualBlockDatas := window.PopSuccessiveBlockData()
				Expect(actualBlockDatas).To(HaveLen(1))
				Expect(actualBlockDatas[0]).To(Equal(&anyBlockData))
			})

			It("should return consecutive block height data starting from first height up to last consecutive ready data", func() {
				anySize := 5
				initHeight := uint64(0)
				window := syncservice.NewBatchBlocksSlidingWindow(anySize, initHeight)

				anyFirstBlockData := RandomBlockData()
				window.Insert(uint64(0), &anyFirstBlockData)

				anySecondBlockData := RandomBlockData()
				window.Insert(uint64(1), &anySecondBlockData)

				actualBlockDatas := window.PopSuccessiveBlockData()
				Expect(actualBlockDatas).To(HaveLen(2))
				Expect(actualBlockDatas).To(Equal([]*usecase.BlockData{
					&anyFirstBlockData, &anySecondBlockData,
				}))
			})

			It("should not return block height data when their precedent block height data is not ready", func() {
				anySize := 5
				initHeight := uint64(0)
				window := syncservice.NewBatchBlocksSlidingWindow(anySize, initHeight)

				anyBlockData := RandomBlockData()
				window.Insert(uint64(0), &anyBlockData)

				InsertRandomBlockData(window, uint64(2))
				InsertRandomBlockData(window, uint64(3))

				actualBlockDatas := window.PopSuccessiveBlockData()
				Expect(actualBlockDatas).To(HaveLen(1))
				Expect(actualBlockDatas).To(Equal([]*usecase.BlockData{
					&anyBlockData,
				}))
			})

			It("should delete returned block height data", func() {
				anySize := 5
				initHeight := uint64(0)
				window := syncservice.NewBatchBlocksSlidingWindow(anySize, initHeight)

				anyFirstBlockData := RandomBlockData()
				window.Insert(uint64(0), &anyFirstBlockData)

				anySecondBlockData := RandomBlockData()
				window.Insert(uint64(1), &anySecondBlockData)

				actualBlockDatas := window.PopSuccessiveBlockData()
				Expect(actualBlockDatas).To(HaveLen(2))

				blockDataAfterFirstGet := window.PopSuccessiveBlockData()
				Expect(blockDataAfterFirstGet).To(HaveLen(0))
			})

			// It("should free the space after successive block height data is retrieved", func () {
			// 	var err error
			// 	sizeOf3 := uint32(3)
			// 	initHeight := uint64(0)
			// 	windowOf3Space := syncservice.NewBatchBlocksProcessorSlidingWindow(sizeOf3, initHeight)

			// 	InsertRandomBlockData(window, uint64(0))
			// 	InsertRandomBlockData(window, uint64(1))
			// 	InsertRandomBlockData(window, uint64(2))
			// })
		})
	})
})

func InsertRandomBlockData(window *syncservice.BatchBlocksProcessorSlidingWindow, height uint64) {
	anyBlockData := RandomBlockData()
	window.Insert(height, &anyBlockData)
}

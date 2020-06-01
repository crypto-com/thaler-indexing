package syncservice

import (
	"sync"
	"time"

	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/adapter/tendermint"
	tenderminttypes "github.com/crypto-com/chainindex/adapter/tendermint/types"
	"github.com/crypto-com/chainindex/usecase"
)

type BatchBlocksProcessor struct {
	logger usecase.Logger
	client tendermint.Client

	totalWorkingWorker    int
	maxWorker             int
	lastDistributedHeight uint64
}

func NewBatchBlocksProcessor(logger usecase.Logger, client tendermint.Client, maxWorker int) *BatchBlocksProcessor {
	return &BatchBlocksProcessor{
		logger: logger.WithFields(usecase.LogFields{
			"module": "BatchBlocksProcessor",
		}),
		client: client,

		totalWorkingWorker: 0,
		maxWorker:          maxWorker,
	}
}

func (processor *BatchBlocksProcessor) Run(params BlocksProcessorParams) {
	defer func() {
		if r := recover(); r != nil {
			processor.logger.Panicf("panic when running: %v", r)
		}
		processor.logger.Info("shutting down")
		params.OnExitCh <- true
	}()

	processor.logger.Infof("staring from last synchronized height: %d", params.LastSyncHeight)

	processor.lastDistributedHeight = params.LastSyncHeight
	aggregatorBlockDataInputCh := make(chan *usecase.BlockData, processor.maxWorker)
	onWorkerAvailableCh := make(chan int, processor.maxWorker)

	aggregator := NewBatchBlocksAggregator(BatchBlocksAggregatorParams{
		Logger: processor.logger,

		BlockDataInputCh: aggregatorBlockDataInputCh,
		MaxSize:          processor.maxWorker,
		LastSyncHeight:   params.LastSyncHeight,

		BlockDataOutputCh: params.BlockDataCh,

		OnWorkerAvailableCh: onWorkerAvailableCh,
	})
	go aggregator.Run()

	// Start the worker
	onWorkerAvailableCh <- 0

	for {
		availableWorkerSize := <-onWorkerAvailableCh
		processor.totalWorkingWorker -= availableWorkerSize

		processor.DistributeBlocksToWorker(&params, aggregatorBlockDataInputCh)
	}
}

func (processor *BatchBlocksProcessor) DistributeBlocksToWorker(params *BlocksProcessorParams, aggregatorBlockDataCh chan<- *usecase.BlockData) {
	latestTendermintBlockHeight := params.TendermintHeight.Get()
	if processor.lastDistributedHeight == latestTendermintBlockHeight {
		processor.logger.Debug("processor has free worker but is blocked because of no new block")

		<-params.OnTendermintHeightUpdate
		latestTendermintBlockHeight = params.TendermintHeight.Get()
	}

	for processor.lastDistributedHeight < latestTendermintBlockHeight {
		nextHeightToHandle := processor.lastDistributedHeight + 1
		processor.logger.Debugf("trying to distribute block height %d", nextHeightToHandle)

		if processor.totalWorkingWorker == processor.maxWorker {
			processor.logger.Debug("processor worker size already exceed maximum worker allowed")
			return
		}

		worker := NewBatchBlocksWorker(
			processor.logger,
			processor.client,

			nextHeightToHandle,

			aggregatorBlockDataCh,
		)

		go worker.Run()

		processor.totalWorkingWorker += 1
		processor.lastDistributedHeight += 1
	}
}

type BatchBlocksWorker struct {
	logger usecase.Logger
	client tendermint.Client

	height uint64

	blockDataCh chan<- *usecase.BlockData
}

func NewBatchBlocksWorker(logger usecase.Logger, client tendermint.Client, height uint64, blockDataCh chan<- *usecase.BlockData) *BatchBlocksWorker {
	return &BatchBlocksWorker{
		logger: logger.WithFields(usecase.LogFields{
			"module":      "BatchBlocksWorker",
			"blockHeight": height,
		}),
		client: client,

		height: height,

		blockDataCh: blockDataCh,
	}
}

func (worker *BatchBlocksWorker) Run() {
	for {
		err := worker.processBlock()
		if err == nil {
			break
		}

		<-time.After(5 * time.Second)
	}
}

func (worker *BatchBlocksWorker) processBlock() error {
	var err error

	logger := worker.logger.WithFields(usecase.LogFields{
		"blockHeight": worker.height,
	})
	logger.Info("processing block")

	var blockData *usecase.BlockData
	if worker.height == uint64(1) {
		blockData, err = worker.handleGenesisBlock()
	} else {
		blockData, err = worker.handleBlock(worker.height)
	}
	if err != nil {
		logger.Errorf("error processing block: %v", err)
		return err
	}

	logger.WithFields(usecase.LogFields{
		"blockData": blockData,
	}).Debug("processed block data")

	worker.blockDataCh <- blockData

	return nil
}

// TODO: Extract to standalone struct
func (worker *BatchBlocksWorker) handleGenesisBlock() (*usecase.BlockData, error) {
	var err error

	var genesis *tenderminttypes.Genesis
	genesis, err = worker.client.Genesis()
	if err != nil {
		return nil, err
	}

	var block *tenderminttypes.Block
	block, err = worker.client.Block(uint64(1))
	if err != nil {
		return nil, err
	}

	blockData := adapter.ParseGenesisToBlockData(adapter.TendermintGenesisBlockData{
		Genesis: genesis,
		Block:   block,
	})

	return blockData, nil
}

func (worker *BatchBlocksWorker) handleBlock(height uint64) (*usecase.BlockData, error) {
	var err error

	var block *tenderminttypes.Block
	block, err = worker.client.Block(height)
	if err != nil {
		return nil, err
	}
	var blockResults *tenderminttypes.BlockResults
	blockResults, err = worker.client.BlockResults(height)
	if err != nil {
		return nil, err
	}

	blockData := adapter.ParseBlockToBlockData(adapter.TendermintBlockData{
		Block:        block,
		BlockResults: blockResults,
	})

	return blockData, nil
}

type BatchBlocksAggregator struct {
	logger usecase.Logger

	blockDataInputCh <-chan *usecase.BlockData
	slidingWindow    *BatchBlocksProcessorSlidingWindow

	blockDataOutputCh   chan<- *usecase.BlockData
	onWorkerAvailableCh chan<- int
}

func NewBatchBlocksAggregator(params BatchBlocksAggregatorParams) *BatchBlocksAggregator {
	return &BatchBlocksAggregator{
		logger: params.Logger.WithFields(usecase.LogFields{
			"module": "BatchBlocksAggregator",
		}),

		blockDataInputCh: params.BlockDataInputCh,
		slidingWindow:    NewBatchBlocksSlidingWindow(params.MaxSize, params.LastSyncHeight+1),

		blockDataOutputCh:   params.BlockDataOutputCh,
		onWorkerAvailableCh: params.OnWorkerAvailableCh,
	}
}

type BatchBlocksAggregatorParams struct {
	Logger usecase.Logger

	BlockDataInputCh <-chan *usecase.BlockData
	MaxSize          int
	LastSyncHeight   uint64

	BlockDataOutputCh chan<- *usecase.BlockData

	OnWorkerAvailableCh chan<- int
}

func (aggregator *BatchBlocksAggregator) Run() {
	for {
		aggregator.ProcessBlockDataInput()
	}
}

func (aggregator *BatchBlocksAggregator) ProcessBlockDataInput() {
	blockData := <-aggregator.blockDataInputCh
	aggregator.logger.WithFields(usecase.LogFields{
		"blockHeight": blockData.Block.Height,
	}).Debug("received block data from worker")

	aggregator.slidingWindow.Insert(blockData.Block.Height, blockData)
	someBlockData := aggregator.slidingWindow.PopSuccessiveBlockData()
	aggregator.logger.WithFields(usecase.LogFields{
		"blockDataSize": len(someBlockData),
	}).Debug("going to notify worker available")

	aggregator.onWorkerAvailableCh <- len(someBlockData)

	aggregator.logger.Debug("going to send block data to downstream")
	for _, blockData := range someBlockData {
		aggregator.blockDataOutputCh <- blockData
	}
}

type BatchBlocksProcessorSlidingWindow struct {
	sync.RWMutex
	data map[uint64]*usecase.BlockData

	firstHeight uint64

	maxSize int
}

func NewBatchBlocksSlidingWindow(maxSize int, initHeight uint64) *BatchBlocksProcessorSlidingWindow {
	return &BatchBlocksProcessorSlidingWindow{
		data: make(map[uint64]*usecase.BlockData),

		firstHeight: initHeight,

		maxSize: maxSize,
	}
}

func (window *BatchBlocksProcessorSlidingWindow) Insert(height uint64, blockData *usecase.BlockData) {
	window.Lock()
	defer window.Unlock()

	if len(window.data) == window.maxSize {
		panic("error inserting block data into sliding window: already full")
	}

	_, exist := window.data[height]
	if !exist {
		window.data[height] = blockData
	}
}

func (window *BatchBlocksProcessorSlidingWindow) Get(height uint64) (*usecase.BlockData, bool) {
	window.RLock()
	defer window.RUnlock()

	blockData, ok := window.data[height]

	if !ok {
		return nil, false
	}

	return blockData, true
}

func (window *BatchBlocksProcessorSlidingWindow) PopSuccessiveBlockData() []*usecase.BlockData {
	window.Lock()
	defer window.Unlock()

	result := make([]*usecase.BlockData, 0)

	for i := window.firstHeight; ; i += 1 {
		blockData, ok := window.data[i]
		if !ok {
			return result
		}

		result = append(result, blockData)
		delete(window.data, i)
		window.firstHeight += 1
	}
}

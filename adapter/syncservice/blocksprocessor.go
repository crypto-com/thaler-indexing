package syncservice

import (
	"time"

	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/adapter/tendermint"
	tenderminttypes "github.com/crypto-com/chainindex/adapter/tendermint/types"
	"github.com/crypto-com/chainindex/usecase"
)

type BlocksProcessor interface {
	Run(BlocksProcessorParams)
}

type BlocksProcessorParams struct {
	LastSyncHeight   uint64
	TendermintHeight RWSerialUint64

	OnTendermintHeightUpdate <-chan bool

	BlockDataCh chan<- *usecase.BlockData
	OnExitCh    chan<- bool
}

type PeriodicBlocksProcessor struct {
	logger usecase.Logger
	client tendermint.Client

	nextHeightAwaiting uint64
}

func NewPeriodicBlocksWorker(
	logger usecase.Logger,

	client tendermint.Client,
) *PeriodicBlocksProcessor {
	return &PeriodicBlocksProcessor{
		logger: logger.WithFields(usecase.LogFields{
			"module": "PeriodicBlocksWorker",
		}),

		client: client,
	}
}

func (worker *PeriodicBlocksProcessor) Run(params BlocksProcessorParams) {
	defer func() {
		if r := recover(); r != nil {
			worker.logger.Panicf("panic when running: %v", r)
		}
		worker.logger.Info("shutting down")
		params.OnExitCh <- true
	}()

	worker.nextHeightAwaiting = params.LastSyncHeight + 1

	for {
		_ = worker.ProcessNewBlocks(&params)

		<-time.After(5 * time.Second)
	}

}

func (worker *PeriodicBlocksProcessor) ProcessNewBlocks(params *BlocksProcessorParams) error {
	var err error

	for {
		latestBlockHeight := params.TendermintHeight.Get()

		worker.logger.WithFields(usecase.LogFields{
			"nextHeightAwaiting": worker.nextHeightAwaiting,
			"latestBlockHeight":  latestBlockHeight,
		}).Debug("block height status")
		if worker.nextHeightAwaiting > latestBlockHeight {
			break
		}

		for worker.nextHeightAwaiting <= latestBlockHeight {
			logger := worker.logger.WithFields(usecase.LogFields{
				"blockHeight": worker.nextHeightAwaiting,
			})
			logger.Info("processing block")

			var blockData *usecase.BlockData
			if worker.nextHeightAwaiting == uint64(1) {
				blockData, err = worker.handleGenesisBlock()
			} else {
				blockData, err = worker.handleBlock(worker.nextHeightAwaiting)
			}
			if err != nil {
				logger.Errorf("error processing block: %v", err)
				return err
			}

			logger.WithFields(usecase.LogFields{
				"blockData": blockData,
			}).Debug("processed block data")

			params.BlockDataCh <- blockData

			worker.nextHeightAwaiting += 1
		}
	}

	return nil
}

func (worker *PeriodicBlocksProcessor) handleGenesisBlock() (*usecase.BlockData, error) {
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

func (worker *PeriodicBlocksProcessor) handleBlock(height uint64) (*usecase.BlockData, error) {
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

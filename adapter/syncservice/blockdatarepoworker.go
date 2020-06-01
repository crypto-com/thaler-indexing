package syncservice

import (
	"time"

	"github.com/crypto-com/chainindex/usecase"
)

type BlockDataRepoWorker interface {
	Run(params BlockDataRepoWorkerParams)
}

type BlockDataRepoWorkerParams struct {
	BlockDataCh     <-chan *usecase.BlockData
	OnBlockStoredCh chan<- uint64
	OnExitCh        chan<- bool
}

type DefaultBlockDataRepoWorker struct {
	logger        usecase.Logger
	blockDataRepo usecase.BlockDataRepository
}

func NewDefaultBlockDataRepoWorker(
	logger usecase.Logger,
	blockDataRepo usecase.BlockDataRepository,
) *DefaultBlockDataRepoWorker {
	return &DefaultBlockDataRepoWorker{
		logger: logger.WithFields(usecase.LogFields{
			"module": "DefaultBlockDataRepoWorker",
		}),
		blockDataRepo: blockDataRepo,
	}
}

func (worker *DefaultBlockDataRepoWorker) Run(params BlockDataRepoWorkerParams) {
	var err error

	defer func() {
		if r := recover(); r != nil {
			worker.logger.Panicf("panic when running: %v", r)
		}
		worker.logger.Info("shutting down")
		params.OnExitCh <- true
	}()

	for {
		blockData := <-params.BlockDataCh
		for {
			if err = worker.processBlockData(blockData); err != nil {
				<-time.After(5 * time.Second)
				continue
			}

			params.OnBlockStoredCh <- blockData.Block.Height
			break
		}
	}
}

func (worker *DefaultBlockDataRepoWorker) processBlockData(blockData *usecase.BlockData) error {
	logger := worker.logger.WithFields(usecase.LogFields{
		"blockHeight": blockData.Block.Height,
	})
	logger.Info("storing block data")

	err := worker.blockDataRepo.Store(blockData)
	if err != nil {
		logger.Errorf("error storing block data: %v", err)
		return err
	}

	return nil
}

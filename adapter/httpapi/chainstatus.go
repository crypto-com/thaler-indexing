package httpapi

import (
	"net/http"

	"github.com/crypto-com/chainindex/internal/bignum"
	"github.com/crypto-com/chainindex/usecase"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
)

type ChainStatus struct {
	Version                string          `json:"version"`
	SyncBlockHeight        uint64          `json:"sync_block_height"`
	TendermintBlockHeight  uint64          `json:"tendermint_block_height"`
	TransactionCount       uint64          `json:"transaction_count"`
	TotalReward            *bignum.WBigInt `json:"total_reward"`
	TotalRewardMinted      *bignum.WBigInt `json:"total_reward_minted"`
	CouncilNodeCount       uint64          `json:"council_node_count"`
	TotalCouncilNodeStaked *bignum.WBigInt `json:"total_council_node_staked"`
}

type ChainStatusReward struct {
	Minted *bignum.WBigInt `json:"minted"`
	Total  *bignum.WBigInt `json:"total"`
}

type ChainStatusHandler struct {
	logger usecase.Logger

	syncService     usecase.SyncService
	activityView    viewrepo.ActivityViewRepo
	rewardView      viewrepo.RewardViewRepo
	councilNodeView viewrepo.CouncilNodeViewRepo
}

func NewChainStatusHandler(
	logger usecase.Logger,

	syncService usecase.SyncService,
	activityView viewrepo.ActivityViewRepo,
	rewardView viewrepo.RewardViewRepo,
	councilNodeView viewrepo.CouncilNodeViewRepo,
) *ChainStatusHandler {
	return &ChainStatusHandler{
		logger: logger.WithFields(usecase.LogFields{
			"module": "ChainStatusHandler",
		}),

		syncService:     syncService,
		activityView:    activityView,
		rewardView:      rewardView,
		councilNodeView: councilNodeView,
	}
}

func (handler *ChainStatusHandler) GetChainStatus(resp http.ResponseWriter, req *http.Request) {
	var err error
	var chainStatus ChainStatus

	chainStatus.Version = "0.5.2"

	syncStatus := handler.syncService.GetStatus()
	chainStatus.TendermintBlockHeight = syncStatus.TendermintBlockHeight
	chainStatus.SyncBlockHeight = syncStatus.SyncBlockHeight

	chainStatus.TransactionCount, err = handler.activityView.TransactionsCount()
	if err != nil {
		handler.logger.Errorf("error querying transactions count: %v", err)
		InternalServerError(resp)
		return
	}

	chainStatus.TotalRewardMinted, err = handler.rewardView.TotalMinted()
	if err != nil {
		handler.logger.Errorf("error querying transactions count: %v", err)
		InternalServerError(resp)
		return
	}
	chainStatus.TotalReward, err = handler.rewardView.Total()
	if err != nil {
		handler.logger.Errorf("error querying transactions count: %v", err)
		InternalServerError(resp)
		return
	}

	councilNodeStats, err := handler.councilNodeView.Stats()
	if err != nil {
		handler.logger.Errorf("error querying council node status: %v", err)
		InternalServerError(resp)
		return
	}
	chainStatus.CouncilNodeCount = councilNodeStats.Count
	chainStatus.TotalCouncilNodeStaked = councilNodeStats.TotalStaked

	Success(resp, chainStatus)
}

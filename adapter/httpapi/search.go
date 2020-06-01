package httpapi

import (
	"errors"
	"net/http"

	"github.com/crypto-com/chainindex/usecase"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
)

type SearchHandler struct {
	logger usecase.Logger

	activityView       viewrepo.ActivityViewRepo
	blockView          viewrepo.BlockViewRepo
	stakingAccountView viewrepo.StakingAccountViewRepo
	councilNodeView    viewrepo.CouncilNodeViewRepo
}

type SearchAllResult struct {
	Blocks          []viewrepo.Block
	Transactions    []viewrepo.Transaction
	Events          []viewrepo.Event
	StakingAccounts []viewrepo.StakingAccount
	CouncilNodes    []viewrepo.CouncilNode
}

func NewSearchHandler(
	logger usecase.Logger,

	activityView viewrepo.ActivityViewRepo,
	blockView viewrepo.BlockViewRepo,
	stakingAccountView viewrepo.StakingAccountViewRepo,
	councilNodeView viewrepo.CouncilNodeViewRepo,
) *SearchHandler {
	return &SearchHandler{
		logger: logger.WithFields(usecase.LogFields{
			"module": "SearchHandler",
		}),

		activityView:       activityView,
		blockView:          blockView,
		stakingAccountView: stakingAccountView,
		councilNodeView:    councilNodeView,
	}
}

func (handler *SearchHandler) All(resp http.ResponseWriter, req *http.Request) {
	var err error

	keyword := req.URL.Query().Get("query")
	if keyword == "" {
		BadRequest(resp, errors.New("missing query"))
		return
	}

	pagination := viewrepo.NewOffsetPagination(uint64(1), uint64(5))

	blocks, _, err := handler.blockView.Search(keyword, pagination)
	if err != nil {
		handler.logger.Errorf("error searching blocks: %v", err)
		InternalServerError(resp)
		return
	}
	transactions, _, err := handler.activityView.SearchTransactions(keyword, pagination)
	if err != nil {
		handler.logger.Errorf("error searching transactions: %v", err)
		InternalServerError(resp)
		return
	}
	events, _, err := handler.activityView.SearchEvents(keyword, pagination)
	if err != nil {
		handler.logger.Errorf("error searching events: %v", err)
		InternalServerError(resp)
		return
	}
	stakingAccounts, _, err := handler.stakingAccountView.Search(keyword, pagination)
	if err != nil {
		handler.logger.Errorf("error searching staking accounts: %v", err)
		InternalServerError(resp)
		return
	}
	councilNodes, _, err := handler.councilNodeView.Search(keyword, pagination)
	if err != nil {
		handler.logger.Errorf("error searching counil nodes: %v", err)
		InternalServerError(resp)
		return
	}

	Success(resp, SearchAllResult{
		Blocks:          blocks,
		Transactions:    transactions,
		Events:          events,
		StakingAccounts: stakingAccounts,
		CouncilNodes:    councilNodes,
	})
}

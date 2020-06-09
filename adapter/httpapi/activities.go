package httpapi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/crypto-com/chainindex"
	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/usecase"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
)

type ActivitiesHandler struct {
	logger usecase.Logger

	routePath    RoutePath
	activityView viewrepo.ActivityViewRepo
}

func NewActivitiesHandler(logger usecase.Logger, routePath RoutePath, activityView viewrepo.ActivityViewRepo) *ActivitiesHandler {
	return &ActivitiesHandler{
		logger: logger.WithFields(usecase.LogFields{
			"module": "ActivitiesHandler",
		}),

		routePath:    routePath,
		activityView: activityView,
	}
}

func (handler *ActivitiesHandler) ListTransactions(resp http.ResponseWriter, req *http.Request) {
	var err error

	pagination, err := ParsePagination(req)
	if err != nil {
		BadRequest(resp, err)
		return
	}

	filter := viewrepo.TransactionFilter{
		MaybeTypes: make([]chainindex.TransactionType, 0),
	}

	filterTypes := req.URL.Query().Get("filter[type]")
	if filterTypes != "" {
		filterTypeInputs := strings.Split(filterTypes, ",")
		for _, input := range filterTypeInputs {
			if !adapter.IsValidTransactionType(input) {
				BadRequest(resp, fmt.Errorf("invalid transaction type filter: %s", input))
				return
			}

			filter.MaybeTypes = append(filter.MaybeTypes, adapter.StringToTransactionType(input))
		}
	}

	filterStakingAccountAddress := req.URL.Query().Get("filter[staking_account_address]")
	if filterStakingAccountAddress != "" {
		filter.MaybeStakingAccountAddress = &filterStakingAccountAddress
	}

	blockTransactions, paginationResult, err := handler.activityView.ListTransactions(filter, pagination)
	if err != nil {
		handler.logger.Errorf("error listing transactions: %v", err)
		InternalServerError(resp)
		return
	}

	SuccessWithPagination(resp, blockTransactions, paginationResult)
}

func (handler *ActivitiesHandler) FindTransactionByTxId(resp http.ResponseWriter, req *http.Request) {
	var err error

	routeVars := handler.routePath.Vars(req)
	txid, ok := routeVars["txid"]
	if !ok {
		BadRequest(resp, errors.New("missing txid path parameter"))
		return
	}

	block, err := handler.activityView.FindTransactionByTxId(txid)
	if err != nil {
		if err == adapter.ErrNotFound {
			NotFound(resp)
			return
		}
		handler.logger.Errorf("error finding transaction by txid: %v", err)
		InternalServerError(resp)
		return
	}

	Success(resp, block)
}

func (handler *ActivitiesHandler) ListEvents(resp http.ResponseWriter, req *http.Request) {
	var err error

	pagination, err := ParsePagination(req)
	if err != nil {
		BadRequest(resp, err)
		return
	}

	filter := viewrepo.EventFilter{
		MaybeTypes: make([]chainindex.TransactionType, 0),
	}
	filterTypes := req.URL.Query().Get("filter[type]")
	if filterTypes != "" {
		filterTypeInputs := strings.Split(filterTypes, ",")
		for _, input := range filterTypeInputs {
			if !adapter.IsValidEventType(input) {
				BadRequest(resp, fmt.Errorf("invalid event type filter: %s", input))
				return
			}

			filter.MaybeTypes = append(filter.MaybeTypes, adapter.StringToEventType(input))
		}
	}

	blockTransactions, paginationResult, err := handler.activityView.ListEvents(filter, pagination)
	if err != nil {
		handler.logger.Errorf("error listing events: %v", err)
		InternalServerError(resp)
		return
	}

	SuccessWithPagination(resp, blockTransactions, paginationResult)
}

func (handler *ActivitiesHandler) FindEventByBlockHeightEventPosition(resp http.ResponseWriter, req *http.Request) {
	var err error

	routeVars := handler.routePath.Vars(req)

	blockHeightVar, ok := routeVars["height"]
	if !ok {
		BadRequest(resp, errors.New("missing id path parameter"))
		return
	}
	blockHeight, err := strconv.ParseUint(blockHeightVar, 10, 64)
	if err != nil {
		BadRequest(resp, errors.New("invalid id path parameter"))
		return
	}

	eventPositionVar, ok := routeVars["position"]
	if !ok {
		BadRequest(resp, errors.New("missing id path parameter"))
		return
	}
	eventPosition, err := strconv.ParseUint(eventPositionVar, 10, 64)
	if err != nil {
		BadRequest(resp, errors.New("invalid id path parameter"))
		return
	}

	block, err := handler.activityView.FindEventByBlockHeightEventPosition(blockHeight, eventPosition)
	if err != nil {
		if err == adapter.ErrNotFound {
			NotFound(resp)
			return
		}
		handler.logger.Errorf("error finding event by id: %v", err)
		InternalServerError(resp)
		return
	}

	Success(resp, block)
}

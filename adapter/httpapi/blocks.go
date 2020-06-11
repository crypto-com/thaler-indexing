package httpapi

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/usecase"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
)

type BlocksHandler struct {
	logger usecase.Logger

	routePath RoutePath
	blockView viewrepo.BlockViewRepo
}

func NewBlocksHandler(logger usecase.Logger, routePath RoutePath, blockView viewrepo.BlockViewRepo) *BlocksHandler {
	return &BlocksHandler{
		logger: logger.WithFields(usecase.LogFields{
			"module": "BlocksHandler",
		}),

		routePath: routePath,
		blockView: blockView,
	}
}

func (handler *BlocksHandler) ListBlocks(resp http.ResponseWriter, req *http.Request) {
	var err error

	pagination, err := ParsePagination(req)
	if err != nil {
		BadRequest(resp, err)
		return
	}

	filter := viewrepo.BlockFilter{
		MaybeProposers: make([]uint64, 0),
	}
	filterProposerIds := req.URL.Query().Get("proposer_ids")
	if filterProposerIds != "" {
		filterProposerIdInputs := strings.Split(filterProposerIds, ",")
		for _, input := range filterProposerIdInputs {
			var id uint64
			id, err = strconv.ParseUint(input, 10, 64)
			if err != nil {
				BadRequest(resp, errors.New("invalid proposer id filter"))
				return
			}

			filter.MaybeProposers = append(filter.MaybeProposers, id)
		}
	}

	blocks, paginationResult, err := handler.blockView.ListBlocks(filter, pagination)
	if err != nil {
		handler.logger.Errorf("error listing blocks: %v", err)
		InternalServerError(resp)
		return
	}

	SuccessWithPagination(resp, blocks, paginationResult)
}

func (handler *BlocksHandler) FindBlock(resp http.ResponseWriter, req *http.Request) {
	var err error

	routeVars := handler.routePath.Vars(req)
	var blockIdentity *viewrepo.BlockIdentity
	if blockIdentity, err = parseBlockIdentity(routeVars); err != nil {
		BadRequest(resp, err)
		return
	}

	block, err := handler.blockView.FindBlock(*blockIdentity)
	if err != nil {
		if err == adapter.ErrNotFound {
			NotFound(resp)
			return
		}
		handler.logger.Errorf("error finding block by height: %v", err)
		InternalServerError(resp)
		return
	}

	Success(resp, block)
}

func (handler *BlocksHandler) ListBlockTransactions(resp http.ResponseWriter, req *http.Request) {
	var err error

	routeVars := handler.routePath.Vars(req)
	var blockIdentity *viewrepo.BlockIdentity
	if blockIdentity, err = parseBlockIdentity(routeVars); err != nil {
		BadRequest(resp, err)
		return
	}

	pagination, err := ParsePagination(req)
	if err != nil {
		BadRequest(resp, err)
		return
	}

	blockTransactions, paginationResult, err := handler.blockView.ListBlockTransactions(
		*blockIdentity, pagination,
	)
	if err != nil {
		if err == adapter.ErrNotFound {
			NotFound(resp)
			return
		}
		handler.logger.Errorf("error listing block transactions: %v", err)
		InternalServerError(resp)
		return
	}

	SuccessWithPagination(resp, blockTransactions, paginationResult)
}

func (handler *BlocksHandler) ListBlockEvents(resp http.ResponseWriter, req *http.Request) {
	var err error

	routeVars := handler.routePath.Vars(req)
	var blockIdentity *viewrepo.BlockIdentity
	if blockIdentity, err = parseBlockIdentity(routeVars); err != nil {
		BadRequest(resp, err)
		return
	}

	pagination, err := ParsePagination(req)
	if err != nil {
		BadRequest(resp, err)
		return
	}

	blockEvents, paginationResult, err := handler.blockView.ListBlockEvents(
		*blockIdentity, pagination,
	)
	if err != nil {
		if err == adapter.ErrNotFound {
			NotFound(resp)
			return
		}
		handler.logger.Errorf("error listing block events: %v", err)
		InternalServerError(resp)
		return
	}

	SuccessWithPagination(resp, blockEvents, paginationResult)
}

func parseBlockIdentity(routeVars map[string]string) (*viewrepo.BlockIdentity, error) {
	identityVar, ok := routeVars["hash_or_height"]
	if !ok {
		return nil, errors.New("missing block identity path parameter")
	}
	var blockIdentity viewrepo.BlockIdentity
	if len(identityVar) == 64 {
		blockIdentity.MaybeHash = &identityVar
	} else if containsNonDigit(identityVar) {
		return nil, errors.New("invalid hash path parameter")
	} else {
		height, err := strconv.ParseUint(identityVar, 10, 64)
		if err != nil {
			return nil, errors.New("invalid height path parameter")
		}

		if height == 0 {
			return nil, errors.New("block height cannot be 0")
		}
		blockIdentity.MaybeHeight = &height
	}

	return &blockIdentity, nil
}

var containsNonDigit = regexp.MustCompile(`[^0-9]+`).MatchString

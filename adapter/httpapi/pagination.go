package httpapi

import (
	"net/http"
	"strconv"

	"github.com/crypto-com/chainindex/usecase/viewrepo"
)

func ParsePagination(req *http.Request) (*viewrepo.Pagination, error) {
	var err error

	var pagination viewrepo.PaginationType
	var page, limit uint64

	pagination = req.URL.Query().Get("pagination")
	if pagination == "" {
		pagination = viewrepo.PAGINATION_OFFSET
	}
	if pagination != viewrepo.PAGINATION_OFFSET {
		return nil, ErrInvalidPagination
	}

	pageQuery := req.URL.Query().Get("page")
	if pageQuery == "" {
		page = uint64(1)
	} else {
		page, err = strconv.ParseUint(pageQuery, 10, 64)
		if err != nil {
			return nil, ErrInvalidPage
		}
		if page == 0 {
			return nil, ErrInvalidPage
		}
	}

	limitQuery := req.URL.Query().Get("limit")
	if limitQuery == "" {
		limit = uint64(20)
	} else {
		limit, err = strconv.ParseUint(limitQuery, 10, 64)
		if err != nil {
			return nil, ErrInvalidPage
		}
	}

	return viewrepo.NewOffsetPagination(page, limit), nil
}

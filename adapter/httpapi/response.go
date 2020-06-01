package httpapi

import (
	"net/http"

	"github.com/crypto-com/chainindex/usecase/viewrepo"
	jsoniter "github.com/json-iterator/go"
)

func Success(resp http.ResponseWriter, result interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	err := jsoniter.NewEncoder(resp).Encode(Response{
		Result: result,
		Err:    "",
	})
	if err != nil {
		InternalServerError(resp)
	}
}

func SuccessWithPagination(
	resp http.ResponseWriter,
	result interface{},
	paginationResult *viewrepo.PaginationResult,
) {
	resp.Header().Set("Content-Type", "application/json")
	err := jsoniter.NewEncoder(resp).Encode(PagedResponse{
		Response: Response{
			Result: result,
			Err:    "",
		},
		OffsetPagination: OptPaginationOffsetResponseFromResult(paginationResult.OffsetResult()),
	})
	if err != nil {
		InternalServerError(resp)
	}
}

func NotFound(
	resp http.ResponseWriter,
) {
	resp.Header().Set("Content-Type", "application/json")
	message, err := jsoniter.Marshal(Response{
		Err: "Record not found",
	})
	if err != nil {
		InternalServerError(resp)
		return
	}

	http.Error(resp, string(message), 404)
}

func BadRequest(resp http.ResponseWriter, errResp error) {
	resp.Header().Set("Content-Type", "application/json")
	message, err := jsoniter.Marshal(Response{
		Err: errResp.Error(),
	})
	if err != nil {
		InternalServerError(resp)
		return
	}

	http.Error(resp, string(message), 400)
}

func InternalServerError(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
	message, _ := jsoniter.Marshal(Response{
		Err: ErrInternalServerError.Error(),
	})
	http.Error(resp, string(message), 500)
}

type PagedResponse struct {
	Response

	OffsetPagination *PaginationOffsetResponse `json:"pagination,omitempty"`
}

type Response struct {
	Result interface{} `json:"result"`
	Err    string      `json:"error,omitempty"`
}

type PaginationOffsetResponse struct {
	TotalRecord uint64 `json:"total_record"`
	TotalPage   uint64 `json:"total_page"`
	CurrentPage uint64 `json:"current_page"`
	Limit       uint64 `json:"limit"`
}

func OptPaginationOffsetResponseFromResult(offsetResult *viewrepo.PaginationOffsetResult) *PaginationOffsetResponse {
	if offsetResult == nil {
		return nil
	}

	return &PaginationOffsetResponse{
		TotalRecord: offsetResult.TotalRecord,
		TotalPage:   offsetResult.TotalPage(),
		CurrentPage: offsetResult.CurrentPage,
		Limit:       offsetResult.Limit,
	}
}

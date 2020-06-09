package adaptertest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
)

type HTTPQueryParams = map[string]string

func NewMockHTTPGetRequest(params HTTPQueryParams) *http.Request {
	urlValues := make(url.Values)
	for key, value := range params {
		urlValues[key] = []string{value}
	}

	return httptest.NewRequest("GET", "/?"+urlValues.Encode(), nil)
}

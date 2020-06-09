package httpapi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chainindex/adapter/httpapi"
	. "github.com/crypto-com/chainindex/adapter/httpapi/test"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
)

var _ = Describe("Pagination", func() {
	Describe("ParsePagination", func() {
		It("should return ErrInvalidPagination when pagination parameter is invalid", func() {
			reqWithInvalidPagination := NewMockHTTPGetRequest(HTTPQueryParams{
				"pagination": "invalid",
				"page":       "1",
				"limit":      "20",
			})

			_, err := httpapi.ParsePagination(reqWithInvalidPagination)
			Expect(err).NotTo(BeNil())
			Expect(err).To(Equal(httpapi.ErrInvalidPagination))
		})

		It("should return ErrInvalidPagination when pagination parameter is unsupported cursor", func() {
			reqWithInvalidPagination := NewMockHTTPGetRequest(HTTPQueryParams{
				"pagination": "cursor",
				"page":       "1",
				"limit":      "20",
			})

			_, err := httpapi.ParsePagination(reqWithInvalidPagination)
			Expect(err).NotTo(BeNil())
			Expect(err).To(Equal(httpapi.ErrInvalidPagination))
		})

		It("should return ErrInvalidPage when page parameter is non-integer", func() {
			reqWithInvalidPagination := NewMockHTTPGetRequest(HTTPQueryParams{
				"pagination": "offset",
				"page":       "invalid",
				"limit":      "20",
			})

			_, err := httpapi.ParsePagination(reqWithInvalidPagination)
			Expect(err).NotTo(BeNil())
			Expect(err).To(Equal(httpapi.ErrInvalidPage))
		})

		It("should return ErrInvalidPage when page parameter is 0", func() {
			reqWithInvalidPagination := NewMockHTTPGetRequest(HTTPQueryParams{
				"pagination": "offset",
				"page":       "0",
				"limit":      "20",
			})

			_, err := httpapi.ParsePagination(reqWithInvalidPagination)
			Expect(err).NotTo(BeNil())
			Expect(err).To(Equal(httpapi.ErrInvalidPage))
		})

		It("should return ErrInvalidPage when page parameter is negative", func() {
			reqWithInvalidPagination := NewMockHTTPGetRequest(HTTPQueryParams{
				"pagination": "offset",
				"page":       "-2",
				"limit":      "20",
			})

			_, err := httpapi.ParsePagination(reqWithInvalidPagination)
			Expect(err).NotTo(BeNil())
			Expect(err).To(Equal(httpapi.ErrInvalidPage))
		})

		It("should return ErrInvalidPage when limit parameter is non-integer", func() {
			reqWithInvalidPagination := NewMockHTTPGetRequest(HTTPQueryParams{
				"pagination": "offset",
				"page":       "1",
				"limit":      "invalid",
			})

			_, err := httpapi.ParsePagination(reqWithInvalidPagination)
			Expect(err).NotTo(BeNil())
			Expect(err).To(Equal(httpapi.ErrInvalidPage))
		})

		It("should return ErrInvalidPage when limit parameter is negative", func() {
			reqWithInvalidPagination := NewMockHTTPGetRequest(HTTPQueryParams{
				"pagination": "offset",
				"page":       "1",
				"limit":      "-20",
			})

			_, err := httpapi.ParsePagination(reqWithInvalidPagination)
			Expect(err).NotTo(BeNil())
			Expect(err).To(Equal(httpapi.ErrInvalidPage))
		})

		It("should return parsed offset pagination when pagination is offset based", func() {
			reqWithInvalidPagination := NewMockHTTPGetRequest(HTTPQueryParams{
				"pagination": "offset",
				"page":       "2",
				"limit":      "15",
			})

			pagination, err := httpapi.ParsePagination(reqWithInvalidPagination)
			Expect(err).To(BeNil())

			Expect(pagination.Type()).To(Equal(viewrepo.PAGINATION_OFFSET))
			offsetParams := pagination.OffsetParams()
			Expect(offsetParams.Page).To(Equal(uint64(2)))
			Expect(offsetParams.Limit).To(Equal(uint64(15)))
		})

		It("should return offset pagination with defaults when none is provided", func() {
			reqWithInvalidPagination := NewMockHTTPGetRequest(HTTPQueryParams{})

			pagination, err := httpapi.ParsePagination(reqWithInvalidPagination)
			Expect(err).To(BeNil())

			Expect(pagination.Type()).To(Equal(viewrepo.PAGINATION_OFFSET))
			offsetParams := pagination.OffsetParams()
			Expect(offsetParams.Page).To(Equal(uint64(1)))
			Expect(offsetParams.Limit).To(Equal(uint64(20)))
		})
	})
})

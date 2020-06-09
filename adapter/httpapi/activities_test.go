package httpapi_test

import (
	"net/http/httptest"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/adapter/httpapi"
	. "github.com/crypto-com/chainindex/adapter/httpapi/test"
	. "github.com/crypto-com/chainindex/adapter/httpapi/test/mock"
	. "github.com/crypto-com/chainindex/usecase/test/fake"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
	. "github.com/crypto-com/chainindex/usecase/viewrepo/test/mock"
)

var _ = Describe("Activities", func() {
	var mockActivityViewRepo *MockActivityViewRepo
	var mockRoutePath *MockRoutePath
	var mockHandler *httpapi.ActivitiesHandler

	BeforeEach(func() {
		fakeLogger := &FakeLogger{}
		mockActivityViewRepo = &MockActivityViewRepo{}
		mockRoutePath = &MockRoutePath{}

		mockHandler = httpapi.NewActivitiesHandler(fakeLogger, mockRoutePath, mockActivityViewRepo)
	})

	Describe("ListTransactions", func() {
		It("should return BadRequest when pagination is missing", func() {
			reqWithInvalidPage := NewMockHTTPGetRequest(HTTPQueryParams{
				"page": "invalid",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListTransactions(respSpy, reqWithInvalidPage)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when type filter has invalid type", func() {
			reqWithInvalidFilter := NewMockHTTPGetRequest(HTTPQueryParams{
				"filter[type]": "invalid",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListTransactions(respSpy, reqWithInvalidFilter)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when type filter has mixed valid and invalid types", func() {
			reqWithInvalidFilter := NewMockHTTPGetRequest(HTTPQueryParams{
				"filter[type]": "transfer,invalid",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListTransactions(respSpy, reqWithInvalidFilter)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when type filter is not using command as separator", func() {
			reqWithInvalidFilter := NewMockHTTPGetRequest(HTTPQueryParams{
				"filter[type]": "transfer;reward",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListTransactions(respSpy, reqWithInvalidFilter)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})
	})

	Describe("FindTransactionByTxId", func() {
		It("should return BadRequest when txid is missing", func() {
			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{})

			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockHandler.FindTransactionByTxId(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when txid is missing", func() {
			anyTxID := "any-txid"
			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"txid": anyTxID,
			})
			mockActivityViewRepo.On(
				"FindTransactionByTxId", anyTxID,
			).Return((*viewrepo.Transaction)(nil), adapter.ErrNotFound)

			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockHandler.FindTransactionByTxId(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(404))
		})
	})

	Describe("ListEvents", func() {
		It("should return BadRequest when pagination is missing", func() {
			reqWithInvalidPage := NewMockHTTPGetRequest(HTTPQueryParams{
				"page": "invalid",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListEvents(respSpy, reqWithInvalidPage)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when type filter has invalid type", func() {
			reqWithInvalidFilter := NewMockHTTPGetRequest(HTTPQueryParams{
				"filter[type]": "invalid",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListEvents(respSpy, reqWithInvalidFilter)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when type filter has mixed valid and invalid types", func() {
			reqWithInvalidFilter := NewMockHTTPGetRequest(HTTPQueryParams{
				"filter[type]": "transfer,invalid",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListEvents(respSpy, reqWithInvalidFilter)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when type filter is not using command as separator", func() {
			reqWithInvalidFilter := NewMockHTTPGetRequest(HTTPQueryParams{
				"filter[type]": "transfer;reward",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListEvents(respSpy, reqWithInvalidFilter)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})
	})

	Describe("FindEventByBlockHeightEventPosition", func() {
		It("should return BadRequest when height is missing", func() {
			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"position": "0",
			})

			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockHandler.FindEventByBlockHeightEventPosition(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when height is invalid", func() {
			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"height":   "invalid",
				"position": "0",
			})

			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockHandler.FindEventByBlockHeightEventPosition(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when position is missing", func() {
			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"height": "1000",
			})

			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockHandler.FindEventByBlockHeightEventPosition(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when position is invalid", func() {
			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"height":   "1000",
				"position": "invalid",
			})

			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockHandler.FindEventByBlockHeightEventPosition(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when event does not exist", func() {
			anyHeight := uint64(1000)
			anyPosition := uint64(1)
			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"height":   strconv.FormatUint(anyHeight, 10),
				"position": strconv.FormatUint(anyPosition, 10),
			})
			mockActivityViewRepo.On(
				"FindEventByBlockHeightEventPosition", anyHeight, anyPosition,
			).Return((*viewrepo.Event)(nil), adapter.ErrNotFound)

			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockHandler.FindEventByBlockHeightEventPosition(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(404))
		})
	})
})

package httpapi_test

import (
	"net/http/httptest"

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

var _ = Describe("Blocks", func() {
	var mockBlockViewRepo *MockBlockViewRepo
	var mockRoutePath *MockRoutePath
	var mockHandler *httpapi.BlocksHandler

	BeforeEach(func() {
		fakeLogger := &FakeLogger{}
		mockBlockViewRepo = &MockBlockViewRepo{}
		mockRoutePath = &MockRoutePath{}

		mockHandler = httpapi.NewBlocksHandler(fakeLogger, mockRoutePath, mockBlockViewRepo)
	})

	Describe("ListBlocks", func() {
		It("should return BadRequest when pagination is missing", func() {
			reqWithInvalidPage := NewMockHTTPGetRequest(HTTPQueryParams{
				"page": "invalid",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListBlocks(respSpy, reqWithInvalidPage)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when proposer id filter has invalid type", func() {
			reqWithInvalidFilter := NewMockHTTPGetRequest(HTTPQueryParams{
				"filter[proposer_id]": "invalid",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListBlocks(respSpy, reqWithInvalidFilter)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when proposer id filter is non-positive-integer", func() {
			reqWithInvalidFilter := NewMockHTTPGetRequest(HTTPQueryParams{
				"filter[proposer_id]": "-1",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListBlocks(respSpy, reqWithInvalidFilter)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when proposer id filter is non-integer", func() {
			reqWithInvalidFilter := NewMockHTTPGetRequest(HTTPQueryParams{
				"filter[proposer_id]": "1.1",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListBlocks(respSpy, reqWithInvalidFilter)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when proposer id filter is not separated by comma", func() {
			reqWithInvalidFilter := NewMockHTTPGetRequest(HTTPQueryParams{
				"filter[proposer_id]": "1;2",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.ListBlocks(respSpy, reqWithInvalidFilter)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})
	})

	Describe("FindBlock", func() {
		It("should return BadRequest when block identity is missing", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{})

			mockHandler.FindBlock(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when block identity is not block hash nor height number", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "invalid",
			})

			mockHandler.FindBlock(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when block identity is block hash of invalid length", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "blockHashOfInvalidLength",
			})

			mockHandler.FindBlock(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when block identity is block height of negative number", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "-10",
			})

			mockHandler.FindBlock(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when block identity is block height of 0", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "0",
			})

			mockHandler.FindBlock(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return NotFound when block does not exist", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "10",
			})

			mockBlockViewRepo.On("FindBlock", mock.Anything).Return((*viewrepo.Block)(nil), adapter.ErrNotFound)

			mockHandler.FindBlock(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(404))
		})
	})

	Describe("ListBlockTransactions", func() {
		It("should return BadRequest when block identity is missing", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{})

			mockHandler.ListBlockTransactions(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when block identity is not block hash nor height number", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "invalid",
			})

			mockHandler.ListBlockTransactions(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when block identity is block hash of invalid length", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "blockHashOfInvalidLength",
			})

			mockHandler.ListBlockTransactions(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when block identity is block height of negative number", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "-10",
			})

			mockHandler.ListBlockTransactions(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when block identity is block height of 0", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "0",
			})

			mockHandler.ListBlockTransactions(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return NotFound when block does not exist", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "10",
			})

			mockBlockViewRepo.On(
				"ListBlockTransactions",
				mock.Anything,
				mock.Anything,
			).Return(([]viewrepo.Transaction)(nil), (*viewrepo.PaginationResult)(nil), adapter.ErrNotFound)

			mockHandler.ListBlockTransactions(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(404))
		})

		It("should return BadRequest when pagination is missing", func() {
			reqWithInvalidPage := NewMockHTTPGetRequest(HTTPQueryParams{
				"page": "invalid",
			})
			respSpy := httptest.NewRecorder()
			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "10",
			})

			mockHandler.ListBlockTransactions(respSpy, reqWithInvalidPage)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})
	})

	Describe("ListBlockEvents", func() {
		It("should return BadRequest when block identity is missing", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{})

			mockHandler.ListBlockEvents(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when block identity is not block hash nor height number", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "invalid",
			})

			mockHandler.ListBlockEvents(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when block identity is block hash of invalid length", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "blockHashOfInvalidLength",
			})

			mockHandler.ListBlockEvents(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when block identity is block height of negative number", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "-10",
			})

			mockHandler.ListBlockEvents(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return BadRequest when block identity is block height of 0", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "0",
			})

			mockHandler.ListBlockEvents(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

		It("should return NotFound when block does not exist", func() {
			anyReq := NewMockHTTPGetRequest(HTTPQueryParams{})
			respSpy := httptest.NewRecorder()

			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "10",
			})

			mockBlockViewRepo.On(
				"ListBlockEvents",
				mock.Anything,
				mock.Anything,
			).Return(([]viewrepo.BlockEvent)(nil), (*viewrepo.PaginationResult)(nil), adapter.ErrNotFound)

			mockHandler.ListBlockEvents(respSpy, anyReq)

			Expect(respSpy.Result().StatusCode).To(Equal(404))
		})

		It("should return BadRequest when pagination is missing", func() {
			reqWithInvalidPage := NewMockHTTPGetRequest(HTTPQueryParams{
				"page": "invalid",
			})
			respSpy := httptest.NewRecorder()
			mockRoutePath.On("Vars", mock.Anything).Return(map[string]string{
				"hash_or_height": "10",
			})

			mockHandler.ListBlockEvents(respSpy, reqWithInvalidPage)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})

	})
})

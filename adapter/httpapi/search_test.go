package httpapi_test

import (
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chainindex/adapter/httpapi"
	. "github.com/crypto-com/chainindex/adapter/httpapi/test"
	. "github.com/crypto-com/chainindex/usecase/test/fake"
	. "github.com/crypto-com/chainindex/usecase/viewrepo/test/mock"
)

var _ = Describe("Search", func() {
	var mockActivityViewRepo *MockActivityViewRepo
	var mockBlockViewRepo *MockBlockViewRepo
	var mockStakingAccountViewRepo *MockStakingAccountViewRepo
	var mockCouncilNodeViewRepo *MockCouncilNodeViewRepo
	var mockHandler *httpapi.SearchHandler

	BeforeEach(func() {
		fakeLogger := &FakeLogger{}
		mockActivityViewRepo = &MockActivityViewRepo{}
		mockBlockViewRepo = &MockBlockViewRepo{}
		mockStakingAccountViewRepo = &MockStakingAccountViewRepo{}
		mockCouncilNodeViewRepo = &MockCouncilNodeViewRepo{}

		mockHandler = httpapi.NewSearchHandler(
			fakeLogger,
			mockActivityViewRepo,
			mockBlockViewRepo,
			mockStakingAccountViewRepo,
			mockCouncilNodeViewRepo,
		)
	})

	Describe("All", func() {
		It("should return BadRequest when pagination is missing", func() {
			reqWithInvalidPage := NewMockHTTPGetRequest(HTTPQueryParams{
				"page": "invalid",
			})
			respSpy := httptest.NewRecorder()

			mockHandler.All(respSpy, reqWithInvalidPage)

			Expect(respSpy.Result().StatusCode).To(Equal(400))
		})
	})
})

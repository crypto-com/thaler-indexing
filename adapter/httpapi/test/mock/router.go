package httpapimock

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockRoutePath struct {
	mock.Mock
}

func (path *MockRoutePath) Vars(req *http.Request) map[string]string {
	args := path.Called(req)

	return args.Get(0).(map[string]string)
}

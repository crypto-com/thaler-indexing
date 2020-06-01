package mock

import testifymock "github.com/stretchr/testify/mock"

// Genrate a slice of testify.mock.Anything.
// mockObj.On("Method", MockAnythingOfTimes(3)...).Return(true)
func MockAnythingOfTimes(time int) []interface{} {
	args := make([]interface{}, time)
	for i := 0; i < time; i += 1 {
		args[i] = testifymock.Anything
	}

	return args
}

package httpapi_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttpapi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Httpapi Suite")
}

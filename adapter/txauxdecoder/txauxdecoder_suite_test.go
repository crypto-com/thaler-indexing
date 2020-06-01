package txauxdecoder_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTxauxdecoder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Txauxdecoder Suite")
}

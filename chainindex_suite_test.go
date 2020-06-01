package chainindex_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestChainIndex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ChainIndex Suite")
}

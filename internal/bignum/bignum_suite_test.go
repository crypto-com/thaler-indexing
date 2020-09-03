package bignum_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBignum(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bignum Suite")
}

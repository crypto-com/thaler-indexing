package chainindex_test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chainindex"
)

var _ = Describe("Coin", func() {
	Describe("CROToCoin", func() {
		It("should return error when the CRO has more than 8 decimal places", func() {
			invalidCRO, _ := new(big.Float).SetString("0.000000000001")

			_, err := chainindex.CROToCoin(invalidCRO)
			Expect(err).NotTo(BeNil())
		})

		It("should return coin representation of the CRO amount", func() {
			oneCRO, _ := new(big.Float).SetString("0.00000001")

			unit, err := chainindex.CROToCoin(oneCRO)
			Expect(err).To(BeNil())
			Expect(unit).To(Equal(big.NewInt(1)))
		})
	})
})

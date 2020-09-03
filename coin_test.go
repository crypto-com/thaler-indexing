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
			Expect(unit.Cmp(big.NewInt(1))).To(Equal(0))
		})

		// FIXME: big.Float has precision problem
		// It("should return coin representation of a CRO amount", func() {
		// 	testVal, _ := new(big.Float).SetString("0.00000884")

		// 	unit, err := chainindex.CROToCoin(testVal)
		// 	Expect(err).To(BeNil())
		// 	Expect(unit.Cmp(big.NewInt(884))).To(Equal(0))
		// })

		It("should return coin representation of a big CRO amount", func() {
			testVal, _ := new(big.Float).SetString("99999999999.00000884")

			unit, err := chainindex.CROToCoin(testVal)
			Expect(err).To(BeNil())
			expectedUnit, _ := new(big.Int).SetString("9999999999900000884", 10)
			Expect(unit.Cmp(expectedUnit)).To(Equal(0))
		})
	})

	Describe("CROStrToCoin", func() {
		It("should return coin representation of the CRO amount", func() {
			oneCRO := "0.00000001"

			unit, err := chainindex.CROStrToCoin(oneCRO)
			Expect(err).To(BeNil())
			Expect(unit.Cmp(big.NewInt(1))).To(Equal(0))
		})

		It("should return coin representation of a CRO amount", func() {
			unit, err := chainindex.CROStrToCoin("0.00000884")
			Expect(err).To(BeNil())
			Expect(unit.Cmp(big.NewInt(884))).To(Equal(0))
		})
	})
})

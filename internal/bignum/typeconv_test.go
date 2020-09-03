package bignum_test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chainindex/internal/bignum"
)

var _ = Describe("bignum", func() {
	Describe("Atof", func() {
		It("should return error when the float number will result in lost in precision", func() {
			result, err := bignum.Atof("2.797693134862315708145274237317043567981e+308")

			Expect(err).To(BeNil())
			expected, _ := new(big.Float).SetString("2.797693134862315708145274237317043567981e+308")
			Expect(result.Cmp(expected)).To(Equal(0))
		})

		It("should return the big.Float of the float amount", func() {
			bigFloat, err := bignum.Atof("0.00000001")

			Expect(err).To(BeNil())
			floatVal, _ := bigFloat.Float64()
			Expect(floatVal).To(Equal(0.00000001))
		})

		It("should return the big.Float of the float amount", func() {
			bigFloat, err := bignum.Atof("0.00000884")

			Expect(err).To(BeNil())
			floatVal, _ := bigFloat.Float64()
			Expect(floatVal).To(Equal(0.00000884))
		})

		It("should return the big.Float of the int amount with decimal", func() {
			bigFloat, err := bignum.Atof("123456.00000884")

			Expect(err).To(BeNil())
			floatVal, _ := bigFloat.Float64()
			Expect(floatVal).To(Equal(123456.00000884))
		})
	})
})

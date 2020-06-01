package infrastructure_test

import (
	"math/big"

	"github.com/jackc/pgtype"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chainindex/infrastructure"
)

var _ = Describe("PgxRDbTypeConv", func() {
	var pgxRdbTypeConv infrastructure.PgxRDbTypeConv

	BeforeEach(func() {
		//nolint:staticcheck
		pgxRdbTypeConv = infrastructure.PgxRDbTypeConv{}
	})

	Describe("Bton", func() {
		It("should return numeric null when bigInt is nil", func() {
			var expected pgtype.Numeric
			_ = expected.Set(nil)

			actual := pgxRdbTypeConv.Bton(nil)
			Expect(actual).To(Equal(expected))
		})

		It("should return numeric of the bigInt value", func() {
			var expected pgtype.Numeric
			_ = expected.Scan("10")

			actual := pgxRdbTypeConv.Bton(big.NewInt(10))
			Expect(actual).To(Equal(expected))
		})

		It("should work for negative number", func() {
			var expected pgtype.Numeric
			_ = expected.Scan("-10")

			actual := pgxRdbTypeConv.Bton(big.NewInt(-10))
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("Iton", func() {
		It("should return numeric representation of the int", func() {
			var expected pgtype.Numeric
			_ = expected.Set(10)

			Expect(pgxRdbTypeConv.Iton(10)).To(Equal(expected))
		})

		It("should work for negative number", func() {
			var expected pgtype.Numeric
			_ = expected.Set(-10)

			Expect(pgxRdbTypeConv.Iton(-10)).To(Equal(expected))
		})
	})

	Describe("NtobReader", func() {
		It("should return Error for decimal number", func() {
			var n pgtype.Numeric
			_ = n.Set(1.23456)

			ntobReader := pgxRdbTypeConv.NtobReader()
			arg, _ := ntobReader.ScannableArg().(*pgtype.Numeric)
			*arg = n

			_, err := ntobReader.Parse()
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("cannot convert 123456e-5 to bigInt"))
		})

		It("should return nil when numeric is null", func() {
			var n pgtype.Numeric
			_ = n.Set(nil)

			ntobReader := pgxRdbTypeConv.NtobReader()
			arg, _ := ntobReader.ScannableArg().(*pgtype.Numeric)
			*arg = n

			actual, err := ntobReader.Parse()
			Expect(err).To(BeNil())
			Expect(actual).To(BeNil())
		})

		It("should return big.Int of the numeric value", func() {
			var n pgtype.Numeric
			_ = n.Set(10)

			expected := big.NewInt(10)

			ntobReader := pgxRdbTypeConv.NtobReader()
			arg, _ := ntobReader.ScannableArg().(*pgtype.Numeric)
			*arg = n

			actual, err := ntobReader.Parse()
			Expect(err).To(BeNil())
			Expect(*actual).To(Equal(*expected))
		})

		It("should work for negative number", func() {
			var n pgtype.Numeric
			_ = n.Set(-10)

			expected := big.NewInt(-10)

			ntobReader := pgxRdbTypeConv.NtobReader()
			arg, _ := ntobReader.ScannableArg().(*pgtype.Numeric)
			*arg = n

			actual, err := ntobReader.Parse()
			Expect(err).To(BeNil())
			Expect(*actual).To(Equal(*expected))
		})
	})
})

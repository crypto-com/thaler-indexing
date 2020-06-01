package adapter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	random "github.com/brianvoe/gofakeit/v5"

	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/internal/bignum"
	"github.com/crypto-com/chainindex/internal/primptr"
	. "github.com/crypto-com/chainindex/test/factory"
)

var _ = Describe("StakingAccount", func() {
	Describe("IncrementNonce", func() {
		It("should increment nonce of the record", func() {
			accountRow := RandomRDbStakingAccountRow()
			accountRow.Nonce = uint64(1)

			accountRow.IncrementNonce()

			Expect(accountRow.Nonce).To(Equal(uint64(2)))
		})
	})

	Describe("AddBonded", func() {
		It("should add bonded amount to the record", func() {
			accountRow := RandomRDbStakingAccountRow()
			accountRow.Bonded = bignum.Int0()

			accountRow.AddBonded(bignum.Int1())

			Expect(accountRow.Bonded).To(Equal(bignum.Int1()))
		})
	})

	Describe("AddUnbonded", func() {
		It("should add unbonded amount to the record", func() {
			accountRow := RandomRDbStakingAccountRow()
			accountRow.Unbonded = bignum.Int0()

			accountRow.AddUnbonded(bignum.Int1())

			Expect(accountRow.Unbonded).To(Equal(bignum.Int1()))
		})
	})
})

func RandomRDbStakingAccountRow() *adapter.RDbStakingAccountRow {
	return &adapter.RDbStakingAccountRow{
		Address:              RandomTendermintAddress(),
		Nonce:                random.Uint64(),
		Bonded:               RandomCoin(),
		Unbonded:             RandomCoin(),
		UnbondedFrom:         primptr.Time(RandomUTCTime()),
		PunishmentKind:       primptr.String(adapter.PunishmentKindToString(RandomPunishmentKind())),
		JailedUntil:          primptr.Time(RandomUTCTime()),
		CurrentCouncilNodeId: primptr.Uint64(random.Uint64()),
	}
}

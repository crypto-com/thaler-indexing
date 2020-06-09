package usecasevewrepomock

import (
	"github.com/crypto-com/chainindex/internal/bignum"
	"github.com/stretchr/testify/mock"
)

type MockRewardViewRepo struct {
	mock.Mock
}

func (repo *MockRewardViewRepo) TotalMinted() (*bignum.WBigInt, error) {
	args := repo.Called()

	return args.Get(0).(*bignum.WBigInt), args.Error(1)
}

func (repo *MockRewardViewRepo) Total() (*bignum.WBigInt, error) {
	args := repo.Called()

	return args.Get(0).(*bignum.WBigInt), args.Error(1)
}

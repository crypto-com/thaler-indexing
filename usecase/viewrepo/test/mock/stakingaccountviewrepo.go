package usecasevewrepomock

import (
	"github.com/crypto-com/chainindex/usecase/viewrepo"
	"github.com/stretchr/testify/mock"
)

type MockStakingAccountViewRepo struct {
	mock.Mock
}

func (repo *MockStakingAccountViewRepo) Search(
	keyword string,
	pagination *viewrepo.Pagination,
) ([]viewrepo.StakingAccount, *viewrepo.PaginationResult, error) {
	args := repo.Called(keyword, pagination)

	return args.Get(0).([]viewrepo.StakingAccount), args.Get(1).(*viewrepo.PaginationResult), args.Error(2)
}

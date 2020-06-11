package usecasevewrepomock

import (
	"github.com/crypto-com/chainindex/usecase/viewrepo"
	"github.com/stretchr/testify/mock"
)

type MockCouncilNodeViewRepo struct {
	mock.Mock
}

func (repo *MockCouncilNodeViewRepo) ListActivities(
	pagination *viewrepo.Pagination,
) ([]viewrepo.CouncilNodeListItem, *viewrepo.PaginationResult, error) {
	args := repo.Called(pagination)

	return args.Get(0).([]viewrepo.CouncilNodeListItem), args.Get(1).(*viewrepo.PaginationResult), args.Error(2)
}

func (repo *MockCouncilNodeViewRepo) FindById(id uint64) (*viewrepo.CouncilNode, error) {
	args := repo.Called(id)

	return args.Get(0).(*viewrepo.CouncilNode), args.Error(1)
}

func (repo *MockCouncilNodeViewRepo) ListActivitiesById(
	id uint64,
	filter viewrepo.ActivityFilter,
	pagination *viewrepo.Pagination,
) ([]viewrepo.StakingAccountActivity, *viewrepo.PaginationResult, error) {
	args := repo.Called(id, filter, pagination)

	return args.Get(0).([]viewrepo.StakingAccountActivity), args.Get(1).(*viewrepo.PaginationResult), args.Error(2)
}

func (repo *MockCouncilNodeViewRepo) Stats() (*viewrepo.CouncilNodeStats, error) {
	args := repo.Called()

	return args.Get(0).(*viewrepo.CouncilNodeStats), args.Error(1)
}

func (repo *MockCouncilNodeViewRepo) Search(
	keyword string,
	pagination *viewrepo.Pagination,
) ([]viewrepo.CouncilNode, *viewrepo.PaginationResult, error) {
	args := repo.Called(keyword, pagination)

	return args.Get(0).([]viewrepo.CouncilNode), args.Get(1).(*viewrepo.PaginationResult), args.Error(2)
}

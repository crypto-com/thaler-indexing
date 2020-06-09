package usecasevewrepomock

import (
	"github.com/crypto-com/chainindex/usecase/viewrepo"
	"github.com/stretchr/testify/mock"
)

type MockBlockViewRepo struct {
	mock.Mock
}

func (repo *MockBlockViewRepo) LatestBlockHeight() (uint64, error) {
	args := repo.Called()

	return args.Get(0).(uint64), args.Error(1)
}

func (repo *MockBlockViewRepo) ListBlocks(
	filter viewrepo.BlockFilter,
	pagination *viewrepo.Pagination,
) ([]viewrepo.Block, *viewrepo.PaginationResult, error) {
	args := repo.Called(filter, pagination)

	return args.Get(0).([]viewrepo.Block), args.Get(1).(*viewrepo.PaginationResult), args.Error(2)
}

func (repo *MockBlockViewRepo) FindBlock(identity viewrepo.BlockIdentity) (*viewrepo.Block, error) {
	args := repo.Called(identity)

	return args.Get(0).(*viewrepo.Block), args.Error(1)
}

func (repo *MockBlockViewRepo) ListBlockTransactions(
	identity viewrepo.BlockIdentity,
	pagination *viewrepo.Pagination,
) ([]viewrepo.Transaction, *viewrepo.PaginationResult, error) {
	args := repo.Called(identity, pagination)

	return args.Get(0).([]viewrepo.Transaction), args.Get(1).(*viewrepo.PaginationResult), args.Error(2)
}

func (repo *MockBlockViewRepo) ListBlockEvents(
	identity viewrepo.BlockIdentity,
	pagination *viewrepo.Pagination,
) ([]viewrepo.BlockEvent, *viewrepo.PaginationResult, error) {
	args := repo.Called(identity, pagination)

	return args.Get(0).([]viewrepo.BlockEvent), args.Get(1).(*viewrepo.PaginationResult), args.Error(2)
}

func (repo *MockBlockViewRepo) Search(
	keyword string,
	pagination *viewrepo.Pagination,
) ([]viewrepo.Block, *viewrepo.PaginationResult, error) {
	args := repo.Called(keyword, pagination)

	return args.Get(0).([]viewrepo.Block), args.Get(1).(*viewrepo.PaginationResult), args.Error(2)
}

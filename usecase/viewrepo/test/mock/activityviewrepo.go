package usecasevewrepomock

import (
	"github.com/stretchr/testify/mock"

	"github.com/crypto-com/chainindex/usecase/viewrepo"
)

type MockActivityViewRepo struct {
	mock.Mock
}

func (repo *MockActivityViewRepo) ListTransactions(
	filter viewrepo.TransactionFilter,
	pagination *viewrepo.Pagination,
) ([]viewrepo.Transaction, *viewrepo.PaginationResult, error) {
	args := repo.Called(filter, pagination)

	return args.Get(0).([]viewrepo.Transaction), args.Get(1).(*viewrepo.PaginationResult), args.Error(2)
}

func (repo *MockActivityViewRepo) FindTransactionByTxId(txid string) (*viewrepo.Transaction, error) {
	args := repo.Called(txid)

	return args.Get(0).(*viewrepo.Transaction), args.Error(1)
}

func (repo *MockActivityViewRepo) ListEvents(
	filter viewrepo.EventFilter,
	pagination *viewrepo.Pagination,
) ([]viewrepo.Event, *viewrepo.PaginationResult, error) {
	args := repo.Called(filter, pagination)

	return args.Get(0).([]viewrepo.Event), args.Get(1).(*viewrepo.PaginationResult), args.Error(2)
}

func (repo *MockActivityViewRepo) FindEventByBlockHeightEventPosition(
	blockHeight uint64,
	eventPosition uint64,
) (*viewrepo.Event, error) {
	args := repo.Called(blockHeight, eventPosition)

	return args.Get(0).(*viewrepo.Event), args.Error(1)
}

func (repo *MockActivityViewRepo) TransactionsCount() (uint64, error) {
	args := repo.Called()

	return args.Get(0).(uint64), args.Error(1)
}

func (repo *MockActivityViewRepo) SearchTransactions(
	keyword string,
	pagination *viewrepo.Pagination,
) ([]viewrepo.Transaction, *viewrepo.PaginationResult, error) {
	args := repo.Called(keyword, pagination)

	return args.Get(0).([]viewrepo.Transaction), args.Get(1).(*viewrepo.PaginationResult), args.Error(2)
}

func (repo *MockActivityViewRepo) SearchEvents(
	keyword string,
	pagination *viewrepo.Pagination,
) ([]viewrepo.Event, *viewrepo.PaginationResult, error) {
	args := repo.Called(keyword, pagination)

	return args.Get(0).([]viewrepo.Event), args.Get(1).(*viewrepo.PaginationResult), args.Error(2)
}

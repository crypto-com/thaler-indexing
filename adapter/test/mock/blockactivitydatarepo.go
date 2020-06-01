package adaptermock

import (
	"github.com/stretchr/testify/mock"

	"github.com/crypto-com/chainindex"
	"github.com/crypto-com/chainindex/adapter"
)

type MockRDbBlockActivityDataRepo struct {
	mock.Mock
}

func (repo *MockRDbBlockActivityDataRepo) InsertGenesisActivity(tx adapter.RDbTx, activity *chainindex.Activity) error {
	args := repo.Called(tx, activity)
	return args.Error(0)
}
func (repo *MockRDbBlockActivityDataRepo) InsertTransferTransaction(tx adapter.RDbTx, activity *chainindex.Activity) error {
	args := repo.Called(tx, activity)
	return args.Error(0)
}
func (repo *MockRDbBlockActivityDataRepo) InsertDepositTransaction(tx adapter.RDbTx, activity *chainindex.Activity) error {
	args := repo.Called(tx, activity)
	return args.Error(0)
}
func (repo *MockRDbBlockActivityDataRepo) InsertUnbondTransaction(tx adapter.RDbTx, activity *chainindex.Activity) error {
	args := repo.Called(tx, activity)
	return args.Error(0)
}
func (repo *MockRDbBlockActivityDataRepo) InsertWithdrawTransaction(tx adapter.RDbTx, activity *chainindex.Activity) error {
	args := repo.Called(tx, activity)
	return args.Error(0)
}
func (repo *MockRDbBlockActivityDataRepo) InsertNodeJoinTransaction(tx adapter.RDbTx, activity *chainindex.Activity) error {
	args := repo.Called(tx, activity)
	return args.Error(0)
}
func (repo *MockRDbBlockActivityDataRepo) InsertUnjailTransaction(tx adapter.RDbTx, activity *chainindex.Activity) error {
	args := repo.Called(tx, activity)
	return args.Error(0)
}
func (repo *MockRDbBlockActivityDataRepo) InsertRewardEvent(tx adapter.RDbTx, activity *chainindex.Activity) error {
	args := repo.Called(tx, activity)
	return args.Error(0)
}
func (repo *MockRDbBlockActivityDataRepo) InsertSlashEvent(tx adapter.RDbTx, activity *chainindex.Activity) error {
	args := repo.Called(tx, activity)
	return args.Error(0)
}
func (repo *MockRDbBlockActivityDataRepo) InsertJailEvent(tx adapter.RDbTx, activity *chainindex.Activity) error {
	args := repo.Called(tx, activity)
	return args.Error(0)
}

package rdbviewrepo

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/internal/bignum"
)

type RDbRewardViewRepo struct {
	conn adapter.RDbConn

	stmtBuilder sq.StatementBuilderType
	typeConv    adapter.RDbTypeConv
}

func NewRDbRewardViewRepo(
	conn adapter.RDbConn,
	stmtBuilder sq.StatementBuilderType,
	typeConv adapter.RDbTypeConv,
) *RDbRewardViewRepo {
	return &RDbRewardViewRepo{
		conn,

		stmtBuilder,
		typeConv,
	}
}

func (repo *RDbRewardViewRepo) TotalMinted() (*bignum.WBigInt, error) {
	var err error

	sql, sqlArgs, err := repo.stmtBuilder.Select(
		"SUM(minted)",
	).From(
		"block_rewards",
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building reward select SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	totalRewardMintedReader := repo.typeConv.NtobReader()
	err = repo.conn.QueryRow(sql, sqlArgs...).Scan(totalRewardMintedReader.ScannableArg())
	if err != nil {
		return nil, fmt.Errorf("error scanning reward row: %v, %w", err, adapter.ErrRepoQuery)
	}

	totalRewardMinted, err := totalRewardMintedReader.ParseW()
	if err != nil {
		return nil, fmt.Errorf("error parsing total reward minted: %v: %w", err, adapter.ErrRepoQuery)
	}

	return totalRewardMinted, nil
}

func (repo *RDbRewardViewRepo) Total() (*bignum.WBigInt, error) {
	var err error

	sql, sqlArgs, err := repo.stmtBuilder.Select(
		"SUM(bonded)",
	).From(
		"activities",
	).Where(
		"type = 'reward'",
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building reward activity select SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	totalRewardReader := repo.typeConv.NtobReader()
	err = repo.conn.QueryRow(sql, sqlArgs...).Scan(totalRewardReader.ScannableArg())
	if err != nil {
		return nil, fmt.Errorf("error scanning reward activity row: %v, %w", err, adapter.ErrRepoQuery)
	}

	totalReward, err := totalRewardReader.ParseW()
	if err != nil {
		return nil, fmt.Errorf("error parsing total reward: %v: %w", err, adapter.ErrRepoQuery)
	}

	return totalReward, nil
}

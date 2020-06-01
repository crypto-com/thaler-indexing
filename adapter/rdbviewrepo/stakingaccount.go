package rdbviewrepo

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
)

type RDbStkaingAccountViewRepo struct {
	conn adapter.RDbConn

	stmtBuilder sq.StatementBuilderType
	typeConv    adapter.RDbTypeConv
}

func NewRDbStkaingAccountViewRepo(
	conn adapter.RDbConn,
	stmtBuilder sq.StatementBuilderType,
	typeConv adapter.RDbTypeConv,
) *RDbStkaingAccountViewRepo {
	return &RDbStkaingAccountViewRepo{
		conn,

		stmtBuilder,
		typeConv,
	}
}

func (repo *RDbStkaingAccountViewRepo) Search(keyword string, pagination *viewrepo.Pagination) ([]viewrepo.StakingAccount, *viewrepo.PaginationResult, error) {
	var err error

	likeKeyword := fmt.Sprint("%", keyword, "%")
	rDbPagination := adapter.NewRDbPaginationBuilder(
		pagination,
		repo.conn,
	).BuildStmt(repo.stmtBuilder.Select(
		"sa.address AS staking_account_address",
		"sa.nonce",
		"sa.bonded",
		"sa.unbonded",
		"sa.unbonded_from",
		"sa.punishment_kind",
		"sa.jailed_until",
		"cn.id",
		"cn.name",
		"cn.security_contact",
		"cn.pubkey_type",
		"cn.pubkey",
		"cn.address",
		"cn.created_at_block_height",
		"cn.last_left_at_block_height",
	).From(
		"staking_accounts sa",
	).LeftJoin(
		"council_nodes cn ON sa.current_council_node_id = cn.id",
	).OrderBy(
		"sa.address",
	).Where(
		"sa.address = ? OR cn.name LIKE ? OR cn.security_contact LIKE ? OR cn.pubkey = ? OR cn.address = ?",
		keyword, likeKeyword, likeKeyword, keyword, keyword,
	))

	sql, sqlArgs, err := rDbPagination.ToStmtBuilder().ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building staking account search SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	rowsResult, err := repo.conn.Query(sql, sqlArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing staking account search SQL: %v: %w", err, adapter.ErrRepoQuery)
	}

	stakingAccounts := make([]viewrepo.StakingAccount, 0)
	for rowsResult.Next() {
		var stakingAccount viewrepo.StakingAccount
		var councilNode viewrepo.StakingAccountCouncilNode

		bondedReader := repo.typeConv.NtobReader()
		unbondedReader := repo.typeConv.NtobReader()
		unbondedFromReader := repo.typeConv.NtotReader()
		jailedUntilReader := repo.typeConv.NtotReader()
		if err = rowsResult.Scan(
			&stakingAccount.Address,
			&stakingAccount.Nonce,
			bondedReader.ScannableArg(),
			unbondedReader.ScannableArg(),
			unbondedFromReader.ScannableArg(),
			&stakingAccount.MaybePunishmentKind,
			jailedUntilReader.ScannableArg(),
			&councilNode.MaybeId,
			&councilNode.MaybeName,
			&councilNode.MaybeSecurityContact,
			&councilNode.MaybePubKeyType,
			&councilNode.MaybePubKey,
			&councilNode.MaybeAddress,
			&councilNode.MaybeCreatedAtBlockHeight,
			&councilNode.MaybeLastLeftAtBlockHeight,
		); err != nil {
			return nil, nil, fmt.Errorf("error scanning council node row: %v: %w", err, adapter.ErrRepoQuery)
		}

		if stakingAccount.Bonded, err = bondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing bonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if stakingAccount.MaybeUnbonded, err = unbondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if stakingAccount.MaybeUnbondedFrom, err = unbondedFromReader.Parse(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded from: %v: %w", err, adapter.ErrRepoQuery)
		}
		if stakingAccount.MaybeJailedUntil, err = jailedUntilReader.Parse(); err != nil {
			return nil, nil, fmt.Errorf("error parsing jailed until: %v: %w", err, adapter.ErrRepoQuery)
		}

		if councilNode.MaybeId != nil {
			stakingAccount.MaybeCurrentCouncilNode = &councilNode
		}

		stakingAccounts = append(stakingAccounts, stakingAccount)
	}

	paginationResult, err := rDbPagination.Result()
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing pagination result: %v", err)
	}

	return stakingAccounts, paginationResult, nil
}

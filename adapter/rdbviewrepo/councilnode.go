package rdbviewrepo

import (
	"fmt"
	"math/big"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
	jsoniter "github.com/json-iterator/go"
)

type RDbCouncilNodeViewRepo struct {
	conn adapter.RDbConn

	stmtBuilder sq.StatementBuilderType
	typeConv    adapter.RDbTypeConv
}

func NewRDbCouncilNodeViewRepo(
	conn adapter.RDbConn,
	stmtBuilder sq.StatementBuilderType,
	typeConv adapter.RDbTypeConv,
) *RDbCouncilNodeViewRepo {
	return &RDbCouncilNodeViewRepo{
		conn,

		stmtBuilder,
		typeConv,
	}
}

func (repo *RDbCouncilNodeViewRepo) ListActivities(pagination *viewrepo.Pagination) ([]viewrepo.CouncilNodeListItem, *viewrepo.PaginationResult, error) {
	var err error

	baseStmtBuilder := repo.stmtBuilder.Select().From(
		"council_nodes c",
	).LeftJoin(
		"staking_accounts sa ON c.id = sa.current_council_node_id",
	).Where(
		"c.last_left_at_block_height IS NULL",
	).OrderBy(
		"sa.bonded DESC",
	)

	totalCumulativeSql, totalCumulativeSqlArgs, err := repo.stmtBuilder.Select(
		"SUM(sa.bonded)",
	).From(
		"council_nodes c",
	).LeftJoin(
		"staking_accounts sa ON c.id = sa.current_council_node_id",
	).Where(
		"c.last_left_at_block_height IS NULL",
	).ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building council nodes total cumulative bonded SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}
	totalCumulativeReader := repo.typeConv.NtobReader()
	if err = repo.conn.QueryRow(
		totalCumulativeSql, totalCumulativeSqlArgs...,
	).Scan(totalCumulativeReader.ScannableArg()); err != nil {
		return nil, nil, fmt.Errorf("error executing council nodes total cumulative bonded SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}
	totalCumulative, err := totalCumulativeReader.Parse()
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing council nodes total cumulative bonded: %v, %w", err, adapter.ErrRepoQuery)
	}
	totalCumulativeFloat := new(big.Float).SetInt(totalCumulative)

	var cumulativeStmtBuilder sq.SelectBuilder
	if pagination.Type() == viewrepo.PAGINATION_OFFSET {
		cumulativeStmtBuilder = repo.stmtBuilder.Select(
			"CASE WHEN COUNT(*) = 0 THEN 0 ELSE SUM(ssa.bonded) END",
		).FromSelect(
			baseStmtBuilder.Columns(
				"sa.bonded",
			).Offset(0).Limit(pagination.OffsetParams().Offset()),
			"ssa",
		)
	} else {
		panic(fmt.Sprintf("unsupported pagination type: %s", pagination.Type()))
	}
	cumulativeSql, cumulativeSqlArgs, err := cumulativeStmtBuilder.ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building council nodes cumulative bonded SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}
	cumulativeReader := repo.typeConv.NtobReader()
	if err = repo.conn.QueryRow(
		cumulativeSql, cumulativeSqlArgs...,
	).Scan(cumulativeReader.ScannableArg()); err != nil {
		return nil, nil, fmt.Errorf("error executing council nodes cumulative bonded SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}
	cumulative, err := cumulativeReader.Parse()
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing council nodes cumulative bonded: %v, %w", err, adapter.ErrRepoQuery)
	}
	cumulativeFloat := new(big.Float).SetInt(cumulative)

	rDbPagination := adapter.NewRDbPaginationBuilder(
		pagination,
		repo.conn,
	).BuildStmt(baseStmtBuilder.Columns(
		"c.id",
		"c.name",
		"c.security_contact",
		"c.pubkey_type",
		"c.pubkey",
		"c.address",
		"sa.address AS staking_account_address",
		"sa.nonce",
		"sa.bonded",
		"sa.unbonded",
		"sa.unbonded_from",
		"sa.punishment_kind",
		"sa.jailed_until",
		"c.created_at_block_height",
		"c.last_left_at_block_height",
	))
	sql, sqlArgs, err := rDbPagination.ToStmtBuilder().ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building council nodes select SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}
	rowsResult, err := repo.conn.Query(sql, sqlArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing council nodes select SQL: %v: %w", err, adapter.ErrRepoQuery)
	}
	councilNodes := make([]viewrepo.CouncilNodeListItem, 0)
	for rowsResult.Next() {
		var councilNode viewrepo.CouncilNodeListItem
		var stakingAccount viewrepo.CouncilNodeStakingAccount

		bondedReader := repo.typeConv.NtobReader()
		unbondedReader := repo.typeConv.NtobReader()
		unbondedFromReader := repo.typeConv.NtotReader()
		jailedUntilReader := repo.typeConv.NtotReader()
		if err = rowsResult.Scan(
			&councilNode.Id,
			&councilNode.Name,
			&councilNode.MaybeSecurityContact,
			&councilNode.PubKeyType,
			&councilNode.PubKey,
			&councilNode.Address,
			&stakingAccount.MaybeAddress,
			&stakingAccount.MaybeNonce,
			bondedReader.ScannableArg(),
			unbondedReader.ScannableArg(),
			unbondedFromReader.ScannableArg(),
			&stakingAccount.MaybePunishmentKind,
			jailedUntilReader.ScannableArg(),
			&councilNode.CreatedAtBlockHeight,
			&councilNode.MaybeLastLeftAtBlockHeight,
		); err != nil {
			return nil, nil, fmt.Errorf("error scanning council node row: %v: %w", err, adapter.ErrRepoQuery)
		}

		if stakingAccount.MaybeAddress != nil {
			if stakingAccount.MaybeBonded, err = bondedReader.ParseW(); err != nil {
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

			councilNode.StakingAccount = &stakingAccount
		}

		bondedFloat := new(big.Float).SetInt(councilNode.StakingAccount.MaybeBonded.Int)

		councilNode.SharePercentage, _ = new(big.Float).Quo(
			bondedFloat, totalCumulativeFloat,
		).Float64()

		cumulativeFloat = cumulativeFloat.Add(cumulativeFloat, bondedFloat)
		councilNode.CumulativeSharePercentage, _ = new(big.Float).Quo(
			cumulativeFloat, totalCumulativeFloat,
		).Float64()

		councilNodes = append(councilNodes, councilNode)
	}

	paginationResult, err := rDbPagination.Result()
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing pagination result: %v", err)
	}

	return councilNodes, paginationResult, nil
}

func (repo *RDbCouncilNodeViewRepo) FindById(id uint64) (*viewrepo.CouncilNode, error) {
	var err error

	sql, sqlArgs, err := repo.stmtBuilder.Select(
		"c.id",
		"c.name",
		"c.security_contact",
		"c.pubkey_type",
		"c.pubkey",
		"c.address",
		"sa.address AS staking_account_address",
		"sa.nonce",
		"sa.bonded",
		"sa.unbonded",
		"sa.unbonded_from",
		"sa.punishment_kind",
		"sa.jailed_until",
		"c.created_at_block_height",
		"c.last_left_at_block_height",
		"c.last_left_at_block_height IS NULL AS is_active",
	).From(
		"council_nodes c",
	).LeftJoin(
		"staking_accounts sa ON c.id = sa.current_council_node_id",
	).Where(
		"c.id = ?", id,
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building council node select SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	var councilNode viewrepo.CouncilNode
	var stakingAccount viewrepo.CouncilNodeStakingAccount

	bondedReader := repo.typeConv.NtobReader()
	unbondedReader := repo.typeConv.NtobReader()
	unbondedFromReader := repo.typeConv.NtotReader()
	jailedUntilReader := repo.typeConv.NtotReader()
	if err = repo.conn.QueryRow(sql, sqlArgs...).Scan(
		&councilNode.Id,
		&councilNode.Name,
		&councilNode.MaybeSecurityContact,
		&councilNode.PubKeyType,
		&councilNode.PubKey,
		&councilNode.Address,
		&stakingAccount.MaybeAddress,
		&stakingAccount.MaybeNonce,
		bondedReader.ScannableArg(),
		unbondedReader.ScannableArg(),
		unbondedFromReader.ScannableArg(),
		&stakingAccount.MaybePunishmentKind,
		jailedUntilReader.ScannableArg(),
		&councilNode.CreatedAtBlockHeight,
		&councilNode.MaybeLastLeftAtBlockHeight,
		&councilNode.IsActive,
	); err != nil {
		if err == adapter.ErrNoRows {
			return nil, adapter.ErrNotFound
		}
		return nil, fmt.Errorf("error scanning council node row: %v: %w", err, adapter.ErrRepoQuery)
	}

	if stakingAccount.MaybeAddress != nil {
		if stakingAccount.MaybeBonded, err = bondedReader.ParseW(); err != nil {
			return nil, fmt.Errorf("error parsing bonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if stakingAccount.MaybeUnbonded, err = unbondedReader.ParseW(); err != nil {
			return nil, fmt.Errorf("error parsing unbonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if stakingAccount.MaybeUnbondedFrom, err = unbondedFromReader.Parse(); err != nil {
			return nil, fmt.Errorf("error parsing unbonded from: %v: %w", err, adapter.ErrRepoQuery)
		}
		if stakingAccount.MaybeJailedUntil, err = jailedUntilReader.Parse(); err != nil {
			return nil, fmt.Errorf("error parsing jailed until: %v: %w", err, adapter.ErrRepoQuery)
		}

		councilNode.StakingAccount = &stakingAccount
	}

	return &councilNode, nil
}

func (repo *RDbCouncilNodeViewRepo) ListActivitiesById(councilNodeId uint64, pagination *viewrepo.Pagination) ([]viewrepo.StakingAccountActivity, *viewrepo.PaginationResult, error) {
	var err error

	rDbPagination := adapter.NewRDbPaginationBuilder(
		pagination,
		repo.conn,
	).BuildStmt(repo.stmtBuilder.Select(
		"a.type",
		"a.block_height",
		"b.time",
		"b.hash",
		"a.txid",
		"a.event_position",
		"a.fee",
		"a.inputs",
		"a.joined_council_node",
		"a.output_count",
		"a.staking_account_address",
		"a.staking_account_nonce",
		"a.bonded",
		"a.unbonded",
		"a.unbonded_from",
		"a.jailed_until",
		"a.punishment_kind",
		"a.affected_council_node",
	).From(
		"activities a",
	).LeftJoin(
		"blocks b ON a.block_height = b.height",
	).Where(
		`a.staking_account_address = (
			SELECT staking_account_address
			FROM activities ja
			LEFT JOIN council_nodes jcn ON ja.joined_council_node_id = jcn.id
			WHERE joined_council_node_id = ? AND block_height = jcn.created_at_block_height
		) AND
		EXISTS
		(
			SELECT 1	
			FROM council_nodes scn
			WHERE scn.id = ? AND 
				a.block_height BETWEEN scn.created_at_block_height AND (
					CASE WHEN scn.last_left_at_block_height IS NULL
					THEN (SELECT MAX(height) FROM blocks)
					ELSE scn.last_left_at_block_height END
				)
		)`, councilNodeId, councilNodeId,
	).OrderBy(
		"block_height DESC",
	))

	sql, sqlArgs, err := rDbPagination.ToStmtBuilder().ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building council node activities select SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	rowsResult, err := repo.conn.Query(sql, sqlArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing council node activities select SQL: %v: %w", err, adapter.ErrRepoQuery)
	}

	activities := make([]viewrepo.StakingAccountActivity, 0)
	for rowsResult.Next() {
		var activity viewrepo.StakingAccountActivity

		blockTimeReader := repo.typeConv.NtotReader()
		feeReader := repo.typeConv.NtobReader()
		var inputsJSON *string
		var joinedCouncilNodeJSON *string
		bondedReader := repo.typeConv.NtobReader()
		unbondedReader := repo.typeConv.NtobReader()
		unbondedFromReader := repo.typeConv.NtotReader()
		jailedUntilReader := repo.typeConv.NtotReader()
		var affectedCouncilNodeJSON *string
		if err = rowsResult.Scan(
			&activity.Type,
			&activity.BlockHeight,
			blockTimeReader.ScannableArg(),
			&activity.BlockHash,
			&activity.MaybeTxID,
			&activity.MaybeEventPosition,
			feeReader.ScannableArg(),
			&inputsJSON,
			&joinedCouncilNodeJSON,
			&activity.MaybeOutputCount,
			&activity.MaybeStakingAccountAddress,
			&activity.MaybeStakingAccountNonce,
			bondedReader.ScannableArg(),
			unbondedReader.ScannableArg(),
			unbondedFromReader.ScannableArg(),
			jailedUntilReader.ScannableArg(),
			&activity.MaybePunishmentKind,
			&affectedCouncilNodeJSON,
		); err != nil {
			return nil, nil, fmt.Errorf("error scanning activity row: %v: %w", err, adapter.ErrRepoQuery)
		}
		var blockTime *time.Time
		blockTime, err = blockTimeReader.Parse()
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing block time: %v: %w", err, adapter.ErrRepoQuery)
		}
		activity.BlockTime = *blockTime
		if activity.MaybeFee, err = feeReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing fee: %v: %w", err, adapter.ErrRepoQuery)
		}
		if activity.MaybeBonded, err = bondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing bonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if activity.MaybeUnbonded, err = unbondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if activity.MaybeUnbondedFrom, err = unbondedFromReader.Parse(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded from: %v: %w", err, adapter.ErrRepoQuery)
		}
		if activity.MaybeJailedUntil, err = jailedUntilReader.Parse(); err != nil {
			return nil, nil, fmt.Errorf("error parsing jailed until: %v: %w", err, adapter.ErrRepoQuery)
		}
		if inputsJSON != nil {
			var inputs []viewrepo.TransactionInput
			if err = jsoniter.Unmarshal([]byte(*inputsJSON), &inputs); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling inputs JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			activity.MaybeInputs = inputs
		}
		if joinedCouncilNodeJSON != nil {
			var joinedCouncilNode viewrepo.ActivityCouncilNode
			if err = jsoniter.Unmarshal([]byte(*joinedCouncilNodeJSON), &joinedCouncilNode); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling joined council node JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			activity.MaybeJoinedCouncilNode = &joinedCouncilNode
		}
		if affectedCouncilNodeJSON != nil {
			var affectedCouncilNode viewrepo.ActivityCouncilNode
			if err = jsoniter.Unmarshal([]byte(*affectedCouncilNodeJSON), &affectedCouncilNode); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling affected council node JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			activity.MaybeAffectedCouncilNode = &affectedCouncilNode
		}

		activities = append(activities, activity)
	}

	paginationResult, err := rDbPagination.Result()
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing pagination result: %v", err)
	}

	return activities, paginationResult, nil
}

func (repo *RDbCouncilNodeViewRepo) Stats() (*viewrepo.CouncilNodeStats, error) {
	var err error

	sql, sqlArgs, err := repo.stmtBuilder.Select(
		"COUNT(*), SUM(sa.bonded + sa.unbonded)",
	).From(
		"council_nodes c",
	).Join(
		"staking_accounts sa ON c.id = sa.current_council_node_id",
	).Where(
		"last_left_at_block_height IS NULL",
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building council node stats SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	var councilNodeStats viewrepo.CouncilNodeStats
	totalStakedReader := repo.typeConv.NtobReader()
	err = repo.conn.QueryRow(sql, sqlArgs...).Scan(&councilNodeStats.Count, totalStakedReader.ScannableArg())
	if err != nil {
		return nil, fmt.Errorf("error scanning council node stats row: %v, %w", err, adapter.ErrRepoQuery)
	}

	councilNodeStats.TotalStaked, err = totalStakedReader.ParseW()
	if err != nil {
		return nil, fmt.Errorf("error parsing total staked: %v: %w", err, adapter.ErrRepoQuery)
	}

	return &councilNodeStats, nil
}

func (repo *RDbCouncilNodeViewRepo) Search(keyword string, pagination *viewrepo.Pagination) ([]viewrepo.CouncilNode, *viewrepo.PaginationResult, error) {
	var err error

	likeKeyword := "%" + keyword + "%"
	rDbPagination := adapter.NewRDbPaginationBuilder(
		pagination,
		repo.conn,
	).BuildStmt(repo.stmtBuilder.Select(
		"c.id",
		"c.name",
		"c.security_contact",
		"c.pubkey_type",
		"c.pubkey",
		"c.address",
		"sa.address AS staking_account_address",
		"sa.nonce",
		"sa.bonded",
		"sa.unbonded",
		"sa.unbonded_from",
		"sa.punishment_kind",
		"sa.jailed_until",
		"c.created_at_block_height",
		"c.last_left_at_block_height",
		"c.last_left_at_block_height IS NULL AS is_active",
	).From(
		"council_nodes c",
	).LeftJoin(
		"staking_accounts sa ON c.id = sa.current_council_node_id",
	).OrderBy(
		"(CASE WHEN c.last_left_at_block_height IS NULL THEN 1 ELSE 2 END) ASC, c.id DESC",
	).Where(
		"c.name LIKE ? OR c.security_contact LIKE ? OR c.pubkey = ? OR c.address = ? OR sa.address = ?",
		likeKeyword, likeKeyword, keyword, keyword, keyword,
	))

	sql, sqlArgs, err := rDbPagination.ToStmtBuilder().ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building council nodes select SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	rowsResult, err := repo.conn.Query(sql, sqlArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing council nodes select SQL: %v: %w", err, adapter.ErrRepoQuery)
	}

	councilNodes := make([]viewrepo.CouncilNode, 0)
	for rowsResult.Next() {
		var councilNode viewrepo.CouncilNode
		var stakingAccount viewrepo.CouncilNodeStakingAccount

		bondedReader := repo.typeConv.NtobReader()
		unbondedReader := repo.typeConv.NtobReader()
		unbondedFromReader := repo.typeConv.NtotReader()
		jailedUntilReader := repo.typeConv.NtotReader()
		if err = rowsResult.Scan(
			&councilNode.Id,
			&councilNode.Name,
			&councilNode.MaybeSecurityContact,
			&councilNode.PubKeyType,
			&councilNode.PubKey,
			&councilNode.Address,
			&stakingAccount.MaybeAddress,
			&stakingAccount.MaybeNonce,
			bondedReader.ScannableArg(),
			unbondedReader.ScannableArg(),
			unbondedFromReader.ScannableArg(),
			&stakingAccount.MaybePunishmentKind,
			jailedUntilReader.ScannableArg(),
			&councilNode.CreatedAtBlockHeight,
			&councilNode.MaybeLastLeftAtBlockHeight,
			&councilNode.IsActive,
		); err != nil {
			return nil, nil, fmt.Errorf("error scanning council node row: %v: %w", err, adapter.ErrRepoQuery)
		}

		if stakingAccount.MaybeAddress != nil {
			if stakingAccount.MaybeBonded, err = bondedReader.ParseW(); err != nil {
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

			councilNode.StakingAccount = &stakingAccount
		}

		councilNodes = append(councilNodes, councilNode)
	}

	paginationResult, err := rDbPagination.Result()
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing pagination result: %v", err)
	}

	return councilNodes, paginationResult, nil
}

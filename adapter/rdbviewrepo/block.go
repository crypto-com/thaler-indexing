package rdbviewrepo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	jsoniter "github.com/json-iterator/go"

	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
)

type RDbBlockViewRepo struct {
	conn adapter.RDbConn

	stmtBuilder sq.StatementBuilderType
	typeConv    adapter.RDbTypeConv
}

func NewRDbBlockViewRepo(
	conn adapter.RDbConn,
	stmtBuilder sq.StatementBuilderType,
	typeConv adapter.RDbTypeConv,
) *RDbBlockViewRepo {
	return &RDbBlockViewRepo{
		conn,

		stmtBuilder,
		typeConv,
	}
}

func (repo *RDbBlockViewRepo) LatestBlockHeight() (uint64, error) {
	var err error

	sql, _, err := repo.stmtBuilder.Select("MAX(height)").From("blocks").ToSql()
	if err != nil {
		return uint64(0), fmt.Errorf("error building block select SQL: %v: %w", err, adapter.ErrBuildSQLStmt)
	}

	var latestBlockHeight *uint64
	if err = repo.conn.QueryRow(sql).Scan(&latestBlockHeight); err != nil {
		if err == adapter.ErrNoRows {
			return uint64(0), nil
		}
		return uint64(0), fmt.Errorf("error executing query: %v: %w", err, adapter.ErrRepoQuery)
	}

	if latestBlockHeight == nil {
		return uint64(0), nil
	}
	return uint64(*latestBlockHeight), nil
}

func (repo *RDbBlockViewRepo) ListBlocks(filter viewrepo.BlockFilter, pagination *viewrepo.Pagination) ([]viewrepo.Block, *viewrepo.PaginationResult, error) {
	var err error

	stmtBuilder := repo.stmtBuilder.Select(
		"hash", "height", "time", "app_hash",
		"committed_council_nodes->0",
		"SUM(CASE WHEN a.type IN "+RDB_SQL_TRANSACTION_TYPES+" THEN 1 ELSE 0 END) AS transaction_count",
		"SUM(CASE WHEN a.type IN "+RDB_SQL_EVENT_TYPES+" THEN 1 ELSE 0 END) AS event_count",
		"committed_council_nodes",
	).From(
		"blocks b",
	).LeftJoin(
		"activities a ON b.height = a.block_height",
	).GroupBy(
		"b.hash, b.height",
	).OrderBy(
		"b.height DESC",
	)

	if filter.MaybeProposers != nil {
		filterTypesSize := len(filter.MaybeProposers)
		if filterTypesSize != 0 {
			preparedTypesQuery := "committed_council_nodes->0->'id' IN (" + strings.TrimRight(strings.Repeat("?,", filterTypesSize), ",") + ")"
			typeValues := make([]interface{}, 0, filterTypesSize)
			for _, proposerId := range filter.MaybeProposers {
				typeValues = append(typeValues, strconv.FormatUint(proposerId, 10))
			}
			stmtBuilder = stmtBuilder.Where(preparedTypesQuery, typeValues...)
		}
	}

	rDbPagination := adapter.NewRDbPaginationBuilder(
		pagination,
		repo.conn,
	).BuildStmt(stmtBuilder)

	sql, sqlArgs, err := rDbPagination.ToStmtBuilder().ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building blocks select SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	rowsResult, err := repo.conn.Query(sql, sqlArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing blocks select SQL: %v: %w", err, adapter.ErrRepoQuery)
	}

	blocks := make([]viewrepo.Block, 0)
	for rowsResult.Next() {
		var block viewrepo.Block
		var proposerJSON *string
		var commmittedCouncilNodesJSON *string
		timeReader := repo.typeConv.NtotReader()
		if err = rowsResult.Scan(
			&block.Hash,
			&block.Height,
			timeReader.ScannableArg(),
			&block.AppHash,
			&proposerJSON,
			&block.TransactionCount,
			&block.EventCount,
			&commmittedCouncilNodesJSON,
		); err != nil {
			return nil, nil, fmt.Errorf("error scanning block row: %v: %w", err, adapter.ErrRepoQuery)
		}
		var blockTime *time.Time
		blockTime, err = timeReader.Parse()
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing block time: %v: %w", err, adapter.ErrRepoQuery)
		}
		block.Time = *blockTime

		if proposerJSON != nil {
			var proposer *viewrepo.BlockCouncilNode
			if err = jsoniter.Unmarshal([]byte(*proposerJSON), &proposer); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling proposer JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			block.MaybeProposer = proposer
		}

		if commmittedCouncilNodesJSON != nil {
			var committedCouncilNodes []viewrepo.BlockCouncilNode
			if err = jsoniter.Unmarshal([]byte(*commmittedCouncilNodesJSON), &committedCouncilNodes); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling block council nodes JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			block.MaybeCommittedCouncilNodes = committedCouncilNodes
		}

		blocks = append(blocks, block)
	}

	paginationResult, err := rDbPagination.Result()
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing pagination result: %v", err)
	}

	return blocks, paginationResult, nil
}

func (repo *RDbBlockViewRepo) FindBlock(blockIdentity viewrepo.BlockIdentity) (*viewrepo.Block, error) {
	var err error

	selectBuilder := repo.stmtBuilder.Select(
		"hash", "height", "time", "app_hash",
		"committed_council_nodes->0",
		"SUM(CASE WHEN a.type IN "+RDB_SQL_TRANSACTION_TYPES+" THEN 1 ELSE 0 END) AS transaction_count",
		"SUM(CASE WHEN a.type IN "+RDB_SQL_EVENT_TYPES+" THEN 1 ELSE 0 END) AS event_count",
		"committed_council_nodes",
	).From(
		"blocks b",
	).LeftJoin(
		"activities a ON b.height = a.block_height",
	).GroupBy(
		"b.hash, b.height",
	).OrderBy(
		"b.height DESC",
	)
	if blockIdentity.MaybeHash != nil {
		selectBuilder = selectBuilder.Where("b.hash = ?", *blockIdentity.MaybeHash)
	} else {
		selectBuilder = selectBuilder.Where("b.height = ?", *blockIdentity.MaybeHeight)
	}

	sql, sqlArgs, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building blocks select SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	var block viewrepo.Block
	var proposerJSON *string
	var committedCouncilNodesJSON *string
	timeReader := repo.typeConv.NtotReader()
	if err = repo.conn.QueryRow(sql, sqlArgs...).Scan(
		&block.Hash,
		&block.Height,
		timeReader.ScannableArg(),
		&block.AppHash,
		&proposerJSON,
		&block.TransactionCount,
		&block.EventCount,
		&committedCouncilNodesJSON,
	); err != nil {
		if err == adapter.ErrNoRows {
			return nil, adapter.ErrNotFound
		}
		return nil, fmt.Errorf("error scanning block row: %v: %w", err, adapter.ErrRepoQuery)
	}
	blockTime, err := timeReader.Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing block time: %v: %w", err, adapter.ErrRepoQuery)
	}
	block.Time = *blockTime

	if proposerJSON != nil {
		var proposer *viewrepo.BlockCouncilNode
		if err = jsoniter.Unmarshal([]byte(*proposerJSON), &proposer); err != nil {
			return nil, fmt.Errorf("error unmarshalling proposer JSON: %v: %w", err, adapter.ErrRepoQuery)
		}

		block.MaybeProposer = proposer
	}

	if committedCouncilNodesJSON != nil {
		var committedCouncilNodes []viewrepo.BlockCouncilNode
		if err = jsoniter.Unmarshal([]byte(*committedCouncilNodesJSON), &committedCouncilNodes); err != nil {
			return nil, fmt.Errorf("error unmarshalling block council nodes JSON: %v: %w", err, adapter.ErrRepoQuery)
		}

		block.MaybeCommittedCouncilNodes = committedCouncilNodes
	}

	return &block, nil
}

func (repo *RDbBlockViewRepo) ListBlockTransactions(
	blockIdentity viewrepo.BlockIdentity,
	pagination *viewrepo.Pagination,
) ([]viewrepo.Transaction, *viewrepo.PaginationResult, error) {
	var err error

	selectBuilder := repo.stmtBuilder.Select(
		"a.type",
		"a.block_height",
		"b.time",
		"b.hash",
		"a.txid",
		"a.fee",
		"a.inputs",
		"a.joined_council_node",
		"a.output_count",
		"a.staking_account_address",
		"a.staking_account_nonce",
		"a.bonded",
		"a.unbonded",
		"a.unbonded_from",
		"a.affected_council_node",
	).From(
		"activities a",
	).LeftJoin(
		"blocks b ON a.block_height = b.height",
	).Where(
		"a.type IN " + RDB_SQL_TRANSACTION_TYPES,
	).OrderBy(
		"a.id",
	)
	if blockIdentity.MaybeHash != nil {
		selectBuilder = selectBuilder.Where("b.hash = ?", *blockIdentity.MaybeHash)
	} else {
		selectBuilder = selectBuilder.Where("b.height = ?", *blockIdentity.MaybeHeight)
	}

	rDbPagination := adapter.NewRDbPaginationBuilder(
		pagination,
		repo.conn,
	).BuildStmt(selectBuilder)

	sql, sqlArgs, err := rDbPagination.ToStmtBuilder().ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building block transactions select SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	rowsResult, err := repo.conn.Query(sql, sqlArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing block transactions select SQL: %v: %w", err, adapter.ErrRepoQuery)
	}

	// TODO:
	blockTransactions := make([]viewrepo.Transaction, 0)
	for rowsResult.Next() {
		var blockTransaction viewrepo.Transaction

		blockTimeReader := repo.typeConv.NtotReader()
		feeReader := repo.typeConv.NtobReader()
		var inputsJSON *string
		var joinedCouncilNodeJSON *string
		var affectedCouncilNodeJSON *string
		bondedReader := repo.typeConv.NtobReader()
		unbondedReader := repo.typeConv.NtobReader()
		unbondedFromReader := repo.typeConv.NtotReader()
		if err = rowsResult.Scan(
			&blockTransaction.Type,
			&blockTransaction.BlockHeight,
			blockTimeReader.ScannableArg(),
			&blockTransaction.BlockHash,
			&blockTransaction.MaybeTxID,
			feeReader.ScannableArg(),
			&inputsJSON,
			&joinedCouncilNodeJSON,
			&blockTransaction.MaybeOutputCount,
			&blockTransaction.MaybeStakingAccountAddress,
			&blockTransaction.MaybeStakingAccountNonce,
			bondedReader.ScannableArg(),
			unbondedReader.ScannableArg(),
			unbondedFromReader.ScannableArg(),
			&affectedCouncilNodeJSON,
		); err != nil {
			return nil, nil, fmt.Errorf("error scanning transaction row: %v: %w", err, adapter.ErrRepoQuery)
		}
		var blockTime *time.Time
		blockTime, err = blockTimeReader.Parse()
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing block time: %v: %w", err, adapter.ErrRepoQuery)
		}
		blockTransaction.BlockTime = *blockTime
		if blockTransaction.MaybeFee, err = feeReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing fee: %v: %w", err, adapter.ErrRepoQuery)
		}
		if blockTransaction.MaybeBonded, err = bondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing bonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if blockTransaction.MaybeUnbonded, err = unbondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if blockTransaction.MaybeUnbondedFrom, err = unbondedFromReader.Parse(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded from: %v: %w", err, adapter.ErrRepoQuery)
		}
		if inputsJSON != nil {
			var inputs []viewrepo.TransactionInput
			if err = jsoniter.Unmarshal([]byte(*inputsJSON), &inputs); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling inputs JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			blockTransaction.MaybeInputs = inputs
		}
		if joinedCouncilNodeJSON != nil {
			var joinedCouncilNode viewrepo.ActivityCouncilNode
			if err = jsoniter.Unmarshal([]byte(*joinedCouncilNodeJSON), &joinedCouncilNode); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling joined council node JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			blockTransaction.MaybeJoinedCouncilNode = &joinedCouncilNode
		}
		if affectedCouncilNodeJSON != nil {
			var affectedCouncilNode viewrepo.ActivityCouncilNode
			if err = jsoniter.Unmarshal([]byte(*affectedCouncilNodeJSON), &affectedCouncilNode); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling affected council node JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			blockTransaction.MaybeJoinedCouncilNode = &affectedCouncilNode
		}

		blockTransactions = append(blockTransactions, blockTransaction)
	}

	paginationResult, err := rDbPagination.Result()
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing pagination result: %v", err)
	}

	return blockTransactions, paginationResult, nil
}

func (repo *RDbBlockViewRepo) ListBlockEvents(
	blockIdentity viewrepo.BlockIdentity,
	pagination *viewrepo.Pagination,
) ([]viewrepo.BlockEvent, *viewrepo.PaginationResult, error) {
	var err error

	var whereClause string
	var whereArgs interface{}
	if blockIdentity.MaybeHash != nil {
		whereClause = "b.hash = ?"
		whereArgs = *blockIdentity.MaybeHash
	} else {
		whereClause = "a.block_height = ?"
		whereArgs = *blockIdentity.MaybeHeight
	}

	rDbPagination := adapter.NewRDbPaginationBuilder(
		pagination,
		repo.conn,
	).BuildSQL(""+
		//nolint:gosec
		`SELECT *
FROM (
	SELECT
		a.type,
		a.block_height,
		b.time,
		b.hash,
		NULL AS event_position,
		NULL AS staking_account_address,
		NULL AS staking_account_nonce,
		NULL AS bonded,
		NULL AS unbonded,
		br.minted AS reward_minted,
		json_agg(json_build_object(
			'event_position',a.event_position,
			'staking_address',a.staking_account_address,
			'bonded',a.bonded::text,
			'affected_council_node',a.affected_council_node
		)) AS reward_distribution,
		NULL AS jailed_until,
		NULL AS punishment_kind,
		NULL AS affected_council_node
	FROM activities a
	LEFT JOIN blocks b ON a.block_height = b.height
	LEFT JOIN block_rewards br ON a.block_height = br.block_height
	WHERE `+strings.Replace(whereClause, "?", "$1", 1)+` AND a.type = 'reward'
	GROUP BY a.type, a.block_height, b.time, b.hash, br.minted

	UNION ALL

	SELECT
		a.type,
		a.block_height,
		b.time,
		b.hash,
		a.event_position,
		a.staking_account_address,
		a.staking_account_nonce,
		a.bonded,
		a.unbonded,
		NULL AS reward_minted,
		NULL AS reward_distribution,
		a.jailed_until,
		a.punishment_kind,
		a.affected_council_node
	FROM activities a
	LEFT JOIN blocks b ON a.block_height = b.height
	WHERE `+strings.Replace(whereClause, "?", "$2", 1)+" AND a.type IN "+RDB_SQL_NON_REWARD_EVENT_TYPES+`
) t
ORDER BY block_height`,
		whereArgs, whereArgs,
	)

	sql, sqlArgs := rDbPagination.ToSQL()
	if err != nil {
		return nil, nil, fmt.Errorf("error building block events select SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	rowsResult, err := repo.conn.Query(sql, sqlArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing block events select SQL: %v: %w", err, adapter.ErrRepoQuery)
	}

	blockEvents := make([]viewrepo.BlockEvent, 0)
	for rowsResult.Next() {
		var blockEvent viewrepo.BlockEvent

		blockTimeReader := repo.typeConv.NtotReader()
		bondedReader := repo.typeConv.NtobReader()
		unbondedReader := repo.typeConv.NtobReader()
		rewardMintedReader := repo.typeConv.NtobReader()
		jailedUntilReader := repo.typeConv.NtotReader()
		var rewardDistributionJSON *string
		var affectedCouncilNodeJSON *string
		if err = rowsResult.Scan(
			&blockEvent.Type,
			&blockEvent.BlockHeight,
			blockTimeReader.ScannableArg(),
			&blockEvent.BlockHash,
			&blockEvent.MaybeEventPosition,
			&blockEvent.MaybeStakingAccountAddress,
			&blockEvent.MaybeStakingAccountNonce,
			bondedReader.ScannableArg(),
			unbondedReader.ScannableArg(),
			rewardMintedReader.ScannableArg(),
			&rewardDistributionJSON,
			jailedUntilReader.ScannableArg(),
			&blockEvent.MaybePunishmentKind,
			&affectedCouncilNodeJSON,
		); err != nil {
			return nil, nil, fmt.Errorf("error scanning event row: %v: %w", err, adapter.ErrRepoQuery)
		}
		var blockTime *time.Time
		blockTime, err = blockTimeReader.Parse()
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing block time: %v: %w", err, adapter.ErrRepoQuery)
		}
		blockEvent.BlockTime = *blockTime
		if blockEvent.MaybeBonded, err = bondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing bonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if blockEvent.MaybeUnbonded, err = unbondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if blockEvent.MaybeRewardMinted, err = rewardMintedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing reward minted: %v: %w", err, adapter.ErrRepoQuery)
		}
		if blockEvent.MaybeJailedUntil, err = jailedUntilReader.Parse(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded from: %v: %w", err, adapter.ErrRepoQuery)
		}
		if rewardDistributionJSON != nil {
			var rewardDistribution []viewrepo.BlockRewardRecord
			if err = jsoniter.Unmarshal([]byte(*rewardDistributionJSON), &rewardDistribution); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling reward distribution JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			blockEvent.MaybeRewardDistribution = rewardDistribution
		}
		if affectedCouncilNodeJSON != nil {
			var affectedCouncilNode viewrepo.ActivityCouncilNode
			if err = jsoniter.Unmarshal([]byte(*affectedCouncilNodeJSON), &affectedCouncilNode); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling affected council node JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			blockEvent.MaybeAffectedCouncilNode = &affectedCouncilNode
		}

		blockEvents = append(blockEvents, blockEvent)
	}

	paginationResult, err := rDbPagination.Result()
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing pagination result: %v", err)
	}

	return blockEvents, paginationResult, nil
}

func (repo *RDbBlockViewRepo) Search(keyword string, pagination *viewrepo.Pagination) ([]viewrepo.Block, *viewrepo.PaginationResult, error) {
	var err error

	rDbPagination := adapter.NewRDbPaginationBuilder(
		pagination,
		repo.conn,
	).BuildStmt(repo.stmtBuilder.Select(
		"hash", "height", "time", "app_hash",
		"committed_council_nodes->0",
		"SUM(CASE WHEN a.type IN "+RDB_SQL_TRANSACTION_TYPES+" THEN 1 ELSE 0 END) AS transaction_count",
		"SUM(CASE WHEN a.type IN "+RDB_SQL_EVENT_TYPES+" THEN 1 ELSE 0 END) AS event_count",
		"committed_council_nodes",
	).From(
		"blocks b",
	).LeftJoin(
		"activities a ON b.height = a.block_height",
	).Where(
		"b.height::text = ? OR b.hash = ?", keyword, keyword,
	).GroupBy(
		"b.hash, b.height",
	).OrderBy(
		"b.height DESC",
	))

	sql, sqlArgs, err := rDbPagination.ToStmtBuilder().ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building blocks search SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	rowsResult, err := repo.conn.Query(sql, sqlArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing blocks search SQL: %v: %w", err, adapter.ErrRepoQuery)
	}

	blocks := make([]viewrepo.Block, 0)
	for rowsResult.Next() {
		var block viewrepo.Block
		var proposerJSON *string
		var committedCouncilNodesJSON *string
		timeReader := repo.typeConv.NtotReader()
		if err = rowsResult.Scan(
			&block.Hash,
			&block.Height,
			timeReader.ScannableArg(),
			&block.AppHash,
			&proposerJSON,
			&block.TransactionCount,
			&block.EventCount,
			&committedCouncilNodesJSON,
		); err != nil {
			return nil, nil, fmt.Errorf("error scanning block row: %v: %w", err, adapter.ErrRepoQuery)
		}
		var blockTime *time.Time
		blockTime, err = timeReader.Parse()
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing block time: %v: %w", err, adapter.ErrRepoQuery)
		}
		block.Time = *blockTime

		if proposerJSON != nil {
			var proposer *viewrepo.BlockCouncilNode
			if err = jsoniter.Unmarshal([]byte(*proposerJSON), &proposer); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling proposer JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			block.MaybeProposer = proposer
		}

		if committedCouncilNodesJSON != nil {
			var committedCouncilNodes []viewrepo.BlockCouncilNode
			if err = jsoniter.Unmarshal([]byte(*committedCouncilNodesJSON), &committedCouncilNodes); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling block council nodes JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			block.MaybeCommittedCouncilNodes = committedCouncilNodes
		}

		blocks = append(blocks, block)
	}

	paginationResult, err := rDbPagination.Result()
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing pagination result: %v", err)
	}

	return blocks, paginationResult, nil
}

package rdbviewrepo

import (
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
	jsoniter "github.com/json-iterator/go"
)

type RDbActivityViewRepo struct {
	conn adapter.RDbConn

	stmtBuilder sq.StatementBuilderType
	typeConv    adapter.RDbTypeConv
}

func NewRDbActivityViewRepo(
	conn adapter.RDbConn,
	stmtBuilder sq.StatementBuilderType,
	typeConv adapter.RDbTypeConv,
) *RDbActivityViewRepo {
	return &RDbActivityViewRepo{
		conn,

		stmtBuilder,
		typeConv,
	}
}

func (repo *RDbActivityViewRepo) ListTransactions(filter viewrepo.TransactionFilter, pagination *viewrepo.Pagination) ([]viewrepo.Transaction, *viewrepo.PaginationResult, error) {
	var err error

	stmtBuilder := repo.stmtBuilder.Select(
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
		"a.id DESC",
	)

	if filter.MaybeTypes != nil {
		filterTypesSize := len(filter.MaybeTypes)
		if filterTypesSize != 0 {
			preparedTypesQuery := "a.type IN (" + strings.TrimRight(strings.Repeat("?,", filterTypesSize), ",") + ")"
			typeValues := make([]interface{}, 0, filterTypesSize)
			for _, t := range filter.MaybeTypes {
				typeValues = append(typeValues, adapter.TransactionTypeToString(t))
			}
			stmtBuilder = stmtBuilder.Where(preparedTypesQuery, typeValues...)
		}
	}

	if filter.MaybeStakingAccountAddress != nil {
		stmtBuilder = stmtBuilder.Where("a.staking_account_address = ?", *filter.MaybeStakingAccountAddress)
	}

	rDbPagination := adapter.NewRDbPaginationBuilder(
		pagination,
		repo.conn,
	).BuildStmt(stmtBuilder)

	sql, sqlArgs, err := rDbPagination.ToStmtBuilder().ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building transactions select SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	rowsResult, err := repo.conn.Query(sql, sqlArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing transactions select SQL: %v: %w", err, adapter.ErrRepoQuery)
	}

	transactions := make([]viewrepo.Transaction, 0)
	for rowsResult.Next() {
		var transaction viewrepo.Transaction

		blockTimeReader := repo.typeConv.NtotReader()
		feeReader := repo.typeConv.NtobReader()
		var transferInputsJSON *string
		var joinedCouncilNodeJSON *string
		var affectedCouncilNodeJSON *string
		bondedReader := repo.typeConv.NtobReader()
		unbondedReader := repo.typeConv.NtobReader()
		unbondedFromReader := repo.typeConv.NtotReader()
		if err = rowsResult.Scan(
			&transaction.Type,
			&transaction.BlockHeight,
			blockTimeReader.ScannableArg(),
			&transaction.BlockHash,
			&transaction.MaybeTxID,
			feeReader.ScannableArg(),
			&transferInputsJSON,
			&joinedCouncilNodeJSON,
			&transaction.MaybeOutputCount,
			&transaction.MaybeStakingAccountAddress,
			&transaction.MaybeStakingAccountNonce,
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
		transaction.BlockTime = *blockTime
		if transaction.MaybeFee, err = feeReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing fee: %v: %w", err, adapter.ErrRepoQuery)
		}
		if transaction.MaybeBonded, err = bondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing bonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if transaction.MaybeUnbonded, err = unbondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if transaction.MaybeUnbondedFrom, err = unbondedFromReader.Parse(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded from: %v: %w", err, adapter.ErrRepoQuery)
		}
		if transferInputsJSON != nil {
			var inputs []viewrepo.TransactionInput
			if err = jsoniter.Unmarshal([]byte(*transferInputsJSON), &inputs); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling inputs JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			transaction.MaybeInputs = inputs
		}
		if joinedCouncilNodeJSON != nil {
			var joinedCouncilNode viewrepo.ActivityCouncilNode
			if err = jsoniter.Unmarshal([]byte(*joinedCouncilNodeJSON), &joinedCouncilNode); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling joined council node JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			transaction.MaybeJoinedCouncilNode = &joinedCouncilNode
		}
		if affectedCouncilNodeJSON != nil {
			var affectedCouncilNode viewrepo.ActivityCouncilNode
			if err = jsoniter.Unmarshal([]byte(*affectedCouncilNodeJSON), &affectedCouncilNode); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling affected council node JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			transaction.MaybeAffectedCouncilNode = &affectedCouncilNode
		}

		transactions = append(transactions, transaction)
	}

	paginationResult, err := rDbPagination.Result()
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing pagination result: %v", err)
	}

	return transactions, paginationResult, nil
}

func (repo *RDbActivityViewRepo) FindTransactionByTxId(txid string) (*viewrepo.Transaction, error) {
	var err error

	sql, sqlArgs, err := repo.stmtBuilder.Select(
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
		"a.txid = ?", txid,
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building transaction query SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	var transaction viewrepo.Transaction

	blockTimeReader := repo.typeConv.NtotReader()
	feeReader := repo.typeConv.NtobReader()
	var transferInputsJSON *string
	var joinedCouncilNodeJSON *string
	var affectedCouncilNodeJSON *string
	bondedReader := repo.typeConv.NtobReader()
	unbondedReader := repo.typeConv.NtobReader()
	unbondedFromReader := repo.typeConv.NtotReader()
	if err = repo.conn.QueryRow(sql, sqlArgs...).Scan(
		&transaction.Type,
		&transaction.BlockHeight,
		blockTimeReader.ScannableArg(),
		&transaction.BlockHash,
		&transaction.MaybeTxID,
		feeReader.ScannableArg(),
		&transferInputsJSON,
		&joinedCouncilNodeJSON,
		&transaction.MaybeOutputCount,
		&transaction.MaybeStakingAccountAddress,
		&transaction.MaybeStakingAccountNonce,
		bondedReader.ScannableArg(),
		unbondedReader.ScannableArg(),
		unbondedFromReader.ScannableArg(),
		&affectedCouncilNodeJSON,
	); err != nil {
		if err != adapter.ErrNoRows {
			return nil, adapter.ErrNotFound
		}
		return nil, fmt.Errorf("error scanning transaction row: %v: %w", err, adapter.ErrRepoQuery)
	}
	blockTime, err := blockTimeReader.Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing block time: %v: %w", err, adapter.ErrRepoQuery)
	}
	transaction.BlockTime = *blockTime
	if transaction.MaybeFee, err = feeReader.ParseW(); err != nil {
		return nil, fmt.Errorf("error parsing fee: %v: %w", err, adapter.ErrRepoQuery)
	}
	if transaction.MaybeBonded, err = bondedReader.ParseW(); err != nil {
		return nil, fmt.Errorf("error parsing bonded: %v: %w", err, adapter.ErrRepoQuery)
	}
	if transaction.MaybeUnbonded, err = unbondedReader.ParseW(); err != nil {
		return nil, fmt.Errorf("error parsing unbonded: %v: %w", err, adapter.ErrRepoQuery)
	}
	if transaction.MaybeUnbondedFrom, err = unbondedFromReader.Parse(); err != nil {
		return nil, fmt.Errorf("error parsing unbonded from: %v: %w", err, adapter.ErrRepoQuery)
	}
	if transferInputsJSON != nil {
		var inputs []viewrepo.TransactionInput
		if err = jsoniter.Unmarshal([]byte(*transferInputsJSON), &inputs); err != nil {
			return nil, fmt.Errorf("error unmarshalling inputs JSON: %v: %w", err, adapter.ErrRepoQuery)
		}

		transaction.MaybeInputs = inputs
	}
	if joinedCouncilNodeJSON != nil {
		var joinedCouncilNode viewrepo.ActivityCouncilNode
		if err = jsoniter.Unmarshal([]byte(*joinedCouncilNodeJSON), &joinedCouncilNode); err != nil {
			return nil, fmt.Errorf("error unmarshalling joined council node JSON: %v: %w", err, adapter.ErrRepoQuery)
		}

		transaction.MaybeJoinedCouncilNode = &joinedCouncilNode
	}
	if affectedCouncilNodeJSON != nil {
		var affectedCouncilNode viewrepo.ActivityCouncilNode
		if err = jsoniter.Unmarshal([]byte(*affectedCouncilNodeJSON), &affectedCouncilNode); err != nil {
			return nil, fmt.Errorf("error unmarshalling affected council node JSON: %v: %w", err, adapter.ErrRepoQuery)
		}

		transaction.MaybeAffectedCouncilNode = &affectedCouncilNode
	}

	return &transaction, nil
}

func (repo *RDbActivityViewRepo) ListEvents(filter viewrepo.EventFilter, pagination *viewrepo.Pagination) ([]viewrepo.Event, *viewrepo.PaginationResult, error) {
	var err error

	stmtBuilder := repo.stmtBuilder.Select(
		"a.type",
		"a.block_height",
		"b.time",
		"b.hash",
		"a.event_position",
		"a.staking_account_address",
		"a.staking_account_nonce",
		"a.bonded",
		"a.unbonded",
		"a.jailed_until",
		"a.punishment_kind",
		"a.affected_council_node",
	).From(
		"activities a",
	).LeftJoin(
		"blocks b ON a.block_height = b.height",
	).Where(
		"a.type IN " + RDB_SQL_EVENT_TYPES,
	).OrderBy(
		"a.id DESC",
	)

	if filter.MaybeTypes != nil {
		filterTypesSize := len(filter.MaybeTypes)
		if filterTypesSize != 0 {
			preparedTypesQuery := "a.type IN (" + strings.TrimRight(strings.Repeat("?,", filterTypesSize), ",") + ")"
			typeValues := make([]interface{}, 0, filterTypesSize)
			for _, t := range filter.MaybeTypes {
				typeValues = append(typeValues, adapter.EventTypeToString(t))
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
		return nil, nil, fmt.Errorf("error building events list SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	rowsResult, err := repo.conn.Query(sql, sqlArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing events list SQL: %v: %w", err, adapter.ErrRepoQuery)
	}

	events := make([]viewrepo.Event, 0)
	for rowsResult.Next() {
		var event viewrepo.Event

		blockTimeReader := repo.typeConv.NtotReader()
		bondedReader := repo.typeConv.NtobReader()
		unbondedReader := repo.typeConv.NtobReader()
		jailedUntilReader := repo.typeConv.NtotReader()
		var affectedCouncilNodeJSON *string
		if err = rowsResult.Scan(
			&event.Type,
			&event.BlockHeight,
			blockTimeReader.ScannableArg(),
			&event.BlockHash,
			&event.EventPosition,
			&event.StakingAccountAddress,
			&event.MaybeStakingAccountNonce,
			bondedReader.ScannableArg(),
			unbondedReader.ScannableArg(),
			jailedUntilReader.ScannableArg(),
			&event.MaybePunishmentKind,
			&affectedCouncilNodeJSON,
		); err != nil {
			return nil, nil, fmt.Errorf("error scanning event row: %v: %w", err, adapter.ErrRepoQuery)
		}
		var blockTime *time.Time
		blockTime, err = blockTimeReader.Parse()
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing block time: %v: %w", err, adapter.ErrRepoQuery)
		}
		event.BlockTime = *blockTime
		if event.MaybeBonded, err = bondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing bonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if event.MaybeUnbonded, err = unbondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if event.MaybeJailedUntil, err = jailedUntilReader.Parse(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded from: %v: %w", err, adapter.ErrRepoQuery)
		}
		if affectedCouncilNodeJSON != nil {
			var affectedCouncilNode viewrepo.ActivityCouncilNode
			if err = jsoniter.Unmarshal([]byte(*affectedCouncilNodeJSON), &affectedCouncilNode); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling affected council node JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			event.MaybeAffectedCouncilNode = &affectedCouncilNode
		}

		events = append(events, event)
	}

	paginationResult, err := rDbPagination.Result()
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing pagination result: %v", err)
	}

	return events, paginationResult, nil
}

func (repo *RDbActivityViewRepo) FindEventByBlockHeightEventPosition(blockHeight uint64, eventPosition uint64) (*viewrepo.Event, error) {
	var err error

	sql, sqlArgs, err := repo.stmtBuilder.Select(
		"a.type",
		"a.block_height",
		"b.time",
		"b.hash",
		"a.event_position",
		"a.staking_account_address",
		"a.staking_account_nonce",
		"a.bonded",
		"a.unbonded",
		"a.jailed_until",
		"a.punishment_kind",
		"a.affected_council_node",
	).From(
		"activities a",
	).LeftJoin(
		"blocks b ON a.block_height = b.height",
	).Where(
		"a.block_height = ? AND a.event_position = ?", blockHeight, eventPosition,
	).ToSql()
	if err != nil {
		if err != adapter.ErrNoRows {
			return nil, adapter.ErrNotFound
		}
		return nil, fmt.Errorf("error building event select SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	var event viewrepo.Event

	blockTimeReader := repo.typeConv.NtotReader()
	bondedReader := repo.typeConv.NtobReader()
	unbondedReader := repo.typeConv.NtobReader()
	jailedUntilReader := repo.typeConv.NtotReader()
	var affectedCouncilNodeJSON *string
	if err = repo.conn.QueryRow(sql, sqlArgs...).Scan(
		&event.Type,
		&event.BlockHeight,
		blockTimeReader.ScannableArg(),
		&event.BlockHash,
		&event.EventPosition,
		&event.StakingAccountAddress,
		&event.MaybeStakingAccountNonce,
		bondedReader.ScannableArg(),
		unbondedReader.ScannableArg(),
		jailedUntilReader.ScannableArg(),
		&event.MaybePunishmentKind,
		&affectedCouncilNodeJSON,
	); err != nil {
		return nil, fmt.Errorf("error scanning event row: %v: %w", err, adapter.ErrRepoQuery)
	}
	blockTime, err := blockTimeReader.Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing block time: %v: %w", err, adapter.ErrRepoQuery)
	}
	event.BlockTime = *blockTime
	if event.MaybeBonded, err = bondedReader.ParseW(); err != nil {
		return nil, fmt.Errorf("error parsing bonded: %v: %w", err, adapter.ErrRepoQuery)
	}
	if event.MaybeUnbonded, err = unbondedReader.ParseW(); err != nil {
		return nil, fmt.Errorf("error parsing unbonded: %v: %w", err, adapter.ErrRepoQuery)
	}
	if event.MaybeJailedUntil, err = jailedUntilReader.Parse(); err != nil {
		return nil, fmt.Errorf("error parsing unbonded from: %v: %w", err, adapter.ErrRepoQuery)
	}
	if affectedCouncilNodeJSON != nil {
		var affectedCouncilNode viewrepo.ActivityCouncilNode
		if err = jsoniter.Unmarshal([]byte(*affectedCouncilNodeJSON), &affectedCouncilNode); err != nil {
			return nil, fmt.Errorf("error unmarshalling affected council node JSON: %v: %w", err, adapter.ErrRepoQuery)
		}

		event.MaybeAffectedCouncilNode = &affectedCouncilNode
	}

	return &event, nil
}

func (repo *RDbActivityViewRepo) TransactionsCount() (uint64, error) {
	var err error

	sql, sqlArgs, err := repo.stmtBuilder.Select(
		"COUNT(*)",
	).From(
		"activities",
	).Where(
		"type IN " + RDB_SQL_TRANSACTION_TYPES,
	).ToSql()
	if err != nil {
		return uint64(0), fmt.Errorf("error building transactions count SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	var transactionsCount uint64
	err = repo.conn.QueryRow(sql, sqlArgs...).Scan(&transactionsCount)
	if err != nil {
		return uint64(0), fmt.Errorf("error scanning transactions count row: %v, %w", err, adapter.ErrRepoQuery)
	}

	return transactionsCount, nil
}

func (repo *RDbActivityViewRepo) SearchTransactions(keyword string, pagination *viewrepo.Pagination) ([]viewrepo.Transaction, *viewrepo.PaginationResult, error) {
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
		"a.type IN "+RDB_SQL_TRANSACTION_TYPES+` AND (
			a.type::text = ? OR a.txid = ? OR b.height::text = ? OR b.hash = ? OR a.staking_account_address = ?
		)`, keyword, keyword, keyword, keyword, keyword,
	).OrderBy(
		"a.id DESC",
	))

	sql, sqlArgs, err := rDbPagination.ToStmtBuilder().ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building transactions search SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	rowsResult, err := repo.conn.Query(sql, sqlArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing transactions search SQL: %v: %w", err, adapter.ErrRepoQuery)
	}

	transactions := make([]viewrepo.Transaction, 0)
	for rowsResult.Next() {
		var transaction viewrepo.Transaction

		blockTimeReader := repo.typeConv.NtotReader()
		feeReader := repo.typeConv.NtobReader()
		var transferInputsJSON *string
		var joinedCouncilNodeJSON *string
		var affectedCouncilNodeJSON *string
		bondedReader := repo.typeConv.NtobReader()
		unbondedReader := repo.typeConv.NtobReader()
		unbondedFromReader := repo.typeConv.NtotReader()
		if err = rowsResult.Scan(
			&transaction.Type,
			&transaction.BlockHeight,
			blockTimeReader.ScannableArg(),
			&transaction.BlockHash,
			&transaction.MaybeTxID,
			feeReader.ScannableArg(),
			&transferInputsJSON,
			&joinedCouncilNodeJSON,
			&transaction.MaybeOutputCount,
			&transaction.MaybeStakingAccountAddress,
			&transaction.MaybeStakingAccountNonce,
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
		transaction.BlockTime = *blockTime
		if transaction.MaybeFee, err = feeReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing fee: %v: %w", err, adapter.ErrRepoQuery)
		}
		if transaction.MaybeBonded, err = bondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing bonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if transaction.MaybeUnbonded, err = unbondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if transaction.MaybeUnbondedFrom, err = unbondedFromReader.Parse(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded from: %v: %w", err, adapter.ErrRepoQuery)
		}
		if transferInputsJSON != nil {
			var inputs []viewrepo.TransactionInput
			if err = jsoniter.Unmarshal([]byte(*transferInputsJSON), &inputs); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling inputs JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			transaction.MaybeInputs = inputs
		}
		if joinedCouncilNodeJSON != nil {
			var joinedCouncilNode viewrepo.ActivityCouncilNode
			if err = jsoniter.Unmarshal([]byte(*joinedCouncilNodeJSON), &joinedCouncilNode); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling joined council node JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			transaction.MaybeJoinedCouncilNode = &joinedCouncilNode
		}
		if affectedCouncilNodeJSON != nil {
			var affectedCouncilNode viewrepo.ActivityCouncilNode
			if err = jsoniter.Unmarshal([]byte(*affectedCouncilNodeJSON), &affectedCouncilNode); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling affected council node JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			transaction.MaybeAffectedCouncilNode = &affectedCouncilNode
		}

		transactions = append(transactions, transaction)
	}

	paginationResult, err := rDbPagination.Result()
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing pagination result: %v", err)
	}

	return transactions, paginationResult, nil
}

func (repo *RDbActivityViewRepo) SearchEvents(keyword string, pagination *viewrepo.Pagination) ([]viewrepo.Event, *viewrepo.PaginationResult, error) {
	var err error

	rDbPagination := adapter.NewRDbPaginationBuilder(
		pagination,
		repo.conn,
	).BuildStmt(repo.stmtBuilder.Select(
		"a.type",
		"a.block_height",
		"b.time",
		"b.hash",
		"a.event_position",
		"a.staking_account_address",
		"a.staking_account_nonce",
		"a.bonded",
		"a.unbonded",
		"a.jailed_until",
		"a.punishment_kind",
		"a.affected_council_node",
	).From(
		"activities a",
	).LeftJoin(
		"blocks b ON a.block_height = b.height",
	).Where(
		"a.type IN "+RDB_SQL_EVENT_TYPES+" AND a.type::text = ?", keyword,
	).OrderBy(
		"a.id DESC",
	))

	sql, sqlArgs, err := rDbPagination.ToStmtBuilder().ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("error building events search SQL: %v, %w", err, adapter.ErrBuildSQLStmt)
	}

	rowsResult, err := repo.conn.Query(sql, sqlArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("error executing events search SQL: %v: %w", err, adapter.ErrRepoQuery)
	}

	events := make([]viewrepo.Event, 0)
	for rowsResult.Next() {
		var event viewrepo.Event

		blockTimeReader := repo.typeConv.NtotReader()
		bondedReader := repo.typeConv.NtobReader()
		unbondedReader := repo.typeConv.NtobReader()
		jailedUntilReader := repo.typeConv.NtotReader()
		var affectedCouncilNodeJSON *string
		if err = rowsResult.Scan(
			&event.Type,
			&event.BlockHeight,
			blockTimeReader.ScannableArg(),
			&event.BlockHash,
			&event.EventPosition,
			&event.StakingAccountAddress,
			&event.MaybeStakingAccountNonce,
			bondedReader.ScannableArg(),
			unbondedReader.ScannableArg(),
			jailedUntilReader.ScannableArg(),
			&event.MaybePunishmentKind,
			&affectedCouncilNodeJSON,
		); err != nil {
			return nil, nil, fmt.Errorf("error scanning event row: %v: %w", err, adapter.ErrRepoQuery)
		}
		var blockTime *time.Time
		if blockTime, err = blockTimeReader.Parse(); err != nil {
			return nil, nil, fmt.Errorf("error parsing block time: %v: %w", err, adapter.ErrRepoQuery)
		}
		event.BlockTime = *blockTime
		if event.MaybeBonded, err = bondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing bonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if event.MaybeUnbonded, err = unbondedReader.ParseW(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded: %v: %w", err, adapter.ErrRepoQuery)
		}
		if event.MaybeJailedUntil, err = jailedUntilReader.Parse(); err != nil {
			return nil, nil, fmt.Errorf("error parsing unbonded from: %v: %w", err, adapter.ErrRepoQuery)
		}
		if affectedCouncilNodeJSON != nil {
			var affectedCouncilNode viewrepo.ActivityCouncilNode
			if err = jsoniter.Unmarshal([]byte(*affectedCouncilNodeJSON), &affectedCouncilNode); err != nil {
				return nil, nil, fmt.Errorf("error unmarshalling affected council node JSON: %v: %w", err, adapter.ErrRepoQuery)
			}

			event.MaybeAffectedCouncilNode = &affectedCouncilNode
		}

		events = append(events, event)
	}

	paginationResult, err := rDbPagination.Result()
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing pagination result: %v", err)
	}

	return events, paginationResult, nil
}

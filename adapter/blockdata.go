package adapter

import (
	"fmt"
	"sort"

	sq "github.com/Masterminds/squirrel"
	jsoniter "github.com/json-iterator/go"

	"github.com/crypto-com/chainindex"
	"github.com/crypto-com/chainindex/usecase"
)

type RDbBlockDataRepo struct {
	conn        RDbConn
	stmtBuilder sq.StatementBuilderType
	typeConv    RDbTypeConv

	activityDataRepo RDbBlockActivityDataRepo
}

func NewRDbBlockDataRepo(
	conn RDbConn,
	stmtBuilder sq.StatementBuilderType,
	typeConv RDbTypeConv,

	activityDataRepo RDbBlockActivityDataRepo,
) *RDbBlockDataRepo {
	return &RDbBlockDataRepo{
		conn,
		stmtBuilder,
		typeConv,

		activityDataRepo,
	}
}

func (repo *RDbBlockDataRepo) Store(blockData *usecase.BlockData) error {
	// FIXME: Persist block data into event store and create projection
	var err error

	tx, err := repo.conn.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v: %w", err, ErrRepoOpen)
	}
	defer func() {
		// Calling rollback on committed transaction has no effect
		_ = tx.Rollback()
	}()

	var committedCouncilNodes []RDbBlockCommittedCouncilNodeRow
	if committedCouncilNodes, err = repo.parseSignaturesToCommittedCouncilNodeRows(tx, blockData.Block.Height, blockData.Signatures); err != nil {
		return err
	}

	if err = repo.insertBlock(tx, &blockData.Block, committedCouncilNodes); err != nil {
		return err
	}

	if err = repo.insertBlockCommittedCouncilNodes(tx, committedCouncilNodes); err != nil {
		return err
	}

	if err = repo.storeActivities(tx, blockData.Activities); err != nil {
		return err
	}

	if blockData.Reward != nil {
		if err = repo.insertReward(tx, blockData.Reward); err != nil {
			return err
		}
	}

	if err = repo.storeCouncilNodeUpdates(tx, blockData.CouncilNodeUpdates, blockData.Block.Height); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error commiting block data: %v: %w", err, ErrRepoWrite)
	}

	return nil
}

func (repo *RDbBlockDataRepo) parseSignaturesToCommittedCouncilNodeRows(
	tx RDbTx,
	blockHeight uint64,
	signatures []chainindex.BlockSignature,
) ([]RDbBlockCommittedCouncilNodeRow, error) {
	if signatures == nil {
		return nil, nil
	}

	sort.SliceStable(signatures, func(i, j int) bool {
		if signatures[i].IsProposer {
			return true
		} else if signatures[j].IsProposer {
			return false
		} else {
			return i < j
		}
	})

	rows := make([]RDbBlockCommittedCouncilNodeRow, 0, len(signatures))
	for _, signature := range signatures {
		councilNodeId, councilNodeName, err := repo.findLatestCouncilNodeByAddress(tx, signature.CouncilNodeAddress)
		if err != nil {
			return nil, fmt.Errorf("error querying council node by signature address: %v", err)
		}
		rows = append(rows, RDbBlockCommittedCouncilNodeRow{
			BlockHeight:        blockHeight,
			ID:                 councilNodeId,
			Name:               councilNodeName,
			CouncilNodeAddress: signature.CouncilNodeAddress,
			Signature:          signature.Signature,
			IsProposer:         signature.IsProposer,
		})
	}

	return rows, nil
}

func (repo *RDbBlockDataRepo) insertBlock(tx RDbTx, block *chainindex.Block, committedCouncilNodes []RDbBlockCommittedCouncilNodeRow) error {
	var err error

	sql, _, err := repo.stmtBuilder.Insert(
		"blocks",
	).Columns(
		"height",
		"hash",
		"time",
		"app_hash",
		"committed_council_nodes",
	).Values("?", "?", "?", "?", "?").ToSql()
	if err != nil {
		return fmt.Errorf("error building block insert SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	var committedCoundcilNodesJSON *string
	if len(committedCouncilNodes) > 0 {
		var jsonStr string
		if jsonStr, err = jsoniter.MarshalToString(committedCouncilNodes); err != nil {
			return fmt.Errorf("error building committed council nodes JSON for insertion: %v: %w", err, ErrBuildSQLStmt)
		}
		committedCoundcilNodesJSON = &jsonStr
	}

	result, err := tx.Exec(sql,
		block.Height,
		block.Hash,
		repo.typeConv.Tton(&block.Time),
		block.AppHash,
		committedCoundcilNodesJSON,
	)
	if err != nil {
		return fmt.Errorf("error inserting block into the table: %v: %w", err, ErrRepoWrite)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("error inserting block into the table: %w", err)
	}

	return nil
}

func (repo *RDbBlockDataRepo) insertBlockCommittedCouncilNodes(tx RDbTx, rows []RDbBlockCommittedCouncilNodeRow) error {
	for _, row := range rows {
		var err error

		// nolint:gosec,scopelint
		if err = repo.insertBlockCommittedCouncilNodeRow(tx, &row); err != nil {
			return err
		}
	}

	return nil
}

func (repo *RDbBlockDataRepo) insertBlockCommittedCouncilNodeRow(tx RDbTx, row *RDbBlockCommittedCouncilNodeRow) error {
	var err error

	sql, _, err := repo.stmtBuilder.Insert(
		"block_committed_council_nodes",
	).Columns(
		"block_height",
		"council_node_id",
		"signature",
		"is_proposer",
	).Values("?", "?", "?", "?").ToSql()
	if err != nil {
		return fmt.Errorf("error building block signature insert SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	result, err := tx.Exec(sql, row.BlockHeight, row.ID, row.Signature, row.IsProposer)
	if err != nil {
		return fmt.Errorf("error inserting block signature into the table: %v: %w", err, ErrRepoWrite)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("error inserting block signature into the table: no row inserted: %w", ErrRepoWrite)
	}

	return nil
}

func (repo *RDbBlockDataRepo) storeActivities(tx RDbTx, activities []chainindex.Activity) error {
	var err error
	for _, activity := range activities {
		if err = repo.storeActivity(tx, activity); err != nil {
			return fmt.Errorf("error storing activity into database: %w", err)
		}
	}

	return nil
}

func (repo *RDbBlockDataRepo) storeActivity(tx RDbTx, activity chainindex.Activity) error {
	switch activity.Type {
	case chainindex.ACTIVITY_GENESIS:
		return repo.activityDataRepo.InsertGenesisActivity(tx, &activity)
	case chainindex.ACTIVITY_TRANSFER:
		return repo.activityDataRepo.InsertTransferTransaction(tx, &activity)
	case chainindex.ACTIVITY_DEPOSIT:
		return repo.activityDataRepo.InsertDepositTransaction(tx, &activity)
	case chainindex.ACTIVITY_UNBOND:
		return repo.activityDataRepo.InsertUnbondTransaction(tx, &activity)
	case chainindex.ACTIVITY_WITHDRAW:
		return repo.activityDataRepo.InsertWithdrawTransaction(tx, &activity)
	case chainindex.ACTIVITY_NODEJOIN:
		return repo.activityDataRepo.InsertNodeJoinTransaction(tx, &activity)
	case chainindex.ACTIVITY_UNJAIL:
		return repo.activityDataRepo.InsertUnjailTransaction(tx, &activity)
	case chainindex.ACTIVITY_REWARD:
		return repo.activityDataRepo.InsertRewardEvent(tx, &activity)
	case chainindex.ACTIVITY_SLASH:
		return repo.activityDataRepo.InsertSlashEvent(tx, &activity)
	case chainindex.ACTIVITY_JAIL:
		return repo.activityDataRepo.InsertJailEvent(tx, &activity)
	default:
		return fmt.Errorf("Unrecognized activity type: %v", activity.Type)
	}
}

func (repo *RDbBlockDataRepo) insertReward(tx RDbTx, reward *chainindex.BlockReward) error {
	var err error

	sql, _, err := repo.stmtBuilder.Insert(
		"block_rewards",
	).Columns(
		"block_height",
		"minted",
	).Values("?", "?").ToSql()
	if err != nil {
		return fmt.Errorf("error building block reward insert SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	result, err := tx.Exec(sql, reward.BlockHeight, repo.typeConv.Bton(reward.Minted))
	if err != nil {
		return fmt.Errorf("error inserting block reward into the table: %v: %w", err, ErrRepoWrite)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("error inserting block reward into the table: no row inserted: %w", ErrRepoWrite)
	}

	return nil
}

func (repo *RDbBlockDataRepo) storeCouncilNodeUpdates(tx RDbTx, updates []chainindex.CouncilNodeUpdate, blockHeight uint64) error {
	for _, update := range updates {
		var err error
		// TODO: insert node kicked activity into database

		if update.Type != chainindex.COUNCIL_NODE_UPDATE_TYPE_LEFT {
			continue
		}

		councilNodeId, _, err := repo.findLatestCouncilNodeByAddress(tx, update.Address)
		if err != nil {
			return fmt.Errorf("error querying council node by Tendermint address: %v", err)
		}

		if err = repo.updateCouncilNodeLastUpdatedAtBlockHeight(tx, blockHeight, councilNodeId); err != nil {
			return err
		}

		if err = repo.removeStakingAccountCurrentCouncilNodeId(tx, councilNodeId); err != nil {
			return err
		}
	}

	return nil
}

func (repo *RDbBlockDataRepo) findLatestCouncilNodeByAddress(tx RDbTx, address string) (uint64, string, error) {
	var err error

	sql, _, err := repo.stmtBuilder.Select(
		"id",
		"name",
	).From(
		"council_nodes",
	).Where(
		"address = ?",
	).OrderBy(
		"id DESC",
	).ToSql()
	if err != nil {
		return uint64(0), "", fmt.Errorf("error building council node query SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	var councilNodeId uint64
	var councilNodeName string
	if err = tx.QueryRow(sql, address).Scan(&councilNodeId, &councilNodeName); err != nil {
		if err == ErrNoRows {
			return uint64(0), "", nil
		}
	}

	return councilNodeId, councilNodeName, nil
}

func (repo *RDbBlockDataRepo) removeStakingAccountCurrentCouncilNodeId(tx RDbTx, councilNodeId uint64) error {
	var err error

	sql, sqlArgs, err := repo.stmtBuilder.Update(
		"staking_accounts",
	).Set(
		"current_council_node_id", nil,
	).Where(
		"current_council_node_id = ?", councilNodeId,
	).ToSql()
	if err != nil {
		return fmt.Errorf("error building remove council node id from staking account update SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	execResult, err := tx.Exec(sql, sqlArgs...)
	if err != nil {
		return fmt.Errorf("error removing council node id from staking account: %v: %w", err, ErrRepoWrite)
	}
	if execResult.RowsAffected() != 1 {
		return fmt.Errorf("error removing council node id from staking account: no row updated: %w", ErrRepoWrite)
	}

	return nil
}

func (repo *RDbBlockDataRepo) updateCouncilNodeLastUpdatedAtBlockHeight(tx RDbTx, blockHeight uint64, councilNodeId uint64) error {

	var err error

	sql, _, err := repo.stmtBuilder.Update(
		"council_nodes",
	).Set(
		"last_left_at_block_height", "?",
	).Where(
		"id = ?",
	).ToSql()
	if err != nil {
		return fmt.Errorf("error building council nodes update SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	result, err := tx.Exec(sql, blockHeight, councilNodeId)
	if err != nil {
		return fmt.Errorf("error updating council node into the table: %v: %w", err, ErrRepoWrite)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("error updating council node into the table: no row updated: %w", ErrRepoWrite)
	}

	return nil
}

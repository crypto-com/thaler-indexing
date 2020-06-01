package adapter

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	jsoniter "github.com/json-iterator/go"

	"github.com/crypto-com/chainindex"
	"github.com/crypto-com/chainindex/internal/bignum"
	"github.com/crypto-com/chainindex/internal/primptr"
)

type RDbBlockActivityDataRepo interface {
	InsertGenesisActivity(tx RDbTx, activity *chainindex.Activity) error
	InsertTransferTransaction(tx RDbTx, activity *chainindex.Activity) error
	InsertDepositTransaction(tx RDbTx, activity *chainindex.Activity) error
	InsertUnbondTransaction(tx RDbTx, activity *chainindex.Activity) error
	InsertWithdrawTransaction(tx RDbTx, activity *chainindex.Activity) error
	InsertNodeJoinTransaction(tx RDbTx, activity *chainindex.Activity) error
	InsertUnjailTransaction(tx RDbTx, activity *chainindex.Activity) error
	InsertRewardEvent(tx RDbTx, activity *chainindex.Activity) error
	InsertSlashEvent(tx RDbTx, activity *chainindex.Activity) error
	InsertJailEvent(tx RDbTx, activity *chainindex.Activity) error
}

type DefaultRDbBlockActivityDataRepo struct {
	stmtBuilder sq.StatementBuilderType

	typeConv RDbTypeConv
}

func NewDefaultRDbBlockActivityDataRepo(stmtBuilder sq.StatementBuilderType, typeConv RDbTypeConv) *DefaultRDbBlockActivityDataRepo {
	return &DefaultRDbBlockActivityDataRepo{
		stmtBuilder,

		typeConv,
	}
}

func (repo *DefaultRDbBlockActivityDataRepo) InsertGenesisActivity(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	var insertedCouncilNodeId *uint64
	if activity.MaybeCouncilNodeMeta != nil {
		councilNodeRow := &RDbCouncilNodeRow{
			ID:                         nil,
			Name:                       activity.MaybeCouncilNodeMeta.Name,
			MaybeSecurityContact:       activity.MaybeCouncilNodeMeta.MaybeSecurityContact,
			PubKeyType:                 PubKeyTypeToString(activity.MaybeCouncilNodeMeta.PubKeyType),
			PubKey:                     activity.MaybeCouncilNodeMeta.PubKey,
			Address:                    activity.MaybeCouncilNodeMeta.Address,
			CreatedAtBlockHeight:       activity.BlockHeight,
			MaybeLastLeftAtBlockHeight: nil,
		}
		insertedCouncilNodeId, err = repo.insertCouncilNode(tx, councilNodeRow)
		if err != nil {
			return err
		}
	}

	bonded, unbonded := bignum.Int0(), bignum.Int0()
	if activity.MaybeBonded != nil {
		bonded = activity.MaybeBonded
	}
	if activity.MaybeUnbonded != nil {
		unbonded = activity.MaybeUnbonded
	}
	accountRow := &RDbStakingAccountRow{
		Address:              *activity.MaybeStakingAccountAddress,
		Nonce:                uint64(0),
		Bonded:               bonded,
		Unbonded:             unbonded,
		UnbondedFrom:         nil,
		PunishmentKind:       nil,
		JailedUntil:          nil,
		CurrentCouncilNodeId: insertedCouncilNodeId,
	}
	if err = repo.insertStakingAccount(tx, accountRow); err != nil {
		return err
	}

	activity.MaybeStakingAccountNonce = primptr.Uint64(accountRow.Nonce)
	if activity.MaybeCouncilNodeMeta != nil {
		activity.MaybeCouncilNodeMeta.Id = insertedCouncilNodeId
	}

	if err = repo.insertActivity(tx, activity); err != nil {
		return err
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) InsertTransferTransaction(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	if err = repo.insertActivity(tx, activity); err != nil {
		return err
	}

	if err = repo.updateUsedPrevTxOutputs(tx, activity); err != nil {
		return err
	}

	if err = repo.insertTransactionOutputs(tx, activity); err != nil {
		return err
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) updateUsedPrevTxOutputs(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	sql, _, err := repo.stmtBuilder.Update(
		"transaction_outputs",
	).Set(
		"spent_at_txid", "?",
	).Where(
		"txid = ? AND index = ?",
	).ToSql()
	if err != nil {
		return fmt.Errorf("error building transaction outputs update SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	txInputs := activity.MaybeTxInputs
	for i, l := 0, len(txInputs); i < l; i += 1 {
		result, err := tx.Exec(sql,
			activity.MaybeTxID,
			txInputs[i].TxId,
			txInputs[i].Index,
		)
		if err != nil {
			return fmt.Errorf("error updating transaction input into table: %v: %w", err, ErrRepoWrite)
		}
		if result.RowsAffected() != 1 {
			return fmt.Errorf("error updating transaction input into table: no row updated: %w", ErrRepoWrite)
		}
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) insertTransactionOutputs(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	sql, _, err := repo.stmtBuilder.Insert(
		"transaction_outputs",
	).Columns(
		"txid",
		"index",
	).Values("?", "?").ToSql()
	if err != nil {
		return fmt.Errorf("error building transaction outputs insert SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	for i, l := uint32(0), *activity.MaybeOutputCount; i < l; i += 1 {
		result, err := tx.Exec(sql, activity.MaybeTxID, i)
		if err != nil {
			return fmt.Errorf("error inserting transaction output into table: %v: %w", err, ErrRepoWrite)
		}
		if result.RowsAffected() != 1 {
			return fmt.Errorf("error insertion transaction output into table: no row updated: %w", ErrRepoWrite)
		}
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) InsertDepositTransaction(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	var accountRow *RDbStakingAccountRow

	accountRow, err = repo.findStakingAccount(tx, activity.MaybeStakingAccountAddress)
	if err != nil {
		return err
	}

	if accountRow == nil {
		accountRow = &RDbStakingAccountRow{
			Address:              *activity.MaybeStakingAccountAddress,
			Nonce:                uint64(0),
			Bonded:               activity.MaybeBonded,
			Unbonded:             bignum.Int0(),
			UnbondedFrom:         nil,
			PunishmentKind:       nil,
			JailedUntil:          nil,
			CurrentCouncilNodeId: nil,
		}
		if err = repo.insertStakingAccount(tx, accountRow); err != nil {
			return err
		}
	} else {
		// Deposit transaction does not increment nonce
		accountRow.AddBonded(activity.MaybeBonded)

		if err = repo.updateStakingAccount(tx, accountRow); err != nil {
			return err
		}
	}

	activity.MaybeStakingAccountNonce = primptr.Uint64(accountRow.Nonce)

	if err = repo.appendAffectedCouncilNodeToActivity(tx, activity); err != nil {
		return err
	}

	if err = repo.insertActivity(tx, activity); err != nil {
		return err
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) insertStakingAccount(tx RDbTx, accountRow *RDbStakingAccountRow) error {
	var err error

	sql, _, err := repo.stmtBuilder.Insert(
		"staking_accounts",
	).Columns(
		"address",
		"nonce",
		"bonded",
		"unbonded",
		"unbonded_from",
		"jailed_until",
		"punishment_kind",
		"current_council_node_id",
	).Values("?", "?", "?", "?", "?", "?", "?", "?").ToSql()
	if err != nil {
		return fmt.Errorf("error building staking account insert SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	result, err := tx.Exec(sql,
		accountRow.Address,
		accountRow.Nonce,
		repo.typeConv.Bton(accountRow.Bonded),
		repo.typeConv.Bton(accountRow.Unbonded),
		accountRow.UnbondedFrom,
		accountRow.JailedUntil,
		accountRow.PunishmentKind,
		accountRow.CurrentCouncilNodeId,
	)
	if err != nil {
		return fmt.Errorf("error inserting transaction output into table: %v: %w", err, ErrRepoWrite)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("error insertion transaction output into table: no row inserted: %w", ErrRepoWrite)
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) InsertUnbondTransaction(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	var accountRow *RDbStakingAccountRow

	accountRow, err = repo.findStakingAccount(tx, activity.MaybeStakingAccountAddress)
	if err != nil {
		return err
	}
	if accountRow == nil {
		panic("staking account of unbond activity does not exist")
	}

	accountRow.IncrementNonce()
	accountRow.AddBonded(activity.MaybeBonded)
	accountRow.AddUnbonded(activity.MaybeUnbonded)

	if err = repo.updateStakingAccount(tx, accountRow); err != nil {
		return err
	}

	activity.MaybeStakingAccountNonce = primptr.Uint64(accountRow.Nonce)

	if err = repo.appendAffectedCouncilNodeToActivity(tx, activity); err != nil {
		return err
	}

	if err = repo.insertActivity(tx, activity); err != nil {
		return err
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) InsertWithdrawTransaction(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	var accountRow *RDbStakingAccountRow

	accountRow, err = repo.findStakingAccount(tx, activity.MaybeStakingAccountAddress)
	if err != nil {
		return err
	}
	if accountRow == nil {
		panic("staking account of withdraw activity does not exist")
	}

	accountRow.IncrementNonce()
	accountRow.AddUnbonded(activity.MaybeUnbonded)

	if err = repo.updateStakingAccount(tx, accountRow); err != nil {
		return err
	}

	activity.MaybeStakingAccountNonce = primptr.Uint64(accountRow.Nonce)

	if err = repo.insertActivity(tx, activity); err != nil {
		return err
	}

	if err = repo.appendAffectedCouncilNodeToActivity(tx, activity); err != nil {
		return err
	}

	if err = repo.insertTransactionOutputs(tx, activity); err != nil {
		return err
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) InsertNodeJoinTransaction(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	accountRow, err := repo.findStakingAccount(tx, activity.MaybeStakingAccountAddress)
	if err != nil {
		return err
	}
	if accountRow == nil {
		panic("staking account of node join activity does not exist")
	}

	lastAddressCouncilNode, err := repo.findLastCouncilNodeByAddress(tx, activity.MaybeCouncilNodeMeta.Address)
	if err != nil {
		return err
	}
	lastStakingAccountCouncilNode, err := repo.findLastCouncilNodeByStakingAccountAddress(tx, *activity.MaybeStakingAccountAddress)
	if err != nil {
		return err
	}

	var councilNodeRow *RDbCouncilNodeRow
	var councilNodeId *uint64
	if lastAddressCouncilNode != nil &&
		lastStakingAccountCouncilNode != nil &&
		*lastAddressCouncilNode.ID == *lastStakingAccountCouncilNode.ID &&
		lastAddressCouncilNode.Name == activity.MaybeCouncilNodeMeta.Name {

		councilNodeRow = lastAddressCouncilNode
		councilNodeId = lastAddressCouncilNode.ID

		if err = repo.clearCouncilNodeLastLeftAtBlockHeight(tx, *councilNodeId); err != nil {
			return err
		}
	} else {
		councilNodeRow = &RDbCouncilNodeRow{
			ID:                         nil,
			Name:                       activity.MaybeCouncilNodeMeta.Name,
			MaybeSecurityContact:       activity.MaybeCouncilNodeMeta.MaybeSecurityContact,
			PubKeyType:                 PubKeyTypeToString(activity.MaybeCouncilNodeMeta.PubKeyType),
			PubKey:                     activity.MaybeCouncilNodeMeta.PubKey,
			Address:                    activity.MaybeCouncilNodeMeta.Address,
			CreatedAtBlockHeight:       activity.BlockHeight,
			MaybeLastLeftAtBlockHeight: nil,
		}
		councilNodeId, err = repo.insertCouncilNode(tx, councilNodeRow)
		if err != nil {
			return err
		}

		councilNodeRow.ID = councilNodeId
	}

	accountRow.IncrementNonce()
	accountRow.CurrentCouncilNodeId = councilNodeId

	if err = repo.updateStakingAccount(tx, accountRow); err != nil {
		return err
	}

	activity.MaybeStakingAccountNonce = primptr.Uint64(accountRow.Nonce)
	activity.MaybeCouncilNodeMeta = RDbCouncilNodeRowToCouncilNode(councilNodeRow)
	activity.MaybeCouncilNodeMeta.Id = councilNodeId

	if err = repo.insertActivity(tx, activity); err != nil {
		return err
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) insertCouncilNode(tx RDbTx, nodeRow *RDbCouncilNodeRow) (*uint64, error) {
	var err error

	sql, _, err := repo.stmtBuilder.Insert(
		"council_nodes",
	).Columns(
		"name",
		"security_contact",
		"pubkey_type",
		"pubkey",
		"address",
		"created_at_block_height",
		"last_left_at_block_height",
	).Values(
		"?", "?", "?", "?", "?", "?", "?",
	).Suffix(
		"RETURNING id",
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building council node insertion SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	var councilNodeId uint64
	if err := tx.QueryRow(sql,
		nodeRow.Name,
		nodeRow.MaybeSecurityContact,
		nodeRow.PubKeyType,
		nodeRow.PubKey,
		nodeRow.Address,
		nodeRow.CreatedAtBlockHeight,
		nodeRow.MaybeLastLeftAtBlockHeight,
	).Scan(&councilNodeId); err != nil {
		return nil, fmt.Errorf("error inserting council node into table: %v: %w", err, ErrRepoWrite)
	}

	return &councilNodeId, nil
}

func (repo *DefaultRDbBlockActivityDataRepo) clearCouncilNodeLastLeftAtBlockHeight(tx RDbTx, councilNodeId uint64) error {
	var err error

	sql, sqlArgs, err := repo.stmtBuilder.Update(
		"council_nodes",
	).Set(
		"last_left_at_block_height", nil,
	).Where(
		"id = ?", councilNodeId,
	).ToSql()
	if err != nil {
		return fmt.Errorf("error building council node update SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	result, err := tx.Exec(sql, sqlArgs...)
	if err != nil {
		return fmt.Errorf("error clearing council node last left at block height into table: %v: %w", err, ErrRepoWrite)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("error clearing council node last left at block height into table: no row updated: %w", ErrRepoWrite)
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) findLastCouncilNodeByAddress(tx RDbTx, address string) (*RDbCouncilNodeRow, error) {
	var err error

	var councilNode RDbCouncilNodeRow
	sql, sqlArgs, err := repo.stmtBuilder.Select(
		"id",
		"name",
		"security_contact",
		"pubkey_type",
		"pubkey",
		"address",
		"created_at_block_height",
		"last_left_at_block_height",
	).From(
		"council_nodes",
	).Where(
		"address = ?", address,
	).OrderBy(
		"id DESC",
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building council nodes query SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	if err = tx.QueryRow(sql, sqlArgs...).Scan(
		&councilNode.ID,
		&councilNode.Name,
		&councilNode.MaybeSecurityContact,
		&councilNode.PubKeyType,
		&councilNode.PubKey,
		&councilNode.Address,
		&councilNode.CreatedAtBlockHeight,
		&councilNode.MaybeLastLeftAtBlockHeight,
	); err != nil {
		if err == ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying council node: %v: %w", err, ErrRepoQuery)
	}

	return &councilNode, nil
}

func (repo *DefaultRDbBlockActivityDataRepo) findLastCouncilNodeByStakingAccountAddress(tx RDbTx, stakingAccountAddress string) (*RDbCouncilNodeRow, error) {
	var err error

	sql, sqlArgs, err := repo.stmtBuilder.Select(
		"c.id",
		"c.name",
		"c.security_contact",
		"c.pubkey_type",
		"c.pubkey",
		"c.address",
		"c.created_at_block_height",
		"c.last_left_at_block_height",
	).From(
		"council_nodes c",
	).Join(
		"activities a ON c.id = a.joined_council_node_id",
	).Where(
		"a.type IN ('genesis', 'nodejoin') AND a.staking_account_address = ?", stakingAccountAddress,
	).OrderBy(
		"a.id DESC",
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building council nodes query SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	var councilNode RDbCouncilNodeRow
	if err = tx.QueryRow(sql, sqlArgs...).Scan(
		&councilNode.ID,
		&councilNode.Name,
		&councilNode.MaybeSecurityContact,
		&councilNode.PubKeyType,
		&councilNode.PubKey,
		&councilNode.Address,
		&councilNode.CreatedAtBlockHeight,
		&councilNode.MaybeLastLeftAtBlockHeight,
	); err != nil {
		if err == ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying council node: %v: %w", err, ErrRepoQuery)
	}

	return &councilNode, nil
}

func (repo *DefaultRDbBlockActivityDataRepo) InsertUnjailTransaction(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	accountRow, err := repo.findStakingAccount(tx, activity.MaybeStakingAccountAddress)
	if err != nil {
		return err
	}
	if accountRow == nil {
		panic("staking account of unjail activity does not exist")
	}

	accountRow.IncrementNonce()
	accountRow.PunishmentKind = nil
	accountRow.JailedUntil = nil

	if err = repo.updateStakingAccount(tx, accountRow); err != nil {
		return err
	}

	activity.MaybeStakingAccountNonce = primptr.Uint64(accountRow.Nonce)

	if err = repo.appendAffectedCouncilNodeToActivity(tx, activity); err != nil {
		return err
	}

	if err = repo.insertActivity(tx, activity); err != nil {
		return err
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) InsertRewardEvent(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	var accountRow *RDbStakingAccountRow

	accountRow, err = repo.findStakingAccount(tx, activity.MaybeStakingAccountAddress)
	if err != nil {
		return err
	}
	if accountRow == nil {
		panic("staking account of unjail activity does not exist")
	}

	accountRow.AddBonded(activity.MaybeBonded)

	if err = repo.updateStakingAccount(tx, accountRow); err != nil {
		return err
	}

	if err = repo.appendAffectedCouncilNodeToActivity(tx, activity); err != nil {
		return err
	}

	if err = repo.insertActivity(tx, activity); err != nil {
		return err
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) InsertSlashEvent(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	var accountRow *RDbStakingAccountRow

	accountRow, err = repo.findStakingAccount(tx, activity.MaybeStakingAccountAddress)
	if err != nil {
		return err
	}
	if accountRow == nil {
		panic("staking account of unjail activity does not exist")
	}

	accountRow.AddBonded(activity.MaybeBonded)
	accountRow.AddUnbonded(activity.MaybeUnbonded)

	if err = repo.updateStakingAccount(tx, accountRow); err != nil {
		return err
	}

	if err = repo.appendAffectedCouncilNodeToActivity(tx, activity); err != nil {
		return err
	}

	if err = repo.insertActivity(tx, activity); err != nil {
		return err
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) InsertJailEvent(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	var accountRow *RDbStakingAccountRow

	accountRow, err = repo.findStakingAccount(tx, activity.MaybeStakingAccountAddress)
	if err != nil {
		return err
	}
	if accountRow == nil {
		panic("staking account of unjail activity does not exist")
	}

	accountRow.JailedUntil = activity.MaybeJailedUntil
	accountRow.PunishmentKind = OptPunishmentKindToString(activity.MaybePunishmentKind)

	if err = repo.updateStakingAccount(tx, accountRow); err != nil {
		return err
	}

	if err = repo.appendAffectedCouncilNodeToActivity(tx, activity); err != nil {
		return err
	}

	if err = repo.insertActivity(tx, activity); err != nil {
		return err
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) findCurrentCouncilNodeByStakingAccountAddress(tx RDbTx, stakingAccountAddress string) (*RDbCouncilNodeRow, error) {
	var err error

	sql, sqlArgs, err := repo.stmtBuilder.Select(
		"c.id",
		"c.name",
		"c.security_contact",
		"c.pubkey_type",
		"c.pubkey",
		"c.address",
		"c.created_at_block_height",
		"c.last_left_at_block_height",
	).From(
		"council_nodes c",
	).Join(
		"staking_accounts sa ON c.id = sa.current_council_node_id",
	).Where(
		"c.last_left_at_block_height IS NULL AND sa.address = ?", stakingAccountAddress,
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building staking account current council node query SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	var councilNode RDbCouncilNodeRow
	if err = tx.QueryRow(sql, sqlArgs...).Scan(
		&councilNode.ID,
		&councilNode.Name,
		&councilNode.MaybeSecurityContact,
		&councilNode.PubKeyType,
		&councilNode.PubKey,
		&councilNode.Address,
		&councilNode.CreatedAtBlockHeight,
		&councilNode.MaybeLastLeftAtBlockHeight,
	); err != nil {
		if err == ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying staking account current council node: %v: %w", err, ErrRepoQuery)
	}

	return &councilNode, nil
}

func (repo *DefaultRDbBlockActivityDataRepo) findStakingAccount(tx RDbTx, address *string) (*RDbStakingAccountRow, error) {
	var err error

	var stakingAccountRow RDbStakingAccountRow
	sql, _, err := repo.stmtBuilder.Select(
		"address,nonce,bonded,unbonded,unbonded_from,punishment_kind,jailed_until,current_council_node_id",
	).From(
		"staking_accounts",
	).Where(
		"address = ?",
	).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building staking account query SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	bondedReader := repo.typeConv.NtobReader()
	unbondedReader := repo.typeConv.NtobReader()
	unbondedFromReader := repo.typeConv.NtotReader()
	jailedUntilReader := repo.typeConv.NtotReader()
	if err = tx.QueryRow(sql, address).Scan(
		&stakingAccountRow.Address,
		&stakingAccountRow.Nonce,
		bondedReader.ScannableArg(),
		unbondedReader.ScannableArg(),
		unbondedFromReader.ScannableArg(),
		&stakingAccountRow.PunishmentKind,
		jailedUntilReader.ScannableArg(),
		&stakingAccountRow.CurrentCouncilNodeId,
	); err != nil {
		if err == ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying staking account: %v: %w", err, ErrRepoQuery)
	}

	if stakingAccountRow.Bonded, err = bondedReader.Parse(); err != nil {
		return nil, fmt.Errorf("error parsing bonded: %v", err)
	}
	if stakingAccountRow.Unbonded, err = unbondedReader.Parse(); err != nil {
		return nil, fmt.Errorf("error parsing unbonded: %v", err)
	}
	if stakingAccountRow.JailedUntil, err = jailedUntilReader.Parse(); err != nil {
		return nil, fmt.Errorf("error parsing jailed until: %v", err)
	}

	return &stakingAccountRow, nil
}

func (repo *DefaultRDbBlockActivityDataRepo) updateStakingAccount(tx RDbTx, accountRow *RDbStakingAccountRow) error {
	var err error

	sql, sqlArgs, err := repo.stmtBuilder.Update(
		"staking_accounts",
	).SetMap(sq.Eq{
		"nonce":                   accountRow.Nonce,
		"bonded":                  repo.typeConv.Bton(accountRow.Bonded),
		"unbonded":                repo.typeConv.Bton(accountRow.Unbonded),
		"unbonded_from":           repo.typeConv.Tton(accountRow.UnbondedFrom),
		"punishment_kind":         accountRow.PunishmentKind,
		"jailed_until":            repo.typeConv.Tton(accountRow.JailedUntil),
		"current_council_node_id": accountRow.CurrentCouncilNodeId,
	}).Where(
		"address = ?", accountRow.Address,
	).ToSql()
	if err != nil {
		return fmt.Errorf("error building staking account update SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	result, err := tx.Exec(sql, sqlArgs...)
	if err != nil {
		return fmt.Errorf("error updating staking account into table: %v: %w", err, ErrRepoWrite)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("error updating staking account into table: no row updated: %w", ErrRepoWrite)
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) insertActivity(tx RDbTx, activity *chainindex.Activity) error {
	var err error

	sql, _, err := repo.stmtBuilder.Insert(
		"activities",
	).Columns(
		"block_height",
		"type",
		"txid",
		"event_position",
		"fee",
		"inputs",
		"output_count",
		"staking_account_address",
		"staking_account_nonce",
		"bonded",
		"unbonded",
		"unbonded_from",
		"joined_council_node",
		"joined_council_node_id",
		"affected_council_node",
		"affected_council_node_id",
		"jailed_until",
		"punishment_kind",
	).Values("?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?").ToSql()
	if err != nil {
		return fmt.Errorf("error building activity insertion SQL: %v: %w", err, ErrBuildSQLStmt)
	}

	var transferInputsJSON *string
	if activity.MaybeTxInputs != nil && len(activity.MaybeTxInputs) > 0 {
		var jsonStr string
		jsonStr, err = jsoniter.MarshalToString(TxInputsToRDbTransferInputs(activity.MaybeTxInputs))
		if err != nil {
			return fmt.Errorf("error marshalling transfer inputs: %v", err)
		}

		transferInputsJSON = &jsonStr
	}
	var joinedCouncilNodeJSON *string
	var joinedCouncilNodeId *uint64
	if activity.MaybeCouncilNodeMeta != nil {
		var jsonStr string
		jsonStr, err = jsoniter.MarshalToString(CouncilNodeToRDbCouncilNodeRow(activity.MaybeCouncilNodeMeta))
		if err != nil {
			return fmt.Errorf("error marshalling joined council node: %v", err)
		}

		joinedCouncilNodeJSON = &jsonStr
		joinedCouncilNodeId = activity.MaybeCouncilNodeMeta.Id
	}
	var affectedCouncilNode *string
	var affectedCouncilNodeId *uint64
	if activity.MaybeAffectedCouncilNode != nil {
		var jsonStr string
		jsonStr, err = jsoniter.MarshalToString(CouncilNodeToRDbCouncilNodeRow(activity.MaybeAffectedCouncilNode))
		if err != nil {
			return fmt.Errorf("error marshalling affected council node: %v", err)
		}

		affectedCouncilNode = &jsonStr
		affectedCouncilNodeId = activity.MaybeAffectedCouncilNode.Id
	}
	result, err := tx.Exec(sql,
		activity.BlockHeight,
		ActivityTypeToString(activity.Type),
		activity.MaybeTxID,
		activity.MaybeEventPosition,
		bignum.OptItoa(activity.MaybeFee),
		transferInputsJSON,
		activity.MaybeOutputCount,
		activity.MaybeStakingAccountAddress,
		activity.MaybeStakingAccountNonce,
		bignum.OptItoa(activity.MaybeBonded),
		bignum.OptItoa(activity.MaybeUnbonded),
		repo.typeConv.Tton(activity.MaybeUnbondedFrom),
		joinedCouncilNodeJSON,
		joinedCouncilNodeId,
		affectedCouncilNode,
		affectedCouncilNodeId,
		repo.typeConv.Tton(activity.MaybeJailedUntil),
		OptPunishmentKindToString(activity.MaybePunishmentKind),
	)
	if err != nil {
		return fmt.Errorf("error inserting activity into table: %v: %w", err, ErrRepoWrite)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("error inserting activity outputs into table: no row updated: %w", ErrRepoWrite)
	}

	return nil
}

func (repo *DefaultRDbBlockActivityDataRepo) appendAffectedCouncilNodeToActivity(tx RDbTx, mutActivity *chainindex.Activity) error {
	var err error
	var stakingAccountCurrentCouncilNode *RDbCouncilNodeRow

	if stakingAccountCurrentCouncilNode, err = repo.findCurrentCouncilNodeByStakingAccountAddress(
		tx, *mutActivity.MaybeStakingAccountAddress,
	); err != nil {
		return err
	}
	if stakingAccountCurrentCouncilNode != nil {
		mutActivity.MaybeAffectedCouncilNode = RDbCouncilNodeRowToCouncilNode(stakingAccountCurrentCouncilNode)
	}

	return nil
}

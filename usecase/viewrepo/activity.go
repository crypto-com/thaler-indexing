package viewrepo

import (
	"time"

	"github.com/crypto-com/chainindex"
	"github.com/crypto-com/chainindex/internal/bignum"
)

type ActivityViewRepo interface {
	ListTransactions(filter TransactionFilter, pagination *Pagination) ([]Transaction, *PaginationResult, error)
	FindTransactionByTxId(txid string) (*Transaction, error)
	ListEvents(filter EventFilter, pagination *Pagination) ([]Event, *PaginationResult, error)
	FindEventByBlockHeightEventPosition(blockHeight uint64, eventPosition uint64) (*Event, error)

	TransactionsCount() (uint64, error)

	SearchTransactions(keyword string, pagination *Pagination) ([]Transaction, *PaginationResult, error)
	SearchEvents(keyword string, pagination *Pagination) ([]Event, *PaginationResult, error)
}

type TransactionFilter struct {
	MaybeTypes                 []chainindex.TransactionType
	MaybeStakingAccountAddress *string
}

type Transaction struct {
	Type                       string               `json:"type"`
	BlockHeight                uint64               `json:"block_height"`
	BlockTime                  time.Time            `json:"block_time"`
	BlockHash                  string               `json:"block_hash"`
	MaybeTxID                  *string              `json:"txid"`
	MaybeFee                   *bignum.WBigInt      `json:"fee"`
	MaybeInputs                []TransactionInput   `json:"inputs"`
	MaybeOutputCount           *uint64              `json:"output_count"`
	MaybeStakingAccountAddress *string              `json:"staking_account_address"`
	MaybeStakingAccountNonce   *uint64              `json:"staking_account_nonce"`
	MaybeBonded                *bignum.WBigInt      `json:"bonded"`
	MaybeUnbonded              *bignum.WBigInt      `json:"unbonded"`
	MaybeUnbondedFrom          *time.Time           `json:"unbonded_from"`
	MaybeJoinedCouncilNode     *ActivityCouncilNode `json:"joined_council_node"`
	MaybeAffectedCouncilNode   *ActivityCouncilNode `json:"affected_council_node"`
}

type TransactionInput struct {
	PrevTxId        string `json:"prev_txid"`
	PrevOutputIndex uint64 `json:"prev_output_index"`
}

type EventFilter struct {
	MaybeTypes                 []chainindex.EventType
	MaybeStakingAccountAddress *string
}

type Event struct {
	Type                     string               `json:"type"`
	BlockHeight              uint64               `json:"block_height"`
	BlockTime                time.Time            `json:"block_time"`
	BlockHash                string               `json:"block_hash"`
	EventPosition            uint64               `json:"event_position"`
	StakingAccountAddress    string               `json:"staking_account_address"`
	MaybeStakingAccountNonce *uint64              `json:"staking_account_nonce"`
	MaybeBonded              *bignum.WBigInt      `json:"bonded"`
	MaybeUnbonded            *bignum.WBigInt      `json:"unbonded"`
	MaybeRewardMinted        *bignum.WBigInt      `json:"reward_minted"`
	MaybeJailedUntil         *time.Time           `json:"jailed_until"`
	MaybePunishmentKind      *string              `json:"punishment_kind"`
	MaybeAffectedCouncilNode *ActivityCouncilNode `json:"affected_council_node"`
}

type ActivityCouncilNode struct {
	ID                         uint64  `json:"id"`
	Name                       string  `json:"name"`
	MaybeSecurityContact       *string `json:"security_contact"`
	PubKeyType                 string  `json:"pubkey_type"`
	PubKey                     string  `json:"pubkey"`
	Address                    string  `json:"address"`
	CreatedAtBlockHeight       uint64  `json:"created_at_block_height"`
	MaybeLastLeftAtBlockHeight *uint64 `json:"last_left_at_block_height"`
}

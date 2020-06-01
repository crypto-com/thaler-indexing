package viewrepo

import (
	"time"

	"github.com/crypto-com/chainindex/internal/bignum"
)

type BlockViewRepo interface {
	LatestBlockHeight() (uint64, error)
	ListBlocks(filter BlockFilter, pagination *Pagination) ([]Block, *PaginationResult, error)
	FindBlock(BlockIdentity) (*Block, error)
	ListBlockTransactions(BlockIdentity, *Pagination) ([]Transaction, *PaginationResult, error)
	ListBlockEvents(BlockIdentity, *Pagination) ([]BlockEvent, *PaginationResult, error)

	Search(keyword string, pagination *Pagination) ([]Block, *PaginationResult, error)
}

type BlockIdentity struct {
	MaybeHash   *string
	MaybeHeight *uint64
}

type BlockFilter struct {
	MaybeProposers []uint64
}

type Block struct {
	Hash                       string             `json:"hash"`
	Height                     uint64             `json:"height"`
	Time                       time.Time          `json:"time"`
	AppHash                    string             `json:"app_hash"`
	MaybeProposer              *BlockCouncilNode  `json:"proposer"`
	TransactionCount           uint64             `json:"transaction_count"`
	EventCount                 uint64             `json:"event_count"`
	MaybeCommittedCouncilNodes []BlockCouncilNode `json:"committed_council_nodes"`
}

type BlockCouncilNode struct {
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Signature  string `json:"signature"`
	IsProposer bool   `json:"is_proposer"`
}

type BlockEvent struct {
	Type                       string               `json:"type"`
	BlockHeight                uint64               `json:"block_height"`
	BlockTime                  time.Time            `json:"block_time"`
	BlockHash                  string               `json:"block_hash"`
	MaybeEventPosition         *uint64              `json:"event_position"`
	MaybeStakingAccountAddress *string              `json:"staking_account_address"`
	MaybeStakingAccountNonce   *uint64              `json:"staking_account_nonce"`
	MaybeBonded                *bignum.WBigInt      `json:"bonded"`
	MaybeUnbonded              *bignum.WBigInt      `json:"unbonded"`
	MaybeRewardMinted          *bignum.WBigInt      `json:"reward_minted"`
	MaybeRewardDistribution    []BlockRewardRecord  `json:"reward_distribution"`
	MaybeJailedUntil           *time.Time           `json:"jailed_until"`
	MaybePunishmentKind        *string              `json:"punishment_kind"`
	MaybeAffectedCouncilNode   *ActivityCouncilNode `json:"affected_council_node"`
}

type BlockRewardRecord struct {
	EventPosition       uint64              `json:"event_position"`
	StakingAddress      string              `json:"staking_address"`
	MaybeBonded         *string             `json:"bonded"`
	AffectedCouncilNode ActivityCouncilNode `json:"affected_council_node"`
}

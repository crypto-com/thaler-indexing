package viewrepo

import (
	"time"

	"github.com/crypto-com/chainindex/internal/bignum"
)

type CouncilNodeViewRepo interface {
	ListActivities(pagination *Pagination) ([]CouncilNodeListItem, *PaginationResult, error)
	FindById(id uint64) (*CouncilNode, error)
	ListActivitiesById(id uint64, pagination *Pagination) ([]StakingAccountActivity, *PaginationResult, error)

	Stats() (*CouncilNodeStats, error)

	Search(keyword string, pagination *Pagination) ([]CouncilNode, *PaginationResult, error)
}

type CouncilNodeListItem struct {
	Id                         uint64                     `json:"id"`
	Name                       string                     `json:"name"`
	MaybeSecurityContact       *string                    `json:"security_contact"`
	PubKeyType                 string                     `json:"pubkey_type"`
	PubKey                     string                     `json:"pubkey"`
	Address                    string                     `json:"address"`
	StakingAccount             *CouncilNodeStakingAccount `json:"staking_account"`
	CreatedAtBlockHeight       uint64                     `json:"created_at_block_height"`
	MaybeLastLeftAtBlockHeight *uint64                    `json:"last_left_at_block_height"`
	SharePercentage            float64                    `json:"share_percentage"`
	CumulativeSharePercentage  float64                    `json:"cumulative_share_percentage"`
}

type CouncilNode struct {
	Id                         uint64                     `json:"id"`
	Name                       string                     `json:"name"`
	MaybeSecurityContact       *string                    `json:"security_contact"`
	PubKeyType                 string                     `json:"pubkey_type"`
	PubKey                     string                     `json:"pubkey"`
	Address                    string                     `json:"address"`
	StakingAccount             *CouncilNodeStakingAccount `json:"staking_account"`
	CreatedAtBlockHeight       uint64                     `json:"created_at_block_height"`
	MaybeLastLeftAtBlockHeight *uint64                    `json:"last_left_at_block_height"`
	IsActive                   bool                       `json:"is_active"`
}

type CouncilNodeStakingAccount struct {
	MaybeAddress        *string         `json:"address"`
	MaybeNonce          *uint64         `json:"nonce"`
	MaybeBonded         *bignum.WBigInt `json:"bonded"`
	MaybeUnbonded       *bignum.WBigInt `json:"unbonded"`
	MaybeUnbondedFrom   *time.Time      `json:"unbonded_from"`
	MaybePunishmentKind *string         `json:"punishment_kind"`
	MaybeJailedUntil    *time.Time      `json:"jailed_until"`
}

type CouncilNodeStats struct {
	Count       uint64          `json:"count"`
	TotalStaked *bignum.WBigInt `json:"total_staked"`
}

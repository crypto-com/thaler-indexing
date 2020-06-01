package viewrepo

import (
	"time"

	"github.com/crypto-com/chainindex/internal/bignum"
)

type StakingAccountViewRepo interface {
	Search(keyword string, pagination *Pagination) ([]StakingAccount, *PaginationResult, error)
}

type StakingAccount struct {
	Address                 string                     `json:"staking_account_address"`
	Nonce                   uint64                     `json:"staking_account_nonce"`
	Bonded                  *bignum.WBigInt            `json:"bonded"`
	MaybeUnbonded           *bignum.WBigInt            `json:"unbonded"`
	MaybeUnbondedFrom       *time.Time                 `json:"unbonded_until"`
	MaybeJailedUntil        *time.Time                 `json:"jailed_until"`
	MaybePunishmentKind     *string                    `json:"punishment_kind"`
	MaybeCurrentCouncilNode *StakingAccountCouncilNode `json:"current_council_node"`
}

type StakingAccountCouncilNode struct {
	MaybeId                    *uint64 `json:"id"`
	MaybeName                  *string `json:"name"`
	MaybeSecurityContact       *string `json:"security_contact"`
	MaybePubKeyType            *string `json:"pubkey_type"`
	MaybePubKey                *string `json:"pubkey"`
	MaybeAddress               *string `json:"address"`
	MaybeCreatedAtBlockHeight  *uint64 `json:"created_at_block_height"`
	MaybeLastLeftAtBlockHeight *uint64 `json:"last_left_at_block_height"`
}

type StakingAccountActivity struct {
	Type                       string               `json:"type"`
	BlockHeight                uint64               `json:"block_height"`
	BlockTime                  time.Time            `json:"block_time"`
	BlockHash                  string               `json:"block_hash"`
	MaybeTxID                  *string              `json:"txid"`
	MaybeFee                   *bignum.WBigInt      `json:"fee"`
	MaybeEventPosition         *uint64              `json:"event_position"`
	MaybeInputs                []TransactionInput   `json:"inputs"`
	MaybeOutputCount           *uint64              `json:"output_count"`
	MaybeStakingAccountAddress *string              `json:"staking_account_address"`
	MaybeStakingAccountNonce   *uint64              `json:"staking_account_nonce"`
	MaybeBonded                *bignum.WBigInt      `json:"bonded"`
	MaybeUnbonded              *bignum.WBigInt      `json:"unbonded"`
	MaybeUnbondedFrom          *time.Time           `json:"unbonded_from"`
	MaybeJoinedCouncilNode     *ActivityCouncilNode `json:"joined_council_node"`
	MaybeRewardMinted          *bignum.WBigInt      `json:"reward_minted"`
	MaybeRewardDistribution    []BlockRewardRecord  `json:"reward_distribution"`
	MaybeJailedUntil           *time.Time           `json:"jailed_until"`
	MaybePunishmentKind        *string              `json:"punishment_kind"`
	MaybeAffectedCouncilNode   *ActivityCouncilNode `json:"affected_council_node"`
}

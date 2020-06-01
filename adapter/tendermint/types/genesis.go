package types

import (
	"time"
)

type GenesisResp struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Genesis struct {
			GenesisTime     time.Time `json:"genesis_time"`
			ChainID         string    `json:"chain_id"`
			ConsensusParams struct {
				Block struct {
					MaxBytes   string `json:"max_bytes"`
					MaxGas     string `json:"max_gas"`
					TimeIotaMs string `json:"time_iota_ms"`
				} `json:"block"`
				Evidence struct {
					MaxAgeNumBlocks string `json:"max_age_num_blocks"`
					MaxAgeDuration  string `json:"max_age_duration"`
				} `json:"evidence"`
				Validator struct {
					PubKeyTypes []string `json:"pub_key_types"`
				} `json:"validator"`
			} `json:"consensus_params"`
			Validators []struct {
				Address string `json:"address"`
				PubKey  struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"pub_key"`
				Power string `json:"power"`
				Name  string `json:"name"`
			} `json:"validators"`
			AppHash  string `json:"app_hash"`
			AppState struct {
				CouncilNodes  map[string][]interface{} `json:"council_nodes"`
				Distribution  map[string][]string      `json:"distribution"`
				NetworkParams struct {
					InitialFeePolicy struct {
						Coefficient int `json:"coefficient"`
						Constant    int `json:"constant"`
					} `json:"initial_fee_policy"`
					JailingConfig struct {
						BlockSigningWindow   int `json:"block_signing_window"`
						MissedBlockThreshold int `json:"missed_block_threshold"`
					} `json:"jailing_config"`
					MaxValidators            int    `json:"max_validators"`
					RequiredCouncilNodeStake string `json:"required_council_node_stake"`
					RewardsConfig            struct {
						MonetaryExpansionCap   string `json:"monetary_expansion_cap"`
						MonetaryExpansionDecay int    `json:"monetary_expansion_decay"`
						MonetaryExpansionR0    int    `json:"monetary_expansion_r0"`
						MonetaryExpansionTau   int64  `json:"monetary_expansion_tau"`
						RewardPeriodSeconds    int    `json:"reward_period_seconds"`
					} `json:"rewards_config"`
					SlashingConfig struct {
						ByzantineSlashPercent string `json:"byzantine_slash_percent"`
						LivenessSlashPercent  string `json:"liveness_slash_percent"`
					} `json:"slashing_config"`
					UnbondingPeriod int `json:"unbonding_period"`
				} `json:"network_params"`
			} `json:"app_state"`
		} `json:"genesis"`
	} `json:"result"`
}

type RawGenesisDistributionType string

var (
	RAW_GENESIS_DISTRIBUTION_TYPE_BONDED   = "Bonded"
	RAW_GENESIS_DISTRIBUTION_TYPE_UNBONDED = "UnbondedFromGenesis"
)

type Genesis struct {
	GenesisTime time.Time
	ChainID     string
	AppHash     string
	AppState    GenesisAppState
}
type GenesisAppState struct {
	CouncilNodes []GenesisCouncilNode
	Distribution []GenesisDistribution
}
type GenesisCouncilNode struct {
	StakingAccountAddress string
	Address               string
	Name                  string
	SecurityContact       string
	PubKeyType            string
	PubKey                string
}
type GenesisDistribution struct {
	StakingAccountAddress string
	Bonded                *string
	Unbonded              *string
}

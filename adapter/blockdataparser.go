package adapter

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/crypto-com/chainindex"
	"github.com/crypto-com/chainindex/adapter/tendermint"
	tenderminttypes "github.com/crypto-com/chainindex/adapter/tendermint/types"
	"github.com/crypto-com/chainindex/adapter/txauxdecoder"
	"github.com/crypto-com/chainindex/internal/bignum"
	"github.com/crypto-com/chainindex/internal/primptr"
	"github.com/crypto-com/chainindex/usecase"
	jsoniter "github.com/json-iterator/go"
)

func ParseGenesisToBlockData(rawBlockData TendermintGenesisBlockData) *usecase.BlockData {
	blockData := usecase.BlockData{
		Block: chainindex.Block{
			Height:  uint64(1),
			Hash:    rawBlockData.Block.Hash,
			Time:    rawBlockData.Genesis.GenesisTime,
			AppHash: rawBlockData.Genesis.AppHash,
		},
		Activities:         parseGenesisActivities(rawBlockData.Genesis.AppState),
		Reward:             nil,
		CouncilNodeUpdates: nil,
	}
	return &blockData
}

type TendermintGenesisBlockData struct {
	Genesis *tenderminttypes.Genesis
	Block   *tenderminttypes.Block
}

func parseGenesisActivities(appState tenderminttypes.GenesisAppState) []chainindex.Activity {
	blockHeight := uint64(1)

	councilNodes := make(map[string]tenderminttypes.GenesisCouncilNode)
	for _, councilNode := range appState.CouncilNodes {
		councilNodes[councilNode.StakingAccountAddress] = councilNode
	}

	activities := make([]chainindex.Activity, 0, len(appState.Distribution))
	for _, entry := range appState.Distribution {
		var bonded, unbonded *big.Int
		if entry.Bonded != nil {
			bonded = bignum.MustAtoi(*entry.Bonded)
		}
		if entry.Unbonded != nil {
			unbonded = bignum.MustAtoi(*entry.Unbonded)
		}

		var councilNodeMeta *chainindex.CouncilNode
		councilNode, hasCouncilNode := councilNodes[entry.StakingAccountAddress]
		if hasCouncilNode {
			var securityContact *string
			if councilNode.SecurityContact != "" {
				securityContact = &councilNode.SecurityContact
			}
			councilNodeMeta = &chainindex.CouncilNode{
				Id:                         nil,
				Name:                       councilNode.Name,
				MaybeSecurityContact:       securityContact,
				PubKeyType:                 chainindex.PUBKEY_TYPE_ED25519,
				PubKey:                     councilNode.PubKey,
				Address:                    councilNode.Address,
				CreatedAtBlockHeight:       blockHeight,
				MaybeLastLeftAtBlockHeight: nil,
			}
		}

		clonedStakingAccountAddress := entry.StakingAccountAddress
		activities = append(activities, chainindex.Activity{
			BlockHeight:                blockHeight,
			Type:                       chainindex.ACTIVITY_GENESIS,
			MaybeTxID:                  nil,
			MaybeEventPosition:         nil,
			MaybeFee:                   nil,
			MaybeTxInputs:              nil,
			MaybeOutputCount:           nil,
			MaybeStakingAccountAddress: &clonedStakingAccountAddress,
			MaybeBonded:                bonded,
			MaybeUnbonded:              unbonded,
			MaybeCouncilNodeMeta:       councilNodeMeta,
			MaybeAffectedCouncilNode:   nil,
			MaybeJailedUntil:           nil,
			MaybePunishmentKind:        nil,
		})
	}

	return activities
}

func ParseBlockToBlockData(rawBlockData TendermintBlockData) *usecase.BlockData {
	var blockData usecase.BlockData

	blockData.Block = chainindex.Block{
		Height:  rawBlockData.Block.Height,
		Hash:    rawBlockData.Block.Hash,
		Time:    rawBlockData.Block.Time,
		AppHash: rawBlockData.Block.AppHash,
	}

	blockData.Signatures = parseSignatures(
		rawBlockData.Block.Height,
		rawBlockData.Block.PropserAddress,
		rawBlockData.Block.Signatures,
	)

	activities := make([]chainindex.Activity, 0)

	if rawBlockData.BlockResults.TxsEvents != nil {
		activities = parseTransactions(
			rawBlockData.Block.Height,
			rawBlockData.BlockResults.TxsEvents,
			rawBlockData.Block.Txs,
		)
	}

	eventActivities, reward := parseBeginBlockEvents(
		rawBlockData.Block.Height,
		rawBlockData.BlockResults.BeginBlockEvents,
	)
	if eventActivities != nil {
		activities = append(activities, eventActivities...)
	}
	if reward != nil {
		blockData.Reward = reward
	}

	if len(activities) != 0 {
		blockData.Activities = activities
	}
	blockData.CouncilNodeUpdates = parseValidatorUpdates(rawBlockData.BlockResults.ValidatorUpdates)

	return &blockData
}

func parseSignatures(blockHeight uint64, proposerAddress string, signatures []tenderminttypes.BlockSignature) []chainindex.BlockSignature {
	if signatures == nil {
		return nil
	}

	parsedSignatures := make([]chainindex.BlockSignature, 0, len(signatures))
	for _, signature := range signatures {
		parsedSignatures = append(parsedSignatures, chainindex.BlockSignature{
			BlockHeight:        blockHeight,
			CouncilNodeAddress: signature.ValidatorAddress,
			Signature:          signature.Signature,
			IsProposer:         signature.ValidatorAddress == proposerAddress,
		})
	}

	return parsedSignatures
}

type TendermintBlockData struct {
	Block        *tenderminttypes.Block
	BlockResults *tenderminttypes.BlockResults
}

func parseTransactions(blockHeight uint64, txsEvents [][]tenderminttypes.BlockResultsEvent, rawTxs []string) []chainindex.Activity {
	activities := make([]chainindex.Activity, 0, len(txsEvents))
	for i, txEvents := range txsEvents {
		var activity chainindex.Activity
		activity.BlockHeight = blockHeight

		decodedTx := txauxdecoder.MustDecodeBase64(rawTxs[i])
		switch decodedTx.TxType {
		case "Transfer":
			activity.Type = chainindex.ACTIVITY_TRANSFER
			txInputs := make([]chainindex.TxInput, 0)
			for _, txInput := range decodedTx.Inputs {
				txInputs = append(txInputs, chainindex.TxInput{
					TxId:  txInput.ID,
					Index: txInput.Index,
				})
			}
			activity.MaybeTxInputs = txInputs
			activity.MaybeOutputCount = decodedTx.OutputCount
		case "Deposit":
			activity.Type = chainindex.ACTIVITY_DEPOSIT
		case "Unbond":
			activity.Type = chainindex.ACTIVITY_UNBOND
		case "Withdraw":
			activity.Type = chainindex.ACTIVITY_WITHDRAW
			activity.MaybeOutputCount = decodedTx.OutputCount
		case "NodeJoin":
			activity.Type = chainindex.ACTIVITY_NODEJOIN
		case "Unjail":
			activity.Type = chainindex.ACTIVITY_UNJAIL
		}

		for _, event := range txEvents {
			switch event.Type {
			case "valid_txs":
				for _, attribute := range event.Attributes {
					value, err := base64DecodeString(attribute.Value)
					if err != nil {
						panic(fmt.Sprintf("error base64 decoding valid_txs event: %v", err))
					}

					switch attribute.Key {
					case ATTRIBUTE_TXID:
						activity.MaybeTxID = &value
					case ATTRIBUTE_FEE:
						croFee := bignum.MustAtof(value)
						activity.MaybeFee = chainindex.MustCROToCoin(croFee)
					}
				}
			case "staking_change":
				for _, attribute := range event.Attributes {
					value, err := base64DecodeString(attribute.Value)
					if err != nil {
						panic(fmt.Sprintf("error base64 decoding staking_change event: %v", err))
					}

					switch attribute.Key {
					case ATTRIBUTE_SATKING_ADDRESS:
						activity.MaybeStakingAccountAddress = &value
					case ATTRIBUTE_STAKING_DIFF:
						var stakingDiffs StakingDiffs
						err := jsoniter.Unmarshal([]byte(value), &stakingDiffs)
						if err != nil {
							panic(fmt.Sprintf("error deserializing staking_diff: %v", err))
						}

						for _, kvPair := range stakingDiffs {
							switch kvPair.Key {
							case "Bonded":
								activity.MaybeBonded, err = bignum.Atoi(kvPair.Value.(string))
								if err != nil {
									panic("error converting staking_diff bonded amount to big.Int")
								}
							case "Unbonded":
								activity.MaybeUnbonded, err = bignum.Atoi(kvPair.Value.(string))
								if err != nil {
									panic("error converting staking_diff unbonded amount to big.Int")
								}
							case "UnbondedFrom":
								unix, _ := kvPair.Value.(float64)
								activity.MaybeUnbondedFrom = primptr.Time(time.Unix(int64(unix), 0).UTC())
							case "CouncilNode":
								councilNode, _ := kvPair.Value.(map[string]interface{})

								var securityContact *string
								if councilNode["security_contact"] != nil {
									*securityContact, _ = councilNode["security_contact"].(string)
								}
								consensusPubKey, _ := councilNode["consensus_pubkey"].(map[string]interface{})
								pubKey, _ := consensusPubKey["value"].(string)
								activity.MaybeCouncilNodeMeta = &chainindex.CouncilNode{
									Id:                         nil,
									Name:                       councilNode["name"].(string),
									MaybeSecurityContact:       securityContact,
									PubKeyType:                 chainindex.PUBKEY_TYPE_ED25519,
									PubKey:                     pubKey,
									Address:                    tendermint.AddressFromPubKey(pubKey),
									CreatedAtBlockHeight:       blockHeight,
									MaybeLastLeftAtBlockHeight: nil,
								}
							}
						}
					}
				}
			}
		}

		activities = append(activities, activity)
	}

	return activities
}

func parseBeginBlockEvents(blockHeight uint64, beginBlockEvents []tenderminttypes.BlockResultsEvent) ([]chainindex.Activity, *chainindex.BlockReward) {
	activities := make([]chainindex.Activity, 0)
	var reward *chainindex.BlockReward

	for position, event := range beginBlockEvents {
		switch event.Type {
		case "reward":
			for _, attribute := range event.Attributes {
				value, err := base64DecodeString(attribute.Value)
				if err != nil {
					panic(fmt.Sprintf("error base64 decoding reward event value: %v", err))
				}

				switch attribute.Key {
				case ATTRIBUTE_MINTED:
					// FIXME: v0.5 reward event minted amount has unnecessary double quotes
					minted, err := bignum.Atoi(strings.Trim(value, "\""))
					if err != nil {
						panic("error converting staking_diff bonded amount to big.Int")
					}

					reward = new(chainindex.BlockReward)
					reward.BlockHeight = blockHeight
					reward.Minted = minted
				}
			}
		case "staking_change":
			var activity chainindex.Activity
			activity.BlockHeight = blockHeight
			activity.MaybeEventPosition = primptr.Uint32(uint32(position))

			for _, attribute := range event.Attributes {
				value, err := base64DecodeString(attribute.Value)
				if err != nil {
					panic(fmt.Sprintf("error base64 decoding staking_change event value: %v", err))
				}

				switch attribute.Key {
				case ATTRIBUTE_SATKING_ADDRESS:
					activity.MaybeStakingAccountAddress = &value
				case ATTRIBUTE_STAKING_OPTYPE:
					switch value {
					case "reward":
						activity.Type = chainindex.ACTIVITY_REWARD
					case "slash":
						activity.Type = chainindex.ACTIVITY_SLASH
					case "jail":
						activity.Type = chainindex.ACTIVITY_JAIL
					}
				case ATTRIBUTE_STAKING_OPREASON:
					switch value {
					case "NonLive":
						activity.MaybePunishmentKind = primptr.Uint8(chainindex.PUNISHMENT_KIND_NON_LIVE)
					case "ByzantineFault":
						activity.MaybePunishmentKind = primptr.Uint8(chainindex.PUNISHMENT_KIND_BYZANTINE_FAULT)
					}
				case ATTRIBUTE_STAKING_DIFF:
					var stakingDiffs StakingDiffs
					err := jsoniter.Unmarshal([]byte(value), &stakingDiffs)
					if err != nil {
						panic(fmt.Sprintf("error deserializing staking_diff: %v", err))
					}

					for _, kvPair := range stakingDiffs {
						switch kvPair.Key {
						case "Bonded":
							activity.MaybeBonded, err = bignum.Atoi(kvPair.Value.(string))
							if err != nil {
								panic("error converting staking_diff bonded amount to big.Int")
							}
						case "Unbonded":
							activity.MaybeUnbonded, err = bignum.Atoi(kvPair.Value.(string))
							if err != nil {
								panic("error converting staking_diff unbonded amount to big.Int")
							}
						case "JailedUntil":
							unix, _ := kvPair.Value.(float64)
							activity.MaybeJailedUntil = primptr.Time(time.Unix(int64(unix), 0).UTC())
						}

					}
				}
			}

			activities = append(activities, activity)
		}
	}

	return activities, reward
}

func parseValidatorUpdates(validatorUpdates []tenderminttypes.BlockResultsValidator) []chainindex.CouncilNodeUpdate {
	if validatorUpdates == nil {
		return nil
	}

	councilNodeUpdates := make([]chainindex.CouncilNodeUpdate, 0)
	for _, validatorUpdate := range validatorUpdates {
		// nolint:gosec,scopelint
		if isValidatorKicked(&validatorUpdate) {
			councilNodeUpdates = append(councilNodeUpdates, chainindex.CouncilNodeUpdate{
				Address: validatorUpdate.PubKey.Address,
				Type:    chainindex.COUNCIL_NODE_UPDATE_TYPE_LEFT,
			})
		}
	}

	return councilNodeUpdates
}

func isValidatorKicked(validatorUpdate *tenderminttypes.BlockResultsValidator) bool {
	return validatorUpdate.Power == nil
}

var (
	ATTRIBUTE_FEE              = base64.StdEncoding.EncodeToString([]byte("fee"))
	ATTRIBUTE_TXID             = base64.StdEncoding.EncodeToString([]byte("txid"))
	ATTRIBUTE_SATKING_ADDRESS  = base64.StdEncoding.EncodeToString([]byte("staking_address"))
	ATTRIBUTE_STAKING_OPTYPE   = base64.StdEncoding.EncodeToString([]byte("staking_optype"))
	ATTRIBUTE_STAKING_OPREASON = base64.StdEncoding.EncodeToString([]byte("staking_opreason"))
	ATTRIBUTE_STAKING_DIFF     = base64.StdEncoding.EncodeToString([]byte("staking_diff"))
	ATTRIBUTE_MINTED           = base64.StdEncoding.EncodeToString([]byte("minted"))
)

type StakingDiffs = []StakingDiffKVPair
type StakingDiffKVPair struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func base64DecodeString(s string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

type StakingDiffCouncilNode struct {
	Name             string                      `json:"name"`
	SecurityContact  string                      `json:"security_contact"`
	ConsensusPubKey  StakingDiffConsensusPubKey  `json:"consensus_pubkey"`
	ConfidentialInit StakingDiffConfidentialInit `json:"confidential_init"`
}

type StakingDiffConsensusPubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type StakingDiffConfidentialInit struct {
	Cert string `json:"cert"`
}

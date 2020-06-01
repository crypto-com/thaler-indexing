package tendermint_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	tendermintadapter "github.com/crypto-com/chainindex/adapter/tendermint"
	"github.com/crypto-com/chainindex/adapter/tendermint/types"
	"github.com/crypto-com/chainindex/infrastructure/tendermint"
	"github.com/crypto-com/chainindex/internal/primptr"
)

var _ = Describe("HTTPClient", func() {
	var server *ghttp.Server

	BeforeEach(func() {
		server = ghttp.NewServer()
	})

	It("should implement Client", func() {
		var _ tendermintadapter.Client = tendermint.NewHTTPClient("http://localhost:26657")
	})

	Describe("Genesis", func() {
		It("should return parsed genesis from Tendermint", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/genesis"),
					ghttp.RespondWith(http.StatusOK, GENESIS_JSON),
				),
			)

			client := tendermint.NewHTTPClient(server.URL())

			genesis, err := client.Genesis()
			Expect(err).To(BeNil())

			genesisTime, _ := time.Parse("2006-01-02T15:04:05.000000Z", "2020-05-01T12:09:01.568951Z")
			Expect(*genesis).To(Equal(types.Genesis{
				GenesisTime: genesisTime,
				ChainID:     "testnet-thaler-crypto-com-chain-42",
				AppHash:     "F62DDB49D7EB8ED0883C735A0FB7DE7F2A3FA322FCD2AA832F452A62B38607D5",
				AppState: types.GenesisAppState{
					CouncilNodes: []types.GenesisCouncilNode{
						{
							StakingAccountAddress: "0x6fc1e3124a7ed07f3710378b68f7046c7300179d",
							Address:               "7570B2D23A4C7B638BEFE02EB4FC7927BFDED6B7",
							Name:                  "canadacentral_validator_1",
							SecurityContact:       "security@crypto.com",
							PubKeyType:            "tendermint/PubKeyEd25519",
							PubKey:                "QMegiWt9+5K1b1ZVd7zOJZxhTnbAtWzvGhViiElAlaw=",
						},
						{
							StakingAccountAddress: "0xb8c6886da09e12db8aebfc8108c67ce2ba086ac6",
							Address:               "D527DAECDE0501CF2E785A8DC0D9F4A64760F0BB",
							Name:                  "uksouth_validator_1",
							SecurityContact:       "security@crypto.com",
							PubKeyType:            "tendermint/PubKeyEd25519",
							PubKey:                "tDLheZJwsA8oYEwarR6/X+zAmNKMLHTVkh/fvcLqcwA=",
						},
						{
							StakingAccountAddress: "0x6dbd5b8fe0dad494465aa7574defba711c184102",
							Address:               "FA7B721B5704DF98EF3ECD3796DDEF6AA2A80257",
							Name:                  "eastus_validator_1",
							SecurityContact:       "security@crypto.com",
							PubKeyType:            "tendermint/PubKeyEd25519",
							PubKey:                "/SvfTeO4Du4oR/VYTjm7IgObc14zzddEAyFb4nU8E3Q=",
						},
					},
					Distribution: []types.GenesisDistribution{
						{
							StakingAccountAddress: "0x4ae85b35597fcb61c6c47b1fe0bdd7eed8421cdd",
							Bonded:                primptr.String("6000000000000000"),
							Unbonded:              nil,
						},
						{
							StakingAccountAddress: "0x4b75f275dde0a8c8e70fb84243adc97a3afb78f2",
							Bonded:                nil,
							Unbonded:              primptr.String("7946000000000000000"),
						},
						{
							StakingAccountAddress: "0x4fd8162521f2e628adced7c1baa39384a08b4a3d",
							Bonded:                primptr.String("6000000000000000"),
							Unbonded:              nil,
						},
						{
							StakingAccountAddress: "0x6c2be7846219eab3086a66f873558b73d8f4a0d4",
							Bonded:                primptr.String("6000000000000000"),
							Unbonded:              nil,
						},
						{
							StakingAccountAddress: "0x6dbd5b8fe0dad494465aa7574defba711c184102",
							Bonded:                primptr.String("6000000000000000"),
							Unbonded:              nil,
						},
						{
							StakingAccountAddress: "0x6fc1e3124a7ed07f3710378b68f7046c7300179d",
							Bonded:                primptr.String("6000000000000000"),
							Unbonded:              nil,
						},
						{
							StakingAccountAddress: "0x9baa6de71cbc6274275eece4b1be15f545897f37",
							Bonded:                primptr.String("6000000000000000"),
							Unbonded:              nil,
						},
						{
							StakingAccountAddress: "0xa9528abb92709370600d2cef41f1677374278337",
							Bonded:                primptr.String("6000000000000000"),
							Unbonded:              nil,
						},
						{
							StakingAccountAddress: "0xb328a39002ede64c33bb60f1dc43f5df9eb47043",
							Bonded:                primptr.String("6000000000000000"),
							Unbonded:              nil,
						},
						{
							StakingAccountAddress: "0xb8c6886da09e12db8aebfc8108c67ce2ba086ac6",
							Bonded:                primptr.String("6000000000000000"),
							Unbonded:              nil,
						},
					},
				},
			}))
		})
	})

	Describe("LatestBlockHeight", func() {
		It("should return block height", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/block_results"),
					ghttp.RespondWith(http.StatusOK, `{
	"jsonrpc": "2.0",
	"id": -1,
	"result": {
		"height": "27012",
		"txs_results": null,
		"begin_block_events": null,
		"end_block_events": null,
		"validator_updates": null,
		"consensus_param_updates": null
	}
}`),
				),
			)

			client := tendermint.NewHTTPClient(server.URL())

			blockHeight, err := client.LatestBlockHeight()
			Expect(err).To(BeNil())
			Expect(blockHeight).To(Equal(uint64(27012)))
		})
	})

	Describe("BlockResults", func() {
		It("should return nil Events when there are no transactions nor events", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/block_results", "height=1"),
					ghttp.RespondWith(http.StatusOK, BLOCK_RESULTS_EMPTY_EVENTS_JSON),
				),
			)

			client := tendermint.NewHTTPClient(server.URL())

			anyBlockHeight := uint64(1)
			blockResults, err := client.BlockResults(anyBlockHeight)
			Expect(err).To(BeNil())
			Expect(*blockResults).To(Equal(types.BlockResults{
				Height:           anyBlockHeight,
				TxsEvents:        nil,
				BeginBlockEvents: nil,
				ValidatorUpdates: nil,
			}))
		})

		It("should return parsed block results", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/block_results", "height=3813"),
					ghttp.RespondWith(http.StatusOK, BLOCK_RESULTS_JSON),
				),
			)

			client := tendermint.NewHTTPClient(server.URL())

			anyBlockHeight := uint64(3813)
			blockResults, err := client.BlockResults(anyBlockHeight)
			Expect(err).To(BeNil())
			Expect(*blockResults).To(Equal(types.BlockResults{
				Height: anyBlockHeight,
				TxsEvents: [][]types.BlockResultsEvent{
					{
						{
							Type: "valid_txs",
							Attributes: []types.BlockResultsEventAttribute{
								{
									Key:   "ZmVl",
									Value: "MC4wMDAwMDQ2OQ==",
								},
								{
									Key:   "dHhpZA==",
									Value: "YmEyMjMwYTA0OTIyZDNmMDFkNDE5OTljNmFkYmUwNmZjOGE5ODQxM2IyZDU1YWM3ZjlhYzMwZmVjMzlmYzdiMg==",
								},
							},
						},
					},
					{
						{
							Type: "valid_txs",
							Attributes: []types.BlockResultsEventAttribute{
								{
									Key:   "ZmVl",
									Value: "MC4wMDAwMDQ2OQ==",
								},
								{
									Key:   "dHhpZA==",
									Value: "N2I1YWMzNmY2M2ZmYjZlZTEzMjg5ZDQ5ZDljZmRhY2UzZjhkMjM0NDY5ZWIwNTc1NWY1ZjVhYzgxYjNhNDVhNg==",
								},
							},
						},
					},
				},
				BeginBlockEvents: []types.BlockResultsEvent{
					{
						Type: "staking_change",
						Attributes: []types.BlockResultsEventAttribute{
							{
								Key:   "c3Rha2luZ19hZGRyZXNz",
								Value: "MHg2ZGJkNWI4ZmUwZGFkNDk0NDY1YWE3NTc0ZGVmYmE3MTFjMTg0MTAy",
							},
							{
								Key:   "c3Rha2luZ19vcHR5cGU=",
								Value: "cmV3YXJk",
							},
							{
								Key:   "c3Rha2luZ19kaWZm",
								Value: "W3sia2V5IjoiQm9uZGVkIiwidmFsdWUiOiI0ODU4OTkwMjAzMzMzIn1d",
							},
						},
					},
					{
						Type: "staking_change",
						Attributes: []types.BlockResultsEventAttribute{
							{
								Key:   "c3Rha2luZ19hZGRyZXNz",
								Value: "MHg2ZmMxZTMxMjRhN2VkMDdmMzcxMDM3OGI2OGY3MDQ2YzczMDAxNzlk",
							},
							{
								Key:   "c3Rha2luZ19vcHR5cGU=",
								Value: "cmV3YXJk",
							},
							{
								Key:   "c3Rha2luZ19kaWZm",
								Value: "W3sia2V5IjoiQm9uZGVkIiwidmFsdWUiOiI0ODU4OTkwMjAzMzMzIn1d",
							},
						},
					},
					{
						Type: "staking_change",
						Attributes: []types.BlockResultsEventAttribute{
							{
								Key:   "c3Rha2luZ19hZGRyZXNz",
								Value: "MHhiOGM2ODg2ZGEwOWUxMmRiOGFlYmZjODEwOGM2N2NlMmJhMDg2YWM2",
							},
							{
								Key:   "c3Rha2luZ19vcHR5cGU=",
								Value: "cmV3YXJk",
							},
							{
								Key:   "c3Rha2luZ19kaWZm",
								Value: "W3sia2V5IjoiQm9uZGVkIiwidmFsdWUiOiI0ODU4OTkwMjAzMzMzIn1d",
							},
						},
					},
					{
						Type: "reward",
						Attributes: []types.BlockResultsEventAttribute{
							{
								Key:   "bWludGVk",
								Value: "IjE0NTc2OTcwNjEwMDAwIg==",
							},
						},
					},
				},
				ValidatorUpdates: []types.BlockResultsValidator{
					{
						PubKey: types.BlockResultsValidatorPubKey{
							Type:    "ed25519",
							PubKey:  "rXhu7xhqYBtJftVLKxvKN0XnpyOzxFnUEfAhD1dEF/8=",
							Address: "34C725CABA703269B3F1D1A907A84DE5FEE96469",
						},
						Power: primptr.String("60000000"),
					},
					{
						PubKey: types.BlockResultsValidatorPubKey{
							Type:    "ed25519",
							PubKey:  "tDLheZJwsA8oYEwarR6/X+zAmNKMLHTVkh/fvcLqcwA=",
							Address: "D527DAECDE0501CF2E785A8DC0D9F4A64760F0BB",
						},
						Power: primptr.String("80000000"),
					},
				},
			}))
		})
	})

	Describe("Block", func() {
		It("should return nil Txs when there are no transactions", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/block", "height=1"),
					ghttp.RespondWith(http.StatusOK, BLOCK_EMPTY_TX_RESULT),
				),
			)

			client := tendermint.NewHTTPClient(server.URL())

			anyBlockHeight := uint64(1)
			block, err := client.Block(anyBlockHeight)
			Expect(err).To(BeNil())
			blockTime, _ := time.Parse("2006-01-02T15:04:05.000000Z", "2020-05-01T12:09:01.568951Z")
			Expect(*block).To(Equal(types.Block{
				Height:         anyBlockHeight,
				Hash:           "BD3D9499EF527035BAE14E7F4932BD1778C1E90141A7A76A1FE67B7ECCC2663E",
				Time:           blockTime,
				AppHash:        "F62DDB49D7EB8ED0883C735A0FB7DE7F2A3FA322FCD2AA832F452A62B38607D5",
				PropserAddress: "D527DAECDE0501CF2E785A8DC0D9F4A64760F0BB",
				Txs:            nil,
				Signatures:     nil,
			}))
		})

		It("should return parsed block", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/block", "height=3510"),
					ghttp.RespondWith(http.StatusOK, BLOCK_JSON),
				),
			)

			client := tendermint.NewHTTPClient(server.URL())

			anyBlockHeight := uint64(3510)
			block, err := client.Block(anyBlockHeight)
			Expect(err).To(BeNil())
			blockTime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "2020-05-07T10:00:25.582998694Z")
			Expect(*block).To(Equal(types.Block{
				Height:         anyBlockHeight,
				Hash:           "2CD22EA622D190B9ABCAB797E2F60F6F4FCFC19CC0A67642E5E7856CEAD78163",
				Time:           blockTime,
				AppHash:        "4AC923568B9DE0AA67A04FC60CC4EA10ECC69D636AACFC519911BB448FEED222",
				PropserAddress: "7570B2D23A4C7B638BEFE02EB4FC7927BFDED6B7",
				Txs: []string{
					"AAAELi1Eyx7NNI3yHk9foI+GtqHLPr5/dflC/DyOdCDniaIBAAIAAAAAAAAAAADNM4VutAy2aB2/DoBlBJC/VJxu+WEpaq4jt0BMOpUgxNEiL4lCUGhvlc3FeFCZ99Fvp9Xz9BOCqZ11XxENUFP4moIHvLZMsrohp3ElmtrJFQpZbeS+Ik0FCbT0erhj/wLMptWgxVY7sD/KggQKJjygJb+jdlre2kZ0f8C+nuYW7Q7k7zR1l/D3zFRsLZqbwkBHGchFt3rc3GKiKN2KhNnUoOKG6EvCDQ/HLNwMg0VA3xAW108NTwMnq3prMfjFcZ8QGOBGHWVjDb3GOL/3vngoNr0Jq98oRVEol+sbtvx3s2eO5mGbApiaDK1/oIh9WBaXxEwtXeSPGrw66/8sa9UoaJWBnNOSTesKls/SP1ZQGT1zlrxAZs9MFwjYDbGgukbWletVqgfBcYt3mrvj/G8LtBUH90A/vrJkbvZ6jvTQXqVVHk1jXuU=",
					"AAAEBuKX03ws+y30Jc0oeZTup1aLw0mEB9W8J+LDAaSaOQUBAAIAAAAAAAAAAACILvI8Gcck+IwmflRlBHvHQzIQdFoMSrOXl9a/ZgETkCyFWpwZl8tmJ0mhjGn2bPB7O/srnVyUKxV58j0C1BKFskIFJIr+lmr75vuM9m2DnyL1g2+TDu9nlZWvl+JHR49h0ToA36f0XshYVjyWG6r6O322w4HHP1ApOvuMpq458Hfy3HHnlgaKbcFrwUEbTDrgRslwpgR24jRvVY38IhkEIEHDgXVl8rsXsc04/UsaLQyYlaDi0zqks3qIuPguTQkQ5++/qDvpK6OromSAdCVHhnfdzf6UKsFFlKv9jMcUSeyg5Obx9u/76C2Gfww/iJ0dTZ/Ob0gVQGp+tWyyfa+2OMac++nhGckwru+mw5RoosoFiLMvmeMQrcHR3mcPIaBBx9p4dndcp0JIUEhzFbo29cm1Z8vsbd5nQHbNOxSN5/daIoEJUcQ=",
				},
				Signatures: []types.BlockSignature{
					{
						ValidatorAddress: "34C725CABA703269B3F1D1A907A84DE5FEE96469",
						Signature:        "W6wNdtDCsUbdodju57/2BlkmUjo6U0PH+Cf19u4RSlaBYS7svNpZOgdQqmXJnUInRjGJp8opE7a9FnnHe3oTAw==",
					},
					{
						ValidatorAddress: "7570B2D23A4C7B638BEFE02EB4FC7927BFDED6B7",
						Signature:        "W2Pkxi/AqAklmEMjsd3P9yd//1ZGxBqRMaHebcGYhlVZxcbZ02dzwgxD7c/BOOMh+kJGYPfuYNiLHD3Kts+pDA==",
					},
					{
						ValidatorAddress: "D527DAECDE0501CF2E785A8DC0D9F4A64760F0BB",
						Signature:        "zCjxVQLbEhBdBBm1VHyLP+4aFH81ke26pC0e+g+pvAzLBvWxrmzdgh347MmOVHWiW6lS9nb8xs+6bkdKRMx5Dg==",
					},
					{
						ValidatorAddress: "FA7B721B5704DF98EF3ECD3796DDEF6AA2A80257",
						Signature:        "TuflcSDZPgVTM618J9JF/tlMFwM8Z/eWrHPjixzeWukIlkFHsNMprRRPnHUKZlu+yDdSwgj6eJ0PrqRm7y/eDw==",
					},
				},
			}))
		})
	})
})

const (
	GENESIS_JSON = `
{
	"jsonrpc": "2.0",
	"id": -1,
	"result": {
		"genesis": {
			"genesis_time": "2020-05-01T12:09:01.568951Z",
			"chain_id": "testnet-thaler-crypto-com-chain-42",
			"consensus_params": {
				"block": {
					"max_bytes": "22020096",
					"max_gas": "-1",
					"time_iota_ms": "1000"
				},
				"evidence": {
					"max_age_num_blocks": "200",
					"max_age_duration": "5400000000000"
				},
				"validator": {
					"pub_key_types": [
						"ed25519"
					]
				}
			},
			"validators": [
				{
					"address": "FA7B721B5704DF98EF3ECD3796DDEF6AA2A80257",
					"pub_key": {
						"type": "tendermint/PubKeyEd25519",
						"value": "/SvfTeO4Du4oR/VYTjm7IgObc14zzddEAyFb4nU8E3Q="
					},
					"power": "60000000",
					"name": "eastus_validator_1"
				},
				{
					"address": "7570B2D23A4C7B638BEFE02EB4FC7927BFDED6B7",
					"pub_key": {
						"type": "tendermint/PubKeyEd25519",
						"value": "QMegiWt9+5K1b1ZVd7zOJZxhTnbAtWzvGhViiElAlaw="
					},
					"power": "60000000",
					"name": "canadacentral_validator_1"
				},
				{
					"address": "D527DAECDE0501CF2E785A8DC0D9F4A64760F0BB",
					"pub_key": {
						"type": "tendermint/PubKeyEd25519",
						"value": "tDLheZJwsA8oYEwarR6/X+zAmNKMLHTVkh/fvcLqcwA="
					},
					"power": "60000000",
					"name": "uksouth_validator_1"
				}
			],
			"app_hash": "F62DDB49D7EB8ED0883C735A0FB7DE7F2A3FA322FCD2AA832F452A62B38607D5",
			"app_state": {
				"council_nodes": {
					"0x6dbd5b8fe0dad494465aa7574defba711c184102": [
						"eastus_validator_1",
						"security@crypto.com",
						{
							"type": "tendermint/PubKeyEd25519",
							"value": "/SvfTeO4Du4oR/VYTjm7IgObc14zzddEAyFb4nU8E3Q="
						},
						{
							"cert": "ABCD"
						}
					],
					"0x6fc1e3124a7ed07f3710378b68f7046c7300179d": [
						"canadacentral_validator_1",
						"security@crypto.com",
						{
							"type": "tendermint/PubKeyEd25519",
							"value": "QMegiWt9+5K1b1ZVd7zOJZxhTnbAtWzvGhViiElAlaw="
						},
						{
							"cert": "ABCD"
						}
					],
					"0xb8c6886da09e12db8aebfc8108c67ce2ba086ac6": [
						"uksouth_validator_1",
						"security@crypto.com",
						{
							"type": "tendermint/PubKeyEd25519",
							"value": "tDLheZJwsA8oYEwarR6/X+zAmNKMLHTVkh/fvcLqcwA="
						},
						{
							"cert": "ABCD"
						}
					]
				},
				"distribution": {
					"0x4ae85b35597fcb61c6c47b1fe0bdd7eed8421cdd": [
						"Bonded",
						"6000000000000000"
					],
					"0x4b75f275dde0a8c8e70fb84243adc97a3afb78f2": [
						"UnbondedFromGenesis",
						"7946000000000000000"
					],
					"0x4fd8162521f2e628adced7c1baa39384a08b4a3d": [
						"Bonded",
						"6000000000000000"
					],
					"0x6c2be7846219eab3086a66f873558b73d8f4a0d4": [
						"Bonded",
						"6000000000000000"
					],
					"0x6dbd5b8fe0dad494465aa7574defba711c184102": [
						"Bonded",
						"6000000000000000"
					],
					"0x6fc1e3124a7ed07f3710378b68f7046c7300179d": [
						"Bonded",
						"6000000000000000"
					],
					"0x9baa6de71cbc6274275eece4b1be15f545897f37": [
						"Bonded",
						"6000000000000000"
					],
					"0xa9528abb92709370600d2cef41f1677374278337": [
						"Bonded",
						"6000000000000000"
					],
					"0xb328a39002ede64c33bb60f1dc43f5df9eb47043": [
						"Bonded",
						"6000000000000000"
					],
					"0xb8c6886da09e12db8aebfc8108c67ce2ba086ac6": [
						"Bonded",
						"6000000000000000"
					]
				},
				"network_params": {
					"initial_fee_policy": {
						"coefficient": 1250,
						"constant": 1100
					},
					"jailing_config": {
						"block_signing_window": 720,
						"missed_block_threshold": 360
					},
					"max_validators": 50,
					"required_council_node_stake": "5000000000000000",
					"rewards_config": {
						"monetary_expansion_cap": "2000000000000000000",
						"monetary_expansion_decay": 999860,
						"monetary_expansion_r0": 350,
						"monetary_expansion_tau": 999999999999999999,
						"reward_period_seconds": 86400
					},
					"slashing_config": {
						"byzantine_slash_percent": "0.200",
						"liveness_slash_percent": "0.100"
					},
					"unbonding_period": 5400
				}
			}
		}
	}
}
`
	BLOCK_RESULTS_JSON = `
{
	"jsonrpc": "2.0",
	"id": -1,
	"result": {
		"height": "3813",
		"txs_results": [
			{
				"code": 0,
				"data": null,
				"log": "",
				"info": "",
				"gasWanted": "0",
				"gasUsed": "0",
				"events": [
					{
						"type": "valid_txs",
						"attributes": [
						{
							"key": "ZmVl",
							"value": "MC4wMDAwMDQ2OQ=="
						},
						{
							"key": "dHhpZA==",
							"value": "YmEyMjMwYTA0OTIyZDNmMDFkNDE5OTljNmFkYmUwNmZjOGE5ODQxM2IyZDU1YWM3ZjlhYzMwZmVjMzlmYzdiMg=="
						}
						]
					}
				],
				"codespace": ""
			},
			{
				"code": 0,
				"data": null,
				"log": "",
				"info": "",
				"gasWanted": "0",
				"gasUsed": "0",
				"events": [
					{
						"type": "valid_txs",
						"attributes": [
						{
							"key": "ZmVl",
							"value": "MC4wMDAwMDQ2OQ=="
						},
						{
							"key": "dHhpZA==",
							"value": "N2I1YWMzNmY2M2ZmYjZlZTEzMjg5ZDQ5ZDljZmRhY2UzZjhkMjM0NDY5ZWIwNTc1NWY1ZjVhYzgxYjNhNDVhNg=="
						}
						]
					}
				],
				"codespace": ""
			}
		],
		"begin_block_events": [
			{
				"type": "staking_change",
				"attributes": [
					{
						"key": "c3Rha2luZ19hZGRyZXNz",
						"value": "MHg2ZGJkNWI4ZmUwZGFkNDk0NDY1YWE3NTc0ZGVmYmE3MTFjMTg0MTAy"
					},
					{
						"key": "c3Rha2luZ19vcHR5cGU=",
						"value": "cmV3YXJk"
					},
					{
						"key": "c3Rha2luZ19kaWZm",
						"value": "W3sia2V5IjoiQm9uZGVkIiwidmFsdWUiOiI0ODU4OTkwMjAzMzMzIn1d"
					}
				]
			},
			{
				"type": "staking_change",
				"attributes": [
					{
						"key": "c3Rha2luZ19hZGRyZXNz",
						"value": "MHg2ZmMxZTMxMjRhN2VkMDdmMzcxMDM3OGI2OGY3MDQ2YzczMDAxNzlk"
					},
					{
						"key": "c3Rha2luZ19vcHR5cGU=",
						"value": "cmV3YXJk"
					},
					{
						"key": "c3Rha2luZ19kaWZm",
						"value": "W3sia2V5IjoiQm9uZGVkIiwidmFsdWUiOiI0ODU4OTkwMjAzMzMzIn1d"
					}
				]
			},
			{
				"type": "staking_change",
				"attributes": [
					{
						"key": "c3Rha2luZ19hZGRyZXNz",
						"value": "MHhiOGM2ODg2ZGEwOWUxMmRiOGFlYmZjODEwOGM2N2NlMmJhMDg2YWM2"
					},
					{
						"key": "c3Rha2luZ19vcHR5cGU=",
						"value": "cmV3YXJk"
					},
					{
						"key": "c3Rha2luZ19kaWZm",
						"value": "W3sia2V5IjoiQm9uZGVkIiwidmFsdWUiOiI0ODU4OTkwMjAzMzMzIn1d"
					}
				]
			},
			{
				"type": "reward",
				"attributes": [
					{
						"key": "bWludGVk",
						"value": "IjE0NTc2OTcwNjEwMDAwIg=="
					}
				]
			}
		],
		"end_block_events": [
			{
				"type": "block_filter",
				"attributes": [
					{
						"key": "ZXRoYmxvb20=",
						"value": "AAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
					}
				]
			}
		],
		"validator_updates": [
			{
				"pub_key": {
					"type": "ed25519",
					"data": "rXhu7xhqYBtJftVLKxvKN0XnpyOzxFnUEfAhD1dEF/8="
				},
				"power": "60000000"
			},
			{
				"pub_key": {
					"type": "ed25519",
					"data": "tDLheZJwsA8oYEwarR6/X+zAmNKMLHTVkh/fvcLqcwA="
				},
				"power": "80000000"
			}
		],
		"consensus_param_updates": null
	}
}`
	BLOCK_RESULTS_EMPTY_EVENTS_JSON = `
{
	"jsonrpc": "2.0",
	"id": -1,
	"result": {
		"height": "1",
		"txs_results": null,
		"begin_block_events": null,
		"end_block_events": null,
		"validator_updates": null,
		"consensus_param_updates": null
	}
}`
	BLOCK_JSON = `
{
	"jsonrpc":"2.0",
	"id":-1,
	"result":{
		"block_id":{
			"hash":"2CD22EA622D190B9ABCAB797E2F60F6F4FCFC19CC0A67642E5E7856CEAD78163",
			"parts":{
				"total":"1",
				"hash":"AABD95EF4748FE1741CB80D3EB4A1CC74EAA40353AC6A0FC150461C4D3F9A7CC"
			}
		},
		"block":{
			"header":{
				"version":{
					"block":"10",
					"app":"1"
				},
				"chain_id":"testnet-thaler-crypto-com-chain-42",
				"height":"3510",
				"time":"2020-05-07T10:00:25.582998694Z",
				"last_block_id":{
					"hash":"BE33BBB643E19D90DDE8FEA7334F43DF8479F3D41ACD91F3B89FF97742ED2C1A",
					"parts":{
						"total":"1",
						"hash":"D5BDA989B8BEFE50B51C89E813E1303869565533AFEAA17CAD0AC559B4C1C331"
					}
				},
				"last_commit_hash":"740754AA358EAA4684F832A06EA6D843CAE836B8C7B168F92659C781E78C1690",
				"data_hash":"071B8A20B69168CE4E6DB8EC8D6137CE13A20FE9AFC956ECDAED9B1F36A77732",
				"validators_hash":"53E9BA193E306C5EEA55F233FB0F46EFC00826B645FF3A5C6BE7428D6E2970CD",
				"next_validators_hash":"53E9BA193E306C5EEA55F233FB0F46EFC00826B645FF3A5C6BE7428D6E2970CD",
				"consensus_hash":"048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
				"app_hash":"4AC923568B9DE0AA67A04FC60CC4EA10ECC69D636AACFC519911BB448FEED222",
				"last_results_hash":"6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D",
				"evidence_hash":"4F25A690A2C7F0961EC427D5F5AF9BED0CDFA2A16BAB10AB1A3A71DF5E60155B",
				"proposer_address":"7570B2D23A4C7B638BEFE02EB4FC7927BFDED6B7"
			},
			"data":{
				"txs": [
					"AAAELi1Eyx7NNI3yHk9foI+GtqHLPr5/dflC/DyOdCDniaIBAAIAAAAAAAAAAADNM4VutAy2aB2/DoBlBJC/VJxu+WEpaq4jt0BMOpUgxNEiL4lCUGhvlc3FeFCZ99Fvp9Xz9BOCqZ11XxENUFP4moIHvLZMsrohp3ElmtrJFQpZbeS+Ik0FCbT0erhj/wLMptWgxVY7sD/KggQKJjygJb+jdlre2kZ0f8C+nuYW7Q7k7zR1l/D3zFRsLZqbwkBHGchFt3rc3GKiKN2KhNnUoOKG6EvCDQ/HLNwMg0VA3xAW108NTwMnq3prMfjFcZ8QGOBGHWVjDb3GOL/3vngoNr0Jq98oRVEol+sbtvx3s2eO5mGbApiaDK1/oIh9WBaXxEwtXeSPGrw66/8sa9UoaJWBnNOSTesKls/SP1ZQGT1zlrxAZs9MFwjYDbGgukbWletVqgfBcYt3mrvj/G8LtBUH90A/vrJkbvZ6jvTQXqVVHk1jXuU=",
					"AAAEBuKX03ws+y30Jc0oeZTup1aLw0mEB9W8J+LDAaSaOQUBAAIAAAAAAAAAAACILvI8Gcck+IwmflRlBHvHQzIQdFoMSrOXl9a/ZgETkCyFWpwZl8tmJ0mhjGn2bPB7O/srnVyUKxV58j0C1BKFskIFJIr+lmr75vuM9m2DnyL1g2+TDu9nlZWvl+JHR49h0ToA36f0XshYVjyWG6r6O322w4HHP1ApOvuMpq458Hfy3HHnlgaKbcFrwUEbTDrgRslwpgR24jRvVY38IhkEIEHDgXVl8rsXsc04/UsaLQyYlaDi0zqks3qIuPguTQkQ5++/qDvpK6OromSAdCVHhnfdzf6UKsFFlKv9jMcUSeyg5Obx9u/76C2Gfww/iJ0dTZ/Ob0gVQGp+tWyyfa+2OMac++nhGckwru+mw5RoosoFiLMvmeMQrcHR3mcPIaBBx9p4dndcp0JIUEhzFbo29cm1Z8vsbd5nQHbNOxSN5/daIoEJUcQ="
				]
			},
			"evidence":{
				"evidence":[
					{
						"type":"tendermint/DuplicateVoteEvidence",
						"value":{
							"PubKey":{
								"type":"tendermint/PubKeyEd25519",
								"value":"rXhu7xhqYBtJftVLKxvKN0XnpyOzxFnUEfAhD1dEF/8="
							},
							"VoteA":{
								"type":1,
								"height":"3509",
								"round":"0",
								"block_id":{
									"hash":"2C278F10E96EB892FF456D91EB48CFFA31679692102A6FC6F600DC51ABAAE989",
									"parts":{
										"total":"1",
										"hash":"C5E150678DF167F29717238E3A5AC90DCEAD3873BA27C8ED44D96238441767BE"
									}
								},
								"timestamp":"2020-05-07T10:00:25.395605367Z",
								"validator_address":"34C725CABA703269B3F1D1A907A84DE5FEE96469",
								"validator_index":"0",
								"signature":"tBWO1Wf9lf39WXpU2aJ1hzZF8nXl+D0izGt1FT1acA/nu/ezE2pVeQDgrU85b16ENTuYy375p2hdaXvyjESCBw=="
							},
							"VoteB":{
								"type":1,
								"height":"3509",
								"round":"0",
								"block_id":{
									"hash":"BE33BBB643E19D90DDE8FEA7334F43DF8479F3D41ACD91F3B89FF97742ED2C1A",
									"parts":{
										"total":"1",
										"hash":"D5BDA989B8BEFE50B51C89E813E1303869565533AFEAA17CAD0AC559B4C1C331"
									}
								},
								"timestamp":"2020-05-07T10:00:25.330002296Z",
								"validator_address":"34C725CABA703269B3F1D1A907A84DE5FEE96469",
								"validator_index":"0",
								"signature":"LqTXlX4msVieutlPzQ9sQCa/ClvY4VYYDjfAZQBH7mT2cJjUS3r7hK8q4VUCktT0Regryrbb40Bts2SJaaYbBw=="
							}
						}
					},
					{
						"type":"tendermint/DuplicateVoteEvidence",
						"value":{
							"PubKey":{
								"type":"tendermint/PubKeyEd25519",
								"value":"rXhu7xhqYBtJftVLKxvKN0XnpyOzxFnUEfAhD1dEF/8="
							},
							"VoteA":{
								"type":2,
								"height":"3509",
								"round":"0",
								"block_id":{
									"hash":"",
									"parts":{
										"total":"0",
										"hash":""
									}
								},
								"timestamp":"2020-05-07T10:00:25.809310628Z",
								"validator_address":"34C725CABA703269B3F1D1A907A84DE5FEE96469",
								"validator_index":"0",
								"signature":"q3abt8718Z1jZcVPjrkqhxmrYM40HEIEzg/1Cu0NbwFkXJBhz1gNRueLJ8jZ1RTu0hbTTLoHUApE0lOdv4efAg=="
							},
							"VoteB":{
								"type":2,
								"height":"3509",
								"round":"0",
								"block_id":{
									"hash":"BE33BBB643E19D90DDE8FEA7334F43DF8479F3D41ACD91F3B89FF97742ED2C1A",
									"parts":{
										"total":"1",
										"hash":"D5BDA989B8BEFE50B51C89E813E1303869565533AFEAA17CAD0AC559B4C1C331"
									}
								},
								"timestamp":"2020-05-07T10:00:25.617739872Z",
								"validator_address":"34C725CABA703269B3F1D1A907A84DE5FEE96469",
								"validator_index":"0",
								"signature":"W6wNdtDCsUbdodju57/2BlkmUjo6U0PH+Cf19u4RSlaBYS7svNpZOgdQqmXJnUInRjGJp8opE7a9FnnHe3oTAw=="
							}
						}
					}
				]
			},
			"last_commit":{
				"height":"3509",
				"round":"0",
				"block_id":{
					"hash":"BE33BBB643E19D90DDE8FEA7334F43DF8479F3D41ACD91F3B89FF97742ED2C1A",
					"parts":{
						"total":"1",
						"hash":"D5BDA989B8BEFE50B51C89E813E1303869565533AFEAA17CAD0AC559B4C1C331"
					}
				},
				"signatures":[
					{
						"block_id_flag": 1,
						"validator_address": "",
						"timestamp": "0001-01-01T00:00:00Z",
						"signature": null
					},
					{
						"block_id_flag":2,
						"validator_address":"34C725CABA703269B3F1D1A907A84DE5FEE96469",
						"timestamp":"2020-05-07T10:00:25.617739872Z",
						"signature":"W6wNdtDCsUbdodju57/2BlkmUjo6U0PH+Cf19u4RSlaBYS7svNpZOgdQqmXJnUInRjGJp8opE7a9FnnHe3oTAw=="
					},
					{
						"block_id_flag":2,
						"validator_address":"7570B2D23A4C7B638BEFE02EB4FC7927BFDED6B7",
						"timestamp":"2020-05-07T10:00:25.582998694Z",
						"signature":"W2Pkxi/AqAklmEMjsd3P9yd//1ZGxBqRMaHebcGYhlVZxcbZ02dzwgxD7c/BOOMh+kJGYPfuYNiLHD3Kts+pDA=="
					},
					{
						"block_id_flag":2,
						"validator_address":"D527DAECDE0501CF2E785A8DC0D9F4A64760F0BB",
						"timestamp":"2020-05-07T10:00:25.622478335Z",
						"signature":"zCjxVQLbEhBdBBm1VHyLP+4aFH81ke26pC0e+g+pvAzLBvWxrmzdgh347MmOVHWiW6lS9nb8xs+6bkdKRMx5Dg=="
					},
					{
						"block_id_flag":2,
						"validator_address":"FA7B721B5704DF98EF3ECD3796DDEF6AA2A80257",
						"timestamp":"2020-05-07T10:00:25.574877035Z",
						"signature":"TuflcSDZPgVTM618J9JF/tlMFwM8Z/eWrHPjixzeWukIlkFHsNMprRRPnHUKZlu+yDdSwgj6eJ0PrqRm7y/eDw=="
					}
				]
			}
		}
	}
}`
	BLOCK_EMPTY_TX_RESULT = `
{
	"jsonrpc":"2.0",
	"id":-1,
	"result":{
		"block_id":{
			"hash":"BD3D9499EF527035BAE14E7F4932BD1778C1E90141A7A76A1FE67B7ECCC2663E",
			"parts":{
				"total":"1",
				"hash":"47A69EDDDDC80CF7DF46C96A76B2C500C91BC5380439D15E8F8968AFDBB7BC7C"
			}
		},
		"block":{
			"header":{
				"version":{
					"block":"10",
					"app":"1"
				},
				"chain_id":"testnet-thaler-crypto-com-chain-42",
				"height":"1",
				"time":"2020-05-01T12:09:01.568951Z",
				"last_block_id":{
					"hash":"",
					"parts":{
						"total":"0",
						"hash":""
					}
				},
				"last_commit_hash":"",
				"data_hash":"",
				"validators_hash":"FA88CACC192B10984F1F34FD1FE32A1E7CA909ADD18457E2C00B0265DC8FAA07",
				"next_validators_hash":"FA88CACC192B10984F1F34FD1FE32A1E7CA909ADD18457E2C00B0265DC8FAA07",
				"consensus_hash":"048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
				"app_hash":"F62DDB49D7EB8ED0883C735A0FB7DE7F2A3FA322FCD2AA832F452A62B38607D5",
				"last_results_hash":"",
				"evidence_hash":"",
				"proposer_address":"D527DAECDE0501CF2E785A8DC0D9F4A64760F0BB"
			},
			"data":{
				"txs":null
			},
			"evidence":{
				"evidence":null
			},
			"last_commit":{
				"height":"0",
				"round":"0",
				"block_id":{
					"hash":"",
					"parts":{
						"total":"0",
						"hash":""
					}
				},
				"signatures":null
			}
		}
	}
}`
)

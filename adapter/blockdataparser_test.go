package adapter_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chainindex"
	. "github.com/crypto-com/chainindex/adapter"
	tenderminttypes "github.com/crypto-com/chainindex/adapter/tendermint/types"
	"github.com/crypto-com/chainindex/internal/bignum"
	"github.com/crypto-com/chainindex/internal/primptr"
	"github.com/crypto-com/chainindex/usecase"
)

var _ = Describe("Blockdataparser", func() {
	Describe("ParseGenesis", func() {
		It("should return parsed BlockData from the genesis", func() {
			genesis := sampleGenesis()
			block := sampleGenesisBlock()
			rawBlockData := TendermintGenesisBlockData{
				Genesis: &genesis,
				Block:   &block,
			}

			actualBlockData := ParseGenesisToBlockData(rawBlockData)
			Expect(*actualBlockData).To(Equal(usecase.BlockData{
				Block: chainindex.Block{
					Height:  block.Height,
					Hash:    block.Hash,
					Time:    block.Time,
					AppHash: block.AppHash,
				},
				Activities: []chainindex.Activity{
					{
						BlockHeight:                uint64(1),
						Type:                       chainindex.ACTIVITY_GENESIS,
						MaybeStakingAccountAddress: primptr.String("0x4ae85b35597fcb61c6c47b1fe0bdd7eed8421cdd"),
						MaybeBonded:                bignum.MustAtoi("6000000000000000"),
					},
					{
						BlockHeight:                uint64(1),
						Type:                       chainindex.ACTIVITY_GENESIS,
						MaybeStakingAccountAddress: primptr.String("0x4b75f275dde0a8c8e70fb84243adc97a3afb78f2"),
						MaybeUnbonded:              bignum.MustAtoi("7946000000000000000"),
					},
					{
						BlockHeight:                uint64(1),
						Type:                       chainindex.ACTIVITY_GENESIS,
						MaybeStakingAccountAddress: primptr.String("0x4fd8162521f2e628adced7c1baa39384a08b4a3d"),
						MaybeBonded:                bignum.MustAtoi("6000000000000000"),
					},
					{
						BlockHeight:                uint64(1),
						Type:                       chainindex.ACTIVITY_GENESIS,
						MaybeStakingAccountAddress: primptr.String("0x6c2be7846219eab3086a66f873558b73d8f4a0d4"),
						MaybeBonded:                bignum.MustAtoi("6000000000000000"),
					},
					{
						BlockHeight:                uint64(1),
						Type:                       chainindex.ACTIVITY_GENESIS,
						MaybeStakingAccountAddress: primptr.String("0x6dbd5b8fe0dad494465aa7574defba711c184102"),
						MaybeBonded:                bignum.MustAtoi("6000000000000000"),
						MaybeCouncilNodeMeta: &chainindex.CouncilNode{
							Name:                 "eastus_validator_1",
							MaybeSecurityContact: primptr.String("security@crypto.com"),
							PubKeyType:           chainindex.PUBKEY_TYPE_ED25519,
							PubKey:               "/SvfTeO4Du4oR/VYTjm7IgObc14zzddEAyFb4nU8E3Q=",
							Address:              "FA7B721B5704DF98EF3ECD3796DDEF6AA2A80257",
							CreatedAtBlockHeight: uint64(1),
						},
					},
					{
						BlockHeight:                uint64(1),
						Type:                       chainindex.ACTIVITY_GENESIS,
						MaybeStakingAccountAddress: primptr.String("0x6fc1e3124a7ed07f3710378b68f7046c7300179d"),
						MaybeBonded:                bignum.MustAtoi("6000000000000000"),
						MaybeCouncilNodeMeta: &chainindex.CouncilNode{
							Name:                 "canadacentral_validator_1",
							MaybeSecurityContact: primptr.String("security@crypto.com"),
							PubKeyType:           chainindex.PUBKEY_TYPE_ED25519,
							PubKey:               "QMegiWt9+5K1b1ZVd7zOJZxhTnbAtWzvGhViiElAlaw=",
							Address:              "7570B2D23A4C7B638BEFE02EB4FC7927BFDED6B7",
							CreatedAtBlockHeight: uint64(1),
						},
					},
					{
						BlockHeight:                uint64(1),
						Type:                       chainindex.ACTIVITY_GENESIS,
						MaybeStakingAccountAddress: primptr.String("0x9baa6de71cbc6274275eece4b1be15f545897f37"),
						MaybeBonded:                bignum.MustAtoi("6000000000000000"),
					},
					{
						BlockHeight:                uint64(1),
						Type:                       chainindex.ACTIVITY_GENESIS,
						MaybeStakingAccountAddress: primptr.String("0xa9528abb92709370600d2cef41f1677374278337"),
						MaybeBonded:                bignum.MustAtoi("6000000000000000"),
					},
					{
						BlockHeight:                uint64(1),
						Type:                       chainindex.ACTIVITY_GENESIS,
						MaybeStakingAccountAddress: primptr.String("0xb328a39002ede64c33bb60f1dc43f5df9eb47043"),
						MaybeBonded:                bignum.MustAtoi("6000000000000000"),
					},
					{
						BlockHeight:                uint64(1),
						Type:                       chainindex.ACTIVITY_GENESIS,
						MaybeStakingAccountAddress: primptr.String("0xb8c6886da09e12db8aebfc8108c67ce2ba086ac6"),
						MaybeBonded:                bignum.MustAtoi("6000000000000000"),
						MaybeCouncilNodeMeta: &chainindex.CouncilNode{
							Name:                 "uksouth_validator_1",
							MaybeSecurityContact: primptr.String("security@crypto.com"),
							PubKeyType:           chainindex.PUBKEY_TYPE_ED25519,
							PubKey:               "tDLheZJwsA8oYEwarR6/X+zAmNKMLHTVkh/fvcLqcwA=",
							Address:              "D527DAECDE0501CF2E785A8DC0D9F4A64760F0BB",
							CreatedAtBlockHeight: uint64(1),
						},
					},
				},
			}))
		})
	})

	Describe("ParseBlock", func() {
		It("should parse signatures", func() {
			anyBlockHeight := uint64(33339)
			anyBlockTime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "2020-05-15T13:19:07.783916012Z")
			block := tenderminttypes.Block{
				Height:         anyBlockHeight,
				Hash:           "42E2A7C6AA135D2652ED8C0BEEB446BFC2B4A54B679FE07109CD249F42EC853C",
				Time:           anyBlockTime,
				AppHash:        "DCAD7CEAD2B0A6A668E1D8458E5C1E2B1B790AFFD2D56B36EE1690825BD9818C",
				PropserAddress: "64161C75F5F2A78806267AB3A4E2BD3C04944AA0",
				Txs:            nil,
				Signatures: []tenderminttypes.BlockSignature{
					{
						ValidatorAddress: "34C725CABA703269B3F1D1A907A84DE5FEE96469",
						Signature:        "waN7MvfcTUwA8hD5aueM5XXsZ4hkwpXCP5MMO0xr/njryxNx1hrfyPD3z07DXPAasVFVrD4mwjkbowzM4T+mCg==",
					},
					{
						ValidatorAddress: "64161C75F5F2A78806267AB3A4E2BD3C04944AA0",
						Signature:        "N9PoH1tgBTLtPxqxfDNWKMrAPZVBpNvrzQEXQpfZMO2TXvBwKW4Gmw31b0bTUMsiyLJJkSAI+UF5zRrHEu2FDQ==",
					},
					{
						ValidatorAddress: "7570B2D23A4C7B638BEFE02EB4FC7927BFDED6B7",
						Signature:        "t3b8S0yakf2xdCGRzsXskGegAv0xu406T7jbBoH+fqfcO3wnbvaD8xUs5A1zOKPEVzj5aylzym/w074K+omoAg==",
					},
					{
						ValidatorAddress: "D527DAECDE0501CF2E785A8DC0D9F4A64760F0BB",
						Signature:        "x9YjNUFpi+W8WcdJXJdUgsfzYUmrqmC7NYDo0vpYpqfRh33VrUWjxkFG6auB0uliEy5yBTEZdivR+Qxb964gCg==",
					},
					{
						ValidatorAddress: "FA7B721B5704DF98EF3ECD3796DDEF6AA2A80257",
						Signature:        "jMos+ST8x69X16GJhRqzpcG0Bxk9aRt0jifL6gJuZGqX50i9iGD5DGNZSkDfTv1YnHL3YdSsOpfXqmbuVCCuAw==",
					},
				},
			}
			blockResults := tenderminttypes.BlockResults{
				Height:           anyBlockHeight,
				TxsEvents:        nil,
				BeginBlockEvents: nil,
				ValidatorUpdates: nil,
			}

			actualBlockData := ParseBlockToBlockData(TendermintBlockData{
				Block:        &block,
				BlockResults: &blockResults,
			})
			Expect(*actualBlockData).To(Equal(usecase.BlockData{
				Block: chainindex.Block{
					Height:  anyBlockHeight,
					Hash:    "42E2A7C6AA135D2652ED8C0BEEB446BFC2B4A54B679FE07109CD249F42EC853C",
					Time:    anyBlockTime,
					AppHash: "DCAD7CEAD2B0A6A668E1D8458E5C1E2B1B790AFFD2D56B36EE1690825BD9818C",
				},
				Signatures: []chainindex.BlockSignature{
					{
						BlockHeight:        anyBlockHeight,
						CouncilNodeAddress: "34C725CABA703269B3F1D1A907A84DE5FEE96469",
						Signature:          "waN7MvfcTUwA8hD5aueM5XXsZ4hkwpXCP5MMO0xr/njryxNx1hrfyPD3z07DXPAasVFVrD4mwjkbowzM4T+mCg==",
						IsProposer:         false,
					},
					{
						BlockHeight:        anyBlockHeight,
						CouncilNodeAddress: "64161C75F5F2A78806267AB3A4E2BD3C04944AA0",
						Signature:          "N9PoH1tgBTLtPxqxfDNWKMrAPZVBpNvrzQEXQpfZMO2TXvBwKW4Gmw31b0bTUMsiyLJJkSAI+UF5zRrHEu2FDQ==",
						IsProposer:         true,
					},
					{
						BlockHeight:        anyBlockHeight,
						CouncilNodeAddress: "7570B2D23A4C7B638BEFE02EB4FC7927BFDED6B7",
						Signature:          "t3b8S0yakf2xdCGRzsXskGegAv0xu406T7jbBoH+fqfcO3wnbvaD8xUs5A1zOKPEVzj5aylzym/w074K+omoAg==",
						IsProposer:         false,
					},
					{
						BlockHeight:        anyBlockHeight,
						CouncilNodeAddress: "D527DAECDE0501CF2E785A8DC0D9F4A64760F0BB",
						Signature:          "x9YjNUFpi+W8WcdJXJdUgsfzYUmrqmC7NYDo0vpYpqfRh33VrUWjxkFG6auB0uliEy5yBTEZdivR+Qxb964gCg==",
						IsProposer:         false,
					},
					{
						BlockHeight:        anyBlockHeight,
						CouncilNodeAddress: "FA7B721B5704DF98EF3ECD3796DDEF6AA2A80257",
						Signature:          "jMos+ST8x69X16GJhRqzpcG0Bxk9aRt0jifL6gJuZGqX50i9iGD5DGNZSkDfTv1YnHL3YdSsOpfXqmbuVCCuAw==",
						IsProposer:         false,
					},
				},
			}))
		})
		It("should parse transfer activity", func() {
			anyBlockHeight := uint64(32168)
			anyBlockTime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "2020-05-15T05:35:05.012038715Z")
			block := tenderminttypes.Block{
				Height:         anyBlockHeight,
				Hash:           "964738311B5657215DC8857FD95FD270C561D1FF93967F9AFCBC398DD52D768A",
				Time:           anyBlockTime,
				AppHash:        "0C2269A48DFE4BCA4CBE1008103D71F133B14C4B7D25C7F34441EAC9C83731E8",
				PropserAddress: "7570B2D23A4C7B638BEFE02EB4FC7927BFDED6B7",
				Txs: []string{
					"AAAEYz5CTye1JIFcys8yZeq4/qzMWr51J21h5SfGoZ7ZKXoBAAIAAAAAAAAAAACQ0AKWJ05yBACTC9ZlBCukZk4EUptxeu8RvUrtX9j4Bxce21WL3c/KYqVPVLolJR9+bDZdNdBaUZ0ff1KyG1KQCK/NHKScILSF/fYxeDjf7rTgYIlqgHpvUDqEIRTeU8Y1dZbecUBT7r9H2IJ64Uou+CEWJpbX8LCgCVH5ropvE/0SsT4fmskhENDqYEE1Mujy2WNgsB8CJ5R40ZLvoVNysN/TvlhfJmcoWEBAWdHpIMnROM+LD33RijZ6okymabCIKQfd4nmlQtRvIU16mmbGfWqkh77rhqgK6AOH5qcQHq13X+Onmt4ScCRsdSfTlCp1OXuuK/s5K5UpkML90Lomn6LTZnHYeTvm/caOPOruMzCWEe0eLzacVqm3aVTFECMNFpOZRexzz7+QhAducX1i7B3cphBvr4mPe/4UIg93nkTdVERFBd4=",
					"AAAErQfrO6yrbZC596XchZcG5znjfJkHEeBTpeqKe7cSNEMBAAIAAAAAAAAAAADvZYMveh0B14aX4TllBHzVJYe9empAx1AGVDADLAZRG3jKclBQlV5XHu8E2qztbmsKUQ3q88ItIaYQJtRly7nFiep/CzscfVlHiw5Zt5ut1/ZQovnNCGyBGe9eR0aV7V42QIrqLSr5D5Ck/qiLfh4xaziCKvdRMMiX/L4qT30UyhEXT3ix3v775qwfMdyURI4Fb2Z6wZpuL3bqBVv1xh5G2aBzhcjM24f8AlqccW6y/ZAGNvejUbh/L4D0wAHpQtI5BnBEtFUnxZzrBA5+xwDByXdZko+PPdsa8UaNkQwjcb8Spm2v7R0Ylz+6qou45jmalSBCxeCiiKymb7wYNJfjAtr/xySPaLxTMClmSvmhd5aUXzX5Nx7GUUMIGyOEskk/MgwtUc6yyPBfkstVOiZIZ40OO4bz5XRogxZtecUQTwW4T7KRGZo=",
				},
				Signatures: nil,
			}
			blockResults := tenderminttypes.BlockResults{
				Height: anyBlockHeight,
				TxsEvents: [][]tenderminttypes.BlockResultsEvent{
					{
						{
							Type: "valid_txs",
							Attributes: []tenderminttypes.BlockResultsEventAttribute{
								{
									Key:   "ZmVl",
									Value: "MC4wMDAwMDQ2OQ==",
								},
								{
									Key:   "dHhpZA==",
									Value: "Y2ZiZjkwODQwNzZlNzE3ZDYyZWMxZGRjYTYxMDZmYWY4OThmN2JmZTE0MjIwZjc3OWU0NGRkNTQ0NDQ1MDVkZQ==",
								},
							},
						},
					},
					{
						{
							Type: "valid_txs",
							Attributes: []tenderminttypes.BlockResultsEventAttribute{
								{
									Key:   "ZmVl",
									Value: "MC4wMDAwMDQ2OQ==",
								},
								{
									Key:   "dHhpZA==",
									Value: "YzhmMDVmOTJjYjU1M2EyNjQ4Njc4ZDBlM2I4NmYzZTU3NDY4ODMxNjZkNzljNTEwNGYwNWI4NGZiMjkxMTk5YQ==",
								},
							},
						},
					},
				},
				BeginBlockEvents: nil,
				ValidatorUpdates: nil,
			}

			actualBlockData := ParseBlockToBlockData(TendermintBlockData{
				Block:        &block,
				BlockResults: &blockResults,
			})
			Expect(*actualBlockData).To(Equal(usecase.BlockData{
				Block: chainindex.Block{
					Height:  anyBlockHeight,
					Hash:    "964738311B5657215DC8857FD95FD270C561D1FF93967F9AFCBC398DD52D768A",
					Time:    anyBlockTime,
					AppHash: "0C2269A48DFE4BCA4CBE1008103D71F133B14C4B7D25C7F34441EAC9C83731E8",
				},
				Activities: []chainindex.Activity{
					{
						BlockHeight: anyBlockHeight,
						Type:        chainindex.ACTIVITY_TRANSFER,
						MaybeTxID:   primptr.String("cfbf9084076e717d62ec1ddca6106faf898f7bfe14220f779e44dd54444505de"),
						MaybeFee:    bignum.MustAtoi("469"),
						MaybeTxInputs: []chainindex.TxInput{
							{
								TxId:  "633e424f27b524815ccacf3265eab8feaccc5abe75276d61e527c6a19ed9297a",
								Index: uint32(1),
							},
						},
						MaybeOutputCount: primptr.Uint32(2),
					},
					{
						BlockHeight: anyBlockHeight,
						Type:        chainindex.ACTIVITY_TRANSFER,
						MaybeTxID:   primptr.String("c8f05f92cb553a2648678d0e3b86f3e5746883166d79c5104f05b84fb291199a"),
						MaybeFee:    bignum.MustAtoi("469"),
						MaybeTxInputs: []chainindex.TxInput{
							{
								TxId:  "ad07eb3bacab6d90b9f7a5dc859706e739e37c990711e053a5ea8a7bb7123443",
								Index: uint32(1),
							},
						},
						MaybeOutputCount: primptr.Uint32(2),
					},
				},
			}))
		})

		It("should parse deposit activity", func() {
			anyBlockHeight := uint64(29220)
			anyBlockTime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "2020-05-14T09:13:45.866764762Z")
			block := tenderminttypes.Block{
				Height: anyBlockHeight,
				Hash:   "77C918D6719285702468F87A017444B3766C79F77540211C1091E083C14BD4D1",
				Time:   anyBlockTime,
				Txs: []string{
					"AAEEEE85Ht83dvrnlKV2rh2byE9Qs1vKzdAOIkUbnE8ja74AAABLdfJ13eCoyOcPuEJDrcl6Ovt48gBCAQAAAAAAAAAAAAAAAAAAAEnuiOdDBdlZ8SU6gNEB0r8zw/QUhtcFC1tdxFoxHO9HdUvDJRMHx6PE7oP22mmbN33PRmJ+kKtNHTCD12ERlQOSOxccDXPygTKqpdZ6Z6M4Q/E/eHhtSBdJHzbfSWY1lHloXJyNLyl9eQjrIDMFozyn+g7umf9NbhGNaBtXrHNsXGcYjt0U2bRlwsXQ6BCSLJ8p1SkmAji45n9BtKfqC2WCcA==",
				},
				Signatures: nil,
			}
			blockResults := tenderminttypes.BlockResults{
				Height: anyBlockHeight,
				TxsEvents: [][]tenderminttypes.BlockResultsEvent{
					{
						{
							Type: "valid_txs",
							Attributes: []tenderminttypes.BlockResultsEventAttribute{
								{
									Key:   "ZmVl",
									Value: "MC4wMDAwMDI5OQ==",
								},
								{
									Key:   "dHhpZA==",
									Value: "MTg4ZWRkMTRkOWI0NjVjMmM1ZDBlODEwOTIyYzlmMjlkNTI5MjYwMjM4YjhlNjdmNDFiNGE3ZWEwYjY1ODI3MA==",
								},
							},
						},
						{
							Type: "staking_change",
							Attributes: []tenderminttypes.BlockResultsEventAttribute{
								{
									Key:   "c3Rha2luZ19hZGRyZXNz",
									Value: "MHg0Yjc1ZjI3NWRkZTBhOGM4ZTcwZmI4NDI0M2FkYzk3YTNhZmI3OGYy",
								},
								{
									Key:   "c3Rha2luZ19vcHR5cGU=",
									Value: "ZGVwb3NpdA==",
								},
								{
									Key:   "c3Rha2luZ19kaWZm",
									Value: "W3sia2V5IjoiQm9uZGVkIiwidmFsdWUiOiIxMDAwMDAwMDAifV0=",
								},
							},
						},
					},
				},
				BeginBlockEvents: nil,
				ValidatorUpdates: nil,
			}

			actualBlockData := ParseBlockToBlockData(TendermintBlockData{
				Block:        &block,
				BlockResults: &blockResults,
			})
			Expect(*actualBlockData).To(Equal(usecase.BlockData{
				Block: chainindex.Block{
					Height: anyBlockHeight,
					Hash:   "77C918D6719285702468F87A017444B3766C79F77540211C1091E083C14BD4D1",
					Time:   anyBlockTime,
				},
				Activities: []chainindex.Activity{
					{
						BlockHeight:                anyBlockHeight,
						Type:                       chainindex.ACTIVITY_DEPOSIT,
						MaybeTxID:                  primptr.String("188edd14d9b465c2c5d0e810922c9f29d529260238b8e67f41b4a7ea0b658270"),
						MaybeFee:                   bignum.MustAtoi("299"),
						MaybeStakingAccountAddress: primptr.String("0x4b75f275dde0a8c8e70fb84243adc97a3afb78f2"),
						MaybeBonded:                bignum.MustAtoi("100000000"),
					},
				},
			}))
		})

		It("should parse unbond transaction", func() {
			anyBlockHeight := uint64(32702)
			anyBlockTime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "2020-05-15T09:17:42.981053198Z")
			block := tenderminttypes.Block{
				Height: anyBlockHeight,
				Hash:   "6400988CEBAABF9AE050E23CB36318372AB725747AF823D27C45452B0780D2D6",
				Time:   anyBlockTime,
				Txs: []string{
					"AQAAS3Xydd3gqMjnD7hCQ63Jejr7ePIDAAAAAAAAAADodkgXAAAAAEIBAAAAAAAAAAAAjCfHg46Jt0tvGoWsLkBukvOGxV3IN9u2esyCc6YHT4IILPXX6T4jGTDZqqgE8NWfk8TJ+Z0AFFVcgb+3tWu2iA==",
				},
				Signatures: nil,
			}
			blockResults := tenderminttypes.BlockResults{
				Height: anyBlockHeight,
				TxsEvents: [][]tenderminttypes.BlockResultsEvent{
					{
						{
							Type: "valid_txs",
							Attributes: []tenderminttypes.BlockResultsEventAttribute{
								{
									Key:   "ZmVl",
									Value: "MC4wMDAwMDE0NQ==",
								},
								{
									Key:   "dHhpZA==",
									Value: "N2Y3YjljNzNkZTlhM2JlMDM3YTRiY2NjNmRmYzEzM2YxNDY0OTNmYWM5ZTdkNWU1OGJhNTk1MjJkOTQ1ZDAyNA==",
								},
							},
						},
						{
							Type: "staking_change",
							Attributes: []tenderminttypes.BlockResultsEventAttribute{
								{
									Key:   "c3Rha2luZ19hZGRyZXNz",
									Value: "MHg0Yjc1ZjI3NWRkZTBhOGM4ZTcwZmI4NDI0M2FkYzk3YTNhZmI3OGYy",
								},
								{
									Key:   "c3Rha2luZ19vcHR5cGU=",
									Value: "dW5ib25k",
								},
								{
									Key:   "c3Rha2luZ19kaWZm",
									Value: "W3sia2V5IjoiQm9uZGVkIiwidmFsdWUiOiItMTAwMDAwMDAwMTQ1In0seyJrZXkiOiJVbmJvbmRlZCIsInZhbHVlIjoiMTAwMDAwMDAwMDAwIn0seyJrZXkiOiJVbmJvbmRlZEZyb20iLCJ2YWx1ZSI6MTU4OTUzOTY2Mn1d",
								},
							},
						},
					},
				},
				BeginBlockEvents: nil,
				ValidatorUpdates: nil,
			}

			actualBlockData := ParseBlockToBlockData(TendermintBlockData{
				Block:        &block,
				BlockResults: &blockResults,
			})
			Expect(*actualBlockData).To(Equal(usecase.BlockData{
				Block: chainindex.Block{
					Height: anyBlockHeight,
					Hash:   "6400988CEBAABF9AE050E23CB36318372AB725747AF823D27C45452B0780D2D6",
					Time:   anyBlockTime,
				},
				Activities: []chainindex.Activity{
					{
						BlockHeight:                anyBlockHeight,
						Type:                       chainindex.ACTIVITY_UNBOND,
						MaybeTxID:                  primptr.String("7f7b9c73de9a3be037a4bccc6dfc133f146493fac9e7d5e58ba59522d945d024"),
						MaybeFee:                   bignum.MustAtoi("145"),
						MaybeStakingAccountAddress: primptr.String("0x4b75f275dde0a8c8e70fb84243adc97a3afb78f2"),
						MaybeBonded:                bignum.MustAtoi("-100000000145"),
						MaybeUnbonded:              bignum.MustAtoi("100000000000"),
						MaybeUnbondedFrom:          primptr.Time(time.Unix(1589539662, 0).UTC()),
					},
				},
			}))
		})

		It("should parse withdraw transaction", func() {
			anyBlockHeight := uint64(32998)
			anyBlockTime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "2020-05-15T09:17:42.981053198Z")
			block := tenderminttypes.Block{
				Height: anyBlockHeight,
				Hash:   "8231F933B562FB64C7AD97BB8A200D681D7EAD0CE6174E60B6E4B8F38994A85E",
				Time:   anyBlockTime,
				Txs: []string{
					"AAIBAAAAW/4Lr7nPBwsjXaN3e3yX/qD+fNV/XQuLt8l9CG7+UhVKQMHLGSANQiCHFQC+6kJYwKN0Ktg2pCyoSuYwJsonmwAAAAAAAAAAbH6NDyezgDjufby45QHPGXU1SINMmemVavh8lOrTnTE/sOl+H3MexFuTX/5RQkN/MEd6sdKvkiXRU+xXU/ERxwwy+Mxn2oqAZEQmTxK0wyFt4Zw9CpJYZOFvxfCfK2+gAemoy2ehFt/5I9xoT36AoVIPjDXoFUSxmx5qYFq767UfESPvWekkMnkBYn/kcHe3iJKeA5t3F7ze3+DgVplirAc7431AbOE=",
				},
				Signatures: nil,
			}
			blockResults := tenderminttypes.BlockResults{
				Height: anyBlockHeight,
				TxsEvents: [][]tenderminttypes.BlockResultsEvent{
					{
						{
							Type: "valid_txs",
							Attributes: []tenderminttypes.BlockResultsEventAttribute{
								{
									Key:   "ZmVl",
									Value: "MC4wMDAwMDMwOA==",
								},
								{
									Key:   "dHhpZA==",
									Value: "MzI3OTAxNjI3ZmU0NzA3N2I3ODg5MjllMDM5Yjc3MTdiY2RlZGZlMGUwNTY5OTYyYWMwNzNiZTM3ZDQwNmNlMQ==",
								},
							},
						},
						{
							Type: "staking_change",
							Attributes: []tenderminttypes.BlockResultsEventAttribute{
								{
									Key:   "c3Rha2luZ19hZGRyZXNz",
									Value: "MHg0Yjc1ZjI3NWRkZTBhOGM4ZTcwZmI4NDI0M2FkYzk3YTNhZmI3OGYy",
								},
								{
									Key:   "c3Rha2luZ19vcHR5cGU=",
									Value: "d2l0aGRyYXc=",
								},
								{
									Key:   "c3Rha2luZ19kaWZm",
									Value: "W3sia2V5IjoiVW5ib25kZWQiLCJ2YWx1ZSI6Ii0xMDAwMDAwMDAwMDAifV0=",
								},
							},
						},
					},
				},
				BeginBlockEvents: nil,
				ValidatorUpdates: nil,
			}

			actualBlockData := ParseBlockToBlockData(TendermintBlockData{
				Block:        &block,
				BlockResults: &blockResults,
			})
			Expect(*actualBlockData).To(Equal(usecase.BlockData{
				Block: chainindex.Block{
					Height: anyBlockHeight,
					Hash:   "8231F933B562FB64C7AD97BB8A200D681D7EAD0CE6174E60B6E4B8F38994A85E",
					Time:   anyBlockTime,
				},
				Activities: []chainindex.Activity{
					{
						BlockHeight:                anyBlockHeight,
						Type:                       chainindex.ACTIVITY_WITHDRAW,
						MaybeTxID:                  primptr.String("327901627fe47077b788929e039b7717bcdedfe0e0569962ac073be37d406ce1"),
						MaybeFee:                   bignum.MustAtoi("308"),
						MaybeOutputCount:           primptr.Uint32(uint32(1)),
						MaybeStakingAccountAddress: primptr.String("0x4b75f275dde0a8c8e70fb84243adc97a3afb78f2"),
						MaybeUnbonded:              bignum.MustAtoi("-100000000000"),
					},
				},
			}))
		})

		It("should parse node join transaction", func() {
			anyBlockHeight := uint64(31)
			anyBlockTime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "2020-05-05T10:09:45.282175368Z")
			block := tenderminttypes.Block{
				Height: anyBlockHeight,
				Hash:   "9A362E7D43B0414461698012A41B75FD0E79D0B336ECDA70A8447F1B5D17A5F1",
				Time:   anyBlockTime,
				Txs: []string{
					"AQIAAAAAAAAAAACzKKOQAu3mTDO7YPHcQ/XfnrRwQwBCAAAAAAAAAAAAZGNhbmFkYWNlbnRyYWxfdmFsaWRhdG9yXzIAAK14bu8YamAbSX7VSysbyjdF56cjs8RZ1BHwIQ9XRBf/FEZJWE1FAAEbC33EXRoay1jVBgObgohxS0Q3NFDj/IprjkML6Vj/+nqbrwYBykRAsVxPFXKB6E+qa6II57Ngb3iStwu4Awfx",
				},
				Signatures: nil,
			}
			blockResults := tenderminttypes.BlockResults{
				Height: anyBlockHeight,
				TxsEvents: [][]tenderminttypes.BlockResultsEvent{
					{
						{
							Type: "valid_txs",
							Attributes: []tenderminttypes.BlockResultsEventAttribute{
								{
									Key:   "ZmVl",
									Value: "MC4wMDAwMDAwMA==",
								},
								{
									Key:   "dHhpZA==",
									Value: "ZmMwODIyMTEyMThiMGE1NWUyNDA1ODI0M2I4ZmNhY2Y2ZjIzYjMxNGYxZTFkY2NkNTU3YmQyNTE4NzgyNmFiMg==",
								},
							},
						},
						{
							Type: "staking_change",
							Attributes: []tenderminttypes.BlockResultsEventAttribute{
								{
									Key:   "c3Rha2luZ19hZGRyZXNz",
									Value: "MHhiMzI4YTM5MDAyZWRlNjRjMzNiYjYwZjFkYzQzZjVkZjllYjQ3MDQz",
								},
								{
									Key:   "c3Rha2luZ19vcHR5cGU=",
									Value: "bm9kZWpvaW4=",
								},
								{
									Key:   "c3Rha2luZ19kaWZm",
									Value: "W3sia2V5IjoiQ291bmNpbE5vZGUiLCJ2YWx1ZSI6eyJuYW1lIjoiY2FuYWRhY2VudHJhbF92YWxpZGF0b3JfMiIsInNlY3VyaXR5X2NvbnRhY3QiOm51bGwsImNvbnNlbnN1c19wdWJrZXkiOnsidHlwZSI6InRlbmRlcm1pbnQvUHViS2V5RWQyNTUxOSIsInZhbHVlIjoiclhodTd4aHFZQnRKZnRWTEt4dktOMFhucHlPenhGblVFZkFoRDFkRUYvOD0ifSwiY29uZmlkZW50aWFsX2luaXQiOnsiY2VydCI6IlJrbFlUVVU9In19fV0=",
								},
							},
						},
					},
				},
				BeginBlockEvents: nil,
				ValidatorUpdates: nil,
			}

			actualBlockData := ParseBlockToBlockData(TendermintBlockData{
				Block:        &block,
				BlockResults: &blockResults,
			})
			Expect(*actualBlockData).To(Equal(usecase.BlockData{
				Block: chainindex.Block{
					Height: anyBlockHeight,
					Hash:   "9A362E7D43B0414461698012A41B75FD0E79D0B336ECDA70A8447F1B5D17A5F1",
					Time:   anyBlockTime,
				},
				Activities: []chainindex.Activity{
					{
						BlockHeight:                anyBlockHeight,
						Type:                       chainindex.ACTIVITY_NODEJOIN,
						MaybeTxID:                  primptr.String("fc082211218b0a55e24058243b8fcacf6f23b314f1e1dccd557bd25187826ab2"),
						MaybeFee:                   bignum.MustAtoi("0"),
						MaybeStakingAccountAddress: primptr.String("0xb328a39002ede64c33bb60f1dc43f5df9eb47043"),
						MaybeCouncilNodeMeta: &chainindex.CouncilNode{
							Name:                 "canadacentral_validator_2",
							MaybeSecurityContact: nil,
							PubKeyType:           chainindex.PUBKEY_TYPE_ED25519,
							PubKey:               "rXhu7xhqYBtJftVLKxvKN0XnpyOzxFnUEfAhD1dEF/8=",
							Address:              "34C725CABA703269B3F1D1A907A84DE5FEE96469",
							CreatedAtBlockHeight: anyBlockHeight,
						},
					},
				},
			}))
		})

		It("should parse unjail transaction", func() {
			anyBlockHeight := uint64(3810)
			anyBlockTime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "2020-05-07T11:33:44.402712005Z")
			block := tenderminttypes.Block{
				Height: anyBlockHeight,
				Hash:   "0E370EFD6E4F9B04D03FEFE104A548B7B49ED50F07869F322B2A090B3DEAF604",
				Time:   anyBlockTime,
				Txs: []string{
					"AQECAAAAAAAAAACzKKOQAu3mTDO7YPHcQ/XfnrRwQwBCAAAAAAAAAAAAAYMvVjvs/jDFA9AUIbDe9I3CLkhx68HGk4DgaUMhmUpNElqrTDOcFjPjDclUZsOpghSpTdrtY+BaVQb44fUlKl0=",
				},
				Signatures: nil,
			}
			blockResults := tenderminttypes.BlockResults{
				Height: anyBlockHeight,
				TxsEvents: [][]tenderminttypes.BlockResultsEvent{
					{
						{
							Type: "valid_txs",
							Attributes: []tenderminttypes.BlockResultsEventAttribute{
								{
									Key:   "ZmVl",
									Value: "MC4wMDAwMDAwMA==",
								},
								{
									Key:   "dHhpZA==",
									Value: "NzFiNzNlZDVhYTI3ZDM5NTQ5YmM5YzU3ODIxZjQxNzE2ZTAxMzlkMDQ0ZDgxZGUzNDBjYjMzNDk3ZDY2M2NhZg==",
								},
							},
						},
						{
							Type: "staking_change",
							Attributes: []tenderminttypes.BlockResultsEventAttribute{
								{
									Key:   "c3Rha2luZ19hZGRyZXNz",
									Value: "MHhiMzI4YTM5MDAyZWRlNjRjMzNiYjYwZjFkYzQzZjVkZjllYjQ3MDQz",
								},
								{
									Key:   "c3Rha2luZ19vcHR5cGU=",
									Value: "dW5qYWls",
								},
							},
						},
					},
				},
				BeginBlockEvents: nil,
				ValidatorUpdates: nil,
			}

			actualBlockData := ParseBlockToBlockData(TendermintBlockData{
				Block:        &block,
				BlockResults: &blockResults,
			})
			Expect(*actualBlockData).To(Equal(usecase.BlockData{
				Block: chainindex.Block{
					Height: anyBlockHeight,
					Hash:   "0E370EFD6E4F9B04D03FEFE104A548B7B49ED50F07869F322B2A090B3DEAF604",
					Time:   anyBlockTime,
				},
				Activities: []chainindex.Activity{
					{
						BlockHeight:                anyBlockHeight,
						Type:                       chainindex.ACTIVITY_UNJAIL,
						MaybeTxID:                  primptr.String("71b73ed5aa27d39549bc9c57821f41716e0139d044d81de340cb33497d663caf"),
						MaybeFee:                   bignum.MustAtoi("0"),
						MaybeStakingAccountAddress: primptr.String("0xb328a39002ede64c33bb60f1dc43f5df9eb47043"),
					},
				},
			}))
		})

		It("should parse reward event", func() {
			anyBlockHeight := uint64(2)
			anyBlockTime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "2020-05-07T11:33:44.402712005Z")
			block := tenderminttypes.Block{
				Height:     anyBlockHeight,
				Hash:       "31378879152B877926D756B9063F7F5C188AD49B0F9B65BC316E4EF501832986",
				Time:       anyBlockTime,
				Txs:        nil,
				Signatures: nil,
			}
			blockResults := tenderminttypes.BlockResults{
				Height:    anyBlockHeight,
				TxsEvents: nil,
				BeginBlockEvents: []tenderminttypes.BlockResultsEvent{
					{
						Type: "staking_change",
						Attributes: []tenderminttypes.BlockResultsEventAttribute{
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
						Attributes: []tenderminttypes.BlockResultsEventAttribute{
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
						Attributes: []tenderminttypes.BlockResultsEventAttribute{
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
						Attributes: []tenderminttypes.BlockResultsEventAttribute{
							{
								Key:   "bWludGVk",
								Value: "IjE0NTc2OTcwNjEwMDAwIg==",
							},
						},
					},
				},
				ValidatorUpdates: nil,
			}

			actualBlockData := ParseBlockToBlockData(TendermintBlockData{
				Block:        &block,
				BlockResults: &blockResults,
			})
			Expect(*actualBlockData).To(Equal(usecase.BlockData{
				Block: chainindex.Block{
					Height: anyBlockHeight,
					Hash:   "31378879152B877926D756B9063F7F5C188AD49B0F9B65BC316E4EF501832986",
					Time:   anyBlockTime,
				},
				Activities: []chainindex.Activity{
					{
						BlockHeight:                anyBlockHeight,
						Type:                       chainindex.ACTIVITY_REWARD,
						MaybeEventPosition:         primptr.Uint32(uint32(0)),
						MaybeStakingAccountAddress: primptr.String("0x6dbd5b8fe0dad494465aa7574defba711c184102"),
						MaybeBonded:                bignum.MustAtoi("4858990203333"),
					},
					{
						BlockHeight:                anyBlockHeight,
						Type:                       chainindex.ACTIVITY_REWARD,
						MaybeEventPosition:         primptr.Uint32(uint32(1)),
						MaybeStakingAccountAddress: primptr.String("0x6fc1e3124a7ed07f3710378b68f7046c7300179d"),
						MaybeBonded:                bignum.MustAtoi("4858990203333"),
					},
					{
						BlockHeight:                anyBlockHeight,
						Type:                       chainindex.ACTIVITY_REWARD,
						MaybeEventPosition:         primptr.Uint32(uint32(2)),
						MaybeStakingAccountAddress: primptr.String("0xb8c6886da09e12db8aebfc8108c67ce2ba086ac6"),
						MaybeBonded:                bignum.MustAtoi("4858990203333"),
					},
				},
				Reward: &chainindex.BlockReward{
					BlockHeight: anyBlockHeight,
					Minted:      bignum.MustAtoi("14576970610000"),
				},
			}))
		})

		It("should parse slash and jail event by byzantine fault", func() {
			anyBlockHeight := uint64(3510)
			anyBlockTime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "2020-05-07T10:00:25.582998694Z")
			block := tenderminttypes.Block{
				Height:     anyBlockHeight,
				Hash:       "2CD22EA622D190B9ABCAB797E2F60F6F4FCFC19CC0A67642E5E7856CEAD78163",
				Time:       anyBlockTime,
				Txs:        nil,
				Signatures: nil,
			}
			blockResults := tenderminttypes.BlockResults{
				Height:    anyBlockHeight,
				TxsEvents: nil,
				BeginBlockEvents: []tenderminttypes.BlockResultsEvent{
					{
						Type: "staking_change",
						Attributes: []tenderminttypes.BlockResultsEventAttribute{
							{
								Key:   "c3Rha2luZ19hZGRyZXNz",
								Value: "MHhiMzI4YTM5MDAyZWRlNjRjMzNiYjYwZjFkYzQzZjVkZjllYjQ3MDQz",
							},
							{
								Key:   "c3Rha2luZ19vcHR5cGU=",
								Value: "c2xhc2g=",
							},
							{
								Key:   "c3Rha2luZ19kaWZm",
								Value: "W3sia2V5IjoiQm9uZGVkIiwidmFsdWUiOiItMTA5NDM2NzkxMjEyMjMxOSJ9LHsia2V5IjoiVW5ib25kZWQiLCJ2YWx1ZSI6Ii0wIn1d",
							},
							{
								Key:   "c3Rha2luZ19vcHJlYXNvbg==",
								Value: "Qnl6YW50aW5lRmF1bHQ=",
							},
						},
					},
					{
						Type: "staking_change",
						Attributes: []tenderminttypes.BlockResultsEventAttribute{
							{
								Key:   "c3Rha2luZ19hZGRyZXNz",
								Value: "MHhiMzI4YTM5MDAyZWRlNjRjMzNiYjYwZjFkYzQzZjVkZjllYjQ3MDQz",
							},
							{
								Key:   "c3Rha2luZ19vcHR5cGU=",
								Value: "amFpbA==",
							},
							{
								Key:   "c3Rha2luZ19kaWZm",
								Value: "W3sia2V5IjoiSmFpbGVkVW50aWwiLCJ2YWx1ZSI6MTU4ODg1MTAyNX1d",
							},
							{
								Key:   "c3Rha2luZ19vcHJlYXNvbg==",
								Value: "Qnl6YW50aW5lRmF1bHQ=",
							},
						},
					},
				},
				ValidatorUpdates: nil,
			}

			actualBlockData := ParseBlockToBlockData(TendermintBlockData{
				Block:        &block,
				BlockResults: &blockResults,
			})
			Expect(*actualBlockData).To(Equal(usecase.BlockData{
				Block: chainindex.Block{
					Height: anyBlockHeight,
					Hash:   "2CD22EA622D190B9ABCAB797E2F60F6F4FCFC19CC0A67642E5E7856CEAD78163",
					Time:   anyBlockTime,
				},
				Activities: []chainindex.Activity{
					{
						BlockHeight:                anyBlockHeight,
						Type:                       chainindex.ACTIVITY_SLASH,
						MaybeEventPosition:         primptr.Uint32(uint32(0)),
						MaybeStakingAccountAddress: primptr.String("0xb328a39002ede64c33bb60f1dc43f5df9eb47043"),
						MaybeBonded:                bignum.MustAtoi("-1094367912122319"),
						MaybeUnbonded:              bignum.MustAtoi("-0"),
						MaybePunishmentKind:        primptr.Uint8(chainindex.PUNISHMENT_KIND_BYZANTINE_FAULT),
					},
					{
						BlockHeight:                anyBlockHeight,
						Type:                       chainindex.ACTIVITY_JAIL,
						MaybeEventPosition:         primptr.Uint32(uint32(1)),
						MaybeStakingAccountAddress: primptr.String("0xb328a39002ede64c33bb60f1dc43f5df9eb47043"),
						MaybePunishmentKind:        primptr.Uint8(chainindex.PUNISHMENT_KIND_BYZANTINE_FAULT),
						MaybeJailedUntil:           primptr.Time(time.Unix(1588851025, 0).UTC()),
					},
				},
			}))
		})

		It("should parse nonlive slash and council node kicked", func() {
			anyBlockHeight := uint64(600)
			anyBlockTime, _ := time.Parse("2006-01-02T15:04:05.000000000Z", "2020-05-05T19:57:10.527286067Z")
			block := tenderminttypes.Block{
				Height:     anyBlockHeight,
				Hash:       "B5B2860C8BE4CB577CBFA9D6007C7F968A238346B3568F1D1982B01D2E3A52B8",
				Time:       anyBlockTime,
				AppHash:    "97CDEB476C160D9E2988A7221B969DA7FDB4F31AD8386D8230E8D40A92FC096E",
				Txs:        nil,
				Signatures: nil,
			}
			blockResults := tenderminttypes.BlockResults{
				Height:    anyBlockHeight,
				TxsEvents: nil,
				BeginBlockEvents: []tenderminttypes.BlockResultsEvent{
					{
						Type: "staking_change",
						Attributes: []tenderminttypes.BlockResultsEventAttribute{
							{
								Key:   "c3Rha2luZ19hZGRyZXNz",
								Value: "MHhiMzI4YTM5MDAyZWRlNjRjMzNiYjYwZjFkYzQzZjVkZjllYjQ3MDQz",
							},
							{
								Key:   "c3Rha2luZ19vcHR5cGU=",
								Value: "c2xhc2g=",
							},
							{
								Key:   "c3Rha2luZ19kaWZm",
								Value: "W3sia2V5IjoiQm9uZGVkIiwidmFsdWUiOiItNjAwMDAwMDAwMDAwMDAwIn0seyJrZXkiOiJVbmJvbmRlZCIsInZhbHVlIjoiLTAifV0=",
							},
							{
								Key:   "c3Rha2luZ19vcHJlYXNvbg==",
								Value: "Tm9uTGl2ZQ==",
							},
						},
					},
				},
				ValidatorUpdates: []tenderminttypes.BlockResultsValidator{
					{
						PubKey: tenderminttypes.BlockResultsValidatorPubKey{
							Type:    "ed25519",
							PubKey:  "rXhu7xhqYBtJftVLKxvKN0XnpyOzxFnUEfAhD1dEF/8=",
							Address: "34C725CABA703269B3F1D1A907A84DE5FEE96469",
						},
						Power: nil,
					},
					{
						PubKey: tenderminttypes.BlockResultsValidatorPubKey{
							Type:    "ed25519",
							PubKey:  "QMegiWt9+5K1b1ZVd7zOJZxhTnbAtWzvGhViiElAlaw=",
							Address: "7570B2D23A4C7B638BEFE02EB4FC7927BFDED6B7",
						},
						Power: nil,
					},
					{
						PubKey: tenderminttypes.BlockResultsValidatorPubKey{
							Type:    "ed25519",
							PubKey:  "tDLheZJwsA8oYEwarR6/X+zAmNKMLHTVkh/fvcLqcwA=",
							Address: "D527DAECDE0501CF2E785A8DC0D9F4A64760F0BB",
						},
						Power: primptr.String("60000000"),
					},
				},
			}

			actualBlockData := ParseBlockToBlockData(TendermintBlockData{
				Block:        &block,
				BlockResults: &blockResults,
			})
			Expect(*actualBlockData).To(Equal(usecase.BlockData{
				Block: chainindex.Block{
					Height:  anyBlockHeight,
					Hash:    "B5B2860C8BE4CB577CBFA9D6007C7F968A238346B3568F1D1982B01D2E3A52B8",
					Time:    anyBlockTime,
					AppHash: "97CDEB476C160D9E2988A7221B969DA7FDB4F31AD8386D8230E8D40A92FC096E",
				},
				Activities: []chainindex.Activity{
					{
						BlockHeight:                anyBlockHeight,
						Type:                       chainindex.ACTIVITY_SLASH,
						MaybeEventPosition:         primptr.Uint32(uint32(0)),
						MaybeStakingAccountAddress: primptr.String("0xb328a39002ede64c33bb60f1dc43f5df9eb47043"),
						MaybeBonded:                bignum.MustAtoi("-600000000000000"),
						MaybeUnbonded:              bignum.MustAtoi("-0"),
						MaybePunishmentKind:        primptr.Uint8(chainindex.PUNISHMENT_KIND_NON_LIVE),
					},
				},
				CouncilNodeUpdates: []chainindex.CouncilNodeUpdate{
					{
						Address: "34C725CABA703269B3F1D1A907A84DE5FEE96469",
						Type:    chainindex.COUNCIL_NODE_UPDATE_TYPE_LEFT,
					},
					{
						Address: "7570B2D23A4C7B638BEFE02EB4FC7927BFDED6B7",
						Type:    chainindex.COUNCIL_NODE_UPDATE_TYPE_LEFT,
					},
				},
			}))
		})
	})
})

func sampleGenesis() tenderminttypes.Genesis {
	genesisTime, _ := time.Parse("2006-01-02T15:04:05.000000Z", "2020-05-01T12:09:01.568951Z")
	return tenderminttypes.Genesis{
		GenesisTime: genesisTime,
		ChainID:     "testnet-thaler-crypto-com-chain-42",
		AppHash:     "F62DDB49D7EB8ED0883C735A0FB7DE7F2A3FA322FCD2AA832F452A62B38607D5",
		AppState: tenderminttypes.GenesisAppState{
			CouncilNodes: []tenderminttypes.GenesisCouncilNode{
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
			Distribution: []tenderminttypes.GenesisDistribution{
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
	}
}

func sampleGenesisBlock() tenderminttypes.Block {
	blockTime, _ := time.Parse("2006-01-02T15:04:05.000000Z", "2020-05-01T12:09:01.568951Z")
	return tenderminttypes.Block{
		Height:         uint64(1),
		Hash:           "BD3D9499EF527035BAE14E7F4932BD1778C1E90141A7A76A1FE67B7ECCC2663E",
		Time:           blockTime,
		AppHash:        "F62DDB49D7EB8ED0883C735A0FB7DE7F2A3FA322FCD2AA832F452A62B38607D5",
		PropserAddress: "D527DAECDE0501CF2E785A8DC0D9F4A64760F0BB",
		Txs:            nil,
		Signatures:     nil,
	}
}

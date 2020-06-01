package adapter_test

import (
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/crypto-com/chainindex"
	"github.com/crypto-com/chainindex/adapter"
	. "github.com/crypto-com/chainindex/adapter/test/fake"
	. "github.com/crypto-com/chainindex/adapter/test/mock"
	"github.com/crypto-com/chainindex/internal/bignum"
	"github.com/crypto-com/chainindex/internal/primptr"
	. "github.com/crypto-com/chainindex/test/factory"
	. "github.com/crypto-com/chainindex/test/mock"
)

const (
	SQL_ACTIVITY_INSERT                                     = "INSERT INTO activities (block_height,type,txid,event_position,fee,inputs,output_count,staking_account_address,staking_account_nonce,bonded,unbonded,unbonded_from,joined_council_node,joined_council_node_id,affected_council_node,affected_council_node_id,jailed_until,punishment_kind) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	SQL_TRANSFER_OUTPUT_UPDATE                              = "UPDATE transaction_outputs SET spent_at_txid = ? WHERE txid = ? AND index = ?"
	SQL_TRANSFER_OUTPUT_INSERT                              = "INSERT INTO transaction_outputs (txid,index) VALUES (?,?)"
	SQL_STAKING_ACCOUNT_INSERT                              = "INSERT INTO staking_accounts (address,nonce,bonded,unbonded,unbonded_from,jailed_until,punishment_kind,current_council_node_id) VALUES (?,?,?,?,?,?,?,?)"
	SQL_STAKING_ACCOUNT_SELECT                              = "SELECT address,nonce,bonded,unbonded,unbonded_from,punishment_kind,jailed_until,current_council_node_id FROM staking_accounts WHERE address = ?"
	SQL_STAKING_ACCOUNT_UPDATE                              = "UPDATE staking_accounts SET bonded = ?, current_council_node_id = ?, jailed_until = ?, nonce = ?, punishment_kind = ?, unbonded = ?, unbonded_from = ? WHERE address = ?"
	SQL_STAKING_ACCOUNT_CURRENT_COUNCIL_NODE_SELECT         = "SELECT c.id, c.name, c.security_contact, c.pubkey_type, c.pubkey, c.address, c.created_at_block_height, c.last_left_at_block_height FROM council_nodes c JOIN staking_accounts sa ON c.id = sa.current_council_node_id WHERE c.last_left_at_block_height IS NULL AND sa.address = ?"
	SQL_COUNCIL_NODE_INSERT                                 = "INSERT INTO council_nodes (name,security_contact,pubkey_type,pubkey,address,created_at_block_height,last_left_at_block_height) VALUES (?,?,?,?,?,?,?) RETURNING id"
	SQL_COUNCIL_NODE_BY_ADDRESS_SELECT                      = "SELECT id, name, security_contact, pubkey_type, pubkey, address, created_at_block_height, last_left_at_block_height FROM council_nodes WHERE address = ? ORDER BY id DESC"
	SQL_COUNCIL_NODE_BY_STAKING_ACCOUNT_ADDRESS_SELECT      = "SELECT c.id, c.name, c.security_contact, c.pubkey_type, c.pubkey, c.address, c.created_at_block_height, c.last_left_at_block_height FROM council_nodes c JOIN activities a ON c.id = a.joined_council_node_id WHERE a.type IN ('genesis', 'nodejoin') AND a.staking_account_address = ? ORDER BY a.id DESC"
	SQL_COUNCIL_NODE_CLEAR_LAST_LEFT_AT_BLOCK_HEIGHT_UPDATE = "UPDATE council_nodes SET last_left_at_block_height = ? WHERE id = ?"
)

var _ = Describe("BlockActivityData", func() {
	Describe("InsertGenesisActivity", func() {
		It("should insert new staking account with the bonded", func() {
			inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
				sq.StatementBuilder,
				new(PrimRDbTypeConv),
			)

			anyGenesisActivity := RandomGenesisActivity()
			anyGenesisActivity.MaybeBonded = bignum.Int10()
			anyGenesisActivity.MaybeUnbonded = nil
			anyGenesisActivity.MaybeCouncilNodeMeta = nil

			tx := new(MockRDbTx)

			execResult := new(MockRDbExecResult)
			execResult.On("RowsAffected").Return(int64(1))
			OnTxInsertAnyActivity(tx).Return(execResult, nil)

			tx.On("Exec",
				SQL_STAKING_ACCOUNT_INSERT,
				*anyGenesisActivity.MaybeStakingAccountAddress, // address
				uint64(0), // nonce
				primptr.String(anyGenesisActivity.MaybeBonded.String()), // bonded
				primptr.String("0"), // unbonded
				(*time.Time)(nil),   // unboned_from
				(*time.Time)(nil),   // jailed_until
				(*string)(nil),      // punishment_kind
				(*uint64)(nil),      // current_council_node_id
			).Once().Return(execResult, nil)

			err := inserter.InsertGenesisActivity(tx, &anyGenesisActivity)
			Expect(err).To(BeNil())
			tx.AssertExpectations(GinkgoT())
		})

		It("should insert new staking account with the unbonded amount", func() {
			inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
				sq.StatementBuilder,
				new(PrimRDbTypeConv),
			)

			anyGenesisActivity := RandomGenesisActivity()
			anyGenesisActivity.MaybeBonded = nil
			anyGenesisActivity.MaybeUnbonded = bignum.Int10()
			anyGenesisActivity.MaybeCouncilNodeMeta = nil

			tx := new(MockRDbTx)

			execResult := new(MockRDbExecResult)
			execResult.On("RowsAffected").Return(int64(1))
			OnTxInsertAnyActivity(tx).Return(execResult, nil)

			tx.On("Exec",
				SQL_STAKING_ACCOUNT_INSERT,
				*anyGenesisActivity.MaybeStakingAccountAddress, // address
				uint64(0),           // nonce
				primptr.String("0"), // bonded
				primptr.String(anyGenesisActivity.MaybeUnbonded.String()), // unbonded
				(*time.Time)(nil), // unboned_from
				(*time.Time)(nil), // jailed_until
				(*string)(nil),    // punishment_kind
				(*uint64)(nil),    // current_council_node_id
			).Once().Return(execResult, nil)

			err := inserter.InsertGenesisActivity(tx, &anyGenesisActivity)
			Expect(err).To(BeNil())
			tx.AssertExpectations(GinkgoT())
		})

		It("should insert council node and insert staking account with the council node id", func() {
			var err error

			inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
				sq.StatementBuilder,
				new(PrimRDbTypeConv),
			)

			anyGenesisActivity := RandomGenesisActivity()
			anyGenesisActivity.MaybeBonded = bignum.Int10()
			anyGenesisActivity.MaybeUnbonded = nil
			anyGenesisActivity.MaybeCouncilNodeMeta.CreatedAtBlockHeight = anyGenesisActivity.BlockHeight
			anyGenesisActivity.MaybeCouncilNodeMeta.MaybeLastLeftAtBlockHeight = nil

			tx := new(MockRDbTx)

			execResult := new(MockRDbExecResult)
			execResult.On("RowsAffected").Return(int64(1))
			OnTxInsertAnyActivity(tx).Return(execResult, nil)

			anyCouncilNodeId := uint64(1)

			rowResult := new(MockRDbRowResult)
			rowResult.On("Scan", mock.MatchedBy(func(councilNodeId *uint64) bool {
				*councilNodeId = anyCouncilNodeId
				return true
			})).Return(nil)
			tx.On("QueryRow",
				SQL_COUNCIL_NODE_INSERT,
				anyGenesisActivity.MaybeCouncilNodeMeta.Name,
				anyGenesisActivity.MaybeCouncilNodeMeta.MaybeSecurityContact,
				adapter.PubKeyTypeToString(anyGenesisActivity.MaybeCouncilNodeMeta.PubKeyType),
				anyGenesisActivity.MaybeCouncilNodeMeta.PubKey,
				anyGenesisActivity.MaybeCouncilNodeMeta.Address,
				anyGenesisActivity.BlockHeight, // created_at_block_height
				(*uint64)(nil),                 // last_left_at_block_height
			).Once().Return(rowResult)

			tx.On("Exec",
				SQL_STAKING_ACCOUNT_INSERT,
				*anyGenesisActivity.MaybeStakingAccountAddress, // address
				uint64(0), // nonce
				primptr.String(anyGenesisActivity.MaybeBonded.String()), // bonded
				primptr.String("0"), // unbonded
				(*time.Time)(nil),   // unboned_from
				(*time.Time)(nil),   // jailed_until
				(*string)(nil),      // punishment_kind
				&anyCouncilNodeId,   // current_council_node_id
			).Once().Return(execResult, nil)

			err = inserter.InsertGenesisActivity(tx, &anyGenesisActivity)
			Expect(err).To(BeNil())
			tx.AssertExpectations(GinkgoT())
		})

		It("should insert activity into the table", func() {
			var err error

			inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
				sq.StatementBuilder,
				new(PrimRDbTypeConv),
			)

			anyGenesisActivity := RandomGenesisActivity()
			anyGenesisActivity.MaybeBonded = bignum.Int10()
			anyGenesisActivity.MaybeUnbonded = nil
			anyGenesisActivity.MaybeCouncilNodeMeta.CreatedAtBlockHeight = anyGenesisActivity.BlockHeight
			anyGenesisActivity.MaybeCouncilNodeMeta.MaybeLastLeftAtBlockHeight = nil

			tx := new(MockRDbTx)

			anyCouncilNodeId := uint64(1)
			OnTxQueryInsertAnyCouncilNodeRowReturnCouncilNodeId(tx, anyCouncilNodeId)

			execResult := new(MockRDbExecResult)
			execResult.On("RowsAffected").Return(int64(1))
			OnTxInsertAnyStakingAccount(tx).Return(execResult, nil)

			joinedCouncilNode := *anyGenesisActivity.MaybeCouncilNodeMeta
			joinedCouncilNode.Id = &anyCouncilNodeId
			tx.On("Exec",
				SQL_ACTIVITY_INSERT,
				anyGenesisActivity.BlockHeight, // block_height
				"genesis",                      // type
				(*string)(nil),                 // txid
				(*uint32)(nil),                 // event_position
				(*string)(nil),                 // fee
				(*string)(nil),                 // inputs
				(*uint32)(nil),                 // output_count
				anyGenesisActivity.MaybeStakingAccountAddress,           // staking_account_address
				primptr.Uint64(uint64(0)),                               // staking_account_nonce
				primptr.String(anyGenesisActivity.MaybeBonded.String()), // bonded
				(*string)(nil), // unbonded
				nil,            // unbdoned_from
				primptr.String(JsonMustMarshal(adapter.CouncilNodeToRDbCouncilNodeRow(
					&joinedCouncilNode,
				))), // joined_council_node
				primptr.Uint64(anyCouncilNodeId), // joined_council_node_id
				(*string)(nil),                   // affected_council_node
				(*uint64)(nil),                   // affected_council_node_id
				nil,                              // jailed_until
				(*string)(nil),                   // punishment_kind
			).Once().Return(execResult, nil)

			err = inserter.InsertGenesisActivity(tx, &anyGenesisActivity)
			Expect(err).To(BeNil())
			tx.AssertExpectations(GinkgoT())
		})
	})

	Describe("InsertTransferTransaction", func() {
		It("should return Error when insertion failed", func() {
			inserter := adapter.NewDefaultRDbBlockActivityDataRepo(sq.StatementBuilder, new(PrimRDbTypeConv))

			tx := new(MockRDbTx)
			anyError := errors.New("any error")
			OnTxInsertAnyActivity(tx).Return(nil, anyError)

			anyTransferActivity := RandomTransferActivity()
			err := inserter.InsertTransferTransaction(tx, &anyTransferActivity)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("any error"))
		})

		Context("When activity is valid", func() {
			var anyTransferActivity chainindex.Activity
			var tx *MockRDbTx
			BeforeEach(func() {
				anyTransferActivity = RandomTransferActivity()
				anyTransferActivity.MaybeTxInputs = RandomTxInputsPtrOfSize(3)
				anyTransferActivity.MaybeOutputCount = primptr.Uint32(3)

				tx = new(MockRDbTx)
			})

			It("should insert transfer activity into table", func() {
				inserter := adapter.NewDefaultRDbBlockActivityDataRepo(sq.StatementBuilder, new(PrimRDbTypeConv))

				execResult := new(MockRDbExecResult)
				execResult.On("RowsAffected").Return(int64(1))
				OnTxUpdateAnyTransaferInput(tx).Return(execResult, nil)
				OnTxInsertAnyTransaferOutput(tx).Return(execResult, nil)

				tx.On("Exec",
					SQL_ACTIVITY_INSERT,
					anyTransferActivity.BlockHeight,              // block_height
					"transfer",                                   // type
					anyTransferActivity.MaybeTxID,                // txid
					(*uint32)(nil),                               // event_position
					bignum.OptItoa(anyTransferActivity.MaybeFee), // fee
					primptr.String(JsonMustMarshal(
						adapter.TxInputsToRDbTransferInputs(anyTransferActivity.MaybeTxInputs),
					)), // inputs
					anyTransferActivity.MaybeOutputCount, // output_count
					(*string)(nil),                       // staking_account_address
					(*uint64)(nil),                       // staking_account_nonce
					(*string)(nil),                       // bonded
					(*string)(nil),                       // unbonded
					nil,                                  // unboned_from
					(*string)(nil),                       // joined_council_node
					(*uint64)(nil),                       // joined_council_node_id
					(*string)(nil),                       // affected_council_node
					(*uint64)(nil),                       // affected_council_node_id
					nil,                                  // jailed_until
					(*string)(nil),                       // punishment_kind
				).Once().Return(execResult, nil)

				err := inserter.InsertTransferTransaction(tx, &anyTransferActivity)
				Expect(err).To(BeNil())
				tx.AssertExpectations(GinkgoT())
			})

			It("should update previous transaction outputs to be spent by this block", func() {
				inserter := adapter.NewDefaultRDbBlockActivityDataRepo(sq.StatementBuilder, new(PrimRDbTypeConv))

				execResult := new(MockRDbExecResult)
				execResult.On("RowsAffected").Return(int64(1))

				OnTxInsertAnyActivity(tx).Return(execResult, nil)
				OnTxInsertAnyTransaferOutput(tx).Return(execResult, nil)

				OnTxUpdateTransaferInput(tx,
					anyTransferActivity.MaybeTxID, anyTransferActivity.MaybeTxInputs[0],
				).Once().Return(execResult, nil)
				OnTxUpdateTransaferInput(tx,
					anyTransferActivity.MaybeTxID, anyTransferActivity.MaybeTxInputs[1],
				).Once().Return(execResult, nil)
				OnTxUpdateTransaferInput(tx,
					anyTransferActivity.MaybeTxID, anyTransferActivity.MaybeTxInputs[2],
				).Once().Return(execResult, nil)

				err := inserter.InsertTransferTransaction(tx, &anyTransferActivity)
				Expect(err).To(BeNil())
				tx.AssertExpectations(GinkgoT())
			})

			It("should insert all transaction outputs into table", func() {
				inserter := adapter.NewDefaultRDbBlockActivityDataRepo(sq.StatementBuilder, new(PrimRDbTypeConv))

				execResult := new(MockRDbExecResult)
				execResult.On("RowsAffected").Return(int64(1))
				OnTxUpdateAnyTransaferInput(tx).Return(execResult, nil)
				OnTxInsertAnyActivity(tx).Return(execResult, nil)

				OnTxInsertTransferOutput(tx,
					anyTransferActivity.MaybeTxID, 0,
				).Once().Return(execResult, nil)
				OnTxInsertTransferOutput(tx,
					anyTransferActivity.MaybeTxID, 1,
				).Once().Return(execResult, nil)
				OnTxInsertTransferOutput(tx,
					anyTransferActivity.MaybeTxID, 2,
				).Once().Return(execResult, nil)

				err := inserter.InsertTransferTransaction(tx, &anyTransferActivity)
				Expect(err).To(BeNil())
				tx.AssertExpectations(GinkgoT())
			})
		})
	})

	Describe("InsertDepositTransaction", func() {
		It("should insert new staking account when the address does not exist before", func() {
			inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
				sq.StatementBuilder,
				new(PrimRDbTypeConv),
			)

			anyDepositActivity := RandomDepositActivity()

			tx := new(MockRDbTx)

			WhenStakingAccountDoesNotExist(tx, anyDepositActivity.MaybeStakingAccountAddress)
			WhenStakingAccountHasNotJoinedCouncilNode(tx, *anyDepositActivity.MaybeStakingAccountAddress)

			execResult := new(MockRDbExecResult)
			execResult.On("RowsAffected").Return(int64(1))
			OnTxInsertAnyActivity(tx).Return(execResult, nil)

			tx.On("Exec",
				SQL_STAKING_ACCOUNT_INSERT,
				*anyDepositActivity.MaybeStakingAccountAddress, // address
				uint64(0), // nonce
				primptr.String(anyDepositActivity.MaybeBonded.String()), // bonded
				primptr.String("0"), // unbonded
				(*time.Time)(nil),   // unboned_from
				(*time.Time)(nil),   // jailed_until
				(*string)(nil),      // punishment_kind
				(*uint64)(nil),      // current_council_node_id
			).Once().Return(execResult, nil)

			err := inserter.InsertDepositTransaction(tx, &anyDepositActivity)
			Expect(err).To(BeNil())
			tx.AssertExpectations(GinkgoT())
		})

		Context("When staking account exists", func() {
			var anyDepositActivity chainindex.Activity
			var tx *MockRDbTx

			BeforeEach(func() {
				anyDepositActivity = RandomDepositActivity()

				tx = new(MockRDbTx)

				stakingAccountRow := adapter.RDbStakingAccountRow{
					Address:              *anyDepositActivity.MaybeStakingAccountAddress,
					Nonce:                uint64(0),
					Bonded:               bignum.Int0(),
					Unbonded:             bignum.Int0(),
					UnbondedFrom:         nil,
					PunishmentKind:       nil,
					JailedUntil:          nil,
					CurrentCouncilNodeId: nil,
				}
				WhenStakingAccountExist(tx, &stakingAccountRow)
			})

			Context("When staking account has not joined council node", func() {
				BeforeEach(func() {
					WhenStakingAccountHasNotJoinedCouncilNode(tx, *anyDepositActivity.MaybeStakingAccountAddress)
				})

				It("should credit staking account bonded balance when the address already existed", func() {
					inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
						sq.StatementBuilder,
						new(PrimRDbTypeConv),
					)

					execResult := new(MockRDbExecResult)
					execResult.On("RowsAffected").Return(int64(1))
					OnTxInsertAnyActivity(tx).Return(execResult, nil)

					tx.On("Exec",
						SQL_STAKING_ACCOUNT_UPDATE,
						primptr.String(anyDepositActivity.MaybeBonded.String()), // bonded
						(*uint64)(nil),      // council_node_id
						nil,                 // jailed_until
						uint64(0),           // nonce
						(*string)(nil),      // punishment_kind
						primptr.String("0"), // unbonded
						nil,                 // ubonded_from
						*anyDepositActivity.MaybeStakingAccountAddress, // WHERE address = ?
					).Once().Return(execResult, nil)

					err := inserter.InsertDepositTransaction(tx, &anyDepositActivity)
					Expect(err).To(BeNil())
					tx.AssertExpectations(GinkgoT())
				})

				It("should insert activity into the table", func() {
					inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
						sq.StatementBuilder,
						new(PrimRDbTypeConv),
					)

					execResult := new(MockRDbExecResult)
					execResult.On("RowsAffected").Return(int64(1))
					OnTxUpdateAnyStakingAccount(tx).Return(execResult, nil)

					tx.On("Exec",
						SQL_ACTIVITY_INSERT,
						anyDepositActivity.BlockHeight,                // block_height
						"deposit",                                     // type
						anyDepositActivity.MaybeTxID,                  // txid
						(*uint32)(nil),                                // event_position
						bignum.OptItoa(anyDepositActivity.MaybeFee),   // fee
						(*string)(nil),                                // inputs
						(*uint32)(nil),                                // output_count
						anyDepositActivity.MaybeStakingAccountAddress, // staking_account_address
						primptr.Uint64(uint64(0)),                     // staking_account_nonce
						primptr.String(anyDepositActivity.MaybeBonded.String()), // bonded
						(*string)(nil), // unbonded
						nil,            // unbonded_from
						(*string)(nil), // joined_council_node
						(*uint64)(nil), // joined_council_node_id
						(*string)(nil), // affected_council_node
						(*uint64)(nil), // affected_council_node_id
						nil,            // jailed_until
						(*string)(nil), // punishment_kind
					).Once().Return(execResult, nil)

					err := inserter.InsertDepositTransaction(tx, &anyDepositActivity)
					Expect(err).To(BeNil())
					tx.AssertExpectations(GinkgoT())
				})
			})

			It("should insert affected_council_node_id when staking account has joined council node", func() {
				inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
					sq.StatementBuilder,
					new(PrimRDbTypeConv),
				)

				prevCouncilNode := RandomCouncilNode()
				prevCouncilNode.Id = primptr.Uint64(uint64(6))
				WhenStakingAccountHasJoinedCouncilNode(tx, *anyDepositActivity.MaybeStakingAccountAddress, prevCouncilNode)

				execResult := new(MockRDbExecResult)
				execResult.On("RowsAffected").Return(int64(1))
				OnTxUpdateAnyStakingAccount(tx).Return(execResult, nil)

				tx.On("Exec",
					SQL_ACTIVITY_INSERT,
					anyDepositActivity.BlockHeight,                // block_height
					"deposit",                                     // type
					anyDepositActivity.MaybeTxID,                  // txid
					(*uint32)(nil),                                // event_position
					bignum.OptItoa(anyDepositActivity.MaybeFee),   // fee
					(*string)(nil),                                // inputs
					(*uint32)(nil),                                // output_count
					anyDepositActivity.MaybeStakingAccountAddress, // staking_account_address
					primptr.Uint64(uint64(0)),                     // staking_account_nonce
					primptr.String(anyDepositActivity.MaybeBonded.String()), // bonded
					(*string)(nil), // unbonded
					nil,            // unbonded_from
					(*string)(nil), // joined_council_node
					(*uint64)(nil), // joined_council_node_id
					primptr.String(JsonMustMarshal(
						adapter.CouncilNodeToRDbCouncilNodeRow(&prevCouncilNode),
					)), // affected_council_node
					primptr.Uint64(*prevCouncilNode.Id), // affected_council_node_id
					nil,                                 // jailed_until
					(*string)(nil),                      // punishment_kind
				).Once().Return(execResult, nil)

				err := inserter.InsertDepositTransaction(tx, &anyDepositActivity)
				Expect(err).To(BeNil())
				tx.AssertExpectations(GinkgoT())
			})
		})
	})

	Describe("InsertUnbondTransaction", func() {
		It("should return Error when the address does not exist before", func() {
			inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
				sq.StatementBuilder,
				new(PrimRDbTypeConv),
			)

			anyUnbondActivity := RandomUnbondActivity()

			tx := new(MockRDbTx)

			WhenStakingAccountDoesNotExist(tx, anyUnbondActivity.MaybeStakingAccountAddress)
			WhenStakingAccountHasNotJoinedCouncilNode(tx, *anyUnbondActivity.MaybeStakingAccountAddress)

			Expect(func() {
				_ = inserter.InsertUnbondTransaction(tx, &anyUnbondActivity)
			}).To(Panic())
		})

		Context("When staking account exist", func() {
			var anyUnbondActivity chainindex.Activity
			var tx *MockRDbTx

			BeforeEach(func() {
				anyUnbondActivity = RandomUnbondActivity()
				anyUnbondActivity.MaybeBonded = bignum.IntN10()
				anyUnbondActivity.MaybeUnbonded = bignum.Int10()

				tx = new(MockRDbTx)

				stakingAccountRow := adapter.RDbStakingAccountRow{
					Address:              *anyUnbondActivity.MaybeStakingAccountAddress,
					Nonce:                uint64(1),
					Bonded:               bignum.Int10(),
					Unbonded:             bignum.Int0(),
					UnbondedFrom:         nil,
					PunishmentKind:       nil,
					JailedUntil:          nil,
					CurrentCouncilNodeId: nil,
				}
				WhenStakingAccountExist(tx, &stakingAccountRow)
			})

			Context("When staking account has not joined any council node", func() {
				BeforeEach(func() {
					WhenStakingAccountHasNotJoinedCouncilNode(tx, *anyUnbondActivity.MaybeStakingAccountAddress)
				})

				It("should deduct staking address bonded balance and credit unbonded balance", func() {
					inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
						sq.StatementBuilder,
						new(PrimRDbTypeConv),
					)

					execResult := new(MockRDbExecResult)
					execResult.On("RowsAffected").Return(int64(1))
					OnTxInsertAnyActivity(tx).Return(execResult, nil)

					tx.On("Exec",
						SQL_STAKING_ACCOUNT_UPDATE,
						primptr.String("0"),  // bonded
						(*uint64)(nil),       // council_node_id
						nil,                  // jailed_until
						uint64(2),            // nonce
						(*string)(nil),       // punishment_kind
						primptr.String("10"), // unbonded
						nil,                  // unbonded_from
						*anyUnbondActivity.MaybeStakingAccountAddress, // WHERE address = ?
					).Once().Return(execResult, nil)

					err := inserter.InsertUnbondTransaction(tx, &anyUnbondActivity)
					Expect(err).To(BeNil())
					tx.AssertExpectations(GinkgoT())
				})

				It("should insert activity into the table", func() {
					inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
						sq.StatementBuilder,
						new(PrimRDbTypeConv),
					)

					execResult := new(MockRDbExecResult)
					execResult.On("RowsAffected").Return(int64(1))

					OnTxUpdateAnyStakingAccount(tx).Return(execResult, nil)

					tx.On("Exec",
						SQL_ACTIVITY_INSERT,
						anyUnbondActivity.BlockHeight,                // block_height
						"unbond",                                     // type
						anyUnbondActivity.MaybeTxID,                  // txid
						(*uint32)(nil),                               // event_position
						bignum.OptItoa(anyUnbondActivity.MaybeFee),   // fee
						(*string)(nil),                               // inputs
						(*uint32)(nil),                               // output_count
						anyUnbondActivity.MaybeStakingAccountAddress, // staking_account_address
						primptr.Uint64(uint64(2)),                    // staking_account_nonce
						primptr.String(anyUnbondActivity.MaybeBonded.String()),   // bonded
						primptr.String(anyUnbondActivity.MaybeUnbonded.String()), // unbonded
						nil,            // unbonded_from
						(*string)(nil), // joined_council_node
						(*uint64)(nil), // joined_council_node_id
						(*string)(nil), // affected_council_node
						(*uint64)(nil), // affected_council_node_id
						nil,            // jailed_until
						(*string)(nil), // punishment_kind
					).Once().Return(execResult, nil)

					err := inserter.InsertUnbondTransaction(tx, &anyUnbondActivity)
					Expect(err).To(BeNil())
					tx.AssertExpectations(GinkgoT())
				})
			})

			It("should insert activity into the table when staking account has joined council node", func() {
				inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
					sq.StatementBuilder,
					new(PrimRDbTypeConv),
				)

				execResult := new(MockRDbExecResult)
				execResult.On("RowsAffected").Return(int64(1))

				OnTxUpdateAnyStakingAccount(tx).Return(execResult, nil)

				prevCouncilNode := RandomCouncilNode()
				prevCouncilNode.Id = primptr.Uint64(uint64(6))
				WhenStakingAccountHasJoinedCouncilNode(tx, *anyUnbondActivity.MaybeStakingAccountAddress, prevCouncilNode)

				tx.On("Exec",
					SQL_ACTIVITY_INSERT,
					anyUnbondActivity.BlockHeight,                // block_height
					"unbond",                                     // type
					anyUnbondActivity.MaybeTxID,                  // txid
					(*uint32)(nil),                               // event_position
					bignum.OptItoa(anyUnbondActivity.MaybeFee),   // fee
					(*string)(nil),                               // inputs
					(*uint32)(nil),                               // output_count
					anyUnbondActivity.MaybeStakingAccountAddress, // staking_account_address
					primptr.Uint64(uint64(2)),                    // staking_account_nonce
					primptr.String(anyUnbondActivity.MaybeBonded.String()),   // bonded
					primptr.String(anyUnbondActivity.MaybeUnbonded.String()), // unbonded
					nil,            // unbonded_from
					(*string)(nil), // joined_council_node
					(*uint64)(nil), // joined_council_node_id
					primptr.String(JsonMustMarshal(
						adapter.CouncilNodeToRDbCouncilNodeRow(&prevCouncilNode),
					)), // affected_council_node
					primptr.Uint64(*prevCouncilNode.Id), // affected_council_node_id
					nil,                                 // jailed_until
					(*string)(nil),                      // punishment_kind
				).Once().Return(execResult, nil)

				err := inserter.InsertUnbondTransaction(tx, &anyUnbondActivity)
				Expect(err).To(BeNil())
				tx.AssertExpectations(GinkgoT())
			})
		})
	})

	Describe("InsertWithdrawTransaction", func() {
		It("should return Error when the address does not exist before", func() {
			inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
				sq.StatementBuilder,
				new(PrimRDbTypeConv),
			)

			anyWithdrawActivity := RandomWithdrawActivity()

			tx := new(MockRDbTx)

			WhenStakingAccountDoesNotExist(tx, anyWithdrawActivity.MaybeStakingAccountAddress)
			WhenStakingAccountHasNotJoinedCouncilNode(tx, *anyWithdrawActivity.MaybeStakingAccountAddress)

			Expect(func() {
				_ = inserter.InsertWithdrawTransaction(tx, &anyWithdrawActivity)
			}).To(Panic())
		})

		Context("When staking account exist", func() {
			var anyWithdrawActivity chainindex.Activity
			var tx *MockRDbTx

			BeforeEach(func() {
				anyWithdrawActivity = RandomWithdrawActivity()
				anyWithdrawActivity.MaybeUnbonded = bignum.IntN10()
				anyWithdrawActivity.MaybeOutputCount = primptr.Uint32(uint32(3))

				tx = new(MockRDbTx)

				stakingAccountRow := adapter.RDbStakingAccountRow{
					Address:              *anyWithdrawActivity.MaybeStakingAccountAddress,
					Nonce:                uint64(1),
					Bonded:               bignum.Int0(),
					Unbonded:             bignum.Int10(),
					UnbondedFrom:         nil,
					PunishmentKind:       nil,
					JailedUntil:          nil,
					CurrentCouncilNodeId: nil,
				}
				WhenStakingAccountExist(tx, &stakingAccountRow)
			})

			Context("When staking account has not joined any council node", func() {
				BeforeEach(func() {
					WhenStakingAccountHasNotJoinedCouncilNode(tx, *anyWithdrawActivity.MaybeStakingAccountAddress)
				})

				It("should deduct staking address unbonded balance", func() {
					inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
						sq.StatementBuilder,
						new(PrimRDbTypeConv),
					)

					execResult := new(MockRDbExecResult)
					execResult.On("RowsAffected").Return(int64(1))
					OnTxInsertAnyActivity(tx).Return(execResult, nil)
					OnTxInsertAnyTransaferOutput(tx).Return(execResult, nil)

					tx.On("Exec",
						SQL_STAKING_ACCOUNT_UPDATE,
						primptr.String("0"), // bonded
						(*uint64)(nil),      // council_node_id
						nil,                 // jailed_until
						uint64(2),           // nonce
						(*string)(nil),      // punishment_kind
						primptr.String("0"), // unbonded
						nil,                 // unbonded_from
						*anyWithdrawActivity.MaybeStakingAccountAddress, // WHERE address = ?
					).Once().Return(execResult, nil)

					err := inserter.InsertWithdrawTransaction(tx, &anyWithdrawActivity)
					Expect(err).To(BeNil())
					tx.AssertExpectations(GinkgoT())
				})

				It("should insert all transaction outputs into table", func() {
					inserter := adapter.NewDefaultRDbBlockActivityDataRepo(sq.StatementBuilder, new(PrimRDbTypeConv))

					execResult := new(MockRDbExecResult)
					execResult.On("RowsAffected").Return(int64(1))
					OnTxUpdateAnyStakingAccount(tx).Return(execResult, nil)
					OnTxInsertAnyActivity(tx).Return(execResult, nil)

					OnTxInsertTransferOutput(tx,
						anyWithdrawActivity.MaybeTxID, 0,
					).Once().Return(execResult, nil)
					OnTxInsertTransferOutput(tx,
						anyWithdrawActivity.MaybeTxID, 1,
					).Once().Return(execResult, nil)
					OnTxInsertTransferOutput(tx,
						anyWithdrawActivity.MaybeTxID, 2,
					).Once().Return(execResult, nil)

					err := inserter.InsertWithdrawTransaction(tx, &anyWithdrawActivity)
					Expect(err).To(BeNil())
					tx.AssertExpectations(GinkgoT())
				})

				It("should insert activity into the table", func() {
					inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
						sq.StatementBuilder,
						new(PrimRDbTypeConv),
					)

					execResult := new(MockRDbExecResult)
					execResult.On("RowsAffected").Return(int64(1))
					OnTxUpdateAnyStakingAccount(tx).Return(execResult, nil)
					OnTxInsertAnyTransaferOutput(tx).Return(execResult, nil)

					tx.On("Exec",
						SQL_ACTIVITY_INSERT,
						anyWithdrawActivity.BlockHeight,                // block_height
						"withdraw",                                     // type
						anyWithdrawActivity.MaybeTxID,                  // txid
						(*uint32)(nil),                                 // event_position
						bignum.OptItoa(anyWithdrawActivity.MaybeFee),   // fee
						(*string)(nil),                                 // inputs
						anyWithdrawActivity.MaybeOutputCount,           // output_count
						anyWithdrawActivity.MaybeStakingAccountAddress, // staking_account_address
						primptr.Uint64(uint64(2)),                      // staking_account_nonce
						(*string)(nil),                                 // bonded
						primptr.String(anyWithdrawActivity.MaybeUnbonded.String()), // unbonded
						nil,            // unbonded_from
						(*string)(nil), // joined_council_node
						(*uint64)(nil), // joined_council_node_id
						(*string)(nil), // affected_council_node
						(*uint64)(nil), // affected_council_node_id
						nil,            // jailed_until
						(*string)(nil), // punishment_kind
					).Once().Return(execResult, nil)

					err := inserter.InsertWithdrawTransaction(tx, &anyWithdrawActivity)
					Expect(err).To(BeNil())
					tx.AssertExpectations(GinkgoT())
				})
			})
		})
	})

	Describe("InsertNodeJoinEvent", func() {
		var anyNodeJoinActivity chainindex.Activity
		var tx *MockRDbTx
		BeforeEach(func() {
			anyNodeJoinActivity = RandomNodeJoinActivity()
			anyNodeJoinActivity.MaybeCouncilNodeMeta.CreatedAtBlockHeight = anyNodeJoinActivity.BlockHeight
			anyNodeJoinActivity.MaybeCouncilNodeMeta.MaybeLastLeftAtBlockHeight = nil

			tx = new(MockRDbTx)

			stakingAccountRow := adapter.RDbStakingAccountRow{
				Address:              *anyNodeJoinActivity.MaybeStakingAccountAddress,
				Nonce:                uint64(1),
				Bonded:               bignum.Int0(),
				Unbonded:             bignum.Int0(),
				UnbondedFrom:         nil,
				PunishmentKind:       nil,
				JailedUntil:          nil,
				CurrentCouncilNodeId: nil,
			}
			WhenStakingAccountExist(tx, &stakingAccountRow)
		})

		It("should continue last council node when staking account and Tendermint address last council node is the same node joining now", func() {
			var err error

			inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
				sq.StatementBuilder,
				new(PrimRDbTypeConv),
			)

			prevCouncilNodeId := uint64(1)

			councilNodeRowResult := new(MockRDbRowResult)
			councilNodeRowResult.On("Scan", mock.MatchedBy(func(id **uint64) bool {
				*id = &prevCouncilNodeId
				return true
			}), mock.MatchedBy(func(name *string) bool {
				*name = anyNodeJoinActivity.MaybeCouncilNodeMeta.Name
				return true
			}), mock.MatchedBy(func(security_contact **string) bool {
				*security_contact = anyNodeJoinActivity.MaybeCouncilNodeMeta.MaybeSecurityContact
				return true
			}), mock.MatchedBy(func(pubkey_type *string) bool {
				*pubkey_type = "ed25519"
				return true
			}), mock.MatchedBy(func(pubkey *string) bool {
				*pubkey = anyNodeJoinActivity.MaybeCouncilNodeMeta.PubKey
				return true
			}), mock.MatchedBy(func(address *string) bool {
				*address = anyNodeJoinActivity.MaybeCouncilNodeMeta.Address
				return true
			}), mock.MatchedBy(func(created_at_block_height *uint64) bool {
				*created_at_block_height = anyNodeJoinActivity.MaybeCouncilNodeMeta.CreatedAtBlockHeight
				return true
			}), mock.MatchedBy(func(last_left_at_block_height **uint64) bool {
				*last_left_at_block_height = anyNodeJoinActivity.MaybeCouncilNodeMeta.MaybeLastLeftAtBlockHeight
				return true
			})).Return(nil)
			OnTxSelectCouncilNodeByAddress(
				tx, anyNodeJoinActivity.MaybeCouncilNodeMeta.Address,
			).Return(councilNodeRowResult)
			OnTxSelectCouncilNodeByStakingAccountAddress(
				tx, *anyNodeJoinActivity.MaybeStakingAccountAddress,
			).Return(councilNodeRowResult)

			execResult := new(MockRDbExecResult)
			execResult.On("RowsAffected").Return(int64(1))
			OnTxUpdateAnyStakingAccount(tx).Return(execResult, nil)
			OnTxInsertAnyActivity(tx).Return(execResult, nil)

			tx.On("Exec",
				SQL_COUNCIL_NODE_CLEAR_LAST_LEFT_AT_BLOCK_HEIGHT_UPDATE,
				nil,
				prevCouncilNodeId,
			).Once().Return(execResult, nil)

			err = inserter.InsertNodeJoinTransaction(tx, &anyNodeJoinActivity)
			Expect(err).To(BeNil())
			tx.AssertExpectations(GinkgoT())
		})

		It("should insert a new council node when it is a new one", func() {
			var err error

			inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
				sq.StatementBuilder,
				new(PrimRDbTypeConv),
			)

			OnTxSelectCouncilNodeByAddress(
				tx, anyNodeJoinActivity.MaybeCouncilNodeMeta.Address,
			).Return(EmptyRowResultOfScanArgsLen(8))
			OnTxSelectCouncilNodeByStakingAccountAddress(
				tx, *anyNodeJoinActivity.MaybeStakingAccountAddress,
			).Return(EmptyRowResultOfScanArgsLen(8))

			rowResult := new(MockRDbRowResult)
			rowResult.On("Scan", mock.MatchedBy(func(councilNodeId *uint64) bool {
				*councilNodeId = uint64(1)
				return true
			})).Return(nil)
			tx.On("QueryRow",
				SQL_COUNCIL_NODE_INSERT,
				anyNodeJoinActivity.MaybeCouncilNodeMeta.Name,
				anyNodeJoinActivity.MaybeCouncilNodeMeta.MaybeSecurityContact,
				adapter.PubKeyTypeToString(anyNodeJoinActivity.MaybeCouncilNodeMeta.PubKeyType),
				anyNodeJoinActivity.MaybeCouncilNodeMeta.PubKey,
				anyNodeJoinActivity.MaybeCouncilNodeMeta.Address,
				anyNodeJoinActivity.BlockHeight, // created_at_block_height
				(*uint64)(nil),                  // last_left_at_block_height
			).Once().Return(rowResult)

			execResult := new(MockRDbExecResult)
			execResult.On("RowsAffected").Return(int64(1))
			OnTxUpdateAnyStakingAccount(tx).Return(execResult, nil)
			OnTxInsertAnyActivity(tx).Return(execResult, nil)

			err = inserter.InsertNodeJoinTransaction(tx, &anyNodeJoinActivity)
			Expect(err).To(BeNil())
			tx.AssertExpectations(GinkgoT())
		})

		It("should update the staking account with the current council node id", func() {
			inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
				sq.StatementBuilder,
				new(PrimRDbTypeConv),
			)

			WhenCouncilNodeJoinIsNew(tx)

			councilNodeId := uint64(1)

			OnTxQueryInsertAnyCouncilNodeRowReturnCouncilNodeId(tx, councilNodeId)

			execResult := new(MockRDbExecResult)
			execResult.On("RowsAffected").Return(int64(1))
			OnTxInsertAnyActivity(tx).Return(execResult, nil)

			tx.On("Exec",
				SQL_STAKING_ACCOUNT_UPDATE,
				primptr.String("0"), // bonded
				&councilNodeId,      // current_council_node_id
				nil,                 // jailed_until
				uint64(2),           // nonce
				(*string)(nil),      // punishment_kind
				primptr.String("0"), // unbonded
				nil,                 // unbonded_from
				*anyNodeJoinActivity.MaybeStakingAccountAddress, // WHERE address = ?
			).Once().Return(execResult, nil)

			err := inserter.InsertNodeJoinTransaction(tx, &anyNodeJoinActivity)
			Expect(err).To(BeNil())
			tx.AssertExpectations(GinkgoT())
		})

		It("should insert activity into the table", func() {
			inserter := adapter.NewDefaultRDbBlockActivityDataRepo(
				sq.StatementBuilder,
				new(PrimRDbTypeConv),
			)

			WhenCouncilNodeJoinIsNew(tx)

			councilNodeId := uint64(1)

			OnTxQueryInsertAnyCouncilNodeRowReturnCouncilNodeId(tx, councilNodeId)

			execResult := new(MockRDbExecResult)
			execResult.On("RowsAffected").Return(int64(1))
			OnTxUpdateAnyStakingAccount(tx).Return(execResult, nil)

			joinedCouncilNode := *anyNodeJoinActivity.MaybeCouncilNodeMeta
			joinedCouncilNode.Id = &councilNodeId
			tx.On("Exec",
				SQL_ACTIVITY_INSERT,
				anyNodeJoinActivity.BlockHeight,                // block_height
				"nodejoin",                                     // type
				anyNodeJoinActivity.MaybeTxID,                  // txid
				(*uint32)(nil),                                 // event_position
				bignum.OptItoa(anyNodeJoinActivity.MaybeFee),   // fee
				(*string)(nil),                                 // inputs
				(*uint32)(nil),                                 // output_count
				anyNodeJoinActivity.MaybeStakingAccountAddress, // staking_account_address
				primptr.Uint64(uint64(2)),                      // staking_account_nonce
				(*string)(nil),                                 // bonded
				(*string)(nil),                                 // unbonded
				nil,                                            // unbdoned_from
				primptr.String(JsonMustMarshal(
					adapter.CouncilNodeToRDbCouncilNodeRow(&joinedCouncilNode),
				)), // joined_council_node
				primptr.Uint64(councilNodeId), // joined_council_node_id
				(*string)(nil),                // affected_council_node
				(*uint64)(nil),                // affected_council_node_id
				nil,                           // jailed_until
				(*string)(nil),                // punishment_kind
			).Once().Return(execResult, nil)

			err := inserter.InsertNodeJoinTransaction(tx, &anyNodeJoinActivity)
			Expect(err).To(BeNil())
			tx.AssertExpectations(GinkgoT())
		})
	})
})

func WhenCouncilNodeJoinIsNew(tx *MockRDbTx) {
	OnTxSelectCouncilNodeByAnyAddress(tx).Return(EmptyRowResultOfScanArgsLen(8))
	OnTxSelectCouncilNodeByAnyStakingAccountAddress(tx).Return(EmptyRowResultOfScanArgsLen(8))
}

func OnTxInsertAnyActivity(tx *MockRDbTx) *mock.Call {
	return tx.On("Exec",
		MockSQLWithAnyArgs(SQL_ACTIVITY_INSERT, 18)...,
	)
}

func OnTxUpdateTransaferInput(tx *MockRDbTx, txId *string, input chainindex.TxInput) *mock.Call {
	return tx.On("Exec",
		SQL_TRANSFER_OUTPUT_UPDATE,
		txId,
		input.TxId,
		input.Index,
	)
}

func OnTxUpdateAnyTransaferInput(tx *MockRDbTx) *mock.Call {
	return tx.On("Exec",
		MockSQLWithAnyArgs(SQL_TRANSFER_OUTPUT_UPDATE, 3)...,
	)
}

func OnTxInsertTransferOutput(tx *MockRDbTx, txId *string, index uint32) *mock.Call {
	return tx.On("Exec",
		SQL_TRANSFER_OUTPUT_INSERT,
		txId,
		index,
	)
}

func OnTxInsertAnyTransaferOutput(tx *MockRDbTx) *mock.Call {
	return tx.On("Exec",
		MockSQLWithAnyArgs(SQL_TRANSFER_OUTPUT_INSERT, 2)...,
	)
}

type MockLike interface {
	On(methodName string, arguments ...interface{}) *mock.Call
}

func WhenStakingAccountDoesNotExist(tx *MockRDbTx, address *string) {
	rowResult := EmptyRowResultOfScanArgsLen(8)
	OnTxQueryStakingAccountRow(tx, address).Return(rowResult)
}

func EmptyRowResultOfScanArgsLen(argsLen int) adapter.RDbRowResult {
	rowResult := new(MockRDbRowResult)
	rowResult.On("Scan", MockAnythingOfTimes(argsLen)...).Return(adapter.ErrNoRows)

	return rowResult
}

func WhenStakingAccountExist(tx *MockRDbTx, stakingAccountRow *adapter.RDbStakingAccountRow) {
	rowResult := new(MockRDbRowResult)
	OnScanStakingAccount(rowResult, stakingAccountRow).Return(nil)
	OnTxQueryStakingAccountRow(tx, &stakingAccountRow.Address).Return(rowResult)

}

// On Scan staking account with PrimRDbTypeConv
func OnScanStakingAccount(subject MockLike, stakingAccountRow *adapter.RDbStakingAccountRow) *mock.Call {
	bondedStr := stakingAccountRow.Bonded.String()
	unbondedStr := stakingAccountRow.Unbonded.String()

	return subject.On("Scan",
		mock.MatchedBy(func(address *string) bool {
			*address = stakingAccountRow.Address
			return true
		}),
		mock.MatchedBy(func(nonce *uint64) bool {
			*nonce = stakingAccountRow.Nonce
			return true
		}),
		mock.MatchedBy(func(bonded *string) bool {
			*bonded = bondedStr
			return true
		}),
		mock.MatchedBy(func(unboned *string) bool {
			*unboned = unbondedStr
			return true
		}),
		mock.MatchedBy(func(unbonded_from *time.Time) bool {
			if stakingAccountRow.UnbondedFrom != nil {
				*unbonded_from = *stakingAccountRow.UnbondedFrom
			}
			return true
		}),
		mock.MatchedBy(func(punishmentKind **string) bool {
			if stakingAccountRow.PunishmentKind != nil {
				**punishmentKind = *stakingAccountRow.PunishmentKind
			}
			return true
		}),
		mock.MatchedBy(func(jailedUntil *time.Time) bool {
			if stakingAccountRow.JailedUntil != nil {
				*jailedUntil = *stakingAccountRow.JailedUntil
			}
			return true
		}),
		mock.MatchedBy(func(currentCouncilNodeId **uint64) bool {
			if stakingAccountRow.CurrentCouncilNodeId != nil {
				**currentCouncilNodeId = *stakingAccountRow.CurrentCouncilNodeId
			}
			return true
		}),
	)
}

func OnTxInsertAnyStakingAccount(tx *MockRDbTx) *mock.Call {
	return tx.On("Exec",
		MockSQLWithAnyArgs(SQL_STAKING_ACCOUNT_INSERT, 8)...,
	)
}

func OnTxUpdateAnyStakingAccount(tx *MockRDbTx) *mock.Call {
	return tx.On("Exec",
		MockSQLWithAnyArgs(SQL_STAKING_ACCOUNT_UPDATE, 8)...,
	)
}

// // On Scan staking account with council node with PrimRDbTypeConv
// func OnScanStakingAccountWithCountilNode(subject MockLike, stakingAccount *chainindex.StakingAccount) *mock.Call {
// 	bondedStr := stakingAccount.Bonded.String()
// 	unbondedStr := stakingAccount.Bonded.String()

// 	return subject.On("Scan", mock.MatchedBy(func(address *string) bool {
// 		*address = stakingAccount.Address
// 		return true
// 	}), mock.MatchedBy(func(bonded *string) bool {
// 		*bonded = bondedStr
// 		return true
// 	}), mock.MatchedBy(func(unboned *string) bool {
// 		*unboned = unbondedStr
// 		return true
// 	}), mock.MatchedBy(func(jailedUntil *time.Time) bool {
// 		if stakingAccount.JailedUntil != nil {
// 			*jailedUntil = *stakingAccount.JailedUntil
// 		}
// 		return true
// 	}), mock.MatchedBy(func(punishmentKind *string) bool {
// 		if stakingAccount.PunishmentKind != nil {
// 			*punishmentKind = adapter.PunishmentKindToString(*stakingAccount.PunishmentKind)
// 		}
// 		return true
// 	}), mock.MatchedBy(func(councilNodeId *uint64) bool {
// 		if stakingAccount.CurrentCouncilNode != nil {
// 			if stakingAccount.CurrentCouncilNode.Id != nil {
// 				*councilNodeId = *stakingAccount.CurrentCouncilNode.Id
// 			}
// 		}
// 		return true
// 	}), mock.MatchedBy(func(councilNodeSecurityContact *string) bool {
// 		if stakingAccount.CurrentCouncilNode != nil {
// 			if stakingAccount.CurrentCouncilNode.SecurityContact != nil {
// 				*councilNodeSecurityContact = *stakingAccount.CurrentCouncilNode.SecurityContact
// 			}
// 		}
// 		return true
// 	}), mock.MatchedBy(func(pubKeyType *string) bool {
// 		if stakingAccount.CurrentCouncilNode != nil {
// 			*pubKeyType = adapter.PubKeyTypeToString(stakingAccount.CurrentCouncilNode.PubKeyType)
// 		}
// 		return true
// 	}), mock.MatchedBy(func(pubKey *string) bool {
// 		if stakingAccount.CurrentCouncilNode != nil {
// 			*pubKey = stakingAccount.CurrentCouncilNode.PubKey
// 		}
// 		return true
// 	}), mock.MatchedBy(func(address *string) bool {
// 		if stakingAccount.CurrentCouncilNode != nil {
// 			*address = stakingAccount.CurrentCouncilNode.Address
// 		}
// 		return true
// 	}), mock.MatchedBy(func(councilNodeCreatedAtBlockHeight *uint64) bool {
// 		if stakingAccount.CurrentCouncilNode != nil {
// 			*councilNodeCreatedAtBlockHeight = stakingAccount.CurrentCouncilNode.CreatedAtBlockHeight
// 		}
// 		return true
// 	}))
// }

func OnTxQueryStakingAccountRow(tx *MockRDbTx, address *string) *mock.Call {
	return tx.On("QueryRow",
		SQL_STAKING_ACCOUNT_SELECT,
		address,
	)
}

func OnTxQueryInsertAnyCouncilNodeRowReturnCouncilNodeId(tx *MockRDbTx, councilNodeId uint64) *mock.Call {
	rowResult := new(MockRDbRowResult)
	rowResult.On("Scan", mock.MatchedBy(func(id *uint64) bool {
		*id = councilNodeId
		return true
	})).Return(nil)

	return tx.On("QueryRow",
		MockSQLWithAnyArgs(SQL_COUNCIL_NODE_INSERT, 7)...,
	).Once().Return(rowResult)
}

func OnTxSelectCouncilNodeByAddress(tx *MockRDbTx, address string) *mock.Call {
	return tx.On("QueryRow",
		SQL_COUNCIL_NODE_BY_ADDRESS_SELECT,
		address,
	)
}

func OnTxSelectCouncilNodeByAnyAddress(tx *MockRDbTx) *mock.Call {
	return tx.On("QueryRow",
		MockSQLWithAnyArgs(SQL_COUNCIL_NODE_BY_ADDRESS_SELECT, 1)...,
	)
}

func OnTxSelectCouncilNodeByStakingAccountAddress(tx *MockRDbTx, stakingAccountAddress string) *mock.Call {
	return tx.On("QueryRow",
		SQL_COUNCIL_NODE_BY_STAKING_ACCOUNT_ADDRESS_SELECT,
		stakingAccountAddress,
	)
}

func OnTxSelectCouncilNodeByAnyStakingAccountAddress(tx *MockRDbTx) *mock.Call {
	return tx.On("QueryRow",
		MockSQLWithAnyArgs(SQL_COUNCIL_NODE_BY_STAKING_ACCOUNT_ADDRESS_SELECT, 1)...,
	)
}

func WhenStakingAccountHasNotJoinedCouncilNode(tx *MockRDbTx, stakingAccountAddress string) {
	OnTxSelectCurrentCouncilNodeByStakingAccountAddress(
		tx, stakingAccountAddress,
	).Return(EmptyRowResultOfScanArgsLen(8))
}

func WhenStakingAccountHasJoinedCouncilNode(
	tx *MockRDbTx,
	stakingAccountAddress string,
	councilNode chainindex.CouncilNode,
) {
	rowResult := new(MockRDbRowResult)
	rowResult.On("Scan", mock.MatchedBy(func(id **uint64) bool {
		*id = councilNode.Id
		return true
	}), mock.MatchedBy(func(name *string) bool {
		*name = councilNode.Name
		return true
	}), mock.MatchedBy(func(security_contact **string) bool {
		*security_contact = councilNode.MaybeSecurityContact
		return true
	}), mock.MatchedBy(func(pubkey_type *string) bool {
		*pubkey_type = "ed25519"
		return true
	}), mock.MatchedBy(func(pubkey *string) bool {
		*pubkey = councilNode.PubKey
		return true
	}), mock.MatchedBy(func(address *string) bool {
		*address = councilNode.Address
		return true
	}), mock.MatchedBy(func(created_at_block_height *uint64) bool {
		*created_at_block_height = councilNode.CreatedAtBlockHeight
		return true
	}), mock.MatchedBy(func(last_left_at_block_height **uint64) bool {
		*last_left_at_block_height = councilNode.MaybeLastLeftAtBlockHeight
		return true
	})).Return(nil)
	OnTxSelectCurrentCouncilNodeByStakingAccountAddress(
		tx, stakingAccountAddress,
	).Return(rowResult)
}

func OnTxSelectCurrentCouncilNodeByStakingAccountAddress(tx *MockRDbTx, stakingAccountAddress string) *mock.Call {
	return tx.On("QueryRow",
		SQL_STAKING_ACCOUNT_CURRENT_COUNCIL_NODE_SELECT,
		stakingAccountAddress,
	)
}

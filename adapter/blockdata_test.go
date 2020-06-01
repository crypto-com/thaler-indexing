package adapter_test

import (
	sq "github.com/Masterminds/squirrel"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	random "github.com/brianvoe/gofakeit/v5"

	"github.com/crypto-com/chainindex"
	"github.com/crypto-com/chainindex/adapter"
	. "github.com/crypto-com/chainindex/adapter/test/fake"
	. "github.com/crypto-com/chainindex/adapter/test/mock"
	"github.com/crypto-com/chainindex/internal/primptr"
	. "github.com/crypto-com/chainindex/test/factory"
	. "github.com/crypto-com/chainindex/usecase/test/factory"
)

const (
	SQL_BLOCK_INSERT                                          = "INSERT INTO blocks (height,hash,time,app_hash,committed_council_nodes) VALUES (?,?,?,?,?)"
	SQL_BLOCK_COMMITTED_COUNCIL_NODES_INSERT                  = "INSERT INTO block_committed_council_nodes (block_height,council_node_id,signature,is_proposer) VALUES (?,?,?,?)"
	SQL_REWARD_INSERT                                         = "INSERT INTO block_rewards (block_height,minted) VALUES (?,?)"
	SQL_COUNCIL_NODE_ID_BY_ADDRESS_SELECT                     = "SELECT id, name FROM council_nodes WHERE address = ? ORDER BY id DESC"
	SQL_COUNCIL_NODE_LAST_LEFT_AT_BLOCK_HEIGHT_UPDATE         = "UPDATE council_nodes SET last_left_at_block_height = ? WHERE id = ?"
	SQL_STAKING_ACCOUNT_REMOVE_CURRENT_COUNCIL_NODE_ID_UPDATE = "UPDATE staking_accounts SET current_council_node_id = ? WHERE current_council_node_id = ?"
)

var _ = Describe("Blockdata", func() {
	var mockConn *MockRDbConn
	var mockTx *MockRDbTx
	var mockActivityRepo *MockRDbBlockActivityDataRepo
	var repo *adapter.RDbBlockDataRepo
	BeforeEach(func() {
		mockTx = new(MockRDbTx)
		mockTx.On("Rollback").Return(nil)
		mockConn = new(MockRDbConn)
		mockConn.On("Begin").Return(mockTx, nil)

		stmtBuilder := sq.StatementBuilder
		mockActivityRepo = new(MockRDbBlockActivityDataRepo)
		fakeTypeConv := new(PrimRDbTypeConv)
		repo = adapter.NewRDbBlockDataRepo(mockConn, stmtBuilder, fakeTypeConv, mockActivityRepo)
	})

	Describe("Store", func() {
		It("should insert Block into table", func() {
			anyBlockData := RandomBlockData()
			anyBlockData.Signatures = RandomBlockSignaturesOfSize(3, anyBlockData.Block.Height)
			anyBlockData.Signatures[0].IsProposer = false
			anyBlockData.Signatures[1].IsProposer = true
			anyBlockData.Signatures[2].IsProposer = false
			anyBlockData.Activities = make([]chainindex.Activity, 0)
			anyBlockData.Reward = nil
			anyBlockData.CouncilNodeUpdates = make([]chainindex.CouncilNodeUpdate, 0)

			mockExecResult := new(MockRDbExecResult)
			mockExecResult.On("RowsAffected").Return(int64(1))
			OnTxInsertAnyBlockCommittedCouncilNode(mockTx).Return(mockExecResult, nil)

			anyCouncilNodeId0 := uint64(1)
			anyCouncilNodeName0 := random.Company()
			OnTxQueryCouncilNodeByAddressRowReturn(mockTx,
				anyBlockData.Signatures[0].CouncilNodeAddress, anyCouncilNodeId0, anyCouncilNodeName0,
			)
			anyCouncilNodeId1 := uint64(2)
			anyCouncilNodeName1 := random.Company()
			OnTxQueryCouncilNodeByAddressRowReturn(mockTx,
				anyBlockData.Signatures[1].CouncilNodeAddress, anyCouncilNodeId1, anyCouncilNodeName1,
			)
			anyCouncilNodeId2 := uint64(3)
			anyCouncilNodeName2 := random.Company()
			OnTxQueryCouncilNodeByAddressRowReturn(mockTx,
				anyBlockData.Signatures[2].CouncilNodeAddress, anyCouncilNodeId2, anyCouncilNodeName2,
			)

			committedCouncilNodes := []adapter.RDbBlockCommittedCouncilNodeRow{
				BlockSignatureToRDbBlockSignatureRow(anyBlockData.Signatures[1], anyCouncilNodeId1, anyCouncilNodeName1),
				BlockSignatureToRDbBlockSignatureRow(anyBlockData.Signatures[0], anyCouncilNodeId0, anyCouncilNodeName0),
				BlockSignatureToRDbBlockSignatureRow(anyBlockData.Signatures[2], anyCouncilNodeId2, anyCouncilNodeName2),
			}
			mockTx.On("Exec",
				SQL_BLOCK_INSERT,
				anyBlockData.Block.Height,
				anyBlockData.Block.Hash,
				&anyBlockData.Block.Time,
				anyBlockData.Block.AppHash,
				primptr.String(JsonMustMarshal(
					committedCouncilNodes,
				)),
			).Return(mockExecResult, nil)

			mockTx.On("Commit").Once().Return(nil)

			err := repo.Store(&anyBlockData)
			Expect(err).To(BeNil())
			mockTx.AssertExpectations(GinkgoT())
		})

		It("should panic when the council node who signed the signature does not exist", func() {
			anyBlockData := RandomBlockData()
			anyBlockData.Signatures = RandomBlockSignaturesOfSize(1, anyBlockData.Block.Height)
			anyBlockData.Activities = make([]chainindex.Activity, 0)
			anyBlockData.Reward = nil
			anyBlockData.CouncilNodeUpdates = make([]chainindex.CouncilNodeUpdate, 0)

			mockExecResult := new(MockRDbExecResult)
			mockExecResult.On("RowsAffected").Return(int64(1))
			OnTxInsertAnyBlock(mockTx).Return(mockExecResult, nil)

			mockRowResult := new(MockRDbRowResult)
			mockRowResult.On("Scan", mock.Anything).Return(adapter.ErrNoRows)
			mockTx.On("QueryRow",
				SQL_COUNCIL_NODE_ID_BY_ADDRESS_SELECT,
				anyBlockData.Signatures[0].CouncilNodeAddress,
			).Once().Return(mockRowResult)

			mockTx.On("Commit").Once().Return(nil)

			Expect(func() {
				_ = repo.Store(&anyBlockData)
			}).To(Panic())
		})

		It("should insert block signature into table", func() {
			anyBlockData := RandomBlockData()
			anyBlockData.Signatures = RandomBlockSignaturesOfSize(3, anyBlockData.Block.Height)
			anyBlockData.Activities = make([]chainindex.Activity, 0)
			anyBlockData.Reward = nil
			anyBlockData.CouncilNodeUpdates = make([]chainindex.CouncilNodeUpdate, 0)

			mockExecResult := new(MockRDbExecResult)
			mockExecResult.On("RowsAffected").Return(int64(1))
			OnTxInsertAnyBlock(mockTx).Return(mockExecResult, nil)

			anyCouncilNodeId0 := uint64(1)
			OnTxQueryCouncilNodeByAddressRowReturn(mockTx,
				anyBlockData.Signatures[0].CouncilNodeAddress, anyCouncilNodeId0, random.Company(),
			)
			anyCouncilNodeId1 := uint64(2)
			OnTxQueryCouncilNodeByAddressRowReturn(mockTx,
				anyBlockData.Signatures[1].CouncilNodeAddress, anyCouncilNodeId1, random.Company(),
			)
			anyCouncilNodeId2 := uint64(3)
			OnTxQueryCouncilNodeByAddressRowReturn(mockTx,
				anyBlockData.Signatures[2].CouncilNodeAddress, anyCouncilNodeId2, random.Company(),
			)

			OnTxInsertBlockCommittedCouncilNode(mockTx,
				&anyBlockData.Signatures[0], anyCouncilNodeId0,
			).Once().Return(mockExecResult, nil)
			OnTxInsertBlockCommittedCouncilNode(mockTx,
				&anyBlockData.Signatures[1], anyCouncilNodeId1,
			).Once().Return(mockExecResult, nil)
			OnTxInsertBlockCommittedCouncilNode(mockTx,
				&anyBlockData.Signatures[2], anyCouncilNodeId2,
			).Once().Return(mockExecResult, nil)

			mockTx.On("Commit").Once().Return(nil)

			err := repo.Store(&anyBlockData)
			Expect(err).To(BeNil())
			mockTx.AssertExpectations(GinkgoT())
		})

		It("should store activity into database", func() {
			anyBlockData := RandomBlockData()
			anyBlockData.Signatures = RandomBlockSignaturesOfSize(0, anyBlockData.Block.Height)
			anyBlockData.Reward = nil
			anyBlockData.CouncilNodeUpdates = make([]chainindex.CouncilNodeUpdate, 0)

			anyTransferActivity := RandomTransferActivity()
			anyDepositActivity := RandomDepositActivity()
			anyUnbondActivity := RandomUnbondActivity()
			anyWithdrawActivity := RandomWithdrawActivity()
			anyNodeJoinActivity := RandomNodeJoinActivity()
			anyUnjailActivity := RandomUnjailActivity()
			anyRewardActivity := RandomRewardActivity()
			anySlashActivity := RandomSlashActivity()
			anyJailActivity := RandomJailActivity()
			anyBlockData.Activities = []chainindex.Activity{
				anyTransferActivity,
				anyDepositActivity,
				anyUnbondActivity,
				anyWithdrawActivity,
				anyNodeJoinActivity,
				anyUnjailActivity,
				anyRewardActivity,
				anySlashActivity,
				anyJailActivity,
			}

			mockActivityRepo.On("InsertTransferTransaction", mock.Anything, &anyTransferActivity).Once().Return(nil)
			mockActivityRepo.On("InsertDepositTransaction", mock.Anything, &anyDepositActivity).Once().Return(nil)
			mockActivityRepo.On("InsertUnbondTransaction", mock.Anything, &anyUnbondActivity).Once().Return(nil)
			mockActivityRepo.On("InsertWithdrawTransaction", mock.Anything, &anyWithdrawActivity).Once().Return(nil)
			mockActivityRepo.On("InsertNodeJoinTransaction", mock.Anything, &anyNodeJoinActivity).Once().Return(nil)
			mockActivityRepo.On("InsertUnjailTransaction", mock.Anything, &anyUnjailActivity).Once().Return(nil)
			mockActivityRepo.On("InsertRewardEvent", mock.Anything, &anyRewardActivity).Once().Return(nil)
			mockActivityRepo.On("InsertSlashEvent", mock.Anything, &anySlashActivity).Once().Return(nil)
			mockActivityRepo.On("InsertJailEvent", mock.Anything, &anyJailActivity).Once().Return(nil)

			mockExecResult := new(MockRDbExecResult)
			mockExecResult.On("RowsAffected").Return(int64(1))
			mockTx.On("Exec",
				SQL_BLOCK_INSERT,
				anyBlockData.Block.Height,
				anyBlockData.Block.Hash,
				&anyBlockData.Block.Time,
				anyBlockData.Block.AppHash,
				(*string)(nil),
			).Return(mockExecResult, nil)

			mockTx.On("Commit").Once().Return(nil)

			err := repo.Store(&anyBlockData)
			Expect(err).To(BeNil())
			mockTx.AssertExpectations(GinkgoT())
			mockActivityRepo.AssertExpectations(GinkgoT())
		})

		It("should insert Reward into table", func() {
			anyBlockData := RandomBlockData()
			anyBlockData.Signatures = RandomBlockSignaturesOfSize(0, anyBlockData.Block.Height)
			anyBlockData.Activities = make([]chainindex.Activity, 0)
			anyBlockData.CouncilNodeUpdates = make([]chainindex.CouncilNodeUpdate, 0)
			reward := RandomBlockReward()
			anyBlockData.Reward = &reward
			anyBlockData.Reward.BlockHeight = anyBlockData.Block.Height

			mockExecResult := new(MockRDbExecResult)
			mockExecResult.On("RowsAffected").Return(int64(1))

			OnTxInsertAnyBlock(mockTx).Return(mockExecResult, nil)

			mockTx.On("Exec",
				SQL_REWARD_INSERT,
				anyBlockData.Block.Height,
				primptr.String(anyBlockData.Reward.Minted.String()),
			).Once().Return(mockExecResult, nil)

			mockTx.On("Commit").Once().Return(nil)

			err := repo.Store(&anyBlockData)
			Expect(err).To(BeNil())
			mockTx.AssertExpectations(GinkgoT())
		})

		It("should update council node last left at block height", func() {
			anyBlockData := RandomBlockData()
			anyBlockData.Signatures = RandomBlockSignaturesOfSize(0, anyBlockData.Block.Height)
			anyBlockData.Activities = make([]chainindex.Activity, 0)
			anyBlockData.Reward = nil
			anyBlockData.CouncilNodeUpdates = RandomCouncilNodeUpdatesOfSize(3)

			mockExecResult := new(MockRDbExecResult)
			mockExecResult.On("RowsAffected").Return(int64(1))

			OnTxInsertAnyBlock(mockTx).Return(mockExecResult, nil)
			OnTxRemoveAnyStakingAccountCurrentCouncilNodeId(mockTx).Return(mockExecResult, nil)

			anyCouncilNodeId0 := uint64(1)
			OnTxQueryCouncilNodeByAddressRowReturn(mockTx,
				anyBlockData.CouncilNodeUpdates[0].Address, anyCouncilNodeId0, random.Company(),
			)
			anyCouncilNodeId1 := uint64(2)
			OnTxQueryCouncilNodeByAddressRowReturn(mockTx,
				anyBlockData.CouncilNodeUpdates[1].Address, anyCouncilNodeId1, random.Company(),
			)
			anyCouncilNodeId2 := uint64(3)
			OnTxQueryCouncilNodeByAddressRowReturn(mockTx,
				anyBlockData.CouncilNodeUpdates[2].Address, anyCouncilNodeId2, random.Company(),
			)

			OnTxUpdateCouncilNodeLastLeftAtBlockHeight(mockTx,
				anyBlockData.Block.Height, anyCouncilNodeId0,
			).Once().Return(mockExecResult, nil)
			OnTxUpdateCouncilNodeLastLeftAtBlockHeight(mockTx,
				anyBlockData.Block.Height, anyCouncilNodeId1,
			).Once().Return(mockExecResult, nil)
			OnTxUpdateCouncilNodeLastLeftAtBlockHeight(mockTx,
				anyBlockData.Block.Height, anyCouncilNodeId2,
			).Once().Return(mockExecResult, nil)

			mockTx.On("Commit").Once().Return(nil)

			err := repo.Store(&anyBlockData)
			Expect(err).To(BeNil())
			mockTx.AssertExpectations(GinkgoT())
		})

		It("should remove current council node id from staking account", func() {
			anyBlockData := RandomBlockData()
			anyBlockData.Signatures = RandomBlockSignaturesOfSize(0, anyBlockData.Block.Height)
			anyBlockData.Activities = make([]chainindex.Activity, 0)
			anyBlockData.Reward = nil
			anyBlockData.CouncilNodeUpdates = RandomCouncilNodeUpdatesOfSize(3)

			mockExecResult := new(MockRDbExecResult)
			mockExecResult.On("RowsAffected").Return(int64(1))

			OnTxInsertAnyBlock(mockTx).Return(mockExecResult, nil)
			OnTxUpdateAnyCouncilNodeLastLeftAtBlockHeight(mockTx).Return(mockExecResult, nil)

			anyCouncilNodeId0 := uint64(1)
			OnTxQueryCouncilNodeByAddressRowReturn(mockTx,
				anyBlockData.CouncilNodeUpdates[0].Address, anyCouncilNodeId0, random.Company(),
			)
			anyCouncilNodeId1 := uint64(2)
			OnTxQueryCouncilNodeByAddressRowReturn(mockTx,
				anyBlockData.CouncilNodeUpdates[1].Address, anyCouncilNodeId1, random.Company(),
			)
			anyCouncilNodeId2 := uint64(3)
			OnTxQueryCouncilNodeByAddressRowReturn(mockTx,
				anyBlockData.CouncilNodeUpdates[2].Address, anyCouncilNodeId2, random.Company(),
			)

			mockTx.On("Exec",
				SQL_STAKING_ACCOUNT_REMOVE_CURRENT_COUNCIL_NODE_ID_UPDATE,
				nil,               // SET current_council_node_id = ?
				anyCouncilNodeId0, // WHERE current_council_node_id = ?,
			).Return(mockExecResult, nil)
			mockTx.On("Exec",
				SQL_STAKING_ACCOUNT_REMOVE_CURRENT_COUNCIL_NODE_ID_UPDATE,
				nil,               // SET current_council_node_id = ?
				anyCouncilNodeId1, // WHERE current_council_node_id = ?,
			).Return(mockExecResult, nil)
			mockTx.On("Exec",
				SQL_STAKING_ACCOUNT_REMOVE_CURRENT_COUNCIL_NODE_ID_UPDATE,
				nil,               // SET current_council_node_id = ?
				anyCouncilNodeId2, // WHERE current_council_node_id = ?,
			).Return(mockExecResult, nil)

			mockTx.On("Commit").Once().Return(nil)

			err := repo.Store(&anyBlockData)
			Expect(err).To(BeNil())
			mockTx.AssertExpectations(GinkgoT())
		})

		It("should roll back all changes whenever there is an error", func() {
			anyBlockData := RandomBlockData()
			anyBlockData.Signatures = RandomBlockSignaturesOfSize(0, anyBlockData.Block.Height)
			anyBlockData.Activities = make([]chainindex.Activity, 0)
			anyBlockData.Reward = nil
			anyBlockData.CouncilNodeUpdates = make([]chainindex.CouncilNodeUpdate, 0)

			OnTxInsertAnyBlock(mockTx).Once().Return(nil, adapter.ErrRepoWrite)

			err := repo.Store(&anyBlockData)
			Expect(err).NotTo(BeNil())
		})
	})
})

func BlockSignatureToRDbBlockSignatureRow(signature chainindex.BlockSignature, councilNodeId uint64, councilNodeName string) adapter.RDbBlockCommittedCouncilNodeRow {
	return adapter.RDbBlockCommittedCouncilNodeRow{
		BlockHeight:        signature.BlockHeight,
		ID:                 councilNodeId,
		Name:               councilNodeName,
		CouncilNodeAddress: signature.CouncilNodeAddress,
		Signature:          signature.Signature,
		IsProposer:         signature.IsProposer,
	}
}

func OnTxInsertAnyBlock(mockTx *MockRDbTx) *mock.Call {
	return mockTx.On("Exec",
		MockSQLWithAnyArgs(SQL_BLOCK_INSERT, 5)...,
	)
}

func OnTxQueryCouncilNodeByAddressRowReturn(mockTx *MockRDbTx, address string, councilNodeId uint64, councilNodeName string) *mock.Call {
	mockRowResult := new(MockRDbRowResult)
	mockRowResult.On("Scan", mock.MatchedBy(func(id *uint64) bool {
		*id = councilNodeId
		return true
	}), mock.MatchedBy(func(name *string) bool {
		*name = councilNodeName
		return true
	})).Return(nil)
	return mockTx.On("QueryRow",
		SQL_COUNCIL_NODE_ID_BY_ADDRESS_SELECT,
		address,
	).Once().Return(mockRowResult)
}

func OnTxInsertBlockCommittedCouncilNode(mockTx *MockRDbTx, signature *chainindex.BlockSignature, councilNodeId uint64) *mock.Call {
	return mockTx.On("Exec",
		SQL_BLOCK_COMMITTED_COUNCIL_NODES_INSERT,
		signature.BlockHeight,
		councilNodeId,
		signature.Signature,
		signature.IsProposer,
	)
}

func OnTxInsertAnyBlockCommittedCouncilNode(mockTx *MockRDbTx) *mock.Call {
	return mockTx.On("Exec",
		MockSQLWithAnyArgs(SQL_BLOCK_COMMITTED_COUNCIL_NODES_INSERT, 5)...,
	)
}

func OnTxUpdateCouncilNodeLastLeftAtBlockHeight(mockTx *MockRDbTx, blockHeight uint64, councilNodeId uint64) *mock.Call {
	return mockTx.On("Exec",
		SQL_COUNCIL_NODE_LAST_LEFT_AT_BLOCK_HEIGHT_UPDATE,
		blockHeight,
		councilNodeId,
	)
}

func OnTxUpdateAnyCouncilNodeLastLeftAtBlockHeight(mockTx *MockRDbTx) *mock.Call {
	return mockTx.On("Exec",
		MockSQLWithAnyArgs(SQL_COUNCIL_NODE_LAST_LEFT_AT_BLOCK_HEIGHT_UPDATE, 2)...,
	)
}

func OnTxRemoveAnyStakingAccountCurrentCouncilNodeId(mockTx *MockRDbTx) *mock.Call {
	return mockTx.On("Exec",
		MockSQLWithAnyArgs(SQL_STAKING_ACCOUNT_REMOVE_CURRENT_COUNCIL_NODE_ID_UPDATE, 2)...,
	)
}

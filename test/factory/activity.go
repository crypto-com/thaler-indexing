package factory

import (
	"encoding/hex"
	"math/big"
	"strings"

	random "github.com/brianvoe/gofakeit/v5"

	"github.com/crypto-com/chainindex"
	"github.com/crypto-com/chainindex/internal/bignum"
	"github.com/crypto-com/chainindex/internal/primptr"
)

// Generate a slice of all possible activities
func RandomAllActivities() []chainindex.Activity {
	anyTransferActivity := RandomTransferActivity()
	anyDepositActivity := RandomDepositActivity()
	anyUnbondActivity := RandomUnbondActivity()
	anyWithdrawActivity := RandomWithdrawActivity()
	anyNodeJoinActivity := RandomNodeJoinActivity()
	anyUnjailActivity := RandomUnjailActivity()
	anyRewardActivity := RandomRewardActivity()
	anySlashActivity := RandomSlashActivity()
	anyJailActivity := RandomJailActivity()
	return []chainindex.Activity{
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
}

func RandomTransferActivity() chainindex.Activity {
	return chainindex.Activity{
		BlockHeight:                random.Uint64(),
		Type:                       chainindex.ACTIVITY_TRANSFER,
		MaybeTxID:                  RandomTxIdPtr(),
		MaybeEventPosition:         nil,
		MaybeFee:                   RandomFee(),
		MaybeTxInputs:              RandomTxInputsPtrOfSize(3),
		MaybeOutputCount:           RandomUint32Ptr(),
		MaybeStakingAccountAddress: nil,
		MaybeBonded:                nil,
		MaybeUnbonded:              nil,
		MaybeCouncilNodeMeta:       nil,
		MaybeAffectedCouncilNode:   nil,
		MaybePunishmentKind:        nil,
		MaybeJailedUntil:           nil,
	}
}

func RandomGenesisActivity() chainindex.Activity {
	return chainindex.Activity{
		BlockHeight:                uint64(1),
		Type:                       chainindex.ACTIVITY_GENESIS,
		MaybeTxID:                  nil,
		MaybeEventPosition:         nil,
		MaybeFee:                   nil,
		MaybeTxInputs:              nil,
		MaybeOutputCount:           nil,
		MaybeStakingAccountAddress: RandomStakingAddressPtr(),
		MaybeBonded:                RandomCoin(),
		MaybeUnbonded:              RandomCoin(),
		MaybeCouncilNodeMeta:       RandomCouncilNodePtr(),
		MaybeAffectedCouncilNode:   nil,
		MaybePunishmentKind:        nil,
		MaybeJailedUntil:           nil,
	}
}

func RandomDepositActivity() chainindex.Activity {
	return chainindex.Activity{
		BlockHeight:                random.Uint64(),
		Type:                       chainindex.ACTIVITY_DEPOSIT,
		MaybeTxID:                  RandomTxIdPtr(),
		MaybeEventPosition:         nil,
		MaybeFee:                   RandomFee(),
		MaybeTxInputs:              nil,
		MaybeOutputCount:           nil,
		MaybeStakingAccountAddress: RandomStakingAddressPtr(),
		MaybeBonded:                RandomCoin(),
		MaybeUnbonded:              nil,
		MaybeCouncilNodeMeta:       nil,
		MaybeAffectedCouncilNode:   nil,
		MaybePunishmentKind:        nil,
		MaybeJailedUntil:           nil,
	}
}

func RandomUnbondActivity() chainindex.Activity {
	unbonded := RandomCoin()
	bonded := new(big.Int).Add(bignum.Int0(), unbonded)
	return chainindex.Activity{
		BlockHeight:                random.Uint64(),
		Type:                       chainindex.ACTIVITY_UNBOND,
		MaybeTxID:                  RandomTxIdPtr(),
		MaybeEventPosition:         nil,
		MaybeFee:                   RandomFee(),
		MaybeTxInputs:              nil,
		MaybeOutputCount:           nil,
		MaybeStakingAccountAddress: RandomStakingAddressPtr(),
		MaybeBonded:                bonded,
		MaybeUnbonded:              unbonded,
		MaybeCouncilNodeMeta:       nil,
		MaybeAffectedCouncilNode:   nil,
		MaybePunishmentKind:        nil,
		MaybeJailedUntil:           nil,
	}
}

func RandomWithdrawActivity() chainindex.Activity {
	return chainindex.Activity{
		BlockHeight:                random.Uint64(),
		Type:                       chainindex.ACTIVITY_WITHDRAW,
		MaybeTxID:                  RandomTxIdPtr(),
		MaybeEventPosition:         nil,
		MaybeFee:                   RandomFee(),
		MaybeTxInputs:              nil,
		MaybeOutputCount:           RandomUint32Ptr(),
		MaybeStakingAccountAddress: RandomStakingAddressPtr(),
		MaybeBonded:                nil,
		MaybeUnbonded:              RandomNegativeCoin(),
		MaybeCouncilNodeMeta:       nil,
		MaybeAffectedCouncilNode:   nil,
		MaybePunishmentKind:        nil,
		MaybeJailedUntil:           nil,
	}
}

func RandomNodeJoinActivity() chainindex.Activity {
	return chainindex.Activity{
		BlockHeight:                random.Uint64(),
		Type:                       chainindex.ACTIVITY_NODEJOIN,
		MaybeTxID:                  RandomTxIdPtr(),
		MaybeEventPosition:         nil,
		MaybeFee:                   bignum.Int0(),
		MaybeTxInputs:              nil,
		MaybeOutputCount:           nil,
		MaybeStakingAccountAddress: RandomStakingAddressPtr(),
		MaybeBonded:                nil,
		MaybeUnbonded:              nil,
		MaybeCouncilNodeMeta:       RandomCouncilNodePtr(),
		MaybeAffectedCouncilNode:   nil,
		MaybePunishmentKind:        nil,
		MaybeJailedUntil:           nil,
	}
}

func RandomUnjailActivity() chainindex.Activity {
	return chainindex.Activity{
		BlockHeight:                random.Uint64(),
		Type:                       chainindex.ACTIVITY_UNJAIL,
		MaybeTxID:                  RandomTxIdPtr(),
		MaybeEventPosition:         nil,
		MaybeFee:                   RandomFee(),
		MaybeTxInputs:              nil,
		MaybeOutputCount:           nil,
		MaybeStakingAccountAddress: RandomStakingAddressPtr(),
		MaybeBonded:                nil,
		MaybeUnbonded:              nil,
		MaybeCouncilNodeMeta:       nil,
		MaybeAffectedCouncilNode:   nil,
		MaybePunishmentKind:        nil,
		MaybeJailedUntil:           nil,
	}
}

func RandomRewardActivity() chainindex.Activity {
	return chainindex.Activity{
		BlockHeight:                random.Uint64(),
		Type:                       chainindex.ACTIVITY_REWARD,
		MaybeTxID:                  nil,
		MaybeEventPosition:         RandomUint32Ptr(),
		MaybeFee:                   nil,
		MaybeTxInputs:              nil,
		MaybeOutputCount:           nil,
		MaybeStakingAccountAddress: RandomStakingAddressPtr(),
		MaybeBonded:                RandomCoin(),
		MaybeUnbonded:              nil,
		MaybeCouncilNodeMeta:       nil,
		MaybeAffectedCouncilNode:   nil,
		MaybePunishmentKind:        nil,
		MaybeJailedUntil:           nil,
	}
}

func RandomSlashActivity() chainindex.Activity {
	return chainindex.Activity{
		BlockHeight:                random.Uint64(),
		Type:                       chainindex.ACTIVITY_SLASH,
		MaybeTxID:                  nil,
		MaybeEventPosition:         RandomUint32Ptr(),
		MaybeFee:                   nil,
		MaybeTxInputs:              nil,
		MaybeOutputCount:           nil,
		MaybeStakingAccountAddress: RandomStakingAddressPtr(),
		MaybeBonded:                nil,
		MaybeUnbonded:              nil,
		MaybeCouncilNodeMeta:       nil,
		MaybeAffectedCouncilNode:   nil,
		MaybePunishmentKind:        nil,
		MaybeJailedUntil:           nil,
	}
}

func RandomJailActivity() chainindex.Activity {
	return chainindex.Activity{
		BlockHeight:                random.Uint64(),
		Type:                       chainindex.ACTIVITY_JAIL,
		MaybeTxID:                  nil,
		MaybeEventPosition:         RandomUint32Ptr(),
		MaybeFee:                   nil,
		MaybeTxInputs:              nil,
		MaybeOutputCount:           nil,
		MaybeStakingAccountAddress: RandomStakingAddressPtr(),
		MaybeBonded:                nil,
		MaybeUnbonded:              nil,
		MaybeCouncilNodeMeta:       nil,
		MaybeAffectedCouncilNode:   nil,
		MaybePunishmentKind:        RandomPunishmentKindPtr(),
		MaybeJailedUntil:           primptr.Time(RandomUTCTime()),
	}
}

func RandomTxInputsPtrOfSize(size int) []chainindex.TxInput {
	inputs := make([]chainindex.TxInput, size)
	for i := 0; i < size; i += 1 {
		inputs[i] = RandomTxInput()
	}

	return inputs
}

func RandomTxInput() chainindex.TxInput {
	return chainindex.TxInput{
		TxId:  RandomTxId(),
		Index: random.Uint32(),
	}
}

func RandomTxIdPtr() *string {
	value := RandomTxId()
	return &value
}

func RandomTxId() string {
	hexTxId := hex.EncodeToString(RandomHex(32))
	return strings.ToUpper(hexTxId)
}

func RandomFee() *big.Int {
	value := random.Number(0, 1_0000_0000)
	return big.NewInt(int64(value))
}

func RandomStakingAddressPtr() *string {
	value := RandomStakingAddress()
	return &value
}

func RandomStakingAddress() string {
	return "0x" + hex.EncodeToString(RandomHex(20))
}

func RandomPunishmentKind() chainindex.PunishmentKind {
	return uint8(random.RandomInt([]int{
		int(chainindex.PUNISHMENT_KIND_NON_LIVE),
		int(chainindex.PUNISHMENT_KIND_BYZANTINE_FAULT),
	}))
}

func RandomPunishmentKindPtr() *chainindex.PunishmentKind {
	value := RandomPunishmentKind()
	return &value
}

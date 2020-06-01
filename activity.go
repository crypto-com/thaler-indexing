package chainindex

import (
	"math/big"
	"time"

	"github.com/luci/go-render/render"
)

type Activity struct {
	BlockHeight                uint64
	Type                       ActivityType
	MaybeTxID                  *string
	MaybeEventPosition         *uint32
	MaybeFee                   *big.Int
	MaybeTxInputs              []TxInput
	MaybeOutputCount           *uint32
	MaybeStakingAccountAddress *string
	MaybeStakingAccountNonce   *uint64
	MaybeBonded                *big.Int
	MaybeUnbonded              *big.Int
	MaybeUnbondedFrom          *time.Time
	MaybeCouncilNodeMeta       *CouncilNode
	MaybeAffectedCouncilNode   *CouncilNode
	MaybeJailedUntil           *time.Time
	MaybePunishmentKind        *PunishmentKind
}

type ActivityType = uint8

const (
	ACTIVITY_GENESIS ActivityType = iota
	ACTIVITY_TRANSFER
	ACTIVITY_DEPOSIT
	ACTIVITY_UNBOND
	ACTIVITY_WITHDRAW
	ACTIVITY_NODEJOIN
	ACTIVITY_UNJAIL
	ACTIVITY_REWARD
	ACTIVITY_SLASH
	ACTIVITY_JAIL
	ACTIVITY_NODEKICKED
)

type TransactionType = int8

const (
	TRANSACTION_TRANSFER TransactionType = iota
	TRANSACTION_DEPOSIT
	TRANSACTION_UNBOND
	TRANSACTION_WITHDRAW
	TRANSACTION_NODEJOIN
	TRANSACTION_UNJAIL
)

type EventType = int8

const (
	EVENT_REWARD EventType = iota
	EVENT_SLASH
	EVENT_JAIL
)

type PunishmentKind = uint8

const (
	PUNISHMENT_KIND_NON_LIVE PunishmentKind = iota
	PUNISHMENT_KIND_BYZANTINE_FAULT
)

func (activity *Activity) String() string {
	return render.Render(activity)
}

type TxInput struct {
	TxId  string
	Index uint32
}

func (input *TxInput) String() string {
	return render.Render(input)
}

type TransactionOutput struct {
	TxId               string
	Index              uint32
	SpentAtBlockHeight uint64
}

func (output *TransactionOutput) String() string {
	return render.Render(output)
}

package adapter

import (
	"github.com/crypto-com/chainindex"
)

type RDbTransferInputRow struct {
	TxID  string `json:"prev_txid"`
	Index int32  `json:"prev_output_index"`
}

func TxInputsToRDbTransferInputs(inputs []chainindex.TxInput) []RDbTransferInputRow {
	if inputs == nil {
		return nil
	}
	transferInputs := make([]RDbTransferInputRow, 0, len(inputs))
	for _, input := range inputs {
		transferInputs = append(transferInputs, RDbTransferInputRow{
			TxID:  input.TxId,
			Index: int32(input.Index),
		})
	}

	return transferInputs
}

func ActivityTypeToString(activityType chainindex.ActivityType) string {
	switch activityType {
	case chainindex.ACTIVITY_GENESIS:
		return "genesis"
	case chainindex.ACTIVITY_TRANSFER:
		return "transfer"
	case chainindex.ACTIVITY_DEPOSIT:
		return "deposit"
	case chainindex.ACTIVITY_UNBOND:
		return "unbond"
	case chainindex.ACTIVITY_WITHDRAW:
		return "withdraw"
	case chainindex.ACTIVITY_NODEJOIN:
		return "nodejoin"
	case chainindex.ACTIVITY_UNJAIL:
		return "unjail"
	case chainindex.ACTIVITY_REWARD:
		return "reward"
	case chainindex.ACTIVITY_SLASH:
		return "slash"
	case chainindex.ACTIVITY_JAIL:
		return "jail"
	case chainindex.ACTIVITY_NODEKICKED:
		return "nodekicked"
	}
	panic("unsupported activity type")
}

func StringToActivityType(transactionType string) chainindex.ActivityType {
	switch transactionType {
	case "transfer":
		return chainindex.ACTIVITY_TRANSFER
	case "deposit":
		return chainindex.ACTIVITY_DEPOSIT
	case "unbond":
		return chainindex.ACTIVITY_UNBOND
	case "withdraw":
		return chainindex.ACTIVITY_WITHDRAW
	case "nodejoin":
		return chainindex.ACTIVITY_NODEJOIN
	case "unjail":
		return chainindex.ACTIVITY_UNJAIL
	case "reward":
		return chainindex.ACTIVITY_REWARD
	case "slash":
		return chainindex.ACTIVITY_SLASH
	case "jail":
		return chainindex.ACTIVITY_JAIL
	}
	panic("unsupported activity type")
}

func IsValidActivityType(transactionType string) bool {
	switch transactionType {
	case "transfer":
		fallthrough
	case "deposit":
		fallthrough
	case "unbond":
		fallthrough
	case "withdraw":
		fallthrough
	case "nodejoin":
		fallthrough
	case "unjail":
		fallthrough
	case "reward":
		fallthrough
	case "slash":
		fallthrough
	case "jail":
		return true
	}
	return false
}

func IsValidTransactionType(transactionType string) bool {
	switch transactionType {
	case "transfer":
		fallthrough
	case "deposit":
		fallthrough
	case "unbond":
		fallthrough
	case "withdraw":
		fallthrough
	case "nodejoin":
		fallthrough
	case "unjail":
		return true
	}
	return false
}

func TransactionTypeToString(transactionType chainindex.TransactionType) string {
	switch transactionType {
	case chainindex.TRANSACTION_TRANSFER:
		return "transfer"
	case chainindex.TRANSACTION_DEPOSIT:
		return "deposit"
	case chainindex.TRANSACTION_UNBOND:
		return "unbond"
	case chainindex.TRANSACTION_WITHDRAW:
		return "withdraw"
	case chainindex.TRANSACTION_NODEJOIN:
		return "nodejoin"
	case chainindex.TRANSACTION_UNJAIL:
		return "unjail"
	}
	panic("unsupported transaction type")
}

func StringToTransactionType(transactionType string) chainindex.TransactionType {
	switch transactionType {
	case "transfer":
		return chainindex.TRANSACTION_TRANSFER
	case "deposit":
		return chainindex.TRANSACTION_DEPOSIT
	case "unbond":
		return chainindex.TRANSACTION_UNBOND
	case "withdraw":
		return chainindex.TRANSACTION_WITHDRAW
	case "nodejoin":
		return chainindex.TRANSACTION_NODEJOIN
	case "unjail":
		return chainindex.TRANSACTION_UNJAIL
	}
	panic("unsupported transaction type")
}

func EventTypeToString(eventType chainindex.EventType) string {
	switch eventType {
	case chainindex.EVENT_REWARD:
		return "reward"
	case chainindex.EVENT_SLASH:
		return "slash"
	case chainindex.EVENT_JAIL:
		return "jail"
	}
	panic("unsupported event type")
}

func StringToEventType(eventType string) chainindex.EventType {
	switch eventType {
	case "reward":
		return chainindex.EVENT_REWARD
	case "slash":
		return chainindex.EVENT_SLASH
	case "jail":
		return chainindex.EVENT_JAIL
	}
	panic("unsupported event type")
}

func IsValidEventType(eventType string) bool {
	switch eventType {
	case "reward":
		fallthrough
	case "slash":
		fallthrough
	case "jail":
		return true
	}
	return false
}

func OptPunishmentKindToString(punishmentKind *chainindex.PunishmentKind) *string {
	if punishmentKind == nil {
		return nil
	}
	str := PunishmentKindToString(*punishmentKind)

	return &str
}

func PunishmentKindToString(punishmentKind chainindex.PunishmentKind) string {
	switch punishmentKind {
	case chainindex.PUNISHMENT_KIND_BYZANTINE_FAULT:
		return "ByzantineFault"
	case chainindex.PUNISHMENT_KIND_NON_LIVE:
		return "NonLive"
	default:
		panic("unsupported punishment kind")
	}
}

func PunishmentKindFromString(punishmentKind string) chainindex.PunishmentKind {
	switch punishmentKind {
	case "ByzantineFault":
		return chainindex.PUNISHMENT_KIND_BYZANTINE_FAULT
	case "NonLive":
		return chainindex.PUNISHMENT_KIND_NON_LIVE
	default:
		panic("unsupported punishment kind")
	}
}

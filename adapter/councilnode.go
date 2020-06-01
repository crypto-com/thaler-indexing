package adapter

import (
	"github.com/crypto-com/chainindex"
	"github.com/luci/go-render/render"
)

type RDbCouncilNodeRow struct {
	ID                         *uint64 `json:"id"`
	Name                       string  `json:"name"`
	MaybeSecurityContact       *string `json:"security_contact"`
	PubKeyType                 string  `json:"pubkey_type"`
	PubKey                     string  `json:"pubkey"`
	Address                    string  `json:"address"`
	CreatedAtBlockHeight       uint64  `json:"created_at_block_height"`
	MaybeLastLeftAtBlockHeight *uint64 `json:"last_left_at_block_height"`
}

func (row *RDbCouncilNodeRow) String() string {
	return render.Render(row)
}

func RDbCouncilNodeRowToCouncilNode(row *RDbCouncilNodeRow) *chainindex.CouncilNode {
	if row == nil {
		return nil
	}
	return &chainindex.CouncilNode{
		Id:                         row.ID,
		Name:                       row.Name,
		MaybeSecurityContact:       row.MaybeSecurityContact,
		PubKeyType:                 PubKeyTypeFromString(row.PubKeyType),
		PubKey:                     row.PubKey,
		Address:                    row.Address,
		CreatedAtBlockHeight:       row.CreatedAtBlockHeight,
		MaybeLastLeftAtBlockHeight: row.MaybeLastLeftAtBlockHeight,
	}
}

func CouncilNodeToRDbCouncilNodeRow(node *chainindex.CouncilNode) *RDbCouncilNodeRow {
	if node == nil {
		return nil
	}
	return &RDbCouncilNodeRow{
		ID:                         node.Id,
		Name:                       node.Name,
		MaybeSecurityContact:       node.MaybeSecurityContact,
		PubKeyType:                 PubKeyTypeToString(node.PubKeyType),
		PubKey:                     node.PubKey,
		Address:                    node.Address,
		CreatedAtBlockHeight:       node.CreatedAtBlockHeight,
		MaybeLastLeftAtBlockHeight: node.MaybeLastLeftAtBlockHeight,
	}
}

func PubKeyTypeToString(pubKeyType chainindex.PubKeyType) string {
	switch pubKeyType {
	case chainindex.PUBKEY_TYPE_ED25519:
		return "ed25519"
	default:
		panic("unsupported pubkey type")
	}
}

func PubKeyTypeFromString(pubKeyType string) chainindex.PubKeyType {
	switch pubKeyType {
	case "ed25519":
		return chainindex.PUBKEY_TYPE_ED25519
	default:
		panic("unsupported pubkey type")
	}
}

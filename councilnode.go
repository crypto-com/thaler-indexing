package chainindex

import "github.com/luci/go-render/render"

type CouncilNode struct {
	Id                         *uint64
	Name                       string
	MaybeSecurityContact       *string
	PubKeyType                 PubKeyType
	PubKey                     string
	Address                    string
	CreatedAtBlockHeight       uint64
	MaybeLastLeftAtBlockHeight *uint64
}

type PubKeyType = int8

const (
	PUBKEY_TYPE_ED25519 PubKeyType = iota
)

func (node *CouncilNode) String() string {
	return render.Render(node)
}

type CouncilNodeUpdate struct {
	Address string
	Type    uint8
}

type CouncilNodeUpdateType = uint8

const (
	COUNCIL_NODE_UPDATE_TYPE_LEFT = iota
)

func (update *CouncilNodeUpdate) String() string {
	return render.Render(update)
}

type CouncilNodeRepository interface {
	Store(*CouncilNode) (id uint32, err error)
	FindById(id uint32) (*CouncilNode, error)
}

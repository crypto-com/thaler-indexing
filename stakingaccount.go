package chainindex

import (
	"math/big"
	"time"

	"github.com/luci/go-render/render"
)

type StakingAccount struct {
	Address                 string
	Nonce                   uint64
	Bonded                  *big.Int
	Unbonded                *big.Int
	MaybeUnbondedFrom       *time.Time
	MaybePunishmentKind     *PunishmentKind
	MaybeJailedUntil        *time.Time
	MaybeCurrentCouncilNode *CouncilNode
}

func (account *StakingAccount) String() string {
	return render.Render(account)
}

type StakingAccountRepository interface {
	Store(*StakingAccount) error
	FindByAddress(string) (*StakingAccount, error)
}

package adapter

import (
	"math/big"
	"time"

	"github.com/luci/go-render/render"
)

type RDbStakingAccountRow struct {
	Address              string
	Nonce                uint64
	Bonded               *big.Int
	Unbonded             *big.Int
	UnbondedFrom         *time.Time
	PunishmentKind       *string
	JailedUntil          *time.Time
	CurrentCouncilNodeId *uint64
}

func (row *RDbStakingAccountRow) IncrementNonce() {
	row.Nonce += 1
}

func (row *RDbStakingAccountRow) AddBonded(value *big.Int) {
	row.Bonded = new(big.Int).Add(row.Bonded, value)
}

func (row *RDbStakingAccountRow) AddUnbonded(value *big.Int) {
	row.Unbonded = new(big.Int).Add(row.Unbonded, value)
}

func (row *RDbStakingAccountRow) String() string {
	return render.Render(row)
}

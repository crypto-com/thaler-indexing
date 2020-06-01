package chainindex

import (
	"fmt"
	"math/big"
)

const (
	MAX_COIN_DECIMALS = 1_0000_0000
)

var (
	MAX_COIN_DECIMALS_BIGFLOAT = big.NewFloat(MAX_COIN_DECIMALS)
)

// TODO: A Coin wrapper struct

func MustCROToCoin(cro *big.Float) *big.Int {
	unit, err := CROToCoin(cro)
	if err != nil {
		panic(err)
	}

	return unit
}

func CROToCoin(cro *big.Float) (*big.Int, error) {
	unit, accuracy := cro.Mul(cro, MAX_COIN_DECIMALS_BIGFLOAT).Int(nil)
	if accuracy != big.Exact {
		return nil, fmt.Errorf("error converting %s CRO to unit: loss in precision", cro.String())
	}

	return unit, nil
}

package factory

import (
	"math/big"

	random "github.com/brianvoe/gofakeit/v5"
)

func RandomCoin() *big.Int {
	value := random.Number(1_0000_0000, 100_000_0000_0000)
	return big.NewInt(int64(value))
}

func RandomNegativeCoin() *big.Int {
	value := 0 - random.Number(1_0000_0000, 100_000_0000_0000)
	return big.NewInt(int64(value))
}

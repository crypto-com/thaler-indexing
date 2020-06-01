package bignum

import "math/big"

func Int0() *big.Int {
	return big.NewInt(int64(0))
}
func Int1() *big.Int {
	return big.NewInt(int64(1))
}
func Int10() *big.Int {
	return big.NewInt(int64(10))
}
func IntN1() *big.Int {
	return big.NewInt(int64(-1))
}
func IntN10() *big.Int {
	return big.NewInt(int64(-10))
}
func Int(i int64) *big.Int {
	return big.NewInt(i)
}

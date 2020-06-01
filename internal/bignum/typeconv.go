package bignum

import (
	"fmt"
	"math/big"
)

func OptItoa(i *big.Int) *string {
	if i == nil {
		return nil
	}
	s := i.String()
	return &s
}

func MustAtoi(s string) *big.Int {
	i, err := Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func Atoi(s string) (*big.Int, error) {
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("error converting string %s to big.Int", s)
	}
	return i, nil
}

func MustAtof(s string) *big.Float {
	f, err := Atof(s)
	if err != nil {
		panic(err)
	}
	return f
}

func Atof(s string) (*big.Float, error) {
	f, ok := new(big.Float).SetString(s)
	if !ok {
		return nil, fmt.Errorf("error converting string %s to big.Float", s)
	}
	return f, nil
}

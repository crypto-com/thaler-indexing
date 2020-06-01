package bignum

import (
	"math/big"

	jsoniter "github.com/json-iterator/go"
)

// JSON friendly big.Int
type WBigInt struct {
	*big.Int
}

func (w *WBigInt) MarshalJSON() ([]byte, error) {
	if w.Int == nil {
		return jsoniter.Marshal(nil)
	}

	s := w.String()
	if s == "<nil>" {
		return jsoniter.Marshal(nil)
	}

	return []byte("\"" + s + "\""), nil
}

func (bigInt *WBigInt) FromBigInt(i *big.Int) *WBigInt {
	bigInt.Int = i
	return bigInt
}

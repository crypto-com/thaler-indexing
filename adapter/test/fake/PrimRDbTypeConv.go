package adapterfake

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/internal/bignum"
)

type PrimRDbTypeConv struct{}

func (conv *PrimRDbTypeConv) Bton(b *big.Int) interface{} {
	if b == nil {
		return nil
	}
	s := b.String()
	return &s
}
func (conv *PrimRDbTypeConv) Iton(i int) interface{} {
	return strconv.Itoa(i)
}
func (conv *PrimRDbTypeConv) NtobReader() adapter.RDbNtobReader {
	return NewPrimRDbNtobReader()
}
func (conv *PrimRDbTypeConv) Tton(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t
}
func (conv *PrimRDbTypeConv) NtotReader() adapter.RDbNtotReader {
	return &PrimRDbNtotReader{}
}

type PrimRDbNtobReader struct {
	str *string
}

func NewPrimRDbNtobReader() *PrimRDbNtobReader {
	var s string
	return &PrimRDbNtobReader{
		str: &s,
	}
}

func (reader *PrimRDbNtobReader) ScannableArg() interface{} {
	return reader.str
}
func (reader *PrimRDbNtobReader) Parse() (*big.Int, error) {
	if reader.str == nil {
		return nil, nil
	}
	b, ok := new(big.Int).SetString(*reader.str, 10)
	if !ok {
		return nil, fmt.Errorf("error converting string %s to big.Int", *reader.str)
	}

	return b, nil
}
func (reader *PrimRDbNtobReader) ParseW() (*bignum.WBigInt, error) {
	var b bignum.WBigInt
	i, err := reader.Parse()
	if err != nil {
		return nil, err
	}
	return b.FromBigInt(i), nil
}

type PrimRDbNtotReader struct {
	t *time.Time
}

func (reader *PrimRDbNtotReader) ScannableArg() interface{} {
	return reader.t
}
func (reader *PrimRDbNtotReader) Parse() (*time.Time, error) {
	return reader.t, nil
}

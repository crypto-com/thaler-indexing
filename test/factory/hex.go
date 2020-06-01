package factory

import "crypto/rand"

func RandomHex(n int) []byte {
	placeholder := make([]byte, n)
	_, err := rand.Read(placeholder)
	if err != nil {
		panic(err)
	}
	return placeholder
}

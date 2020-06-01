package tendermint

import (
	"encoding/base64"

	"github.com/tendermint/tendermint/crypto/ed25519"
)

func AddressFromPubKey(base64PubKey string) string {
	var key ed25519.PubKeyEd25519
	data, _ := base64.StdEncoding.DecodeString(base64PubKey)
	copy(key[:], data[:])

	return key.Address().String()
}
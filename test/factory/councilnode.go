package factory

import (
	"encoding/base64"

	random "github.com/brianvoe/gofakeit/v5"

	"github.com/crypto-com/chainindex"
)

func RandomCouncilNodePtr() *chainindex.CouncilNode {
	node := RandomCouncilNode()
	return &node
}

func RandomCouncilNode() chainindex.CouncilNode {
	return chainindex.CouncilNode{
		Name:                       random.Company(),
		MaybeSecurityContact:       RandomEmailPtr(),
		PubKeyType:                 RandomPubKeyType(),
		PubKey:                     RandomPubKey(),
		Address:                    RandomTendermintAddress(),
		CreatedAtBlockHeight:       random.Uint64(),
		MaybeLastLeftAtBlockHeight: RandomUint64Ptr(),
	}
}

func RandomCouncilNodeLeftUpdate() chainindex.CouncilNodeUpdate {
	return chainindex.CouncilNodeUpdate{
		Address: RandomTendermintAddress(),
		Type:    chainindex.COUNCIL_NODE_UPDATE_TYPE_LEFT,
	}
}

func RandomPubKeyType() chainindex.PubKeyType {
	return chainindex.PUBKEY_TYPE_ED25519
}

func RandomPubKey() string {
	pubKey := RandomHex(20)
	return base64.StdEncoding.EncodeToString([]byte(pubKey))
}

func RandomTendermintAddress() string {
	address := RandomHex(32)
	return base64.StdEncoding.EncodeToString([]byte(address))
}

func RandomTendermintSignature() string {
	address := RandomHex(32)
	return base64.StdEncoding.EncodeToString([]byte(address))
}

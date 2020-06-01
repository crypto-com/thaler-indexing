package factory

import (
	"encoding/hex"
	"strings"

	random "github.com/brianvoe/gofakeit/v5"

	"github.com/crypto-com/chainindex"
	_ "github.com/crypto-com/chainindex/test/factory/init"
)

func RandomBlock() chainindex.Block {
	return chainindex.Block{
		Height:  random.Uint64(),
		Hash:    RandomBlockHash(),
		Time:    RandomUTCTime(),
		AppHash: RandomAppHash(),
	}
}

func RandomBlockHash() string {
	hash := hex.EncodeToString(RandomHex(32))
	return strings.ToUpper(hash)
}

func RandomAppHash() string {
	hash := hex.EncodeToString(RandomHex(32))
	return strings.ToUpper(hash)
}

func RandomBlockSignature() chainindex.BlockSignature {
	return chainindex.BlockSignature{
		BlockHeight:        random.Uint64(),
		CouncilNodeAddress: RandomTendermintAddress(),
		Signature:          RandomTendermintSignature(),
		IsProposer:         random.Bool(),
	}
}

func RandomBlockReward() chainindex.BlockReward {
	return chainindex.BlockReward{
		BlockHeight: random.Uint64(),
		Minted:      RandomCoin(),
	}
}

package viewrepo

import (
	"github.com/crypto-com/chainindex/internal/bignum"
)

type RewardViewRepo interface {
	TotalMinted() (*bignum.WBigInt, error)
	Total() (*bignum.WBigInt, error)
}

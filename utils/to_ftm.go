package utils

import "math/big"

// ToArt number of ART to Wei
func ToArt(art uint64) *big.Int {
	return new(big.Int).Mul(new(big.Int).SetUint64(art), big.NewInt(1e18))
}

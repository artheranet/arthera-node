package utils

import "math/big"

// ToArt number of AA to Wei
func ToArt(art uint64) *big.Int {
	return new(big.Int).Mul(new(big.Int).SetUint64(art), big.NewInt(1e18))
}

func WeiToArt(wei *big.Int) uint64 {
	return new(big.Int).SetUint64(wei.Uint64()).Div(wei, big.NewInt(1e18)).Uint64()
}

package DnsCompare

import (
	"crypto/rand"
	"math/big"
)

/*
const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1
*/

func RandomInt(min, max int) int {
	//TODO: verify min >= MinUint and max <= MaxInt
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)+int64(min)))
	return int(n.Uint64()) + min
}

func RandomFromDOTSlice(list []DOTSettings) DOTSettings {
	idx := RandomInt(0, len(list))
	return list[idx]
}
func RandomFromStringSlice(list []string) string {
	idx := RandomInt(0, len(list))
	return list[idx]
}

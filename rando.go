package DnsCompare

import (
	"crypto/rand"
	"math/big"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func RandomInt(min, max int) int {
	/*if min < MinUint {
		return 0
	}
	if max > MaxInt {
		return 0
	}
	*/
	newMax := big.NewInt(int64(max) - int64(min))
	if max-min <= 0 {
		return 0
	}
	n, _ := rand.Int(rand.Reader, newMax)
	return int(n.Uint64()) + min
}

func RandomFromDOTSlice(list []DOTSettings) DOTSettings {
	idx := RandomInt(0, len(list))
	return list[idx]
}
func RandomFromStringSlice(list []string) string {
	listLen := len(list)
	if listLen < 1 {
		return ""
	} else if listLen == 1 {
		return list[0]
	}
	idx := RandomInt(0, listLen)
	return list[idx]
}

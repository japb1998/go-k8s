package util

import (
	"math/rand/v2"
	"strings"
	"time"
)

var r *rand.Rand
var letter = "abcdefghijklmnopqrstuvwxyz"

func init() {
	src := rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()))
	r = rand.New(src)
}

func RandomInt(min, max int64) int64 {
	return min + r.Int64N(max-min+1) // min -> max
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(letter)

	for i := 0; i < n; i++ {
		c := letter[r.IntN(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[r.IntN(n)]
}

package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func randomInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)

}

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(k)]
		sb.WriteByte(ch)
	}

	return sb.String()
}

func RandomOwner() string {
	return randomString(10)
}

func RandomAmount() int64 {
	return randomInt64(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "GBP", "PLN"}
	return currencies[rand.Intn(len(currencies))]
}

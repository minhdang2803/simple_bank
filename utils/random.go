package utils

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

const alphabet string = "abcdefghijklmnopqrstuvwxyz"

func RandomOwner() string {
	return RandomString(6)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len((alphabet))
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomMoney() int64 {
	return RandomInt(0, 100)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}


// type Account struct {
// 	ID        int64     `json:"id"`
// 	Owner     string    `json:"owner"`
// 	Balance   int64     `json:"balance"`
// 	Currency  string    `json:"currency"`
// 	CreatedAt time.Time `json:"created_at"`
// }
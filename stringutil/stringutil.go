package stringutil

import (
	"math/rand"
)

// RandomString returns a 'random' string. Don't rely on this to have a secure
// level of entropy.
func RandomString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

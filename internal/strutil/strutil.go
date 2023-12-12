package strutil

import (
	"math/rand"
	"time"
)

// RandomNumericString returns a random numeric string with the given length.
func RandomNumericString(length int) string {
	const charset = "0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

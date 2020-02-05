package easyrandom

import (
	"crypto/rand"
	"math/big"
)

// RandomBytes []byte
func RandomBytes(size uint) []byte {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	return data
}

// RandomInt Returns int in range [min, max)
func RandomInt(min, max int64) int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		panic(err)
	}
	return n.Int64() + min
}

// RandomBool true or false
func RandomBool() bool {
	n := RandomInt(0, 2)
	return n == 1
}

package util

import (
	crand "crypto/rand"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcutil/base58"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var (
	randBuffer = [4]byte{}
	randMutex  = sync.Mutex{}
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	return RandomStringWithCharset(length, alphabet)
}

func RandomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateUniqueID() string {
	result := make([]byte, 10)

	randMutex.Lock()
	defer randMutex.Unlock()

	_, _ = io.ReadAtLeast(crand.Reader, randBuffer[:], len(randBuffer))
	copy(result[6:], randBuffer[:])

	PutUint48(result[:6], uint64(time.Now().UnixNano()/1e6))

	return string(base58.Encode(result))
}

func PutUint48(b []byte, v uint64) {
	_ = b[5] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 8)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 24)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 40)
	b[5] = byte(v >> 48)
}

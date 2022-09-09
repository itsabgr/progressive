package span

import (
	crand "crypto/rand"
	"encoding/binary"
	"io"
	srand "math/rand"
)

func must[T any](t T, e any) T {
	throw(e)
	return t
}
func throw(e any) {
	if e != nil {
		panic(e)
	}
}
func randbytes(length int) []byte {
	b := make([]byte, length)
	must(io.ReadFull(crand.Reader, b))
	return b
}

func compare[T number](t1, t2 T) int {
	if t1 < t2 {
		return -1
	}
	if t1 == t2 {
		return 0
	}
	return 1
}

func min[T number](a, b T) T {
	if a > b {
		return b
	}
	return a
}

func randint(max int) int {
	n := binary.BigEndian.Uint64(randbytes(8))
	srand.Seed(int64(n))
	return srand.Intn(max)
}

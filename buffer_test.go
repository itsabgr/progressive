package span

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"io"
	srand "math/rand"
	"testing"
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

func concat[S ~[]E, E any](slices ...S) []E {
	length := 0
	for _, item := range slices {
		length += len(item)
	}
	if length == 0 {
		return nil
	}
	result := make([]E, 0, length)
	for _, item := range slices {
		result = append(result, item...)
	}
	return result
}

func sum[E any](slices [][]E) int {
	s := 0
	for _, slice := range slices {
		s += len(slice)
	}
	return s
}
func randint(max int) int {
	n := binary.BigEndian.Uint64(randbytes(8))
	srand.Seed(int64(n))
	return srand.Intn(max)
}

func shuffle[E any](a []E) {
	for i := len(a) - 1; i > 0; i-- {
		j := randint(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func TestBuffer(t *testing.T) {
	chunks := make([][]byte, 1000)
	for i := range chunks {
		chunks[i] = randbytes(randint(512) + 1)
	}
	buff := new(Buffer)
	go func() {
		order := make([]int, len(chunks))
		for i := range order {
			order[i] = i
		}
		shuffle(order)
		for _, i := range order {
			_, err := buff.WriteAt(chunks[i], sum(chunks[:i]))
			throw(err)
		}
		buff.Close()
	}()
	all := must(io.ReadAll(buff))
	if len(all) == 0 || !bytes.Equal(all, concat(chunks...)) {
		t.FailNow()
	}
}

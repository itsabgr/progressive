package span

import (
	"bytes"
	"io"
	"testing"
)

func TestStream(t *testing.T) {
	data := randbytes(1_000_000)
	buff := new(Stream)
	defer func() { _ = buff.Close() }()
	t.Parallel()
	go func() {
		defer t.Log("writer closed")
		start := 0
		for {
			if start == len(data) {
				return
			}
			end := start + randint(min(512, len(data)-start)) + 1
			_, err := buff.WriteAt(data[start:end], start)
			if err == ErrClosed {
				return
			}
			_, err = buff.WriteAt(data[start:end], start)
			if err == ErrClosed {
				return
			}
			if err != nil {
				t.Log(err)
			}
			start = buff.Missing()
		}
	}()
	all := make([]byte, len(data))
	must(io.ReadFull(buff, all))
	if len(all) == 0 || !bytes.Equal(all, data) {
		t.FailNow()
	}
}

package span

import (
	"io"
	"os"
	"runtime"
	"sync"
)

type Buffer struct {
	mutex  sync.Mutex
	puzzle []Span[int]
	buff   []byte
	off    int
	closed bool
}

func (b *Buffer) WriteAt(data []byte, off int) (int, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.closed {
		return 0, os.ErrClosed
	}
	//
	b.puzzle = Span[int]{off, off + len(data)}.AddTo(b.puzzle)
	if required := (len(data) + off) - len(b.buff); required > 0 {
		b.buff = append(b.buff, make([]byte, required)...)
	}
	return copy(b.buff[off-b.off:], data), nil
}

func (b *Buffer) Close() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.closed {
		return os.ErrClosed
	}
	b.closed = true
	return nil
}

func (b *Buffer) Read(dst []byte) (int, error) {
	for {
		b.mutex.Lock()
		for _, s := range b.puzzle {
			if s.Contains(b.off) {
				n := copy(dst, b.buff[:s.End-b.off])
				b.off += n
				b.buff = b.buff[n:]
				b.puzzle = Span[int]{0, b.off}.SubFrom(b.puzzle)
				b.mutex.Unlock()
				return n, nil
			}
		}
		if b.closed {
			b.mutex.Unlock()
			return 0, io.EOF
		}
		b.mutex.Unlock()
		runtime.Gosched()
	}
}

package span

import (
	"io"
	"runtime"
	"sync"
)

var ErrClosed = io.ErrClosedPipe

type Stream struct {
	mutex  sync.Mutex
	puzzle []span[int]
	buff   []byte
	off    int
	closed bool
}

func (stream *Stream) WriteAt(data []byte, off int) (int, error) {
	stream.mutex.Lock()
	defer stream.mutex.Unlock()
	if stream.closed {
		return 0, ErrClosed
	}
	if len(data) == 0 {
		return 0, nil
	}
	if off < stream.off {
		return len(data), nil
	}
	if required := len(data) + off - len(stream.buff); required > 0 {
		stream.buff = append(stream.buff, make([]byte, required)...)
	}
	stream.puzzle = span[int]{off, off + len(data)}.AddTo(stream.puzzle)
	return copy(stream.buff[off-stream.off:], data), nil
}
func (stream *Stream) Len() int {
	stream.mutex.Lock()
	defer stream.mutex.Unlock()
	if stream.closed {
		return -1
	}
	for _, s := range stream.puzzle {
		if s.Contains(stream.off) {
			return s.End - stream.off
		}
	}
	return 0
}
func (stream *Stream) Close() error {
	stream.mutex.Lock()
	defer stream.mutex.Unlock()
	if stream.closed {
		return ErrClosed
	}
	stream.closed = true
	return nil
}

func (stream *Stream) Missing() int {
	stream.mutex.Lock()
	defer stream.mutex.Unlock()
	if stream.closed {
		return -1
	}
	if len(stream.puzzle) == 0 {
		return 0
	}
	if stream.puzzle[0].Start != 0 {
		return 0
	}
	return stream.puzzle[0].End
}

func (stream *Stream) Read(dst []byte) (int, error) {
	for {
		stream.mutex.Lock()
		for _, s := range stream.puzzle {
			if s.Contains(stream.off) {
				n := copy(dst, stream.buff[:s.End-stream.off])
				stream.off += n
				stream.buff = stream.buff[n:]
				//	stream.puzzle = span[int]{0, stream.off}.SubFrom(stream.puzzle)
				stream.mutex.Unlock()
				return n, nil
			}
		}
		if stream.closed {
			stream.mutex.Unlock()
			return 0, io.EOF
		}
		stream.mutex.Unlock()
		runtime.Gosched()
	}
}

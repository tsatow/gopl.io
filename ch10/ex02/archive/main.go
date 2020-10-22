package archive

import (
	"bufio"
	"io"
	"sync"
	"sync/atomic"
)

type format struct {
	name, magic  string
	decode       func(io.Reader) *io.Reader
}

type peekableReader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

var (
	formatsMu     sync.Mutex
	atomicFormats atomic.Value
)

func RegisterFormat(name, magic string, decode func(io.Reader) *io.Reader) {
	formatsMu.Lock()
	formats, _ := atomicFormats.Load().([]format)
	atomicFormats.Store(append(formats, format{name, magic, decode}))
	formatsMu.Unlock()
}

func asPeekableReader(r io.Reader) peekableReader {
	if rr, ok := r.(peekableReader); ok {
		return rr
	}
	return bufio.NewReader(r)
}

func File() {

}

func match(magic string, b []byte) bool {
	if len(magic) != len(b) {
		return false
	}
	for i, c := range b {
		if magic[i] != c && magic[i] != '?' {
			return false
		}
	}
	return true
}
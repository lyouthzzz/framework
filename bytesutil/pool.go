package bytesutil

import (
	"bytes"
	"sync"
)

type Pool struct {
	p *sync.Pool
}

func (p Pool) Get() *bytes.Buffer {
	buf := p.p.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (p Pool) Put(buf *bytes.Buffer) {
	p.p.Put(buf)
}

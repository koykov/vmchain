package vmchain

import (
	"github.com/koykov/byteconv"
	"github.com/koykov/x2bytes"
)

type builder struct {
	buf []byte
	lc  int
}

func (b *builder) setName(name string) {
	b.reset()
	b.buf = append(b.buf, name...)
}

func (b *builder) setLabel(label, value string) {
	if b.lc == 0 {
		b.buf = append(b.buf, '{')
	} else {
		b.buf = append(b.buf, ',')
	}
	b.buf = append(b.buf, label...)
	b.buf = append(b.buf, `="`...)
	b.buf = append(b.buf, value...)
	b.buf = append(b.buf, '"')
	b.lc++
}

func (b *builder) setAnyLabel(label string, value any) {
	if b.lc == 0 {
		b.buf = append(b.buf, '{')
	} else {
		b.buf = append(b.buf, ',')
	}
	b.buf = append(b.buf, label...)
	b.buf = append(b.buf, `="`...)

	if value != nil {
		var err error
		if b.buf, err = x2bytes.ToBytes(b.buf, value); err != nil {
			b.buf = append(b.buf, err.Error()...)
		}
	} else {
		b.buf = append(b.buf, "<nil>"...)
	}

	b.buf = append(b.buf, '"')
	b.lc++
}

func (b *builder) commit() string {
	if b.lc > 0 {
		b.buf = append(b.buf, '}')
	}
	return byteconv.B2S(b.buf)
}

func (b *builder) reset() {
	b.buf = b.buf[:0]
	b.lc = 0
}

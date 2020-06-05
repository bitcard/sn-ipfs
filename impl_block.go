package ipfs_filestore

import (
	"errors"
	"io"
)

type block struct {
	node
}

func (b block) Type() Type {
	return BLK
}

func (b block) Reader() io.ReadSeeker {
	return &bytesReaderSeeker{b: b.head.Data()}
}

func (b block) Blocks() []Block {
	return []Block{Block(b)}
}

type bytesReaderSeeker struct {
	b     []byte
	index int
}

func (b *bytesReaderSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		b.index = int(offset)
	case io.SeekCurrent:
		if offset+int64(b.index) != 0 {
			return int64(b.index), errors.New("EOF")
		}
		b.index = b.index + int(offset)
	case io.SeekEnd:
		if offset != 0 {
			return int64(b.index), errors.New("EOF")
		}
		b.index = len(b.b)
	}
	return 0, nil
}

func (b *bytesReaderSeeker) Read(p []byte) (n int, err error) {
	dstCap := cap(p)
	src := b.b[b.index:]
	if dstCap >= len(src) {
		copy(p, src)
		return len(src), nil
		b.index = b.index + len(src)
	}
	copy(p, src[:dstCap])
	b.index = b.index + dstCap
	return dstCap, nil
}

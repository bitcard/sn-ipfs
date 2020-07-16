package ipfs

import (
	"errors"
	ipld "github.com/ipfs/go-ipld-format"
	"io"
)

type File interface {
	io.ReadSeeker
	// 获取block数据
	Block(index int) Block
	Blocks() []Block
	Node
}

func newFile(n Node, s Store) *file {
	return &file{
		store: s,
		Node:  n,
	}
}

type file struct {
	store Store
	// 当前读取的位置
	index uint64
	// 缓存
	caches []byte
	Node
}

func (f *file) Read(p []byte) (n int, err error) {
	var maxBlockSize uint64 = 262144
	// 已经读完
	if f.index == f.Size() {
		return 0, io.EOF
	}
	// 初始化缓存
	if f.caches == nil {
		f.caches = make([]byte, maxBlockSize)
	}
	// 计算block数据位置
	blkDataIndex := f.index % maxBlockSize
	// 当前块已经读完,拷贝新内容
	if blkDataIndex == 0 && f.index < f.Size() {
		blkIndex := f.index / maxBlockSize
		n = copy(f.caches, f.Blocks()[blkIndex].Data())
		f.caches = f.caches[:n]
	}
	n = copy(p, f.caches[blkDataIndex:])
	f.index += uint64(n)
	return n, nil
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	var index int64
	switch whence {
	case io.SeekStart:
		if offset < 0 {
			return 0, errors.New("index out of range")
		}
		index = offset
		f.index = uint64(offset)
	case io.SeekCurrent:
		index = int64(f.index) + offset
		if index < 0 || index > int64(f.index) {
			return 0, errors.New("index out of range")
		}
		f.index = uint64(index)
	case io.SeekEnd:
		index = int64(f.Size()) + offset
		if index < 0 || index > int64(f.index) {
			return 0, errors.New("index out of range")
		}
		f.index = uint64(index)
	default:
		return 0, errors.New("unknow whence")
	}
	return index, nil
}

func (f *file) Block(index int) Block {
	node := f.store.Get(f.Links()[index])
	return &block{Node: node}
}

func (f *file) Blocks() []Block {
	links := f.Links()
	if len(links) == 0 {
		links = []*ipld.Link{newLink(f.Cid())}
	}
	var blocks = make([]Block, 0, len(links))
	nodes := f.store.GetMany(links)
	for _, node := range nodes {
		blocks = append(blocks, &block{Node: node})
	}
	return blocks
}

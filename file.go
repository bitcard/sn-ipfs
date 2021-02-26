package ipfs

import (
	"errors"
	"io"
)

type File interface {
	io.ReadSeeker
	// 获取block数据
	Block(index int) Block
	Blocks() []Block
	BaseNode
}

func newFile(n Node, s *store) *file {
	return &file{
		store: s,
		Node:  n,
	}
}

type file struct {
	store *store
	// 当前读取的位置
	index uint64
	// 缓存
	caches []byte
	// block缓存
	children []File
	// TODO:使用currentChild简化寻找file
	currentChild int64
	Node
}

func (f *file) Read(p []byte) (int, error) {
	// 已经读完
	if f.index == f.Size() {
		f.unload()
		return 0, io.EOF
	}
	// 如果当前文件是叶节点
	if f.isLeafNode() {
		if f.caches == nil {
			f.load()
		}
		n := copy(p, f.caches[f.index:])
		f.index = f.index + uint64(n)
		return n, nil
	}
	// TODO：同一级的节点data大小应该是一样的
	// 那么应该可以通过计算直接得到节点位置
	// 如果不是叶节点，寻找正在读取的节点位置
	if f.children == nil {
		f.load()
	}
	var index = f.index
	var currentChild File
	for _, cf := range f.children {
		if index < cf.Size() {
			currentChild = cf
			break
		}
		index = index - cf.Size()
	}
	// 找到block，转化为文件，开始读写
	n, err := currentChild.Read(p)
	if err == nil || err == io.EOF {
		f.index = f.index + uint64(n)
		return n, nil
	}
	return n, err
}

// TODO:调整seek
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

func (f *file) load() {
	// 尝试加载子节点
	nodes := f.Node.Children()
	if len(nodes) != 0 {
		f.children = make([]File, 0, len(nodes))
		for _, v := range nodes {
			cf, err := v.ToFile()
			if err != nil {
				panic(err)
			}
			f.children = append(f.children, cf)
		}
	}
	// 尝试加载数据
	if data := f.Node.(*node).data(); data != nil {
		// 初始化缓存，肯定已经加载完成了。
		if f.caches == nil {
			f.caches = make([]byte, len(data))
			n := copy(f.caches[:], data)
			if n != len(data) {
				panic("无法获取到完整的块数据")
			}
		}
	}
}

func (f *file) isLeafNode() bool {
	return f.Node.(*node).data() != nil
}

// 卸载数据
func (f *file) unload() {
	f.caches = nil
	f.children = nil
}

func (f *file) Block(index int) Block {
	// 首先要查看是否是leafParents
	if f.isLeafNode() {
		return nil
	}
	panic("implement me")
}

func (f *file) Blocks() []Block {
	if f.isLeafNode() {
		return nil
	}
	panic("implement me")
}

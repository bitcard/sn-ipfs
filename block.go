package ipfs

// Block肯定是一个文件,而且肯定有数据
type Block interface {
	BaseNode
	File() File
	Data() []byte
}

type block struct {
	index uint64
	Node
}

func (b *block) Data() []byte {
	return b.Node.(*node).data()
}

func (b *block) File() File {
	f, err := b.Node.ToFile()
	if err != nil {
		panic(err)
	}
	return f
}

func (b *block) Type() Type {
	return BLK
}

package ipfs

// Block肯定是一个文件
type Block interface {
	BaseNode
	Data() []byte
}

type block struct {
	index uint64
	Node
}

func (b *block) Data() []byte {
	return b.Node.(*node).data()
}

func (b *block) Type() Type {
	return BLK
}

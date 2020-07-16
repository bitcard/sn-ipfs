package ipfs

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

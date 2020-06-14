package ipfs_filestore

type Block interface {
	Node
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

package ipfs_filestore

type Block interface {
	Node
}

type block struct {
	Node
}

func (n block) Type() Type {
	return BLK
}

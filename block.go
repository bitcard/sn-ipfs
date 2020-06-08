package ipfs_filestore

type Block interface {
	Node
	ToFile() File
}

type block struct {
	Node
}

func (n block) Type() Type {
	return BLK
}

func (n block) ToFile() File {
	return NewFile(n.Node)
}

package ipfs_filestore

type Dir interface {
	Nodes() []Node
	Node
}

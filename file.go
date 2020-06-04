package ipfs_filestore

import "io"

type File interface {
	io.Reader
	Blocks() []Block
	Node
}

package ipfs_filestore

import "io"

type File interface {
	Reader() io.ReadSeeker
	Blocks() []Block
	Node
}

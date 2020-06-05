package ipfs_filestore

import "io"

type Block interface {
	Reader() io.ReadSeeker
	Node
}

package ipfs_filestore

import "io"

type Block interface {
	io.Reader
	Node
}

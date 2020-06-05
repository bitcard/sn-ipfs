package ipfs_filestore

import (
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-unixfs"
)

type node struct {
	cid   string
	name  string
	tp    Type
	head  *unixfs.FSNode
	links []*ipld.Link
}

func (n node) Name() string {
	return n.name
}

func (n node) Cid() string {
	return n.cid
}

func (n node) Type() Type {
	return n.tp
}

func (n node) Size() uint64 {
	return n.head.FileSize()
}

package ipfs_filestore

import (
	ipld "github.com/ipfs/go-ipld-format"
	"io"
)

type File interface {
	Reader() io.ReadSeeker
	Blocks() []Block
	Node
}

type file struct {
	Node
}

func (f file) Reader() io.ReadSeeker {
	panic("implement me")
}

func (f file) Blocks() []Block {
	links := f.Links()
	if len(links) == 0 {
		links = []*ipld.Link{newLink(f.Cid())}
	}
	var blocks = make([]Block, 0, len(links))
	nodes := Gstore.GetMany(links)
	for _, node := range nodes {
		blocks = append(blocks, block{node})
	}
	return blocks
}

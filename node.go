package ipfs_filestore

import (
	"errors"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-unixfs"
)

type Type uint

const (
	BLK Type = iota
	DIR
	FIL
	UNW
)

type Node interface {
	Name() string
	Cid() string
	Type() Type
	Size() uint64
	RawSize() uint64
	Links() []*ipld.Link
	ToFile() (File, error)
	ToDir() (Dir, error)
}

func NewNode(link *ipld.Link) Node {
	return &node{
		cid:    link.Cid.String(),
		name:   link.Name,
		inited: false,
	}
}

type node struct {
	cid    string
	name   string
	raw    uint64
	tp     Type
	head   *unixfs.FSNode
	links  []*ipld.Link
	inited bool
}

func (n *node) load() {
	if n.inited {
		return
	}
	pn, err := Gstore.(*store).get(newLink(n.cid))
	head, err := unixfs.FSNodeFromBytes(pn.Data())
	if err != nil {
		return
	}
	raw, _ := pn.Marshal()
	n.raw = uint64(len(raw))
	n.head = head
	n.links = pn.Links()
	n.inited = true
}

func (n *node) Name() string {
	return n.name
}

func (n *node) Cid() string {
	return n.cid
}

func (n *node) Type() Type {
	n.load()
	switch n.head.Type() {
	case 1:
		return DIR
	case 2:
		return FIL
	default:
		return UNW
	}
}

func (n *node) Size() uint64 {
	n.load()
	return n.head.FileSize()
}

func (n *node) RawSize() uint64 {
	n.load()
	return n.raw
}

func (n *node) ToFile() (File, error) {
	if n.Type() != FIL && n.Type() != BLK {
		return nil, errors.New("node not a file")
	}
	return file{n}, nil
}

func (n *node) ToDir() (Dir, error) {
	if n.Type() != DIR {
		return nil, errors.New("node not a dir")
	}
	return &dir{nil, n}, nil
}

func (n *node) Links() []*ipld.Link {
	n.load()
	return n.links
}

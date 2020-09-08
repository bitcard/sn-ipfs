package ipfs

import (
	"errors"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-unixfs"
)

type Type int32

const (
	BLK Type = iota
	DIR
	FIL
	UNW
)

func (t Type) String() string {
	switch t {
	case BLK:
		return "Block"
	case DIR:
		return "Directory"
	case FIL:
		return "File"
	default:
		return "Unknown"
	}
}

type BaseNode interface {
	Name() string
	Cid() string
	Type() Type
	Size() uint64
	RawSize() uint64
}

type Node interface {
	BaseNode
	Links() []*ipld.Link
	ToFile() (File, error)
	ToDir() (Dir, error)
}

func newNode(link *ipld.Link, s *store) *node {
	return &node{
		store:  s,
		cid:    link.Cid.String(),
		name:   link.Name,
		inited: false,
	}
}

type node struct {
	store  *store
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
	pn, err := n.store.getProtoNode(newLink(n.cid))
	if pn == nil {
		return
	}
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

func (n *node) data() []byte {
	n.load()
	if n.head == nil {
		return nil
	}
	return n.head.Data()
}

func (n *node) Type() Type {
	n.load()
	if n.head == nil {
		return UNW
	}
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
	if n.head == nil {
		return 0
	}
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
	return newFile(n, n.store), nil
}

func (n *node) ToDir() (Dir, error) {
	if n.Type() != DIR {
		return nil, errors.New("node not a dir")
	}
	return newDir(n, n.store), nil
}

func (n *node) Links() []*ipld.Link {
	n.load()
	return n.links
}

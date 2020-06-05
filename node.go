package ipfs_filestore

type Type uint

const (
	BLK Type = iota
	DIR
	FIL
)

type Node interface {
	Name() string
	Cid() string
	Type() Type
	Size() uint64
}

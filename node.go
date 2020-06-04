package ipfs_filestore

type Type uint

const (
	FIL Type = iota
	DIR
	BLK
)

type Node interface {
	Name() string
	Cid() string
	Type() Type
	Size() int64
}

package ipfs_filestore

type file struct {
	node
}

func (f file) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (f file) Blocks() []Block {
	return
}

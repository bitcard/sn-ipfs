package ipfs_filestore

type Dir interface {
	Nodes() []Node
	Node
}

type dir struct {
	nodes []Node
	Node
}

// 谨慎执行，时间会比较久
func (d *dir) Size() uint64 {
	d.loadNodes()
	var size uint64
	for _, n := range d.nodes {
		switch n.Type() {
		case FIL:
			size = size + n.Size()
		case DIR:
			size = size + newDir(n).Size()
		default:
			continue
		}
	}
	return size
}

func (d *dir) loadNodes() {
	if d.nodes == nil {
		for _, link := range d.Links() {
			d.nodes = append(d.nodes, Gstore.Get(link))
		}
	}
}

func (d *dir) Nodes() []Node {
	d.loadNodes()
	return d.nodes
}

func newDir(n Node) Dir {
	return &dir{nodes: nil, Node: n}
}

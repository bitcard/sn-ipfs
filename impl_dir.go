package ipfs_filestore

type dir struct {
	nodes []Node
	node
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
			size = size + newDir(n.(node)).Size()
		default:
			continue
		}
	}
}

func (d *dir) loadNodes() {
	if d.nodes == nil {
		for _, link := range d.node.links {
			d.nodes = append(d.nodes, GStore.Get(link.Cid.String()))
		}
	}
}

func (d *dir) Nodes() []Node {
	d.loadNodes()
	return d.nodes
}

func newDir(n node) Dir {
	return &dir{nodes: nil, node: n}
}

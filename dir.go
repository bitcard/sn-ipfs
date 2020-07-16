package ipfs

type Dir interface {
	Nodes() []Node
	Node
}

type dir struct {
	store Store
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
			size = size + newDir(n, d.store).Size()
		default:
			continue
		}
	}
	return size
}

func (d *dir) loadNodes() {
	if d.nodes == nil {
		for _, link := range d.Links() {
			d.nodes = append(d.nodes, d.store.Get(link))
		}
	}
}

func (d *dir) Nodes() []Node {
	d.loadNodes()
	return d.nodes
}

func newDir(n Node, s Store) Dir {
	return &dir{store: s, nodes: nil, Node: n}
}

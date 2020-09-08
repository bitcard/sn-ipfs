package ipfs

type Dir interface {
	Nodes() []Node
	AddFile(f File) (Dir, error)
	BaseNode
}

type dir struct {
	store *store
	nodes []Node
	BaseNode
}

func (d *dir) AddFile(f File) (Dir, error) {
	nodes := append(d.nodes, f.(*file).Node)
	node, err := d.store.combine(DIR, nodes)
	if err != nil {
		return nil, err
	}
	return node.ToDir()
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
		node := d.BaseNode.(Node)
		for _, link := range node.Links() {
			d.nodes = append(d.nodes, d.store.get(link))
		}
	}
}

func (d *dir) Nodes() []Node {
	d.loadNodes()
	return d.nodes
}

func newDir(n Node, s *store) Dir {
	return &dir{store: s, nodes: nil, BaseNode: n}
}

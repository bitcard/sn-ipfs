package ipfs_filestore

import (
	client "github.com/ipfs/go-ipfs-api"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	"io"
)

func NewStore(url string) (Store, error) {
	api := client.NewShell(url)
	_, _, err := api.Version()
	if err != nil {
		return nil, err
	}
	return &store{
		api: api,
	}, nil
}

type store struct {
	api *client.Shell
}

func (s *store) AddFromReader(reader io.Reader) {
	s.api.BlockGet()
}

func (s *store) AddFromBytes(bytes []byte) {
	panic("implement me")
}

func (s *store) Pin(node Node) error {
	panic("implement me")
}

func (s *store) Get(link *ipld.Link) (Node, error) {
	cid := link.Cid.String()
	data, err := s.api.BlockGet(cid)
	if err != nil {
		return nil, err
	}
	pn, err := merkledag.DecodeProtobuf(data)
	if err != nil {
		return nil, err
	}
	head, err := unixfs.FSNodeFromBytes(pn.Data())
	if err != nil {
		return nil, err
	}
	return node{
		name:  link.Name,
		cid:   cid,
		head:  head,
		links: pn.Links(),
	}, nil
}

// TODO: 多线程处理
func (s *store) GetMany(links []*ipld.Link) ([]Node, error) {
	var nodes = make([]Node, 0, len(cids))
	for _, link := range links {
		node, err := s.Get(link)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (s *store) Combine(blocks []Block) {
	//cl := newShell()
	//data,_:= cl.BlockGet(cid2)
	//node,_ := merkledag.DecodeProtobuf(data)
	//fack := merkledag.NodeWithData(node.Data())
	//for  _,link := range node.Links() {
	//	fack.AddRawLink(link.Name,link)
	//}
	//cid,err := cl.BlockPut(fack.RawData(),"",mh.Codes[mh.SHA2_256],-1)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(cid)
}

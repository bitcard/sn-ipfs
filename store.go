package ipfs_filestore

import (
	"bytes"
	"errors"
	client "github.com/ipfs/go-ipfs-api"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	mh "github.com/multiformats/go-multihash"
	"io"
)

func init() {
	s, err := NewStore("http://127.0.0.1:5001")
	if err != nil {
		panic(err)
	}
	Gstore = s
}

type Store interface {
	AddFromReader(io.Reader) (File, error)   // node,err 从reader对象中读取创建node
	AddFromBytes(bytes []byte) (File, error) // node,err 从字节数组中读取创建node
	Pin(blk Block) error                     // 固定文件，长期保存，为了更好的存储必须是block
	PinMany(blks []Block) error
	Get(link *ipld.Link) Node // 获取node
	GetMany(links []*ipld.Link) []Node
	Combine([]Block) (File, error) // node 按照顺序组合文件块
}

var Gstore Store

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

func (s *store) AddFromReader(reader io.Reader) (File, error) {
	cid, err := s.api.Add(reader)
	if err != nil {
		return nil, err
	}
	return NewFile(s.Get(newLink(cid))), nil
}

func (s *store) AddFromBytes(data []byte) (File, error) {
	return s.AddFromReader(bytes.NewReader(data))
}

func (s *store) Pin(blk Block) error {
	if blk.Type() != BLK {
		return errors.New("not a block")
	}
	return s.api.Pin(blk.Cid())
}

func (s *store) PinMany(blks []Block) error {
	var err error
	for _, blk := range blks {
		err = s.Pin(blk)
		if err != nil {
			return err
		}
	}
	return err
}

func (s *store) get(link *ipld.Link) (*merkledag.ProtoNode, error) {
	cid := link.Cid.String()
	data, err := s.api.BlockGet(cid)
	if err != nil {
		return nil, err
	}
	return merkledag.DecodeProtobuf(data)
}

func (s *store) Get(link *ipld.Link) Node {
	return NewNode(link)
}

// TODO: 多线程处理
func (s *store) GetMany(links []*ipld.Link) []Node {
	var nodes = make([]Node, 0, len(links))
	for _, link := range links {
		node := s.Get(link)
		nodes = append(nodes, node)
	}
	return nodes
}

func (s *store) Combine(blocks []Block) (File, error) {
	head := unixfs.NewFSNode(2)
	newNode := merkledag.NodeWithData(nil)
	for _, blk := range blocks {
		newNode.AddRawLink("", newFullLink(blk.Cid(), "", blk.RawSize()))
		head.AddBlockSize(blk.Size())
	}
	data, err := head.GetBytes()
	if err != nil {
		return nil, err
	}
	newNode.SetData(data)
	cid, err := s.api.BlockPut(newNode.RawData(), "", mh.Codes[mh.SHA2_256], -1)
	if err != nil {
		panic(err)
	}
	node := s.Get(newLink(cid))
	return NewFile(node), nil
}
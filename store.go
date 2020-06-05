package ipfs_filestore

import (
	ipld "github.com/ipfs/go-ipld-format"
	"io"
)

var GStore Store

type Store interface {
	AddFromReader(io.Reader)           // node,err 从reader对象中读取创建node
	AddFromBytes(bytes []byte)         // node,err 从字节数组中读取创建node
	Pin(Node) error                    // 固定文件，长期保存
	Get(link *ipld.Link) (Node, error) // 获取node
	GetMany(links []*ipld.Link) ([]Node, error)
	Combine([]*ipld.Link) (Node, error) // node 按照顺序组合文件块
}

package ipfs_filestore

import "io"

type Store interface {
	AddFromReader(io.Reader)   // node,err 从reader对象中读取创建node
	AddFromBytes(bytes []byte) // node,err 从字节数组中读取创建node
	Pin(Node) error            // 固定文件，长期保存
	Get(cid string)            // 获取node
	Combine([]Block)           // node 按照顺序组合文件块
}

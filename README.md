# ipfsFileStore

## 约定

1. 文件对象和文件夹对象都是不可变的，只可读，但是可以通过添加文件产生新的对象
2. 能够获取到文件下载链接，同时也能够直接读取内容
3. 可以访问底层文件块，node无权限修改，但是store可以

## 抽象

### store：仓库

```go
type Store interface {
  AddFromReader(io.reader) // node,err 从reader对象中读取创建node
  AddFromBytes([]byte)     // node,err 从字节数组中读取创建node
  Pin(node) err            // 固定文件，长期保存
  Get(cid stirng)          // 获取node
  Combine([]block)         // node 按照顺序组合文件块
}
```



### Node：文件树节点

node是所有底层文件的公共接口，block，file，dir都要满足它，可以通过type来获取底层类型

```go
type Node interface {
  Name() string
  Cid() string
  Type() type
  Size() int64
}
```



### Dir：文件夹类型

文件夹类型不包含任何自己的数据，只是链接的集合。它不能直接被读取，只能获取内部的Node对象，这次node对象只会是file或者dir类型。

```go
type Dir interface {
  Nodes() NodeIterator
  Node
}
```



### File：文件类型

file是可以被识别的文件，这些文件一般是通过store仓库直接添加的。block可以转化为file，进行一些基本操作，但是file不能类型转化为block，但是可以通过

```go
type File interface {
  io.Reader
  Blocks() BlockIterator
  Node
}
```



### Block：底层存储块

block并不是ipfs层面的block，它一定包含数据，是数据存储的最小单元，可以进行拼接操作，可以转化为file。这种类型的对象应该只由file的block接口产生，普通dir和file操作不会产生block类型的node

```go
type Block interface {
  io.Reader
  Node
}
```




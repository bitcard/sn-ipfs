# ipfsFileStore

## 约定

1. 文件对象和文件夹对象都是不可变的，只可读，但是可以通过添加文件产生新的对象
2. 能够获取到文件下载链接，同时也能够直接读取内容
3. 可以访问底层文件块，node无权限修改，但是store可以

## 一些概念

### store：仓库

store提供一些获取node和一些操作node的方法。主要是负责管理ipfs的行为的模块

### Node：文件树节点

node是所有底层文件的公共接口，block，file，dir都要满足它，可以通过type来获取底层类型。

### Dir：文件夹类型

文件夹类型不包含任何自己的数据，只是链接的集合。它不能直接被读取，只能获取内部的Node对象，获取到的node对象只会是file或者dir类型。


### File：文件类型

file是可以被识别的文件，这些文件一般是通过store仓库直接添加的。block可以转化为file，进行一些基本操作。只要file才能被读取。


### Block：底层存储块

block并不是ipfs层面的block，它一定包含数据，是数据存储的最小单元，可以进行拼接操作，可以转化为file。这种类型的对象应该只由file的blocks方法产生，普通dir和file操作不会产生block类型的node

### example
```go
    const(
	apiAddr     = "127.0.0.1:5001"
	gatewayAddr = "127.0.0.1:8080"

	singleFileCid = "QmVtZPoeiqpREqkpTTNMzXkUt74SgQA4JYMG8zPjMVULby"
	emptyDirCid   = "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"
	)
	var file File
	var dir Dir
	// Creat a store
	store, err := NewStore(apiAddr, gatewayAddr)
	if err != nil {
		panic(err)
	}
	// Create a node, notice that it won't be init until being used
	node := store.Get(singleFileCid)
	if file, err = node.ToFile(); err == nil {
		data, _ := ioutil.ReadAll(file)
		fmt.Println(string(data))
	}
	node = store.Get(emptyDirCid)
	if dir, err = node.ToDir(); err == nil {
		fmt.Println("Dir:")
		// Add a file to dir
		dir, err = dir.AddFile(file)
		if err != nil {
			panic(err)
		}
		for _, node := range dir.Nodes() {
			fmt.Println(node.Name(), node.Type())
		}
	}
```
更多的例子可以查看测试代码
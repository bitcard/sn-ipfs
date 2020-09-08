package ipfs

import (
	"fmt"
	client "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	"testing"
)

func TestExample(t *testing.T) {
	const (
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
}

// 测试了文件组合的功能
func Test_store_Combine(t *testing.T) {
	store := testLocalStore()
	cid1 := "QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn"
	//d,_ := client.NewLocalShell().BlockGet(cid1)
	//pn,_ := merkledag.DecodeProtobuf(d)
	//fmt.Println(len(pn.RawData()))
	node := store.Get(cid1)
	file, err := node.ToFile()
	if err != nil {
		panic(err)
	}
	blks := file.Blocks()
	type fields struct {
		api *client.Shell
	}
	type args struct {
		blocks []Block
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "t1",
			fields:  fields{api: client.NewLocalShell()},
			args:    args{blocks: blks},
			want:    cid1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store
			got, err := s.Combine(tt.args.blocks)
			if (err != nil) != tt.wantErr {
				t.Errorf("Combine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Cid() != tt.want {
				t.Errorf("Combine() got = %v, want %v", got.Cid(), tt.want)
			}
		})
	}
}

// 能够测试addbytes和reader
func Test_store_AddFromBytes(t *testing.T) {
	store := testLocalStore()
	data := []byte("hello world")
	file, err := store.AddFromBytes(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(file.Cid())
}

// 通过
func TestStore_Pin(t *testing.T) {
	store := testLocalStore()
	cid1 := "QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn"
	node := store.Get(cid1)
	file, err := node.ToFile()
	if err != nil {
		panic(err)
	}
	err = store.PinMany(file.Blocks())
	if err != nil {
		panic(err)
	}
}

func TestStore_Get(t *testing.T) {
	store := testLocalStore()
	cid1 := "QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn"
	node := store.Get(cid1)
	t.Log(node.Size())
}

// 通过
func TestStore_Unpin(t *testing.T) {
	store := testLocalStore()
	cid1 := "QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn"
	node := store.Get(cid1)
	file, err := node.ToFile()
	if err != nil {
		panic(err)
	}
	err = store.UnpinMany(file.Blocks())
	if err != nil {
		panic(err)
	}
}

func TestDir_Nodes(t *testing.T) {
	store := testLocalStore()
	var (
		err error
		dir Dir
	)
	dirCid := "Qmci4Sm9Cvm4a8fSwzbn6aKYmvmsmkGXpb5f93gBGWgyr9"
	wantNodecid := []map[string]string{
		{"name": "source.list", "cid": "QmRRog1H5wCfzcbCu1BNPas8mW3JuBvRqVBi8xVPFc1KhD"},
	}
	node := store.Get(dirCid)
	// 测试类型转化
	_, err = node.ToFile()
	if err == nil {
		t.Fatal("should be dir but recognized as file")
	}
	dir, err = node.ToDir()
	if err != nil {
		t.Fatal("should be dir but can't be recognized")
	}
	nodes := dir.Nodes()
	// 测试名称和cid是否符合要求
	for i, node := range nodes {
		if node.Cid() != wantNodecid[i]["cid"] {
			t.Fatalf("cid : want %v but got %v", wantNodecid[i]["cid"], node.Cid())
		}
		if node.Name() != wantNodecid[i]["name"] {
			t.Fatalf("name : want %v but got %v", wantNodecid[i]["name"], node.Cid())
		}
	}
}

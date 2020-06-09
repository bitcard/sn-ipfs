package ipfs_filestore

import (
	"fmt"
	client "github.com/ipfs/go-ipfs-api"
	"testing"
)

// 测试了文件组合的功能
func Test_store_Combine(t *testing.T) {
	cid1 := "QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn"
	//d,_ := client.NewLocalShell().BlockGet(cid1)
	//pn,_ := merkledag.DecodeProtobuf(d)
	//fmt.Println(len(pn.RawData()))
	node := Gstore.Get(newLink(cid1))
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
			s := &store{
				api: tt.fields.api,
			}
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
	data := []byte("hello world")
	file, err := Gstore.AddFromBytes(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(file.Cid())
}

// 通过
func TestStore_Pin(t *testing.T) {
	cid1 := "QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn"
	node := Gstore.Get(newLink(cid1))
	file, err := node.ToFile()
	if err != nil {
		panic(err)
	}
	err = Gstore.PinMany(file.Blocks())
	if err != nil {
		panic(err)
	}
}

func TestStore_Get(t *testing.T) {
	cid1 := "QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn"
	node := Gstore.Get(newLink(cid1))
	t.Log(node.Size())
}

// 通过
func TestStore_Unpin(t *testing.T) {
	cid1 := "QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn"
	node := Gstore.Get(newLink(cid1))
	file, err := node.ToFile()
	if err != nil {
		panic(err)
	}
	err = Gstore.UnpinMany(file.Blocks())
	if err != nil {
		panic(err)
	}
}

func TestDir_Nodes(t *testing.T) {
	var (
		err error
		dir Dir
	)
	dirCid := "Qmci4Sm9Cvm4a8fSwzbn6aKYmvmsmkGXpb5f93gBGWgyr9"
	wantNodecid := []map[string]string{
		{"name": "source.list", "cid": "QmRRog1H5wCfzcbCu1BNPas8mW3JuBvRqVBi8xVPFc1KhD"},
	}
	node := Gstore.Get(newLink(dirCid))
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

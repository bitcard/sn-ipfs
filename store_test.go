package ipfs_filestore

import (
	"fmt"
	client "github.com/ipfs/go-ipfs-api"
	"testing"
)

// 测试了文件组合的功能
func Test_store_Combine(t *testing.T) {
	fmt.Println(Gstore)
	cid1 := "QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn"
	//d,_ := client.NewLocalShell().BlockGet(cid1)
	//pn,_ := merkledag.DecodeProtobuf(d)
	//fmt.Println(len(pn.RawData()))
	node := Gstore.Get(newLink(cid1))
	blks := NewFile(node).Blocks()
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

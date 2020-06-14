package ipfs_filestore

import (
	"io"
	"os"
	"testing"
)

// 测试了seeker功能和rader功能
func Test_file_Reader(t *testing.T) {
	var offset int64 = 15552
	cid := "QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn"
	node := Gstore.Get(newLink(cid))
	file, err := node.ToFile()
	if err != nil {
		panic(err)
	}
	f, err := os.Create("qma")
	defer os.Remove("qma")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	file.Seek(offset, 0)
	n, err := io.Copy(f, file)
	wantRead := file.Size() - uint64(offset)
	if uint64(n) != wantRead {
		t.Fatalf("should read %v but read %v", wantRead, n)
	}
}

package ipfs_filestore

import (
	"io"
	"os"
	"testing"
)

func testLocalStore() Store {
	s, err := NewStore("http://127.0.0.1:5001", "http://127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	return s
}

// 测试了seeker功能和rader功能
// 尝试通过一个id打开文件，并进行读取操作
func Test_file_Reader(t *testing.T) {
	store := testLocalStore()
	var offset int64 = 15552
	cid := "QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn"
	node := store.Get(newLink(cid))
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

package ipfs_filestore

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func Test_file_Reader(t *testing.T) {
	offset := 700000
	cid := "QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn"
	node := Gstore.Get(newLink(cid))
	file, err := node.ToFile()
	if err != nil {
		panic(err)
	}
	reader := file.Data()
	f, err := os.Create("qma")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	n, err := io.Copy(f, reader)
	wantRead := file.Size() - uint64(offset)
	if uint64(n) != wantRead {
		t.Fatalf("should read %v but read %v", wantRead, n)
	}
}

func TestHttpClient_Read(t *testing.T) {
	reader := newHttpClient("http://127.0.0.1:8080/ipfs/QmaArqeu69Ss8dhiE9hfZDAMYG8tdoKhgEJREUjQyZLhVn", 7509524)
	reader.open(0)
	f, _ := os.Create("qma")
	defer f.Close()
	var (
		d1 = make([]byte, 1024)
		d2 = make([]byte, 1024)
	)
	n1, err1 := reader.conn.Body.Read(d1)
	n2, err2 := reader.conn.Body.Read(d2)
	io.Copy(f, reader.conn.Body)
	if err1 != nil || err2 != nil {
		t.Fatalf("read error")
	}
	if reflect.DeepEqual(n1, n2) {
		t.Fatalf("dump read")
	}
}

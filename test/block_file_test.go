package test

import (
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"strings"
	"testing"
)

func createShell() *shell.Shell {
	return shell.NewShell("localhost:5001")
}

func Test_File(t *testing.T) {
	sh := createShell()
	cid, err := sh.Add(strings.NewReader("hello world!"))
	// 存在错误，有待研究
	//bt := []byte("hello world")
	//cid,err = sh.DagPutWithOpts(bt)
	if err != nil {
		panic(err)
	}
	var data interface{}
	sh.DagGet(cid, &data)
	fmt.Println(cid, data)
}

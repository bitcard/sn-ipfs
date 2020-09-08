package ipfs

import (
	"fmt"
	"testing"
)

func TestDir_AddFile(t *testing.T) {
	s := testLocalStore()
	node := s.Get("QmeYG2g2LuTnEuekqBYEhWFwUju62D5AinjtRg6kFSv3bz")
	dir, err := node.ToDir()
	if err != nil {
		panic(err)
	}
	for _, node := range dir.Nodes() {
		fmt.Println(node.Name(), node.Type())
	}
}

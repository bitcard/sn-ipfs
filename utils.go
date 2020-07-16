package ipfs

import (
	c "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
)

func newLink(cid string) *ipld.Link {
	_cid, _ := c.Decode(cid)
	return &ipld.Link{Cid: _cid}
}

func newLinkWithName(cid, name string) *ipld.Link {
	_cid, err := c.Decode(cid)
	if err != nil {
		return nil
	}
	return &ipld.Link{Cid: _cid, Name: name}
}

func newFullLink(cid, name string, size uint64) *ipld.Link {
	_cid, err := c.Decode(cid)
	if err != nil {
		return nil
	}
	return &ipld.Link{Cid: _cid, Name: name, Size: size}
}

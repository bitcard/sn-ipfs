package ipfs_filestore

import (
	"errors"
	"fmt"
	ipld "github.com/ipfs/go-ipld-format"
	"io"
	"net/http"
)

type File interface {
	Data() io.ReadCloser
	Blocks() []Block
	Node
}

type file struct {
	Node
}

func (f file) Data() io.ReadCloser {
	url := Gstore.(*store).getGateway() + "/ipfs/" + f.Cid()
	cl := newHttpClient(url, f.Size())
	cl.open(0)
	return cl.raw()
}

func (f file) Blocks() []Block {
	links := f.Links()
	if len(links) == 0 {
		links = []*ipld.Link{newLink(f.Cid())}
	}
	var blocks = make([]Block, 0, len(links))
	nodes := Gstore.GetMany(links)
	for _, node := range nodes {
		blocks = append(blocks, block{node})
	}
	return blocks
}

func newHttpClient(url string, end uint64) *httpClient {
	return &httpClient{url: url, conn: nil, cl: &http.Client{}, end: end}
}

type httpClient struct {
	url   string
	end   uint64
	index uint64
	conn  *http.Response
	cl    *http.Client
}

func (h *httpClient) open(index uint64) error {
	if h.conn != nil {
		h.Close()
		h.conn = nil
	}
	if index > h.end {
		return errors.New("index out of range")
	}
	reqest, err := http.NewRequest("GET", h.url, nil)
	if err != nil {
		return err
	}
	reqest.Header.Add("Range", fmt.Sprintf("bytes=%v-%v", index, h.end))
	h.conn, err = h.cl.Do(reqest)
	if err != nil {
		h.conn = nil
		return err
	}
	return nil
}

func (h *httpClient) raw() io.ReadCloser {
	return h.conn.Body
}

func (h httpClient) Read(p []byte) (n int, err error) {
	if h.conn == nil {
		err = h.open(h.index)
		if err != nil {
			return 0, err
		}
	}
	n, err = h.conn.Body.Read(p)
	h.index += uint64(n)
	if h.index > h.end {
		return n, io.EOF
	}
	return n, err
}

// 暂时废弃，不会使用
func (h httpClient) Seek(offset int64, whence int) (int64, error) {
	var err error
	index := offset + int64(h.index)
	if offset <= 0 {
		return 0, errors.New("read out of range")
	}
	switch whence {
	case io.SeekStart:
		if offset <= 0 {
			return 0, errors.New("read out of range")
		}
		err = h.open(uint64(index))
		return offset, err
	case io.SeekCurrent:
		err = h.open(uint64(index))
		if err != nil {
			return 0, err
		}
		return index, err
	case io.SeekEnd:
		err = h.open(uint64(index))
		if err != nil {
			return 0, err
		}
		return index, err
	default:
		return 0, errors.New("unkonw whence")
	}
}

func (h httpClient) Close() error {
	if h.conn.Close {
		return nil
	}
	return h.conn.Body.Close()
}

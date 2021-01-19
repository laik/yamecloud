package uri

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"reflect"
	"strings"
	"unsafe"
)

const (
	separateKeyword   = "/"
	namespacesKeyword = "namespaces"
)

type uriStruct struct {
	Service      string `json:"service"`
	API          string `json:"api"`
	Group        string `json:"group"`
	Version      string `json:"version"`
	Namespace    string `json:"namespace"`
	Resource     string `json:"resource"`
	ResourceName string `json:"resource_name"`
	Op           string `json:"op"`
}

func (u *uriStruct) String() string {
	b, _ := json.Marshal(u)
	return hackString(b)
}

// URI is only via this project url spec.
type URI struct {
	_index  uint
	_count  int
	_offset int64
	_url    string
	uriStruct
}

// Parser yamecloud URI general interface analysis
type Parser interface {
	ParseOp(uri string) (*URI, error)
}

// NewURIParser return default parser
func NewURIParser() Parser {
	return &parseImplement{}
}

var _ Parser = (*parseImplement)(nil)

type parseImplement struct{}

func (p *parseImplement) ParseOp(url string) (*URI, error) {
	return parse(url)
}

func parse(_url string) (*URI, error) {
	_URL, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	uri := &URI{
		_url:   _URL.Path,
		_count: strings.Count(_URL.Path, separateKeyword),
	}
	if err := uri.parse(); err != nil {
		return nil, err
	}
	return uri, nil
}

func (u *URI) parse() error {
	lastOp := false
	for index := 1; index <= u._count; index++ {
		item, err := u.shift()
		if err != nil {
			return err
		}

		switch index {
		case 1:
			u.Service = item
			continue
		case 2:
			u.API = item
			continue
		}

		if item == "op" {
			lastOp = true
			continue
		}

		if lastOp {
			u.Op = item
			continue
		}

		///workload/api/v1/namespaces/im/pods

		switch u.API {
		case "api":
			switch index {
			case 3:
				u.Version = item
				continue
			case 4:
				u.Resource = item
			case 5:
				if u.Resource == "namespaces" {
					u.ResourceName = item
				}
				u.Namespace = item
				continue
			case 6:
				u.Resource = item
				if u.Resource != "namespaces" {
					u.ResourceName = ""
				}
				continue
			case 7:
				u.ResourceName = item
				continue
			case 8:
				u.Op = item
				continue
			}

		case "apis":
			switch index {
			case 3:
				u.Group = item
				continue
			case 4:
				u.Version = item
				continue
			case 5:
				if item == "namespaces" {
					continue
				}
				u.Resource = item
				continue
			case 6:
				u.Namespace = item
				continue
			case 7:
				u.Resource = item
				continue
			case 8:
				u.ResourceName = item
				continue
			case 9:
				u.Op = item
				continue
			}
		}
	}

	return nil
}

func (u *URI) shift() (item string, err error) {
	itemBytes, err := u.shiftItem()
	if err != nil {
		return "", err
	}
	item = hackString(itemBytes)
	return
}

func (u *URI) shiftItem() (item []byte, err error) {
	item, u._offset, err = readItem(u._url, int64(u._offset))
	if err != nil {
		return nil, err
	}
	return
}

func readItem(uri string, offset int64) (item []byte, nextOffset int64, err error) {
	bytesReader := bytes.NewReader(hackSlice(uri))
	if nextOffset, err = bytesReader.Seek(offset, io.SeekCurrent); err != nil {
		return
	}
	prefix := false
	for {
		b, err := bytesReader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return item, nextOffset, nil
			}
			return nil, nextOffset, err
		}
		if b == '/' && !prefix {
			prefix = true
			nextOffset++
			continue
		}

		if b == '/' && prefix {
			break
		}
		item = append(item, b)
		nextOffset++
	}

	return
}

func hackString(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}

func hackSlice(s string) (b []byte) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pbytes.Data = pstring.Data
	pbytes.Len = pstring.Len
	pbytes.Cap = pstring.Len
	return
}

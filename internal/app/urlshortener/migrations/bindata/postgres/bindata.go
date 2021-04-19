package postgres

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var __1_create_urls_table_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4a\x29\xca\x2f\x50\x28\x49\x4c\xca\x49\x55\xc8\x4c\x53\x48\xad\xc8\x2c\x2e\x29\x56\x28\x2d\xca\x29\xb6\x06\x04\x00\x00\xff\xff\xf1\x7e\xf6\xa9\x1a\x00\x00\x00")

func _1_create_urls_table_down_sql() ([]byte, error) {
	return bindata_read(
		__1_create_urls_table_down_sql,
		"1_create_urls_table.down.sql",
	)
}

var __1_create_urls_table_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\xcd\x41\xaa\xc3\x30\x0c\x04\xd0\xbd\x4f\x31\xcb\x04\xfe\x0d\xfe\x61\x8a\x52\x29\x45\x54\x49\x8c\x2c\xd3\xe4\xf6\xa5\xb5\x29\x04\x4a\xb5\x12\xbc\xd1\xe8\xea\x42\x21\x08\x9a\x4c\xa0\x33\xd6\x2d\x20\xbb\x96\x28\xa8\x6e\x25\x0d\x09\x00\x94\x71\x9a\x49\x6f\x45\x5c\xc9\xfe\xde\x5c\xdd\x4e\x1c\xb2\x47\xdb\x5e\x75\x6b\xb5\x9e\x6b\xcf\xf8\x42\x5d\x43\x17\x29\x41\x4b\xfe\xe4\xc0\x32\x53\xb5\xc0\xba\x3d\x86\xb1\xb7\x67\xfe\x7e\xd5\x98\xc5\xe4\x07\x67\xd7\x85\xfc\xc0\x5d\x0e\x0c\xca\x63\x1a\xff\x53\x7a\x06\x00\x00\xff\xff\x85\xad\x0d\x5b\xf6\x00\x00\x00")

func _1_create_urls_table_up_sql() ([]byte, error) {
	return bindata_read(
		__1_create_urls_table_up_sql,
		"1_create_urls_table.up.sql",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"1_create_urls_table.down.sql": _1_create_urls_table_down_sql,
	"1_create_urls_table.up.sql": _1_create_urls_table_up_sql,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"1_create_urls_table.down.sql": &_bintree_t{_1_create_urls_table_down_sql, map[string]*_bintree_t{
	}},
	"1_create_urls_table.up.sql": &_bintree_t{_1_create_urls_table_up_sql, map[string]*_bintree_t{
	}},
}}

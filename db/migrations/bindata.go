// Code generated by go-bindata. DO NOT EDIT.
// sources:
// 000001_init.down.sql
// 000001_init.up.sql
// 000002_log.down.sql
// 000002_log.up.sql
// 000003_remove_frequency.down.sql
// 000003_remove_frequency.up.sql
// 000004_comment.down.sql
// 000004_comment.up.sql
// migrations.go
// migrations_test.go
package migrations

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __000001_initDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\xc8\x4d\x2d\x29\xca\x4c\xb6\xe6\xc2\x2a\x99\x5c\x82\x4b\x22\x3f\x2f\x2d\x33\xdd\x9a\x0b\x10\x00\x00\xff\xff\xcc\xc4\x2f\x49\x53\x00\x00\x00")

func _000001_initDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_initDownSql,
		"000001_init.down.sql",
	)
}

func _000001_initDownSql() (*asset, error) {
	bytes, err := _000001_initDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_init.down.sql", size: 83, mode: os.FileMode(436), modTime: time.Unix(1599459506, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000001_initUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x90\xcf\x4e\xf3\x30\x10\xc4\xef\x7e\x8a\x39\x26\x52\xdf\xe0\x3b\xf9\x73\x37\x95\x85\xe3\x84\xcd\x5a\xa2\xa7\x2a\x6a\x0d\xb2\x44\xd2\x0a\x0c\xe2\xf1\x51\x12\xfe\x48\x50\x51\x81\x8f\xb3\xbb\xe3\xdf\x8c\x61\xd2\x42\x10\xfd\xdf\x91\xb2\x15\x7c\x23\xa0\x1b\xdb\x49\x87\x21\xe6\x87\xb4\x57\x85\x02\x80\x74\x80\xf5\x42\x1b\xe2\x79\xc5\x07\xe7\xd0\xb2\xad\x35\x6f\x71\x45\xdb\xd5\xbc\x34\xf6\x43\x44\x8e\x2f\x59\x95\xff\x94\xfa\xc1\x7a\x9f\x7f\x63\x9b\xd3\x10\x1f\x73\x3f\x9c\xb0\xd6\x42\x62\x6b\xc2\x9a\x2a\x1d\x9c\xc0\x04\x66\xf2\xb2\x9b\xc4\x4e\x74\xdd\x2e\x17\x0b\xfa\xee\x8c\xfb\x32\x7f\xee\xef\x9f\x22\xd2\x98\xbf\xe8\x55\xc3\x64\x37\x7e\xfa\xba\xf8\xf0\x28\xc1\x54\x11\x93\x37\xf4\x5e\x4a\x91\x0e\xe5\xc5\x8c\xc7\xf1\x36\xdd\xe1\x2d\xe8\x25\xa2\xe3\x29\xcf\xcd\x7d\x07\x3d\x27\x07\x6f\xaf\x03\x7d\x22\xae\xa6\xfb\x72\x1e\x4d\xaf\xf1\x30\x8d\xaf\x9c\x35\x02\xa6\xd6\x69\x43\x7f\xcb\xf7\x1a\x00\x00\xff\xff\x10\x01\xe1\xc0\x1e\x02\x00\x00")

func _000001_initUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_initUpSql,
		"000001_init.up.sql",
	)
}

func _000001_initUpSql() (*asset, error) {
	bytes, err := _000001_initUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_init.up.sql", size: 542, mode: os.FileMode(436), modTime: time.Unix(1599459506, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000002_logDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x90\xc1\x6a\xc3\x30\x10\x44\xef\xfa\x8a\x39\xda\xe0\x3f\xc8\x49\xb5\x57\x41\x54\x96\xcc\x6a\x0d\xcd\x29\x04\xdb\x14\x41\xdc\x96\x56\xed\xf7\x17\xdb\x69\x0d\xa5\xe4\x3a\x33\x3b\x6f\x99\x9a\x49\x0b\x41\xf4\x83\x23\x65\x0d\x7c\x10\xd0\x93\x8d\x12\x31\x64\x55\x28\x00\x48\x23\xac\x17\x3a\x12\xaf\xb6\xef\x9d\x43\xc7\xb6\xd5\x7c\xc2\x23\x9d\xaa\x35\x94\xd3\x3c\x7d\xe4\xcb\xfc\x86\x46\x0b\x89\x6d\x09\x0d\x19\xdd\x3b\x41\xdd\x33\x93\x97\xf3\x22\x46\xd1\x6d\xb7\x5d\xcc\x53\x7e\x4f\xc3\xf9\x9f\xf6\xcd\xff\xba\x5c\x3f\x27\xa4\x97\xfc\x47\x37\x81\xc9\x1e\xfd\x82\x2e\x7e\x3b\x4a\x30\x19\x62\xf2\x35\xc5\x5b\x73\x91\xc6\x52\x95\x07\xa5\xac\x8f\xc4\xb2\x50\x02\x86\x5c\xa4\xb1\xda\xd9\xd5\x86\xa9\xf6\xff\x4b\x44\x72\x54\xdf\x88\xf7\x92\x30\x1c\x5a\x5c\x5f\x9f\x0f\x4a\x35\x1c\xba\x6d\x45\x58\xf3\xb3\xe0\x6a\x7d\x07\x00\x00\xff\xff\xc8\x17\x88\xf2\x62\x01\x00\x00")

func _000002_logDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000002_logDownSql,
		"000002_log.down.sql",
	)
}

func _000002_logDownSql() (*asset, error) {
	bytes, err := _000002_logDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000002_log.down.sql", size: 354, mode: os.FileMode(436), modTime: time.Unix(1599459506, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000002_logUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x90\x41\x6a\xc3\x30\x10\x45\xf7\x3a\xc5\x2c\x6d\xf0\x0d\xb2\x52\x9d\xef\x20\x2a\x4b\xee\x68\x04\xcd\x2a\x84\xd8\x14\x41\xdc\x96\x56\xed\xf9\x8b\xad\x10\x43\x29\x5d\x6a\xe6\x0d\xff\x7d\xb5\x0c\x2d\x20\xd1\x0f\x16\xca\x74\xe4\xbc\x10\x9e\x4d\x90\x40\xd7\xb7\x17\x55\x29\x22\xa2\x34\x92\x71\x82\x03\x78\xdd\xbb\x68\x2d\x0d\x6c\x7a\xcd\x47\x7a\xc4\xb1\x59\xa1\x79\xca\x1f\xe9\x72\xfa\x83\x2d\xfb\xef\xf3\xf5\x6b\xa2\xf4\x9a\x7f\xcd\x73\x9a\xa7\xcf\x7c\x9e\xdf\x69\xaf\x05\x62\x7a\xd0\x1e\x9d\x8e\x56\xa8\x8d\xcc\x70\x72\x5a\x86\x41\x74\x3f\x94\x8b\xce\x33\xcc\xc1\x2d\xd1\xd5\x3d\xb5\x26\x46\x07\x86\x6b\x11\x6e\x2e\x55\x1a\xeb\x72\x11\x9d\x79\x8a\xd8\xe0\xe6\x1e\x5a\xab\x7a\xa7\x94\x71\x01\x2c\x8b\xb8\x5f\x6a\x57\x69\x6c\xb6\x3e\x4d\x51\x6f\x36\xd3\x9a\x02\x2c\xda\x5b\x8b\xff\x48\xea\xd8\xf7\x74\xc9\x3b\xa5\xf6\xec\x87\xf2\xcf\xeb\xfb\x27\x00\x00\xff\xff\x47\x6b\xd9\x83\x79\x01\x00\x00")

func _000002_logUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000002_logUpSql,
		"000002_log.up.sql",
	)
}

func _000002_logUpSql() (*asset, error) {
	bytes, err := _000002_logUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000002_log.up.sql", size: 377, mode: os.FileMode(436), modTime: time.Unix(1599459506, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000003_remove_frequencyDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func _000003_remove_frequencyDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000003_remove_frequencyDownSql,
		"000003_remove_frequency.down.sql",
	)
}

func _000003_remove_frequencyDownSql() (*asset, error) {
	bytes, err := _000003_remove_frequencyDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000003_remove_frequency.down.sql", size: 0, mode: os.FileMode(436), modTime: time.Unix(1599459506, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000003_remove_frequencyUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x71\xf5\x71\x0d\x71\x55\x70\x0b\xf2\xf7\x55\x48\xce\xcf\x4b\xcb\x4c\x57\x08\xf7\x70\x0d\x72\x55\xc8\x2f\x28\x51\xb0\x55\x50\x4f\x2b\x4a\x2d\x2c\x4d\xcd\x4b\xae\x54\xe7\x02\x04\x00\x00\xff\xff\x2b\x11\x3b\x02\x2b\x00\x00\x00")

func _000003_remove_frequencyUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000003_remove_frequencyUpSql,
		"000003_remove_frequency.up.sql",
	)
}

func _000003_remove_frequencyUpSql() (*asset, error) {
	bytes, err := _000003_remove_frequencyUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000003_remove_frequency.up.sql", size: 43, mode: os.FileMode(436), modTime: time.Unix(1599459506, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000004_commentDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xc8\xc9\x4f\x8f\x4f\xce\xcf\xcd\x4d\xcd\x2b\xb1\xe6\x02\x04\x00\x00\xff\xff\x72\xe9\x3e\xa8\x18\x00\x00\x00")

func _000004_commentDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000004_commentDownSql,
		"000004_comment.down.sql",
	)
}

func _000004_commentDownSql() (*asset, error) {
	bytes, err := _000004_commentDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000004_comment.down.sql", size: 24, mode: os.FileMode(420), modTime: time.Unix(1606466873, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __000004_commentUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x8d\xb1\x0a\xc2\x30\x14\x45\xf7\x7c\xc5\x1d\x13\xf0\x0f\x9c\x6a\xb9\x29\xc1\xf0\x8a\xe9\x0b\xe8\xe4\x60\x8b\x14\xac\x5d\x32\xf8\xf9\x42\xaa\x4b\xd7\x73\xce\xe5\xb6\x89\x8d\x12\xda\x9c\x22\x4d\xf0\x90\x5e\xc1\x6b\x18\x74\xc0\x6b\x7d\xde\x1f\xeb\xb2\x4c\xef\x62\xac\x01\x50\xc9\x3c\x22\x88\xb2\x63\xaa\xad\xe4\x18\x0f\x55\xfe\x52\x94\xe9\x53\x76\x2a\x4b\xb8\x64\xda\x6d\xee\x36\xe6\xfb\xc4\xd0\x09\xce\xbc\xfd\x05\x12\x3d\x13\xa5\x65\x3d\xb7\xf3\xe8\x8c\x3b\x9a\x6f\x00\x00\x00\xff\xff\x73\x51\x62\x38\xa3\x00\x00\x00")

func _000004_commentUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000004_commentUpSql,
		"000004_comment.up.sql",
	)
}

func _000004_commentUpSql() (*asset, error) {
	bytes, err := _000004_commentUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000004_comment.up.sql", size: 163, mode: os.FileMode(420), modTime: time.Unix(1606466873, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x90\x3d\x8f\xd4\x30\x10\x86\x6b\xcf\xaf\x18\x52\x9c\x6c\x69\x2f\x2e\xa0\x02\x5d\x01\x84\x02\x09\xb6\x38\x74\xa2\x40\xe8\xe4\x64\x27\x5e\x8b\xc4\x0e\x63\x07\x84\xd0\xfe\x77\x14\xe7\x83\x08\x51\xec\xa5\x89\xc6\x7a\xe7\x79\x46\xef\x60\x9a\x6f\xc6\x12\xf6\xce\xb2\x49\x2e\xf8\x08\xe0\xfa\x21\x70\x42\x09\xa2\xb0\x2e\x9d\xc7\xba\x6c\x42\xaf\x6d\xe8\x8c\xb7\xb7\x73\x90\xf4\xfa\xff\xf1\xa2\x00\xf1\x88\x57\x25\xf5\xc9\x24\x53\x9b\x48\x3a\x7e\xef\x5c\xa2\xe7\x05\x6a\xed\x43\xe7\x7c\xba\x9e\x11\xc3\xc8\x0d\xe9\xd6\x75\x54\x60\xfe\xfe\x32\x6a\xe7\x27\xc5\xd3\x48\x36\x3c\x2e\x7b\x05\x28\x00\xad\xb1\x0a\x1f\xe7\x54\x55\xe3\x40\xdc\x06\xee\x23\x56\x6f\x76\x25\x95\xd0\x8e\xbe\xd9\x07\xe5\xa9\x7e\xb8\xff\x80\x31\xb1\xf3\x56\x21\x31\x07\xc6\xdf\x20\x98\x66\x4d\xc4\x97\x77\xb8\x78\xca\xfb\xe5\x51\xbe\x8e\x91\xd2\xd1\xf4\x14\xa5\x3a\x80\x10\x13\x55\x7a\xd3\xd3\x06\x92\x5f\xbe\xd6\xbf\x12\x1d\x66\xa2\x9a\x90\x42\x30\xa5\x91\x3d\xe6\xed\x1c\x57\x20\xc4\x45\x01\x88\xed\xc4\xca\x24\x93\x97\xf6\xde\xcf\x2e\x9d\xdf\xfb\x98\x8c\x6f\x48\x6e\x97\x29\x10\xae\xcd\xd1\x67\x77\xe8\x5d\x97\x1d\x8b\x82\x98\x41\x5c\x26\xf0\x06\x5b\x1a\x2c\x8f\xf4\x73\xe2\x7d\xca\x90\x8d\x5a\xd8\x70\xbb\xd6\x79\xc0\x7f\xce\xc9\x1d\x5d\xa5\x5b\x12\x93\xae\x7c\x18\xa4\x7a\xb5\x5f\xb8\xb9\x59\xa7\xf5\x96\x77\xcc\xc7\xf0\xf6\x6c\xbc\xa5\xff\xe0\xd6\xd1\xbb\x0e\x2e\xf0\x27\x00\x00\xff\xff\x21\x84\xf6\xe9\xf3\x02\x00\x00")

func migrationsGoBytes() ([]byte, error) {
	return bindataRead(
		_migrationsGo,
		"migrations.go",
	)
}

func migrationsGo() (*asset, error) {
	bytes, err := migrationsGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations.go", size: 755, mode: os.FileMode(436), modTime: time.Unix(1599459506, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations_testGo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\xcb\xb1\x0e\x82\x30\x10\x00\xd0\xb9\xf7\x15\x67\x27\xea\x00\x83\x5b\x0d\x1b\x8e\x6e\xfc\x40\x25\x47\x73\x91\xb6\x78\x3d\x07\x63\xf8\x77\x23\x71\x70\x7f\x6f\x0d\xd3\x3d\x44\xc2\xc4\x51\x82\x72\xc9\x15\x80\xd3\x5a\x44\xb1\x01\x63\x95\xaa\x72\x8e\x16\x1c\xc0\xfc\xcc\x13\x8e\x54\x75\x28\xd7\x5d\xd3\x70\x6b\x14\x8f\x3f\xd3\x8e\x0e\xdf\x60\x78\x46\x12\x41\xdf\xe3\x3f\xb3\xf5\xb1\xb0\xd2\xc9\x77\x9d\x4f\x94\x8a\xbc\xbc\x75\xe7\x5d\x1e\x7a\xcc\xbc\x7c\xab\xd1\xf6\x22\x52\xa4\x21\x11\x07\x66\x83\x0d\x3e\x01\x00\x00\xff\xff\x61\xad\xbb\x6d\x9f\x00\x00\x00")

func migrations_testGoBytes() ([]byte, error) {
	return bindataRead(
		_migrations_testGo,
		"migrations_test.go",
	)
}

func migrations_testGo() (*asset, error) {
	bytes, err := migrations_testGoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations_test.go", size: 159, mode: os.FileMode(436), modTime: time.Unix(1599459506, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
var _bindata = map[string]func() (*asset, error){
	"000001_init.down.sql": _000001_initDownSql,
	"000001_init.up.sql": _000001_initUpSql,
	"000002_log.down.sql": _000002_logDownSql,
	"000002_log.up.sql": _000002_logUpSql,
	"000003_remove_frequency.down.sql": _000003_remove_frequencyDownSql,
	"000003_remove_frequency.up.sql": _000003_remove_frequencyUpSql,
	"000004_comment.down.sql": _000004_commentDownSql,
	"000004_comment.up.sql": _000004_commentUpSql,
	"migrations.go": migrationsGo,
	"migrations_test.go": migrations_testGo,
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
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"000001_init.down.sql": &bintree{_000001_initDownSql, map[string]*bintree{}},
	"000001_init.up.sql": &bintree{_000001_initUpSql, map[string]*bintree{}},
	"000002_log.down.sql": &bintree{_000002_logDownSql, map[string]*bintree{}},
	"000002_log.up.sql": &bintree{_000002_logUpSql, map[string]*bintree{}},
	"000003_remove_frequency.down.sql": &bintree{_000003_remove_frequencyDownSql, map[string]*bintree{}},
	"000003_remove_frequency.up.sql": &bintree{_000003_remove_frequencyUpSql, map[string]*bintree{}},
	"000004_comment.down.sql": &bintree{_000004_commentDownSql, map[string]*bintree{}},
	"000004_comment.up.sql": &bintree{_000004_commentUpSql, map[string]*bintree{}},
	"migrations.go": &bintree{migrationsGo, map[string]*bintree{}},
	"migrations_test.go": &bintree{migrations_testGo, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}


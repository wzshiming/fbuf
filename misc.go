package fbuf

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const perm os.FileMode = 0700

func selfDir() string {
	dir := os.Args[0]
	dir, _ = filepath.Abs(dir)
	dir = filepath.Dir(dir)
	return dir
}

func mkDir(dir string) error {
	return os.MkdirAll(dir, perm)
}

func readFile(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}

func writeFile(name string, data []byte) error {
	return ioutil.WriteFile(name, data, perm)
}

func readGetHttp(name string) ([]byte, error) {
	res, err := http.Get(name)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

func joinPath(elem ...string) string {
	return filepath.Join(elem...)
}

func encodeFilename(name string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(name))
}

func decodeFilename(name string) string {
	dir, _ := base64.RawStdEncoding.DecodeString(name)
	return string(dir)
}

package fbuf

import (
	"os"
	"path/filepath"
)

const perm os.FileMode = 0700

const (
	temp = "/go_build_fbuf/"
)

func tempDir() string {
	return os.TempDir() + temp
}

func selfDir() string {
	dir := os.Args[0]
	dir, _ = filepath.Abs(dir)
	dir = filepath.Dir(dir)
	return dir
}

func joinPath(elem ...string) string {
	return filepath.Join(elem...)
}

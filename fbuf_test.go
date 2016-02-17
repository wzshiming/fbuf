package fbuf

import (
	"fmt"
	"testing"
)

func TestWFile(t *testing.T) {
	err := Defaul.Save("11.TXT", []byte("1111"))
	if err != nil {
		t.Error(err)
	}
}

func TestFile(t *testing.T) {
	b, err := Defaul.Open("c://go/robots.txt")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
}

func TestSave(t *testing.T) {
	dir, err := Defaul.TempDir("http", "http://www.baidu.com")
	fmt.Println(dir, err)
}

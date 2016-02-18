package fbuf

import "testing"

func TestWFile(t *testing.T) {
	err := NewFbuf().Write("temp 11.txt", []byte("1111"))
	if err != nil {
		t.Error(err)
	}
	d, err := NewFbuf().Read("temp 11.txt")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(d))
}

func TestFile(t *testing.T) {
	b, err := NewFbuf().Read("c:/go/robots.txt")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
}

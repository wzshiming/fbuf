package fbuf_http

import (
	"testing"

	"github.com/wzshiming/fbuf"
)

func TestHttp(t *testing.T) {
	b, err := fbuf.NewFbuf().Read("buff http://weibo.com/login.php")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
}

func TestHttps(t *testing.T) {
	b, err := fbuf.NewFbuf().Read("buff https://github.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
}

func TestPostHttps(t *testing.T) {
	b, err := fbuf.NewFbuf().Read("rebuff post http://www.w3school.com.cn/tiy/v.asp", map[string]string{
		"code": `11`,
		"bt":   "",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
}

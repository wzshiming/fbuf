package fbuf_http

import (
	"testing"

	"github.com/wzshiming/fbuf"
)

func TestHttp(t *testing.T) {
	b, err := fbuf.Defaul.Open("http://www.baidu.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
}

func TestHttps(t *testing.T) {
	b, err := fbuf.Defaul.Open("https://github.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
}

func TestPostHttps(t *testing.T) {
	b, err := fbuf.Defaul.Open("http://www.w3school.com.cn/tiy/v.asp", "post", map[string]string{
		"code": `11`,
		"bt":   "",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
}

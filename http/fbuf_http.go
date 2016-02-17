package fbuf_http

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/wzshiming/fbuf"
)

const (
	GET  = "GET"
	POST = "POST"
)

func init() {

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			//DisableCompression: true,
		},
	}
	fbuf.Defaul.RegisterReg("http", `^http://`, `^https://`)
	fbuf.Defaul.RegisterRead("http", func(name string, args ...interface{}) ([]byte, error) {
		met := GET
		arg := url.Values{}
		for _, v := range args {
			switch v.(type) {
			case string:
				v0 := v.(string)
				v0 = strings.ToUpper(v0)
				if v0 == POST {
					met = v0
				}
			case url.Values:
				arg = v.(url.Values)
			case map[string]string:
				v0 := v.(map[string]string)
				for k, v := range v0 {
					arg.Add(k, v)
				}
			}
		}
		var res *http.Response
		var err error
		if met == POST {
			res, err = httpClient.PostForm(name, arg)
		} else {
			n, va := urlParse(name, arg)
			name = n + "?" + va.Encode()
			res, err = httpClient.Get(name)
		}
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(res.Body)
	})
}

func urlParse(que string, args url.Values) (string, url.Values) {
	ra, err := url.Parse(que)
	if err != nil {
		return "", nil
	}
	m := ra.Query()
	if args != nil {
		for k, _ := range args {
			m.Add(k, args.Get(k))
		}
	}
	return ra.Scheme + "://" + ra.Host + ra.Path, m
}

package fbuf_http

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/wzshiming/fbuf"
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

	fbuf.Defaul.RegisterRegexp("get", `^get `)
	fbuf.Defaul.RegisterRead("get", func(name string, args ...interface{}) ([]byte, error) {
		if strings.Index(name, "get ") == 0 {
			name = name[4:]
		}
		return fbuf.Defaul.ReadByMethod("http", name, args...)
	})

	fbuf.Defaul.RegisterRegexp("post", `^post `)
	fbuf.Defaul.RegisterRead("post", func(name string, args ...interface{}) ([]byte, error) {
		if strings.Index(name, "post ") == 0 {
			name = name[5:]
		}
		return fbuf.Defaul.ReadByMethod("http", name, append(args, http.MethodPost)...)
	})

	fbuf.Defaul.RegisterRegexp("http", `^http://`, `^https://`)
	fbuf.Defaul.RegisterRead("http", func(name string, args ...interface{}) ([]byte, error) {
		met := http.MethodGet
		arg := url.Values{}
		for _, v := range args {
			switch v.(type) {
			case string:
				v0 := v.(string)
				v0 = strings.ToUpper(v0)
				if v0 == http.MethodPost {
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
		if met == http.MethodPost {
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

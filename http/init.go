package fbuf_http

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/wzshiming/fbuf"
)

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5741
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

var (
	UserAgent = "Mozilla/4.0 (Windows; MSIE 6.0; Windows NT 5.2)"
)

func init() {

	cookieJar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Jar: cookieJar,
	}

	fbuf.Defaul.RegisterRegexp("get", `^get\s+`)
	fbuf.Defaul.RegisterRead("get", func(name string, args ...interface{}) ([]byte, error) {
		name = regexp.MustCompile(`^get\s+`).ReplaceAllString(name, "")
		return fbuf.Defaul.ReadByMethod("http", name, args...)
	})

	fbuf.Defaul.RegisterRegexp("post", `^post\s+`)
	fbuf.Defaul.RegisterRead("post", func(name string, args ...interface{}) ([]byte, error) {
		name = regexp.MustCompile(`^post\s+`).ReplaceAllString(name, "")
		return fbuf.Defaul.ReadByMethod("http", name, append(args, MethodPost)...)
	})

	fbuf.Defaul.RegisterRegexp("http", `^http://`, `^https://`)
	fbuf.Defaul.RegisterRead("http", func(name string, args ...interface{}) ([]byte, error) {
		met := MethodGet
		arg := url.Values{}
		for _, v := range args {
			switch v.(type) {
			case string:
				v0 := v.(string)
				v0 = strings.ToUpper(v0)
				if v0 == MethodPost {
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

		var req *http.Request
		var err error
		switch met {
		case MethodPost:
			req, err = http.NewRequest(MethodPost, name, strings.NewReader(arg.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		case MethodGet:
			req, err = http.NewRequest(MethodGet, name+"?"+arg.Encode(), nil)
		default:
		}
		req.Header.Set("User-Agent", UserAgent)
		if err != nil {
			return nil, err
		}

		res, err := httpClient.Do(req)
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

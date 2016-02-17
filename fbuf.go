package fbuf

import (
	"encoding/base64"
	"errors"
	"regexp"
)

var dirname = "_gofbuf"

type readMethod func(name string, args ...interface{}) ([]byte, error)
type writeMethod func(name string, data []byte, args ...interface{}) error

type method struct {
	regexp      []*regexp.Regexp
	readMethod  readMethod
	writeMethod writeMethod
}

type fbuf struct {
	methods map[string]*method
}

func newFbuf() *fbuf {
	fb := &fbuf{
		methods: map[string]*method{},
	}
	return fb
}

func (fb *fbuf) RegisterReg(name string, regs ...string) {
	if fb.methods[name] == nil {
		fb.methods[name] = &method{}
	}
	mm := fb.methods[name].regexp
	for _, v := range regs {
		mm = append(mm, regexp.MustCompile(v))
	}
	fb.methods[name].regexp = mm
	return
}

func (fb *fbuf) RegisterRead(name string, met readMethod) {
	if fb.methods[name] == nil {
		fb.methods[name] = &method{}
	}
	fb.methods[name].readMethod = met
	return
}

func (fb *fbuf) RegisterWrite(name string, met writeMethod) {
	if fb.methods[name] == nil {
		fb.methods[name] = &method{}
	}
	fb.methods[name].writeMethod = met
	return
}

func (fb *fbuf) matchMethod(name string) map[string]int {
	r := map[string]int{}
	for k, v := range fb.methods {
		for _, v2 := range v.regexp {
			if v2.MatchString(name) {
				r[k] += 1
			}
		}
	}
	return r
}

func (fb *fbuf) MatchMethod(name string) string {
	ss := fb.matchMethod(name)
	sort := []string{}
	si := 0
	for _, v := range ss {
		if si < v {
			si = v
		}
	}
	for k, v := range ss {
		if si == v {
			sort = append(sort, k)
		}
	}

	switch len(sort) {
	case 0:
		return "file"
	case 1:
		return sort[0]
	default:
		return sort[0] // 如果有多个符合条件的 还没写
	}
	return "file"
}

func (fb *fbuf) Save(name string, data []byte, args ...interface{}) error {
	return fb.SaveByMethod(fb.MatchMethod(name), name, data, args...)
}

func (fb *fbuf) SaveByMethod(method, name string, data []byte, args ...interface{}) error {
	met := fb.methods[method]
	if met == nil {
		return errors.New("SaveByMethod: The method does not exist")
	}
	return met.writeMethod(name, data, args...)
}

func (fb *fbuf) Open(name string, args ...interface{}) ([]byte, error) {
	return fb.OpenByMethod(fb.MatchMethod(name), name, args...)
}

func (fb *fbuf) OpenByMethod(method, name string, args ...interface{}) ([]byte, error) {
	met := fb.methods[method]
	if met == nil {
		return nil, errors.New("OpenByMethod: The method does not exist")
	}
	return met.readMethod(name, args...)
}

func (fb *fbuf) getTempRootDir() string {
	return joinPath(selfDir(), dirname)
}

func (fb *fbuf) TempDir(method, name string) (string, error) {
	na := base64.RawStdEncoding.EncodeToString([]byte(name))
	dir := joinPath(fb.getTempRootDir(), method, na)
	return dir, mkDir(dir)
}

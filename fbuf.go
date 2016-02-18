package fbuf

import (
	"fmt"
	"regexp"
)

type readMethod func(name string, args ...interface{}) ([]byte, error)
type writeMethod func(name string, data []byte, args ...interface{}) error

type method struct {
	regexp []*regexp.Regexp
	read   readMethod
	write  writeMethod
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

func NewFbuf() *fbuf {
	fb := newFbuf()
	for k, v := range Defaul.methods {
		fb.methods[k] = v
	}
	return fb
}

func (fb *fbuf) RegisterRegexp(name string, regs ...string) {
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
	fb.methods[name].read = met
	return
}

func (fb *fbuf) RegisterWrite(name string, met writeMethod) {
	if fb.methods[name] == nil {
		fb.methods[name] = &method{}
	}
	fb.methods[name].write = met
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

func (fb *fbuf) MatchMethod(name string) []string {
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

	if len(sort) == 0 {
		return []string{"file"}
	}
	return sort
}

func (fb *fbuf) Write(name string, data []byte, args ...interface{}) error {
	so := fb.MatchMethod(name)
	var err error
	for _, v := range so {
		err = fb.WriteByMethod(v, name, data, args...)
		if err == nil {
			return nil
		}
	}
	return err
}

func (fb *fbuf) WriteByMethod(method, name string, data []byte, args ...interface{}) error {
	met := fb.methods[method]
	if met == nil || met.write == nil {
		return fmt.Errorf("write %s %s: The method does not exist", method, name)
	}
	return met.write(name, data, args...)
}

func (fb *fbuf) Read(name string, args ...interface{}) ([]byte, error) {
	so := fb.MatchMethod(name)
	var err error
	var data []byte
	for _, v := range so {
		data, err = fb.ReadByMethod(v, name, args...)
		if err == nil {
			return data, nil
		}
	}
	return nil, err
}

func (fb *fbuf) ReadByMethod(method, name string, args ...interface{}) ([]byte, error) {
	met := fb.methods[method]
	if met == nil || met.read == nil {
		return nil, fmt.Errorf("read %s %s: The method does not exist", method, name)
	}
	return met.read(name, args...)
}

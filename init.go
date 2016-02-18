package fbuf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var Defaul = newFbuf()

func init() {

	// 本地文件
	Defaul.RegisterRegexp("file", `^\w:/`, `^/`, `^./`, `^\w+/`)
	Defaul.RegisterRead("file", func(name string, args ...interface{}) ([]byte, error) {
		return ioutil.ReadFile(name)
	})
	Defaul.RegisterWrite("file", func(name string, data []byte, args ...interface{}) error {
		pr := perm
		for _, v := range args {
			switch v.(type) {
			case os.FileMode, uint32:
				pr = v.(os.FileMode)
			}
		}
		err := ioutil.WriteFile(name, data, pr)
		if err != nil {
			os.MkdirAll(filepath.Dir(name), pr)
			return ioutil.WriteFile(name, data, pr)
		}
		return err
	})

	// 临时文件
	mc := regexp.MustCompile("[\\/:*?\"<>|]")
	Defaul.RegisterRegexp("temp", `^temp `)
	Defaul.RegisterRead("temp", func(name string, args ...interface{}) ([]byte, error) {
		if strings.Index(name, "temp ") == 0 {
			name = name[5:]
		}
		// 定位到临时文件夹读取
		name = mc.ReplaceAllString(name, "/")
		name = joinPath(tempDir(), name)
		return Defaul.ReadByMethod("file", name, args...)
	})
	Defaul.RegisterWrite("temp", func(name string, data []byte, args ...interface{}) error {
		if strings.Index(name, "temp ") == 0 {
			name = name[5:]
		}
		// 定位到临时文件夹写入
		name = mc.ReplaceAllString(name, "/")
		name = joinPath(tempDir(), name)
		return Defaul.WriteByMethod("file", name, data, args...)
	})

	// 缓存文件
	Defaul.RegisterRegexp("buff", `^buff `)
	Defaul.RegisterRead("buff", func(name string, args ...interface{}) ([]byte, error) {
		if strings.Index(name, "buff ") == 0 {
			name = name[5:]
		}
		// 从缓存里查找
		d, err := Defaul.ReadByMethod("temp", name)
		if err != nil {
			// 找不到重新加载
			return Defaul.ReadByMethod("rebuff", name, d)
		}
		return d, nil
	})

	// 重新加载缓存文件
	Defaul.RegisterRegexp("rebuff", `^rebuff `)
	Defaul.RegisterRead("rebuff", func(name string, args ...interface{}) ([]byte, error) {
		if strings.Index(name, "rebuff ") == 0 {
			name = name[7:]
		}
		// 加载请求
		d, err := Defaul.Read(name, args...)
		if err != nil {
			return nil, err
		}
		// 写缓存
		Defaul.WriteByMethod("temp", name, d)
		return d, nil
	})
}

package fbuf

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

var Defaul = newFbuf()

func init() {

	// 本地文件
	Defaul.RegisterReg("file", `^\w:/`, `^/`, `^./`, `^\w+/`)
	Defaul.RegisterRead("file", func(name string, args ...interface{}) ([]byte, error) {
		return ioutil.ReadFile(name)
	})
	Defaul.RegisterWrite("file", func(name string, data []byte, args ...interface{}) error {
		pr := perm
		for _, v := range args {
			switch v.(type) {
			case os.FileMode:
				pr = v.(os.FileMode)
			}
		}
		err := ioutil.WriteFile(name, data, pr)
		if err != nil {
			os.MkdirAll(filepath.Dir(name), pr)
			err = ioutil.WriteFile(name, data, pr)
		}
		return err
	})
}

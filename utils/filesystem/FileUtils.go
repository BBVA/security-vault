package filesystem

import (
	"os"
	"io/ioutil"
)

type FileUtils interface {
	Write(file string, content []byte, perm os.FileMode) error
	Read(file string) ([]byte, error)
}

type DefaultFileUtils struct{}

func (*DefaultFileUtils) Write(file string, content []byte, perm os.FileMode) error {
	return ioutil.WriteFile(file, content, perm)
}
func (*DefaultFileUtils) Read(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

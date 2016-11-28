package filesystem

import (
	"io/ioutil"
	"os"
)

type FileUtils interface {
	WriteFile(file string, content []byte, perm os.FileMode) error
	ReadFile(file string) ([]byte, error)
	ReadDir(path string) ([]os.FileInfo, error)
	Remove (name string) error
	Getenv(string) string
}

type DefaultFileUtils struct{}

func (*DefaultFileUtils) WriteFile(file string, content []byte, perm os.FileMode) error {
	return ioutil.WriteFile(file, content, perm)
}
func (*DefaultFileUtils) ReadFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}
func (*DefaultFileUtils) ReadDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(path)
}
func (*DefaultFileUtils) Remove(name string) error {
	return os.Remove(name)
}
func (*DefaultFileUtils) Getenv(env string) string {
	return os.Getenv(env)
}

package filesystem

import (
	"os"
	"io/ioutil"
)

type DirUtils interface {
	Lstat(mountPoint string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
	IsNotExist(err error) bool
	IsExist(err error) bool
	RemoveAll(path string) error
	ReadDir(dir string) ([]os.FileInfo, error)
}

type DefaultDirUtils struct{}

func (d *DefaultDirUtils) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (d *DefaultDirUtils) Lstat(mountPoint string) (os.FileInfo, error) {
	return os.Lstat(mountPoint)
}

func (d *DefaultDirUtils) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (d *DefaultDirUtils) IsExist(err error) bool {
	return os.IsExist(err)
}

func (d *DefaultDirUtils) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (d *DefaultDirUtils) ReadDir(dir string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dir)
}

package filesystem

import "os"

type DirUtils interface {
	Lstat(mountPoint string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
}

type DefaultDirUtils struct{}

func (d DefaultDirUtils) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (d DefaultDirUtils) Lstat(mountPoint string) (os.FileInfo, error) {
	return os.Lstat(mountPoint)
}

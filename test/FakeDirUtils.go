package test

import (
	"os"
	"time"
)

type FakeDirUtils struct {
	lstatError             error
	exist                  bool
	lstatFileInfo          os.FileInfo
	mkdirError             error
	mkdirCalls             int
	mkdirExpectedCalls     int
	removeAllCalls         int
	removeAllExpectedCalls int
	removeAllError         error
	readDirCalls           int
	readDirExpectedCalls   int
	readDirError           error
	readDirFiles           []os.FileInfo
}

func (f *FakeDirUtils) Lstat(mountPoint string) (os.FileInfo, error) {
	return f.lstatFileInfo, f.lstatError
}

func (f *FakeDirUtils) MkdirAll(path string, perm os.FileMode) error {
	return f.mkdirError
}

func (f *FakeDirUtils) IsNotExist(err error) bool {
	return !f.exist
}

func (f *FakeDirUtils) IsExist(err error) bool {
	return f.exist
}

func (f *FakeDirUtils) RemoveAll(path string) error {
	f.removeAllCalls++
	return f.removeAllError
}

func (f *FakeDirUtils) ReadDir(dir string) ([]os.FileInfo, error) {
	f.readDirCalls++
	return f.readDirFiles, f.readDirError
}

type FakeFileInfo struct {
	name string
	isDir bool
}

func (f FakeFileInfo) Name() string {
	return f.name
}

func (f FakeFileInfo) Size() int64 {
	return 0
}

func (f FakeFileInfo) Mode() os.FileMode {
	return 0600
}

func (f FakeFileInfo) ModTime() time.Time {
	return time.Now()
}

func (f FakeFileInfo) IsDir() bool {
	return f.isDir
}

func (f FakeFileInfo) Sys() interface{} {
	return nil
}

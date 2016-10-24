package test

import "os"

type FakeDirUtils struct {
	lstatError     error
	exist          bool
	lstatFileInfo  os.FileInfo
	mkdirError     error
	removeAllCalls int
	removeAllExpectedCalls int
	removeAllError error
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

package main

import (
	"testing"
	"os"
)

type FakeDirUtils struct {
	lstatError    error
	lstatFileInfo os.FileInfo
	mkdirError    error
}

func (f FakeDirUtils) Lstat(mountPoint string) (os.FileInfo, error) {
	return f.lstatFileInfo, f.lstatError
}

func (f FakeDirUtils) MkdirAll(path string, perm os.FileMode) error {
	return f.mkdirError
}

func TestVaultDriver_Mount(t *testing.T) {



}


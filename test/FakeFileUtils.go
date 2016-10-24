package test

import "os"

type FakeFileUtils struct {
	writeCalls         int
	expectedWriteCalls int
	writeBytes         string
	writeError         error
	bytesRead          string
	readCalls          int
	readError          error
}

func (f *FakeFileUtils) Write(file string, content []byte, perm os.FileMode) error {
	f.writeCalls++
	f.writeBytes = string(content[:])
	return f.writeError
}

func (f *FakeFileUtils) Read(file string) ([]byte, error) {
	f.readCalls++
	return []byte(f.bytesRead), f.readError
}

package test

import (
	"os"
)

type WriteFileTestMetrics struct {
	error        error
	writtenBytes string
	MethodCallMetrics
}

type ReadFileTestMetrics struct {
	error   error
	content string
	MethodCallMetrics
}

type ReadEnvTestMetrics struct {
	content map[string]string
	MethodCallMetrics
}

type FakeFileUtils struct {
	writeFile WriteFileTestMetrics
	readFile  ReadFileTestMetrics
	readEnv   ReadEnvTestMetrics
}

func (f *FakeFileUtils) WriteFile(file string, content []byte, perm os.FileMode) error {
	f.writeFile.Call()
	f.writeFile.writtenBytes = string(content[:])
	return f.writeFile.error
}

func (f *FakeFileUtils) ReadFile(file string) ([]byte, error) {
	f.readFile.Call()
	return []byte(f.readFile.content), f.readFile.error
}

func (f *FakeFileUtils) Getenv(env string) string {
	f.readEnv.Call()
	value, ok := f.readEnv.content[env];
	if ok {
		return value
	} else {
		return ""
	}
}

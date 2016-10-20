package main

import (
	"testing"
	"bazil.org/fuse"
	"reflect"
	. "descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
	"errors"
	"os"
	"os/exec"
	"time"
)

type FakeFuseWrapper struct {
	mountError     error
	serveError     error
	connMountError error
	unmountError   error
	waitReady      func()
}

func (f FakeFuseWrapper) Mount(dir string, options ...fuse.MountOption) (*fuse.Conn, error) {
	return nil, f.mountError
}

func (f FakeFuseWrapper) Unmount(dir string) error {
	return f.unmountError
}

func (f FakeFuseWrapper) Serve(conn *fuse.Conn, ff *FS) error {
	return f.serveError
}

func (f FakeFuseWrapper) WaitReady(conn *fuse.Conn) {
	f.waitReady()
}

func (f FakeFuseWrapper) GetError(conn *fuse.Conn) error {
	return f.connMountError
}

func (f FakeFuseWrapper) CloseConnection(conn *fuse.Conn) error {
	return nil
}

func TestFS_Mount(t *testing.T) {

	fixtures := []struct {
		fuse             FakeFuseWrapper
		mountPoint       string
		expectedResponse error
	}{
		{
			fuse: FakeFuseWrapper{
				waitReady: func() {},
			},
			mountPoint: "test",
			expectedResponse: nil,
		},
		{
			fuse: FakeFuseWrapper{
				mountError: errors.New("error"),
				waitReady: func() {},
			},
			mountPoint: "test",
			expectedResponse: errors.New("error"),
		},
		{
			fuse: FakeFuseWrapper{
				connMountError: errors.New("error"),
				waitReady: func() {},
			},
			mountPoint: "test",
			expectedResponse: errors.New("error"),
		},
	}

	for i, fixture := range fixtures {
		f, _ := NewFS(fixture.mountPoint, fixture.fuse)

		err := f.Mount(fixture.mountPoint)

		if !reflect.DeepEqual(err, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, err)
		}
	}

}

func TestFS_MountCrashOnServe(t *testing.T) {

	fixtures := []struct {
		fuse       FakeFuseWrapper
		mountPoint string
	}{
		{
			fuse: FakeFuseWrapper{
				serveError: errors.New("error"),
				waitReady: func() {
					time.Sleep(2000)
				},
			},
			mountPoint: "test",
		},
	}

	wrapperForTestingCrashingFunction(t, "TestFS_MountCrashOnServe", func() {
		fixture := fixtures[0]

		f, _ := NewFS(fixture.mountPoint, fixture.fuse)

		f.Mount(fixture.mountPoint)
	})
}

func wrapperForTestingCrashingFunction(t *testing.T, test string, crasher func()) {
	if os.Getenv("BE_CRASHER") == "1" {
		crasher()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=" + test)
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	// Check that the program exited
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && e.Success() {
		t.Fatalf("Process ran with err %v, want exit status 1", err)
	}

}

func TestFS_Unmount(t *testing.T) {
	fixtures := []struct {
		fuse             FakeFuseWrapper
		mountPoint       string
		expectedResponse error
	}{
		{
			fuse: FakeFuseWrapper{
				waitReady: func() {},
			},
			mountPoint: "test",
			expectedResponse: nil,
		},
		{
			fuse: FakeFuseWrapper{
				unmountError: errors.New("error"),
				waitReady: func() {},
			},
			mountPoint: "test",
			expectedResponse: errors.New("error"),
		},
	}

	for i, fixture := range fixtures {
		f, _ := NewFS(fixture.mountPoint, fixture.fuse)

		err := f.Unmount()

		if !reflect.DeepEqual(err, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, err)
		}
	}
}
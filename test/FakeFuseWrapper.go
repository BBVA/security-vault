package test

import (
	"bazil.org/fuse"
	. "descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
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

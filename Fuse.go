package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"golang.org/x/net/context"
	"log"
	"os"
)

type FS struct {
	mountpoint string
	volumeName string
	conn       *fuse.Conn
	errChan    chan (error)
	server     *fs.Server
	//store      store.SecretStore
	//files      map[string]*File
	//tick       *time.Ticker
}

func NewFS(mountpoint string) (*FS, error) {
	c := make(chan error)
	go func() {
		err := <-c
		log.Fatalf("fs: %s", err.Error())
	}()

	return &FS{
		mountpoint: mountpoint,
		errChan:    c,
	}, nil
}

func (f *FS) Mount(volumeName string) error {
	log.Printf("setting up fuse: volume=%s", volumeName)
	c, err := fuse.Mount(
		f.mountpoint,
		fuse.FSName("vault"),
		fuse.Subtype("vaultfs"),
		fuse.LocalVolume(),
		fuse.VolumeName(volumeName),
		fuse.ReadOnly(),
		//fuse.NoExec(),
	)
	if err != nil {
		return err
	}

	srv := fs.New(c, nil)

	f.server = srv
	f.volumeName = volumeName
	f.conn = c

	go func() {
		err = f.server.Serve(f)
		if err != nil {
			f.errChan <- err
		}
	}()

	// check if the mount process has an error to report
	log.Println("waiting for mount")
	<-c.Ready
	if err := c.MountError; err != nil {
		return err
	}

	return nil

}

func (f *FS) Unmount() error {
	return fuse.Unmount(f.mountpoint)
}

func (FS) Root() (fs.Node, error) {
	return Dir{}, nil
}

type Dir struct{}

func (Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	return nil
}

func (Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == "credential" {
		return File{}, nil
	}
	return nil, fuse.ENOENT
}

var dirDirs = []fuse.Dirent{
	{Inode: 2, Name: "credential", Type: fuse.DT_File},
}

func (Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return dirDirs, nil
}

// File implements both Node and Handle for the hello file.
type File struct{}

const greeting = "hello cloudframe\n"

func (File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(len(greeting))
	return nil
}

func (File) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(greeting), nil
}

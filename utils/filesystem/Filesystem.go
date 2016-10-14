package filesystem

import (
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

type FS struct {
	Mountpoint string
	VolumeId   string
	conn       *fuse.Conn
	ErrChan    chan (error)
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
		Mountpoint: mountpoint,
		ErrChan:    c,
	}, nil
}

func (f *FS) Mount(volumeName string) error {
	log.Printf("setting up fuse: volume=%s", volumeName)
	c, err := fuse.Mount(
		f.Mountpoint,
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
	f.conn = c

	go func() {
		err = f.server.Serve(f)
		if err != nil {
			f.ErrChan <- err
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
	defer f.conn.Close()

	return fuse.Unmount(f.Mountpoint)
}

func (FS) Root() (fs.Node, error) {
	return Dir{}, nil
}

type Dir struct{}

func (Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Println("Fuse Dir Attr")
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	return nil
}

func (Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Println("Fuse Dir Lookup")
	switch name {
	case "cert":
		return &File{
			name:    "cert",
			inode:   2,
			content: []byte("certificadooorr\n"),
			mode:    0444,
		}, nil
	case "private":
		return &File{
			name:    "private",
			inode:   3,
			content: []byte("clave super privada\n"),
			mode:    0444,
		}, nil
	default:
		return nil, fuse.ENOENT

	}
}

var dirDirs = []fuse.Dirent{
	{Inode: 2, Name: "cert", Type: fuse.DT_File},
	{Inode: 3, Name: "private", Type: fuse.DT_File},
}

func (Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("Fuse Dir ReadDirAll")
	return dirDirs, nil
}

// File implements both Node and Handle for the hello file.
type File struct {
	name    string
	inode   uint64
	content []byte
	mode    os.FileMode
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Println("Fuse File Attr")
	a.Inode = f.inode
	a.Mode = f.mode
	a.Size = uint64(len(f.content))
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	log.Println("Fuse File ReadAll")
	return f.content, nil
}

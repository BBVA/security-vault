package fuseutils

import (
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
	"fmt"

	. "descinet.bbva.es/cloudframe-security-vault/SecretApi"
)

type Fuse interface {
	Mount(dir string, options ...fuse.MountOption) (*fuse.Conn, error)
	Unmount(dir string) error
	Serve(*fuse.Conn, *FS) error
	WaitReady(*fuse.Conn)
	GetError(*fuse.Conn) error
	CloseConnection(*fuse.Conn) error
}

type DefaultFuseWrapper struct{}

func (DefaultFuseWrapper) Mount(dir string, options ...fuse.MountOption) (*fuse.Conn, error) {
	return fuse.Mount(dir, options...)
}

func (DefaultFuseWrapper) Unmount(dir string) error {
	return fuse.Unmount(dir)
}

func (DefaultFuseWrapper) Serve(conn *fuse.Conn, f *FS) error {
	return fs.Serve(conn, f)
}

func (DefaultFuseWrapper) WaitReady(conn *fuse.Conn) {
	<-conn.Ready
}

func (DefaultFuseWrapper) GetError(conn *fuse.Conn) error {
	return conn.MountError
}

func (DefaultFuseWrapper) CloseConnection(conn *fuse.Conn) error {
	return conn.Close()
}

type FS struct {
	fuse          Fuse
	Mountpoint    string
	VolumeId      string
	conn          *fuse.Conn
	ErrChan       chan (error)
	server        *fs.Server
	secretHandler SecretApi
	//store      store.SecretStore
	//files      map[string]*File
	//tick       *time.Ticker
}

type Dir struct {
	secretHandler SecretApi
	dir           []fuse.Dirent
}

type File struct {
	Name    string
	Inode   uint64
	Mode    os.FileMode
	Size    int
	Content []byte
}

func NewFS(mountpoint string, fuse Fuse, secretHandler SecretApi) (*FS, error) {
	c := make(chan error)
	go func() {
		err := <-c
		log.Fatalf("fs: %s", err.Error())
	}()

	return &FS{
		fuse:          fuse,
		Mountpoint:    mountpoint,
		ErrChan:       c,
		secretHandler: secretHandler,
	}, nil
}

func (f *FS) Mount(volumeName string) error {
	log.Printf("setting up fuse: volume=%s", volumeName)
	var c *fuse.Conn
	c, err := f.fuse.Mount(
		f.Mountpoint,
		fuse.FSName("vault"),
		fuse.Subtype("vaultfs"),
		fuse.LocalVolume(),
		fuse.VolumeName(volumeName),
		fuse.ReadOnly(),
		fuse.AllowNonEmptyMount(),
		//fuse.NoExec(),
	)
	if err != nil {
		return err
	}

	f.conn = c

	go func() {
		fmt.Println("gorutine")
		err := f.fuse.Serve(c, f)
		if err != nil {
			f.ErrChan <- err
		}
	}()

	// check if the mount process has an error to report
	log.Println("waiting for mount")
	f.fuse.WaitReady(f.conn)
	if err := f.fuse.GetError(f.conn); err != nil {
		return err
	}

	return nil
}

func (f *FS) Unmount() error {
	defer f.fuse.CloseConnection(f.conn)

	return f.fuse.Unmount(f.Mountpoint)
}

func (f *FS) Root() (fs.Node, error) {
	return &Dir{
		secretHandler: f.secretHandler,
	}, nil
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Println("Fuse Dir Attr")
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	return nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Println("Fuse Dir Lookup")
	var inode uint64
	for _, file := range d.dir {
		//Es feo, lo sé, no se me ocurre nada más bonito.
		if name == file.Name {
			inode = file.Inode
			break
		}
	}

	if secret, err := d.secretHandler.GetSecret(name); err != nil {
		return nil,err
	} else {
		return &File{
			Name:    name,
			Inode:   inode,
			Mode:    0444,
			Content: secret,
			Size:    len(secret),
		}, nil
	}
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("Fuse Dir ReadDirAll")
	var dir []fuse.Dirent
	var inode uint64 = 2 // Because inode 1 is always the Dir itself.
	files := d.secretHandler.GetSecretFiles()
	for k := range files {
		dir = append(dir, fuse.Dirent{Inode: inode, Name: k, Type: fuse.DT_File})
		inode++
	}
	d.dir = dir
	return dir, nil
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Println("Fuse File Attr")
	a.Inode = f.Inode
	a.Mode = f.Mode
	a.Size = uint64(len(f.Content))
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	log.Println("Fuse File ReadAll")
	return f.Content, nil
}

package main

type FuseUtils interface {
	Mount(volumeId, mountPoint, volumeName string) error
	Unmount(volumeId string) error
}

type DefaultFuseUtils struct {
	fs map[string]*FS
}

func NewFuseUtils() FuseUtils {
	return DefaultFuseUtils{
		fs: make(map[string]*FS),
	}
}

func (d DefaultFuseUtils) Mount(volumeId, mountPoint, volumeName string) error {
	fs, err := NewFS(mountPoint)
	if err != nil {
		fs.errChan <- err
	}

	d.fs[volumeId] = fs

	return fs.Mount(volumeName)
}

func (d DefaultFuseUtils) Unmount(volumeId string) error {
	return d.fs[volumeId].Unmount()
}

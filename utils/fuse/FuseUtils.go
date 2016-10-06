package fuseutils

import (
	. "descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
)

type FuseUtils interface {
	Mount(volumeId, mountPoint, volumeName string) error
	Unmount(volumeId string) error
	Path(volumeName string) (string, error)
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
		fs.ErrChan <- err
	}
	fs.VolumeId = volumeId
	d.fs[volumeName] = fs

	return fs.Mount(volumeName)
}

func (d DefaultFuseUtils) Unmount(volumeName string) error {
	return d.fs[volumeName].Unmount()
}

func (d DefaultFuseUtils) Path(volumeName string) (string, error){
	fs, ok := d.fs[volumeName]
	if ok {
		return fs.Mountpoint, nil
	}
	return nil, ok
}

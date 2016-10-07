package fuseutils

import (
	. "descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
	"errors"
)

type VolumeOptions map[string]string

type Volume struct {
	Filesystem *FS
	Options    VolumeOptions
}

type FuseUtils interface {
	Mount(volumeId, mountPoint, volumeName string) error
	Unmount(volumeId string) error
	Path(volumeName string) (string, error)
}

type DefaultFuseUtils struct {
	vols map[string]*Volume
}

func NewFuseUtils() FuseUtils {
	return DefaultFuseUtils{
		vols: make(map[string]*Volume),
	}
}

func (d DefaultFuseUtils) Create(name string, options VolumeOptions) error {
	d.vols[name] = &Volume{
		Options: options,
		Filesystem: nil,
	}

	return nil
}

func (d DefaultFuseUtils) Mount(volumeId, mountPoint, volumeName string) error {
	fs, err := NewFS(mountPoint)
	if err != nil {
		fs.ErrChan <- err
	}

	fs.VolumeId = volumeId

	if _, ok := d.vols[volumeName]; !ok {
		d.vols[volumeName] = &Volume{
			Options: make(VolumeOptions),
			Filesystem: nil,
		}
	}

	d.vols[volumeName].Filesystem = fs

	return fs.Mount(volumeName)
}

func (d DefaultFuseUtils) Unmount(volumeName string) error {
	return d.vols[volumeName].Filesystem.Unmount()
}

func (d DefaultFuseUtils) Path(volumeName string) (string, error) {
	vol, ok := d.vols[volumeName]
	if ok {
		return vol.Filesystem.Mountpoint, nil
	}
	return "", errors.New("Volume not found")
}

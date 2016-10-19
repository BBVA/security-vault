package fuseutils

import (
	"errors"

	. "descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
)

type VolumeOptions map[string]string

type Volume struct {
	Filesystem *FS
	Options    VolumeOptions
}

type VolumeData struct {
	Name       string
	Mountpoint string
}

type FuseUtils interface {
	Mount(volumeId, mountPoint, volumeName string) error
	Unmount(volumeId string) error
	Path(volumeName string) (string, error)
	Create(volumeName string, options VolumeOptions) error
	Remove(volumeName string) error
	List() ([]VolumeData, error)
	Get(volumeName string) (VolumeData, error)
}

type DefaultFuseUtils struct {
	vols map[string]*Volume
	fuse Fuse
}

func NewFuseUtils(fuse Fuse) FuseUtils {
	return DefaultFuseUtils{
		vols: make(map[string]*Volume),
		fuse: fuse,
	}
}

func (d DefaultFuseUtils) List() ([]VolumeData, error) {
	var vols []VolumeData
	for name, data := range d.vols {
		vols = append(vols, VolumeData{
			Name:       name,
			Mountpoint: data.Filesystem.Mountpoint,
		})
	}

	return vols, nil
}

func (d DefaultFuseUtils) Create(name string, options VolumeOptions) error {
	d.vols[name] = &Volume{
		Options:    options,
		Filesystem: nil,
	}

	return nil
}

func (d DefaultFuseUtils) Remove(volumeName string) error {
	if _, ok := d.vols[volumeName]; ok {
		delete(d.vols, volumeName)
		return nil
	} else {
		return errors.New("Volume not found")
	}
}

func (d DefaultFuseUtils) Mount(volumeId, mountPoint, volumeName string) error {
	fs, err := NewFS(mountPoint, d.fuse)
	if err != nil {
		fs.ErrChan <- err
	}

	fs.VolumeId = volumeId

	if _, ok := d.vols[volumeName]; !ok {
		d.vols[volumeName] = &Volume{
			Options:    make(VolumeOptions),
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
	if vol, ok := d.vols[volumeName]; ok {
		return vol.Filesystem.Mountpoint, nil
	}
	return "", errors.New("Volume not found")
}

func (d DefaultFuseUtils) Get(volumeName string) (VolumeData, error) {
	if vol, ok := d.vols[volumeName]; ok {
		data := VolumeData{
			Name:       volumeName,
			Mountpoint: vol.Filesystem.Mountpoint,
		}

		return data, nil
	}
	return VolumeData{}, errors.New("Volume not found")
}

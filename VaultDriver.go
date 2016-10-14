package main

import (
	"log"
	"path"

	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
	"descinet.bbva.es/cloudframe-security-vault/utils/fuse"
	"github.com/docker/go-plugins-helpers/volume"
)

type VaultDriver struct {
	VolumePath string
	ServerUrl  string
	VaultToken string
	dirUtils   filesystem.DirUtils
	fuseUtils  fuseutils.FuseUtils
}

func NewVaultDriver(VolumePath string, ServerUrl string, VaultToken string, dirUtils filesystem.DirUtils, fuseUtils fuseutils.FuseUtils) VaultDriver {
	return VaultDriver{
		VolumePath: VolumePath,
		ServerUrl:  ServerUrl,
		VaultToken: VaultToken,
		dirUtils:   dirUtils,
		fuseUtils:  fuseUtils,
	}
}

func (d VaultDriver) Create(r volume.Request) volume.Response {
	log.Println("Create Volume", r)
	if err := d.fuseUtils.Create(r.Name, r.Options); err != nil {
		return volume.Response{Err: err.Error()}
	}
	return volume.Response{}
}

func (d VaultDriver) List(r volume.Request) volume.Response {
	log.Println("List Volumes", r)
	if volumes, err := d.fuseUtils.List(); err == nil {
		var v []*volume.Volume
		for _, vol := range volumes {
			v = append(v, &volume.Volume{
				Name:       vol.Name,
				Mountpoint: vol.Mountpoint,
			})
		}

		return volume.Response{Volumes: v}
	} else {
		return volume.Response{Err: err.Error()}
	}

}

func (d VaultDriver) Get(r volume.Request) volume.Response {
	log.Println("Get Volume", r)
	if vol, err := d.fuseUtils.Get(r.Name); err == nil {
		return volume.Response{
			Volume: &volume.Volume{
				Name:       vol.Name,
				Mountpoint: vol.Mountpoint,
			},
		}
	} else {
		return volume.Response{Err: err.Error()}
	}
}

func (d VaultDriver) Remove(r volume.Request) volume.Response {
	log.Println("Remove Volume", r)
	if err := d.fuseUtils.Remove(r.Name); err != nil {
		return volume.Response{Err: err.Error()}
	}
	return volume.Response{}
}

func (d VaultDriver) Path(r volume.Request) volume.Response {
	log.Println("Volume Path", r)
	mountPoint, err := d.fuseUtils.Path(r.Name)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}
	return volume.Response{Mountpoint: mountPoint}
}

func (d VaultDriver) Mount(r volume.MountRequest) volume.Response {
	log.Println("Mount Volume", r)

	mountPoint := path.Join(d.VolumePath, r.ID, r.Name)

	log.Println("check mountpoint", mountPoint)
	_, err := d.dirUtils.Lstat(mountPoint)

	if d.dirUtils.IsNotExist(err) {
		if err := d.dirUtils.MkdirAll(mountPoint, 0755); err != nil {
			return volume.Response{Err: err.Error()}
		}
	} else if err != nil {
		return volume.Response{Err: err.Error()}
	}

	log.Println("mount volume", mountPoint)

	if err := d.fuseUtils.Mount(r.ID, mountPoint, r.Name); err != nil {
		log.Println(err.Error())
		return volume.Response{Err: err.Error()}

	}

	log.Printf("response: %v\n", mountPoint)
	return volume.Response{Mountpoint: mountPoint}

}

func (d VaultDriver) Unmount(r volume.UnmountRequest) volume.Response {
	log.Println("Unmount Volume", r)

	if err := d.fuseUtils.Unmount(r.Name); err != nil {
		return volume.Response{Err: err.Error()}
	}

	mountPoint := path.Join(d.VolumePath, r.ID)
	if err := d.dirUtils.RemoveAll(mountPoint); err != nil {
		return volume.Response{Err: err.Error()}
	}

	log.Printf("Unmounted: %s\n", mountPoint)
	return volume.Response{}

}

func (d VaultDriver) Capabilities(r volume.Request) volume.Response {
	log.Println("Volume capabilities", r)
	return volume.Response{
		Capabilities: volume.Capability{
			Scope: "local",
		},
	}
}

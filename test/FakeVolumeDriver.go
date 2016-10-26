package test

import "github.com/docker/go-plugins-helpers/volume"

type FakeVolumeDriver struct {
	capabilitiesResponse  volume.Response
	getResponse           volume.Response
	listResponse          volume.Response
	pathResponse          volume.Response
	mountRequests         int
	mountRequestsExpected int
	mountResponse         volume.Response
	unmountResponse       volume.Response
}

func (p *FakeVolumeDriver) Create(r volume.Request) volume.Response {
	return volume.Response{}
}

func (p *FakeVolumeDriver) List(r volume.Request) volume.Response {
	return p.listResponse
}

func (p *FakeVolumeDriver) Get(r volume.Request) volume.Response {
	return p.getResponse
}

func (p *FakeVolumeDriver) Remove(r volume.Request) volume.Response {
	return volume.Response{}
}

func (p *FakeVolumeDriver) Path(r volume.Request) volume.Response {
	return p.pathResponse
}

func (p *FakeVolumeDriver) Mount(r volume.MountRequest) volume.Response {
	p.mountRequests++
	return p.mountResponse
}

func (p *FakeVolumeDriver) Unmount(r volume.UnmountRequest) volume.Response {
	return p.mountResponse
}

func (p *FakeVolumeDriver) Capabilities(r volume.Request) volume.Response {
	return p.capabilitiesResponse
}

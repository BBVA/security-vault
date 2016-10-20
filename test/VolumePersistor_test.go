package main

import (
	"testing"
	"reflect"
	"github.com/docker/go-plugins-helpers/volume"
)

type FakeVolumeDriver struct{}

func (p FakeVolumeDriver) Create(r volume.Request) volume.Response {
	return volume.Response{}
}

func (p FakeVolumeDriver) List(r volume.Request) volume.Response {
	return volume.Response{}
}

func (p FakeVolumeDriver) Get(r volume.Request) volume.Response {
	return volume.Response{}
}

func (p FakeVolumeDriver) Remove(r volume.Request) volume.Response {
	return volume.Response{}
}

func (p FakeVolumeDriver) Path(r volume.Request) volume.Response {
	return volume.Response{}
}

func (p FakeVolumeDriver) Mount(r volume.MountRequest) volume.Response {
	return volume.Response{}
}

func (p FakeVolumeDriver) Unmount(r volume.UnmountRequest) volume.Response {
	return volume.Response{}
}

func (p FakeVolumeDriver) Capabilities(r volume.Request) volume.Response {
	return volume.Response{
		Capabilities: volume.Capability{
			Scope: "local",
		},
	}
}

func TestVolumePersistor_Capabilities(t *testing.T) {

	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		volumeDriver     FakeVolumeDriver
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			volumeDriver: FakeVolumeDriver{},
			expectedResponse: volume.Response{
				Capabilities: volume.Capability{
					Scope: "local",
				},
			},
		},
	}

	for i, fixture := range fixtures {
		driver, _ := NewVolumePersistor("testpath", fixture.volumeDriver, fixture.dirUtils)

		response := driver.Capabilities(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}

}
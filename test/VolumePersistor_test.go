package test

import (
	"testing"
	"reflect"
	"github.com/docker/go-plugins-helpers/volume"
	. "descinet.bbva.es/cloudframe-security-vault"
)

func TestNewVolumePersistor(t *testing.T) {
	t.Skip("Not yet implemented")
}

func TestVolumePersistor_Create(t *testing.T) {
	t.Skip("Not yet implemented")
}

func TestVolumePersistor_List(t *testing.T) {
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
			volumeDriver: FakeVolumeDriver{
				listResponse: volume.Response{},
			},
			expectedResponse: volume.Response{},
		},
		{
			request: volume.Request{},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			volumeDriver: FakeVolumeDriver{
				listResponse: volume.Response{Err: "error"},
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver, _ := NewVolumePersistor("testpath", fixture.volumeDriver, fixture.dirUtils)

		response := driver.List(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}

}

func TestVolumePersistor_Get(t *testing.T) {
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
			volumeDriver: FakeVolumeDriver{
				getResponse: volume.Response{},
			},
			expectedResponse: volume.Response{},
		},
		{
			request: volume.Request{},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			volumeDriver: FakeVolumeDriver{
				getResponse: volume.Response{Err: "error"},
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver, _ := NewVolumePersistor("testpath", fixture.volumeDriver, fixture.dirUtils)

		response := driver.Get(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}

}

func TestVolumePersistor_Remove(t *testing.T) {
	t.Skip("Not yet implemented")
}

func TestVolumePersistor_Path(t *testing.T) {
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
			volumeDriver: FakeVolumeDriver{
				pathResponse: volume.Response{},
			},
			expectedResponse: volume.Response{},
		},
		{
			request: volume.Request{},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			volumeDriver: FakeVolumeDriver{
				pathResponse: volume.Response{Err: "error"},
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver, _ := NewVolumePersistor("testpath", fixture.volumeDriver, fixture.dirUtils)

		response := driver.Path(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}

}

func TestVolumePersistor_Mount(t *testing.T) {
	t.Skip("Not yet implemented")
}

func TestVolumePersistor_Unmount(t *testing.T) {
	t.Skip("Not yet implemented")
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
			volumeDriver: FakeVolumeDriver{
				capabilitiesResponse:volume.Response{},
			},
			expectedResponse: volume.Response{},
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
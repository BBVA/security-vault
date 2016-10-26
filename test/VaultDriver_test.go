package test

import (
	"errors"
	"path"
	"reflect"
	"testing"

	"descinet.bbva.es/cloudframe-security-vault/utils/fuse"
	"github.com/docker/go-plugins-helpers/volume"
	. "descinet.bbva.es/cloudframe-security-vault"
)

func TestVaultDriver_Mount(t *testing.T) {
	fixtures := []struct {
		mountRequest     volume.MountRequest
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
		expectedResponse volume.Response
	}{
		{
			mountRequest: volume.MountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
			},
			expectedResponse: volume.Response{Mountpoint: path.Join("testpath", "abcdef1234567890", "Test_volume")},
		},
		{
			mountRequest: volume.MountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:    errors.New("error"),
				exist:         true,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
			},
			expectedResponse: volume.Response{Err: "error"},
		},
		{
			mountRequest: volume.MountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    errors.New("error"),
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
			},
			expectedResponse: volume.Response{Err: "error"},
		},
		{
			mountRequest: volume.MountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         true,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   errors.New("error"),
				UnmountError: nil,
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		d := NewVaultDriver("testpath", "testserver", "testtoken", &fixture.dirUtils, fixture.fuseUtils)
		response := d.Mount(fixture.mountRequest)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}
}

func TestVaultDriver_Unmount(t *testing.T) {
	fixtures := []struct {
		unmountRequest   volume.UnmountRequest
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
		expectedResponse volume.Response
	}{
		{
			unmountRequest: volume.UnmountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
			},
			expectedResponse: volume.Response{},
		},
		{
			unmountRequest: volume.UnmountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: errors.New("error"),
			},
			expectedResponse: volume.Response{Err: "error"},
		},
		{
			unmountRequest: volume.UnmountRequest{
				Name: "Test_volume",
				ID:   "abcdef1234567890",
			},
			dirUtils: FakeDirUtils{
				lstatError:     nil,
				exist:          false,
				lstatFileInfo:  nil,
				mkdirError:     nil,
				removeAllError: errors.New("error"),
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver := NewVaultDriver("testpath", "testserver", "testtoken", &fixture.dirUtils, fixture.fuseUtils)

		response := driver.Unmount(fixture.unmountRequest)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}

}

func TestVaultDriver_Path(t *testing.T) {
	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{
				Name: "Test_volume",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
				path: Path4FakeFuse{
					path: "Test_volume",
					err:  nil,
				},
			},
			expectedResponse: volume.Response{Mountpoint: "Test_volume"},
		},
		{
			request: volume.Request{
				Name: "Test_volume",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
				path: Path4FakeFuse{
					path: "",
					err:  errors.New("error"),
				},
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver := NewVaultDriver("testpath", "testserver", "testtoken", &fixture.dirUtils, fixture.fuseUtils)

		response := driver.Path(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}
}

func TestVaultDriver_Capabilities(t *testing.T) {
	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
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
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
				path: Path4FakeFuse{
					path: "Test_volume",
					err:  nil,
				},
			},
			expectedResponse: volume.Response{
				Capabilities: volume.Capability{
					Scope: "local",
				},
			},
		},
	}

	for i, fixture := range fixtures {
		driver := NewVaultDriver("testpath", "testserver", "testtoken", &fixture.dirUtils, fixture.fuseUtils)

		response := driver.Capabilities(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}
}

func TestVaultDriver_Create(t *testing.T) {
	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{
				Name: "test_vol",
				Options: map[string]string{
					"key1": "val1",
				},
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
				CreateError:  nil,
				path: Path4FakeFuse{
					path: "Test_volume",
					err:  nil,
				},
			},
			expectedResponse: volume.Response{},
		},
		{
			request: volume.Request{
				Name: "test_vol",
				Options: map[string]string{
					"key1": "val1",
				},
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				MountError:   nil,
				UnmountError: nil,
				CreateError:  errors.New("error"),
				path: Path4FakeFuse{
					path: "Test_volume",
					err:  nil,
				},
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver := NewVaultDriver("testpath", "testserver", "testtoken", &fixture.dirUtils, fixture.fuseUtils)

		response := driver.Create(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}
}

func TestVaultDriver_Remove(t *testing.T) {
	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{
				Name: "test_vol",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				RemoveError: nil,
			},
			expectedResponse: volume.Response{},
		},
		{
			request: volume.Request{
				Name: "test_vol",
				Options: map[string]string{
					"key1": "val1",
				},
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				RemoveError: errors.New("error"),
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver := NewVaultDriver("testpath", "testserver", "testtoken", &fixture.dirUtils, fixture.fuseUtils)

		response := driver.Remove(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}
}

func TestVaultDriver_List(t *testing.T) {
	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
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
			fuseUtils: FakeFuseUtils{
				list: List4FakeFuse{
					volumes: []fuseutils.VolumeData{
						{
							Name:       "test1",
							Mountpoint: "test1",
						},
						{
							Name:       "test2",
							Mountpoint: "test2",
						},
					},
					err: nil,
				},
			},
			expectedResponse: volume.Response{
				Volumes: []*volume.Volume{
					{
						Name:       "test1",
						Mountpoint: "test1",
					},
					{
						Name:       "test2",
						Mountpoint: "test2",
					},
				},
			},
		},
		{
			request: volume.Request{},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				list: List4FakeFuse{
					err: errors.New("error"),
				},
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver := NewVaultDriver("testpath", "testserver", "testtoken", &fixture.dirUtils, fixture.fuseUtils)

		response := driver.List(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}
}

func TestVaultDriver_Get(t *testing.T) {
	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		fuseUtils        FakeFuseUtils
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{
				Name: "test_vol",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				get: Get4FakeFuse{
					volume: fuseutils.VolumeData{
						Name:       "test_vol",
						Mountpoint: "test_mount",
					},
					err: nil,
				},
			},
			expectedResponse: volume.Response{
				Volume: &volume.Volume{
					Name:       "test_vol",
					Mountpoint: "test_mount",
				},
			},
		},
		{
			request: volume.Request{
				Name: "test_vol",
			},
			dirUtils: FakeDirUtils{
				lstatError:    nil,
				exist:         false,
				lstatFileInfo: nil,
				mkdirError:    nil,
			},
			fuseUtils: FakeFuseUtils{
				get: Get4FakeFuse{
					volume: fuseutils.VolumeData{},
					err:    errors.New("error"),
				},
			},
			expectedResponse: volume.Response{
				Err: "error",
			},
		},
	}

	for i, fixture := range fixtures {
		driver := NewVaultDriver("testpath", "testserver", "testtoken", &fixture.dirUtils, fixture.fuseUtils)

		response := driver.Get(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}
}
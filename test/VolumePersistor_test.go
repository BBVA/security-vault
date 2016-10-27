package test

import (
	. "descinet.bbva.es/cloudframe-security-vault"
	"errors"
	"github.com/docker/go-plugins-helpers/volume"
	"os"
	"reflect"
	"testing"
)

func TestNewVolumePersistor(t *testing.T) {
	fixtures := []struct {
		dirUtils         FakeDirUtils
		fileUtils        FakeFileUtils
		volumeDriver     FakeVolumeDriver
		expectedResponse error
	}{
		// Folder exists, single file
		{
			dirUtils: FakeDirUtils{
				lstatError:           nil,
				exist:                true,
				mkdirExpectedCalls:   0,
				readDirExpectedCalls: 1,
				readDirFiles: []os.FileInfo{
					&FakeFileInfo{
						name: "Test.json",
					},
				},
			},
			fileUtils: FakeFileUtils{
				expectedReadCalls: 1,
				bytesRead:         "{\"Name\":\"test\",\"ID\":\"1234567\"}",
			},
			volumeDriver: FakeVolumeDriver{
				mountRequestsExpected: 1,
				mountResponse:         volume.Response{},
			},
			expectedResponse: nil,
		},
		// Folder exists, multiple files file
		{
			dirUtils: FakeDirUtils{
				lstatError:           nil,
				exist:                true,
				mkdirExpectedCalls:   0,
				readDirExpectedCalls: 1,
				readDirFiles: []os.FileInfo{
					&FakeFileInfo{
						name: "Test.json",
					},
					&FakeFileInfo{
						name: "Test2.json",
					},
				},
			},
			fileUtils: FakeFileUtils{
				expectedReadCalls: 2,
				bytesRead:         "{\"Name\":\"test\",\"ID\":\"1234567\"}",
			},
			volumeDriver: FakeVolumeDriver{
				mountRequestsExpected: 2,
				mountResponse:         volume.Response{},
			},
			expectedResponse: nil,
		},
		// Folder exists, no files
		{
			dirUtils: FakeDirUtils{
				lstatError:           nil,
				exist:                true,
				mkdirExpectedCalls:   0,
				readDirExpectedCalls: 1,
				readDirFiles:         []os.FileInfo{},
			},
			fileUtils: FakeFileUtils{
				expectedReadCalls: 0,
				bytesRead:         "",
			},
			volumeDriver: FakeVolumeDriver{
				mountRequestsExpected: 0,
				mountResponse:         volume.Response{},
			},
			expectedResponse: nil,
		},
		// Folder do not exist, no files
		{
			dirUtils: FakeDirUtils{
				lstatError:           nil,
				exist:                false,
				mkdirExpectedCalls:   1,
				readDirExpectedCalls: 0,
				readDirFiles:         []os.FileInfo{},
			},
			fileUtils: FakeFileUtils{
				expectedReadCalls: 0,
				bytesRead:         "",
			},
			volumeDriver: FakeVolumeDriver{
				mountRequestsExpected: 0,
				mountResponse:         volume.Response{},
			},
			expectedResponse: nil,
		},
	}

	for i, fixture := range fixtures {
		_, err := NewVolumePersistor("testpath", &fixture.volumeDriver, &fixture.dirUtils, &fixture.fileUtils)

		if !reflect.DeepEqual(err, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, err)
		}

		if fixture.dirUtils.readDirExpectedCalls != fixture.dirUtils.readDirCalls {
			t.Errorf("%d - Expected %d ReadDir calls, received %d\n", i, fixture.dirUtils.readDirExpectedCalls, fixture.dirUtils.readDirCalls)
		}

		if fixture.fileUtils.expectedReadCalls != fixture.fileUtils.readCalls {
			t.Errorf("%d - Expected %d Read calls, received %d\n", i, fixture.fileUtils.expectedReadCalls, fixture.fileUtils.readCalls)
		}

		if fixture.volumeDriver.mountRequestsExpected != fixture.volumeDriver.mountRequests {
			t.Errorf("%d - Expected %d Mount requests, received %d\n", i, fixture.volumeDriver.mountRequestsExpected, fixture.volumeDriver.mountRequests)
		}

	}
}

func TestVolumePersistor_Create(t *testing.T) {
	t.Skip("Not yet implemented")
}

func TestVolumePersistor_List(t *testing.T) {
	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		fileUtils        FakeFileUtils
		volumeDriver     FakeVolumeDriver
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{},
			volumeDriver: FakeVolumeDriver{
				listResponse: volume.Response{},
			},
			expectedResponse: volume.Response{},
		},
		{
			request: volume.Request{},
			volumeDriver: FakeVolumeDriver{
				listResponse: volume.Response{Err: "error"},
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver, _ := NewVolumePersistor("testpath", &fixture.volumeDriver, &fixture.dirUtils, &fixture.fileUtils)

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
		fileUtils        FakeFileUtils
		volumeDriver     FakeVolumeDriver
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{},
			volumeDriver: FakeVolumeDriver{
				getResponse: volume.Response{},
			},
			expectedResponse: volume.Response{},
		},
		{
			request: volume.Request{},
			volumeDriver: FakeVolumeDriver{
				getResponse: volume.Response{Err: "error"},
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver, _ := NewVolumePersistor("testpath", &fixture.volumeDriver, &fixture.dirUtils, &fixture.fileUtils)

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
		fileUtils        FakeFileUtils
		volumeDriver     FakeVolumeDriver
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{},
			volumeDriver: FakeVolumeDriver{
				pathResponse: volume.Response{},
			},
			expectedResponse: volume.Response{},
		},
		{
			request: volume.Request{},
			volumeDriver: FakeVolumeDriver{
				pathResponse: volume.Response{Err: "error"},
			},
			expectedResponse: volume.Response{Err: "error"},
		},
	}

	for i, fixture := range fixtures {
		driver, _ := NewVolumePersistor("testpath", &fixture.volumeDriver, &fixture.dirUtils, &fixture.fileUtils)

		response := driver.Path(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}

}

func TestVolumePersistor_Mount(t *testing.T) {
	fixtures := []struct {
		request             volume.MountRequest
		dirUtils            FakeDirUtils
		fileUtils           FakeFileUtils
		volumeDriver        FakeVolumeDriver
		expectedResponse    volume.Response
		expectedFileContent string
	}{
		{
			request: volume.MountRequest{
				Name: "test",
				ID:   "1234567",
			},
			fileUtils: FakeFileUtils{
				writeError:         nil,
				expectedWriteCalls: 1,
			},
			volumeDriver: FakeVolumeDriver{
				mountResponse: volume.Response{},
			},
			expectedResponse:    volume.Response{},
			expectedFileContent: "{\"Name\":\"test\",\"ID\":\"1234567\"}",
		},
		{
			request: volume.MountRequest{
				Name: "test",
				ID:   "1234567",
			},
			fileUtils: FakeFileUtils{
				writeError:         errors.New("error"),
				expectedWriteCalls: 1,
			},
			volumeDriver: FakeVolumeDriver{
				mountResponse: volume.Response{},
			},
			expectedResponse:    volume.Response{Err: "Could not persist volume data: error"},
			expectedFileContent: "{\"Name\":\"test\",\"ID\":\"1234567\"}",
		},
	}

	for i, fixture := range fixtures {
		driver, _ := NewVolumePersistor("testpath", &fixture.volumeDriver, &fixture.dirUtils, &fixture.fileUtils)

		response := driver.Mount(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected response %v, received %v\n", i, fixture.expectedResponse, response)
		}

		if fixture.fileUtils.writeCalls != fixture.fileUtils.expectedWriteCalls {
			t.Errorf("%d - Expected %d write calls, received %d\n", i, fixture.fileUtils.expectedWriteCalls, fixture.fileUtils.writeCalls)
		}

		if !reflect.DeepEqual(fixture.fileUtils.writeBytes, fixture.expectedFileContent) {
			t.Errorf("%d - Expected file content %s, received %s\n", i, fixture.expectedFileContent, fixture.fileUtils.writeBytes)
		}
	}
}

func TestVolumePersistor_Unmount(t *testing.T) {
	fixtures := []struct {
		request          volume.UnmountRequest
		dirUtils         FakeDirUtils
		fileUtils        FakeFileUtils
		volumeDriver     FakeVolumeDriver
		expectedResponse volume.Response
	}{
		{
			request: volume.UnmountRequest{
				Name: "test",
				ID:   "1234567",
			},
			dirUtils: FakeDirUtils{
				removeAllExpectedCalls: 1,
			},
			volumeDriver: FakeVolumeDriver{
				unmountResponse: volume.Response{},
			},
			expectedResponse: volume.Response{},
		},
		{
			request: volume.UnmountRequest{
				Name: "test",
				ID:   "1234567",
			},
			dirUtils: FakeDirUtils{
				removeAllExpectedCalls: 1,
				removeAllError:         errors.New("error"),
			},
			volumeDriver: FakeVolumeDriver{
				unmountResponse: volume.Response{},
			},
			expectedResponse: volume.Response{Err: "Error removing volume data: error"},
		},
	}

	for i, fixture := range fixtures {
		driver, _ := NewVolumePersistor("testpath", &fixture.volumeDriver, &fixture.dirUtils, &fixture.fileUtils)

		response := driver.Unmount(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected response %v, received %v\n", i, fixture.expectedResponse, response)
		}

		if fixture.dirUtils.removeAllExpectedCalls != fixture.dirUtils.removeAllCalls {
			t.Errorf("%d - Expected %d removeAll calls, received %d\n", i, fixture.dirUtils.removeAllExpectedCalls, fixture.dirUtils.removeAllCalls)
		}
	}
}

func TestVolumePersistor_Capabilities(t *testing.T) {

	fixtures := []struct {
		request          volume.Request
		dirUtils         FakeDirUtils
		fileUtils        FakeFileUtils
		volumeDriver     FakeVolumeDriver
		expectedResponse volume.Response
	}{
		{
			request: volume.Request{},
			volumeDriver: FakeVolumeDriver{
				capabilitiesResponse: volume.Response{},
			},
			expectedResponse: volume.Response{},
		},
	}

	for i, fixture := range fixtures {
		driver, _ := NewVolumePersistor("testpath", &fixture.volumeDriver, &fixture.dirUtils, &fixture.fileUtils)

		response := driver.Capabilities(fixture.request)

		if !reflect.DeepEqual(response, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, response)
		}
	}

}

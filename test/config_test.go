package test

import (
	"testing"
	"descinet.bbva.es/cloudframe-security-vault/utils/config"
	"github.com/facebookgo/inject"
	"reflect"
	"errors"
)

func TestConfig_ReadConfig(t *testing.T) {

	fixtures := []struct {
		fileUtils        FakeFileUtils
		expectedResponse error
	}{
		{
			fileUtils: FakeFileUtils{
				readEnv: ReadEnvTestMetrics{
					content: map[string]string{
						"VAULT_SERVER": "test",
						"TOKEN_PATH": "test",
						"SECRET_PATH": "test",
						"ROLE": "test",
						"PERSISTENCE_PATH": "test",
					},
					MethodCallMetrics: MethodCallMetrics{
						method: "Readenv",
						expectedCalls: 5,
					},
				},
			},
			expectedResponse: nil,
		},
		{
			fileUtils: FakeFileUtils{
				readEnv: ReadEnvTestMetrics{
					content: map[string]string{
						"VAULT_SERVER": "test",
						"TOKEN_PATH": "test",
						"SECRET_PATH": "test",
						"ROLE": "test",
						"PERSISTENCE_PATH": "",
					},
					MethodCallMetrics: MethodCallMetrics{
						method: "Readenv",
						expectedCalls: 5,
					},
				},
			},
			expectedResponse: errors.New("Undefined configuration: persistencePath"),
		},
		{
			fileUtils: FakeFileUtils{
				readEnv: ReadEnvTestMetrics{
					content: map[string]string{
						"VAULT_SERVER": "test",
						"TOKEN_PATH": "test",
						"SECRET_PATH": "test",
						"ROLE": "test",
					},
					MethodCallMetrics: MethodCallMetrics{
						method: "Readenv",
						expectedCalls: 5,
					},
				},
			},
			expectedResponse: errors.New("Undefined configuration: persistencePath"),
		},
		{
			fileUtils: FakeFileUtils{
				readEnv: ReadEnvTestMetrics{
					content: map[string]string{
						"VAULT_SERVER": "test",
						"TOKEN_PATH": "test",
						"SECRET_PATH": "test",
						"PERSISTENCE_PATH": "test",
					},
					MethodCallMetrics: MethodCallMetrics{
						method: "Readenv",
						expectedCalls: 5,
					},
				},
			},
			expectedResponse: errors.New("Undefined configuration: role"),
		},

	}

	for i, fixture := range fixtures {
		cfg := config.Config{}

		if err := inject.Populate(&cfg, &fixture.fileUtils); err != nil {
			t.Error(err.Error())
		}

		actualResponse := cfg.ReadConfig()

		fixture.fileUtils.readEnv.Report(t, i)

		if !reflect.DeepEqual(actualResponse, fixture.expectedResponse) {
			t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, actualResponse)
		}
	}
}

func TestConfig_GetToken(t *testing.T) {
	fixtures := []struct {
		fileUtils       FakeFileUtils
		expectedContent string
		expectedError   error
	}{
		{
			fileUtils: FakeFileUtils{
				readEnv: ReadEnvTestMetrics{
					content: map[string]string{
						"VAULT_SERVER": "test",
						"TOKEN_PATH": "test",
						"SECRET_PATH": "test",
						"ROLE": "test",
						"PERSISTENCE_PATH": "test",
					},
					MethodCallMetrics: MethodCallMetrics{
						method: "Readenv",
						expectedCalls: 5,
					},
				},
				readFile: ReadFileTestMetrics{
					content: "token",
					MethodCallMetrics: MethodCallMetrics{
						method: "ReadFile",
						expectedCalls: 1,
					},
				},
			},
			expectedContent: "token",
			expectedError: nil,
		},
		{
			fileUtils: FakeFileUtils{
				readEnv: ReadEnvTestMetrics{
					content: map[string]string{
						"VAULT_SERVER": "test",
						"TOKEN_PATH": "test",
						"SECRET_PATH": "test",
						"ROLE": "test",
						"PERSISTENCE_PATH": "test",
					},
					MethodCallMetrics: MethodCallMetrics{
						method: "Readenv",
						expectedCalls: 5,
					},
				},
				readFile: ReadFileTestMetrics{
					content: "",
					error: errors.New("error"),
					MethodCallMetrics: MethodCallMetrics{
						method: "ReadFile",
						expectedCalls: 1,
					},
				},
			},
			expectedContent: "",
			expectedError: errors.New("error"),
		},
	}

	for i, fixture := range fixtures {
		cfg := config.Config{}

		if err := inject.Populate(&cfg, &fixture.fileUtils); err != nil {
			t.Error(err.Error())
		}

		cfg.ReadConfig()
		fixture.fileUtils.readEnv.Report(t, i)

		token, err := cfg.GetToken()
		fixture.fileUtils.readFile.Report(t, i)

		if (token != fixture.expectedContent) {
			t.Errorf("%d - Expected Token %v, Received %v\n", i, fixture.expectedContent, token)
		}

		if !reflect.DeepEqual(err, fixture.expectedError) {
			t.Errorf("%d - Expected %v, Received %v\n", i, fixture.expectedError, err)
		}
	}
}

func TestConfig_Get(t *testing.T) {
	fixtures := []struct {
		fileUtils       FakeFileUtils
		testKey         string
		expectedContent string
		expectedError   error
	}{
		{
			fileUtils: FakeFileUtils{
				readEnv: ReadEnvTestMetrics{
					content: map[string]string{
						"VAULT_SERVER": "test",
						"TOKEN_PATH": "test",
						"SECRET_PATH": "test",
						"ROLE": "test",
						"PERSISTENCE_PATH": "test",
					},
					MethodCallMetrics: MethodCallMetrics{
						method: "Readenv",
						expectedCalls: 5,
					},
				},
			},
			testKey: "role",
			expectedContent: "test",
			expectedError: nil,
		},
		{
			fileUtils: FakeFileUtils{
				readEnv: ReadEnvTestMetrics{
					content: map[string]string{
						"VAULT_SERVER": "test",
						"TOKEN_PATH": "test",
						"SECRET_PATH": "test",
						"ROLE": "test",
						"PERSISTENCE_PATH": "test",
					},
					MethodCallMetrics: MethodCallMetrics{
						method: "Readenv",
						expectedCalls: 5,
					},
				},
			},
			testKey: "undefined",
			expectedContent: "",
			expectedError: errors.New("Missing Key: undefined"),
		},
	}

	for i, fixture := range fixtures {
		cfg := config.Config{}

		if err := inject.Populate(&cfg, &fixture.fileUtils); err != nil {
			t.Error(err.Error())
		}

		cfg.ReadConfig()
		fixture.fileUtils.readEnv.Report(t, i)

		value, err := cfg.Get(fixture.testKey)

		if (value != fixture.expectedContent) {
			t.Errorf("%d - Expected Value %v, Received %v\n", i, fixture.expectedContent, value)
		}

		if !reflect.DeepEqual(err, fixture.expectedError) {
			t.Errorf("%d - Expected %v, Received %v\n", i, fixture.expectedError, err)
		}
	}
}
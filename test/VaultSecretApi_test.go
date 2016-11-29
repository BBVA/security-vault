package test

import (
	"testing"

	"descinet.bbva.es/cloudframe-security-vault/utils/config"
	"github.com/facebookgo/inject"
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
	"net/http/httptest"
	"net/http"
	"reflect"
	"errors"
	"strings"
)

var response = []byte(`{
  "request_id": "1234567898765432",
  "lease_id": "pki/issue/cloudframe-dot-wtf/7ad6cfa5-f04f-c62a-d477-f33210475d05",
  "lease_duration": 21600,
  "renewable": false,
  "data": {
    "certificate": "public",
    "issuing_ca": "issuing-ca",
    "ca_chain": ["ca-chain"],
    "private_key": "private-key",
    "private_key_type": "rsa",
    "serial_number": "39:dd:2e:90:b7:23:1f:8d:d3:7d:31:c5:1b:da:84:d0:5b:65:31:58"
    }
}`)

type GetSecretFilesFixture struct {
	vaultResponse    int
	vaultBody        []byte
	role             string
	secrets          *SecretApi.Secrets
	expectedResponse error
}

func TestFakeExampleSecretApi_GetSecretFiles(t *testing.T) {

	fixtures := []GetSecretFilesFixture{
		{
			vaultBody: response,
			vaultResponse:http.StatusOK,
			role: "test",
			secrets: &SecretApi.Secrets{
				Cacert: "issuing-ca",
				Private: "private-key",
				Public: "public",
				LeaseDuration: 21600,
				LeaseID: "pki/issue/cloudframe-dot-wtf/7ad6cfa5-f04f-c62a-d477-f33210475d05",
				Renewable: false,
			},
			expectedResponse: nil,

		},
		{
			vaultBody: []byte("{}"),
			vaultResponse:http.StatusNotFound,
			role: "test",
			secrets: nil,
			expectedResponse: errors.New("Error making API request."),

		},

	}

	for i, fixture := range fixtures {
		runTest(t, i, fixture)
	}

}

func runTest(t *testing.T, i int, fixture GetSecretFilesFixture) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(fixture.vaultResponse)
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture.vaultBody)
	}))
	defer ts.Close()

	cfg, err := setupConfiguration(ts.URL, "testToken")
	if err != nil {
		t.Error(err.Error())
	}

	api, err := SecretApi.NewVaultSecretApi(cfg)
	if err != nil {
		t.Error(err.Error())
	}

	secrets, err := api.GetSecretFiles(fixture.role)
	if err != nil && !strings.Contains(err.Error(), fixture.expectedResponse.Error()) {
		t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, err)
	}

	if !reflect.DeepEqual(secrets, fixture.secrets) {
		t.Errorf("%d - Expected %v, received %v\n", i, fixture.expectedResponse, err)
	}
}

func setupConfiguration(address, token string) (config.ConfigHandler, error) {
	var cfg config.Config
	fileUtils := FakeFileUtils{
		readEnv: ReadEnvTestMetrics{
			content: map[string]string{
				"VAULT_SERVER": address,
				"TOKEN_PATH": "test",
				"SECRET_PATH": "test",
				"ROLE": "test",
				"PERSISTENCE_PATH": "test",
			},
			MethodCallMetrics: DefaultReadEnvCallMetrics(),
		},
		readFile: ReadFileTestMetrics{
			content: token,
			error: nil,
			MethodCallMetrics: DefaultReadFileCallMetrics(),
		},
	}

	if err := inject.Populate(&cfg, &fileUtils); err != nil {
		return nil, err
	}
	cfg.ReadConfig()

	return &cfg, nil
}

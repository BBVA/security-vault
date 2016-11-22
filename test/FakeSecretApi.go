package test

import (
	"bytes"
)

type GetSecretFilesTestMetrics struct {
	secrets *bytes.Buffer
	error   error
	MethodCallMetrics
}

type DeleteSecretsTestMetrics struct {
	error error
	MethodCallMetrics
}

type FakeSecretApi struct {
	getSecretFilesTestMetrics GetSecretFilesTestMetrics
	deleteSecretsTestMetrics  DeleteSecretsTestMetrics
}

func (f *FakeSecretApi) GetSecretFiles() (*bytes.Buffer, error) {
	f.getSecretFilesTestMetrics.Call()
	return f.getSecretFilesTestMetrics.secrets, f.getSecretFilesTestMetrics.error

}

func (f *FakeSecretApi) DeleteSecrets(containerID string) error {
	f.deleteSecretsTestMetrics.Call()
	return f.deleteSecretsTestMetrics.error
}

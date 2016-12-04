package test

import "descinet.bbva.es/cloudframe-security-vault/SecretApi"

type GetSecretFilesTestMetrics struct {
	secrets SecretApi.Secrets
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

func (f *FakeSecretApi) GetSecretFiles(string) (*SecretApi.Secrets, error) {
	f.getSecretFilesTestMetrics.Call()
	return &f.getSecretFilesTestMetrics.secrets, f.getSecretFilesTestMetrics.error

}

func (f *FakeSecretApi) DeleteSecrets(containerID string) error {
	f.deleteSecretsTestMetrics.Call()
	return f.deleteSecretsTestMetrics.error
}

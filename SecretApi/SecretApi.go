package SecretApi

import "bazil.org/fuse"

type SecretApi interface {
	GetSecret(SecretID string) (string, error)
	LookupSecretsDir() []fuse.Dirent
}

type SecretApiHandler struct {
	GetSecretFunc func()
	LookupSecretsDirFunc func()
}

func NewSecretApi (handle SecretApi) *SecretApiHandler {
	return &SecretApiHandler{
		GetSecretFunc: handle.GetSecret,
		LookupSecretsDirFunc: handle.LookupSecretsDir,
	}
}

package SecretApi

type Secret struct {
	content []byte
	len     int
}

type SecretApi interface {
	GetSecret(SecretID string) ([]byte, error)
	GetSecretFiles() map[string]*secret
}

type SecretApiHandler struct {
	GetSecretFunc      func(SecretID string) ([]byte, error)
	GetSecretFilesFunc func() map[string]*secret
}

func NewSecretApi(handle SecretApi) *SecretApiHandler {
	return &SecretApiHandler{
		GetSecretFunc:      handle.GetSecret,
		GetSecretFilesFunc: handle.GetSecretFiles,
	}
}

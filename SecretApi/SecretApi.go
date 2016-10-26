package SecretApi

type Secret struct {
	content []byte
	len     int
}

type SecretApi interface {
	GetSecret(SecretID string) ([]byte, error)
	GetSecretFiles() map[string]*Secret
}


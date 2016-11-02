package SecretApi

import "bytes"

type Secret struct {
	content []byte
}

type SecretApi interface {
	GetSecret(SecretID string) []byte
	GetSecretFiles(SecretID string) *bytes.Buffer
}


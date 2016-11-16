package SecretApi

import "bytes"

type Secret struct {
	content []byte
}

type SecretApi interface {
	GetSecretFiles(SecretID string, containerID string) (*bytes.Buffer,error)
	DeleteSecrets(containerID string) error
}


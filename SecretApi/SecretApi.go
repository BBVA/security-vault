package SecretApi

import "bytes"

type Secret struct {
	content []byte
}

type SecretApi interface {
	GetSecretFiles(SecretID string) (*bytes.Buffer,error)
}


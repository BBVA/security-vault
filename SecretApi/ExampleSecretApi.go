package SecretApi

import (
	"github.com/pkg/errors"
)

type secret struct {
	content []byte
	len     int

}

type ExampleSecretApi struct {
	secrets map[string]*secret
}


func NewExampleSecretApi() *ExampleSecretApi {

	privateContent := []byte("clave super privada\n")
	certContent := []byte("certificadooorr\n")

	secrets := make(map[string]*secret)
	secrets["private"] = &secret{
		content: privateContent,
		len:     len(privateContent),
	}
	secrets["cert"] = &secret{
		content: certContent,
		len:     len(certContent),
	}

	return &ExampleSecretApi{
		secrets: secrets,
	}
}

func (Api *ExampleSecretApi) GetSecret(SecretID string) ([]byte, error) {

	secret, ok := Api.secrets[SecretID]
	if ok {
		return secret.content, nil
	}
	return nil, errors.New("No secret")

}

func (Api *ExampleSecretApi) GetSecretFiles() map[string]*secret {
	return Api.secrets
}


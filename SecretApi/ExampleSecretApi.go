package SecretApi

import (

	"github.com/pkg/errors"
	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
)

type ExampleSecretApi struct {
	secrets map[string]*Secret
}

func NewExampleSecretApi(cacert string, private string, public string, fileUtilsHandler filesystem.FileUtils) (*ExampleSecretApi, error) {


	privateContent, err := fileUtilsHandler.Read(private)
	if err != nil {
		return nil,err
	}

	publicContent, err := fileUtilsHandler.Read(public)
	if err != nil {
		return nil,err
	}

	caCertContent, err := fileUtilsHandler.Read(cacert)
	if err != nil {
		return nil,err
	}

	secrets := make(map[string]*Secret)
	secrets["private"] = &Secret{
		content: privateContent,
		len:     len(privateContent),
	}
	secrets["cacert"] = &Secret{
		content: caCertContent,
		len:     len(caCertContent),
	}
	secrets["public"] = &Secret{
		content: publicContent,
		len:     len(publicContent),
	}

	return &ExampleSecretApi{
		secrets: secrets,
	},nil
}

func (Api *ExampleSecretApi) GetSecret(SecretID string) ([]byte, error) {

	secret, ok := Api.secrets[SecretID]
	if ok {
		return secret.content, nil
	}
	return nil, errors.New("No secret")

}

func (Api *ExampleSecretApi) GetSecretFiles() map[string]*Secret {
	return Api.secrets
}

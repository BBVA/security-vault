package SecretApi

import (


	"descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
	"github.com/rancher/secrets-bridge/pkg/archive"
	"log"
	"bytes"
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
	}
	secrets["cacert"] = &Secret{
		content: caCertContent,
	}
	secrets["public"] = &Secret{
		content: publicContent,
	}

	return &ExampleSecretApi{
		secrets: secrets,
	},nil
}

func (Api *ExampleSecretApi) GetSecret(SecretID string) ([]byte) {

	secret, ok := Api.secrets[SecretID]
	if ok {
		return secret.content
	}
	return nil

}

func (Api *ExampleSecretApi) GetSecretFiles(SecretID string) *bytes.Buffer {
	//no usamos SecretID para nada porque en este caso no aplica.
	files := []archive.ArchiveFile{}

	for  k := range Api.secrets {
		message := archive.ArchiveFile{
			Name: k,
			Content: string(Api.GetSecret(k)),
		}
		files = append(files,message)
	}

	tarball, err := archive.CreateTarArchive(files)
	if err != nil {
		log.Printf("Failed to create Tar file")
	}
	return tarball
}

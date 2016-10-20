package SecretApi

import (
	"bazil.org/fuse"
	"os"
)

type secret struct {
	content []byte
	len     uint64
	mode    os.FileMode
}

type ExampleSecretApi struct {
	secrets map[string]secret
}

func NewExampleSecretApi() SecretApi {

	privateContent := []byte("clave super privada\n")
	certContent := []byte("certificadooorr\n")

	secrets := make(map[string]secret)
	secrets["private"] = secret{
		content: privateContent,
		mode:    0444,
		len:     len(privateContent),
	}
	secrets["cert"] = secret{
		content: certContent,
		mode:    0444,
		len:     len(certContent),
	}

	return ExampleSecretApi{
		secrets: secrets,
	}
}

func (Api *ExampleSecretApi) getSecret(SecretID string) (string, error) {

	secret, ok := Api.secrets[SecretID]
	if ok {
		return secret.content, nil
	}
	return nil, "No secret found boyyyzzzz\n"

}

func (Api *ExampleSecretApi) lookupSecretsDir() []fuse.Dirent {

	var dir []fuse.Dirent
	var inode = 2 // Because inode 1 is always the Dir itself.
	for k := range Api.secrets {
		append(dir, fuse.Dirent{Inode: inode, Name: k, Type: fuse.DT_File})
	}
	return dir
}

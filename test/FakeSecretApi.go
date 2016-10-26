package test

import (
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
)

type FakeExampleSecretApi struct {
	getSecretContent string
}

func (f *FakeExampleSecretApi) GetSecret(secretID string) ([]byte,error){
	return []byte(f.getSecretContent),nil
}
func (f *FakeExampleSecretApi) GetSecretFiles() map[string]*SecretApi.Secret{
	return map[string]*SecretApi.Secret{}

}

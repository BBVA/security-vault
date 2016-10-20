package SecretApi

import "bazil.org/fuse"

type SecretApi interface {
	getSecret(SecretID string) (string, error)
	lookupSecretsDir() []fuse.Dirent
}

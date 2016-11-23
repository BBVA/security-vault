package SecretApi

type Secret struct {
	content []byte
}

type Secrets struct {
	Public        string
	Private       string
	Cacert        string
	LeaseID       string
	LeaseDuration int
	Renewable     bool
}

type SecretApi interface {
	GetSecretFiles(string) (*Secrets, error)
	DeleteSecrets(string) error
}

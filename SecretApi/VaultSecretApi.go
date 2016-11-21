package SecretApi

import (
	"bytes"
	"descinet.bbva.es/cloudframe-security-vault/utils/config"
	"encoding/json"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/rancher/secrets-bridge/pkg/archive"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type leaseInfo struct {
	LeaseID   string `json:"lease_id"`
	LeaseTime int    `json:"lease_time"`
	Renewable bool   `json:"renewable"`
	Timestamp int64	`json:"timestamp"`
}
type leaseEvent struct {
	eventType   string
	containerID string
	lease       leaseInfo
}

type persistenceObject struct {
	leases map[string]leaseInfo
}

type VaultSecretApi struct {
	client             *vault.Client
	role               string
	leases             map[string]leaseInfo
	persistenceChannel chan leaseEvent
	config             config.Config
}

func NewVaultSecretApi(mainConfig config.Config) (*VaultSecretApi, error) {

	config := vault.DefaultConfig()

	if err := config.ReadEnvironment(); err != nil {
		return nil, err
	}

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	token, err := ioutil.ReadFile(mainConfig["tokenPath"])
	if err != nil {
		return nil, err
	}

	client.SetToken(string(token))
	client.SetAddress(mainConfig["vaultServer"])

	leases := make(map[string]leaseInfo)
	persistenceChannel := make(chan leaseEvent)

	return &VaultSecretApi{
		client:             client,
		role:               mainConfig["role"],
		leases:             leases,
		persistenceChannel: persistenceChannel,
		config:             mainConfig,
	}, nil
}

func (Api *VaultSecretApi) GetSecretFiles(commonName string, containerID string) (*bytes.Buffer, error) {
	fmt.Println("Generating secret\n")
	files := []archive.ArchiveFile{}
	params := make(map[string]interface{})
	params["common_name"] = commonName

	path := filepath.Join("pki/issue/", Api.role)

	secrets, err := Api.client.Logical().Write(path, params)
	if err != nil {
		return nil, err
	}

	files = append(files, archive.ArchiveFile{Name: "private", Content: secrets.Data["private_key"].(string)})
	files = append(files, archive.ArchiveFile{Name: "cacert", Content: secrets.Data["issuing_ca"].(string)})
	files = append(files, archive.ArchiveFile{Name: "public", Content: secrets.Data["certificate"].(string)})

	tarball, err := archive.CreateTarArchive(files)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Unix()

	Api.persistenceChannel <- leaseEvent{
		eventType:   "start",
		containerID: containerID,
		lease: leaseInfo{
			LeaseID:   secrets.LeaseID,
			LeaseTime: secrets.LeaseDuration,
			Renewable: secrets.Renewable,
			Timestamp: timestamp,
		},
	}

	return tarball, nil

}

func (Api *VaultSecretApi) DeleteSecrets(containerID string) error {
	fmt.Println("Deleting secret persistence..")
	event := leaseEvent{
		eventType:   "stop",
		containerID: containerID,
		lease:       leaseInfo{},
	}
	Api.persistenceChannel <- event

	return nil
}

func (Api *VaultSecretApi) PersistenceManager() {

	fmt.Println("Starting persistence goroutine..")
	path := Api.config["persistencePath"]

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		var lease leaseInfo
		filepath := filepath.Join(path, file.Name())
		content, err := ioutil.ReadFile(filepath)
		if err != nil {
			panic(err.Error())
		}
		if err := json.Unmarshal(content, &lease); err != nil {
			panic(err.Error())
		}
		fmt.Printf("Succesfully read persistence information for containerID: %s\nleaseID: %s\nleasetime: %v\nrenewable: %v\ntimestamp: v%\n",file.Name(),lease.LeaseID,lease.LeaseTime,lease.Renewable,lease.Timestamp)

		Api.leases[file.Name()] = lease
	}

	for {
		select {
		case event := <-Api.persistenceChannel:
			fmt.Println("Lease event received\n")

			switch event.eventType {
			case "start":
				fmt.Println("Start event processing")
				Api.leases[event.containerID] = event.lease
				bytes, err := json.Marshal(&event.lease)
				if err != nil {
					panic(err.Error())
				}
				file := filepath.Join(path, event.containerID)
				if err := ioutil.WriteFile(file, bytes, 0777); err != nil {
					panic(err.Error())
				}
				fmt.Printf("Succesfully write persistence information for containerID: %s\nleaseID: %s\nleasetime: %v\nrenewable: %v\ntimestamp: v%\n",event.containerID,event.lease.LeaseID,event.lease.LeaseTime,event.lease.Renewable,event.lease.Timestamp)
			case "stop":
				_, ok := Api.leases[event.containerID]
				if ok {
					delete(Api.leases, event.containerID)

					file := filepath.Join(path, event.containerID)
					if err := os.Remove(file); err != nil {
						panic(err.Error())
					}
					fmt.Printf("Deleted file: %s\n", file)
				}
			}

		}
	}
}

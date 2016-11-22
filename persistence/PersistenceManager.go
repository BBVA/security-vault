package persistence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"descinet.bbva.es/cloudframe-security-vault/utils/config"
)

type LeaseInfo struct {
	LeaseID   string `json:"lease_id"`
	LeaseTime int    `json:"lease_time"`
	Renewable bool   `json:"renewable"`
	Timestamp int64  `json:"timestamp"`
}
type LeaseEvent struct {
	EventType   string
	ContainerID string
	Lease       LeaseInfo
}

type PersistenceObject struct {
	leases map[string]LeaseInfo
}

type PersistenceManager struct {
	config             config.ConfigHandler
	persistenceChannel chan LeaseEvent
	leases             map[string]LeaseInfo
}

func NewPersistenceManager(cfg config.ConfigHandler) (chan LeaseEvent, *PersistenceManager) {

	leases := make(map[string]LeaseInfo)
	persistenceChannel := make(chan LeaseEvent)

	return persistenceChannel, &PersistenceManager{
		leases:             leases,
		persistenceChannel: persistenceChannel,
		config:             cfg,
	}
}

func (api *PersistenceManager) Run() {

	fmt.Println("Starting persistence goroutine..")
	path := api.config.GetPersistencePath()

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		var lease LeaseInfo
		filePath := filepath.Join(path, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err.Error())
		}
		if err := json.Unmarshal(content, &lease); err != nil {
			panic(err.Error())
		}
		fmt.Printf("Succesfully read persistence information for containerID: %s\nleaseID: %s\nleasetime: %v\nrenewable: %v\ntimestamp: v%\n", file.Name(), lease.LeaseID, lease.LeaseTime, lease.Renewable, lease.Timestamp)

		api.leases[file.Name()] = lease
	}

	for {
		select {
		case event := <-api.persistenceChannel:
			fmt.Println("Lease event received\n")

			switch event.EventType {
			case "start":
				fmt.Println("Start event processing")
				api.leases[event.ContainerID] = event.Lease
				bytes, err := json.Marshal(&event.Lease)
				if err != nil {
					panic(err.Error())
				}
				file := filepath.Join(path, event.ContainerID)
				if err := ioutil.WriteFile(file, bytes, 0777); err != nil {
					panic(err.Error())
				}
				fmt.Printf("Succesfully write persistence information for containerID: %s\nleaseID: %s\nleasetime: %v\nrenewable: %v\ntimestamp: v%\n", event.ContainerID, event.Lease.LeaseID, event.Lease.LeaseTime, event.Lease.Renewable, event.Lease.Timestamp)
			case "stop":
				_, ok := api.leases[event.ContainerID]
				if ok {
					delete(api.leases, event.ContainerID)

					file := filepath.Join(path, event.ContainerID)
					if err := os.Remove(file); err != nil {
						panic(err.Error())
					}
					fmt.Printf("Deleted file: %s\n", file)
				}
			}

		}
	}
}

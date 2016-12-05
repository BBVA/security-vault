package persistence

import (
	"encoding/json"
	"fmt"
	"path/filepath"
        "descinet.bbva.es/cloudframe-security-vault/utils/filesystem"
	"descinet.bbva.es/cloudframe-security-vault/utils/config"
	"sync"
)

type LeaseInfo struct {
	LeaseID   string `json:"lease_id"`
	LeaseTime int    `json:"lease_time"`
	Renewable bool   `json:"renewable"`
	Timestamp int64  `json:"timestamp"`
	CommonName string `json:"common_name"`
}
type LeaseEvent struct {
	EventType  string
	Identifier string
	Lease      LeaseInfo
}

type PersistenceObject struct {
	leases map[string]LeaseInfo
}

type PersistenceManager struct {
	config             config.ConfigHandler
	persistenceChannel chan LeaseEvent
	Leases             map[string]LeaseInfo
	LeaseMutex	   sync.RWMutex
	FileUtils          filesystem.FileUtils `inject:""`
}

func NewPersistenceManager(cfg config.ConfigHandler,persistenceCfg *PersistenceManager) (chan LeaseEvent, *PersistenceManager) {

	leases := make(map[string]LeaseInfo)
	persistenceChannel := make(chan LeaseEvent)

	return persistenceChannel, &PersistenceManager{
		Leases:             leases,
		persistenceChannel: persistenceChannel,
		config:             cfg,
		FileUtils:          persistenceCfg.FileUtils,
	}
}

func (p *PersistenceManager) RecoverLeases() error {
	path := p.config.GetPersistencePath()

	files, err := p.FileUtils.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		var lease LeaseInfo
		filePath := filepath.Join(path, file.Name())
		content, err := p.FileUtils.ReadFile(filePath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(content, &lease); err != nil {
			return err
		}
		fmt.Printf("Succesfully read persistence information for containerID: %s\nleaseID: %s\nleasetime: %v\nrenewable: %v\ntimestamp: v%\n", file.Name(), lease.LeaseID, lease.LeaseTime, lease.Renewable, lease.Timestamp)
		p.LeaseMutex.Lock()
		p.Leases[file.Name()] = lease
		p.LeaseMutex.Unlock()
	}

	return nil
}

func (p *PersistenceManager) Run() {
	path := p.config.GetPersistencePath()
	Infinity:
	for {
		select {
		case event := <-p.persistenceChannel:
			fmt.Println("Lease event received\n")

			switch event.EventType {
			case "start":
				fmt.Println("Start event processing")
				p.LeaseMutex.Lock()
				p.Leases[event.Identifier] = event.Lease
				p.LeaseMutex.Unlock()
				bytes, err := json.Marshal(&event.Lease)
				if err != nil {
					panic(err.Error())
				}
				file := filepath.Join(path, event.Identifier)
				if err := p.FileUtils.WriteFile(file, bytes, 0777); err != nil {
					panic(err.Error())
				}
				fmt.Printf("Succesfully write persistence information for containerID: %s\nleaseID: %s\nleasetime: %v\nrenewable: %v\ntimestamp: v%\n", event.Identifier, event.Lease.LeaseID, event.Lease.LeaseTime, event.Lease.Renewable, event.Lease.Timestamp)
			case "stop":
				fmt.Println("Stop event processing")
				p.LeaseMutex.Lock()
				_, ok := p.Leases[event.Identifier]
				if ok {
					delete(p.Leases, event.Identifier)
					p.LeaseMutex.Unlock()
					file := filepath.Join(path, event.Identifier)
					if err := p.FileUtils.Remove(file); err != nil {
						panic(err.Error())
					}
					fmt.Printf("Deleted file: %s\n", file)
				}
			case "dieHard":
				fmt.Printf("Die switch triggered\n. Stopping persistance manager")
				break Infinity
			}
		}
	}
}

package leaseManager

import (
	"descinet.bbva.es/cloudframe-security-vault/persistence"
	"fmt"
	"time"
	"descinet.bbva.es/cloudframe-security-vault/EventConnector"
)

func Run (persistenceInfo *persistence.PersistenceManager,dockerConnector *EventConnector.DockerConnector) {
	fmt.Println("STARTING LEASE MANAGER")
	for {
		now := time.Now().Unix()
		time.Sleep(15 * time.Second)

		persistenceInfo.LeaseMutex.RLock()
		for k,v := range persistenceInfo.Leases {
				if ((int64(v.LeaseTime) + v.Timestamp) < (now + 100)) {
					fmt.Printf("Renew of %s activated for %v",v.CommonName,k)
					if err := dockerConnector.CopySecretsToContainer(v.CommonName,k); err != nil {
						panic(err.Error())
					}
				}
		}
		persistenceInfo.LeaseMutex.RUnlock()
	}
}

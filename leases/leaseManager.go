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
		time.Sleep(300 * time.Second)

		persistenceInfo.LeaseMutex.RLock()
		for k,v := range persistenceInfo.Leases {
				if ((int64(v.LeaseTime) + v.Timestamp) < (now + 300)) {
					fmt.Printf("Renew of %s activated for %v",v.CommonName,k)
					if err := dockerConnector.CopySecretsToContainer(v.CommonName,k); err != nil {
						fmt.Printf("Error generating secrets for container %s: %s",k,err.Error())
					}
				}
		}
		persistenceInfo.LeaseMutex.RUnlock()
	}
}

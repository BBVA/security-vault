package test

import (
	"descinet.bbva.es/cloudframe-security-vault/EventConnector"
	"descinet.bbva.es/cloudframe-security-vault/SecretApi"
	"descinet.bbva.es/cloudframe-security-vault/persistence"
	"encoding/json"
	"github.com/docker/engine-api/types/events"
	"io"
	"testing"
	"time"
)

func TestDockerConnector_Start(t *testing.T) {
	secretHandler := &FakeSecretApi{
		getSecretFilesTestMetrics: GetSecretFilesTestMetrics{
			secrets: SecretApi.Secrets{
				Cacert:        "cacert",
				Private:       "private",
				Public:        "public",
				Renewable:     false,
				LeaseDuration: 0,
				LeaseID:       "1234567",
			},
			MethodCallMetrics: MethodCallMetrics{
				expectedCalls: 1,
			},
		},
	}
	cfg, _ := setupConfiguration("test", "test")

	eventsOut, eventsIn := io.Pipe()
	defer eventsIn.Close()

	cli := &FakeDockerCli{
		events: EventsTestMetrics{
			readCloser: eventsOut,
			error:      nil,
			MethodCallMetrics: MethodCallMetrics{
				expectedCalls: 1,
			},
		},
		copyToContainer: CopyToContainerMetrics{
			MethodCallMetrics: MethodCallMetrics{
				expectedCalls: 1,
			},
		},
	}

	persistenceChannel := make(chan persistence.LeaseEvent)
	defer close(persistenceChannel)

	connector, err := EventConnector.NewConnector(secretHandler, cfg, cli, persistenceChannel)
	if err != nil {
		t.Error(err.Error())
	}

	go connector.Start()
	defer connector.Stop()

	// Actual test

	msg := events.Message{
		Action: "start",
		ID:     "1234567",
		Actor: events.Actor{
			Attributes: map[string]string{
				"common_name": "test",
			},
		},
	}
	message, _ := json.Marshal(msg)
	eventsIn.Write(message)

	// End Test

	time.Sleep(100 * time.Millisecond)

	i := 0

	expectedEventsCalls := cli.events.expectedCalls
	actualEventsCalls := cli.events.actualCalls
	expectedVsActualCalls(t, i, "Client.Events", expectedEventsCalls, actualEventsCalls)

	expectedGetSecretFilesCalls := secretHandler.getSecretFilesTestMetrics.expectedCalls
	actualGetSecretFilesCalls := secretHandler.getSecretFilesTestMetrics.actualCalls
	expectedVsActualCalls(t, i, "SecretApi.GetSecretFiles", expectedGetSecretFilesCalls, actualGetSecretFilesCalls)

	expectedCopyToContainerCalls := cli.copyToContainer.expectedCalls
	actualCopyToContainerCalls := cli.copyToContainer.actualCalls
	expectedVsActualCalls(t, i, "Client.CopyToContainer", expectedCopyToContainerCalls, actualCopyToContainerCalls)

	expectedPersistenceEvent := persistence.LeaseEvent{
		EventType: "start",
		Identifier: "1234567",
		Lease: persistence.LeaseInfo{
			LeaseID: "1234567",
			Renewable: false,
			Timestamp: time.Now().Unix(),
		},
	}
	actualPersistenceEvent := <-persistenceChannel
	checkLeaseEvent(t, i, expectedPersistenceEvent, actualPersistenceEvent)

}

func expectedVsActualCalls(t *testing.T, index int, name string, expected, actual int) {
	if expected != actual {
		t.Errorf("%d - Expected %v Calls %v, Received %v\n", index, name, expected, actual)
	}
}

func checkLeaseEvent(t *testing.T, index int, expected, actual persistence.LeaseEvent) {
	if expected.EventType != actual.EventType {
		t.Errorf("%d - Expected PersistenceEvent %v, Received %v\n", index, expected, actual)
	}
	if expected.Identifier != actual.Identifier {
		t.Errorf("%d - Expected PersistenceEvent %v, Received %v\n", index, expected, actual)
	}
	if expected.Lease.LeaseID != actual.Lease.LeaseID {
		t.Errorf("%d - Expected PersistenceEvent %v, Received %v\n", index, expected, actual)
	}
	if expected.Lease.LeaseTime != actual.Lease.LeaseTime {
		t.Errorf("%d - Expected PersistenceEvent %v, Received %v\n", index, expected, actual)
	}
	if expected.Lease.Renewable != actual.Lease.Renewable {
		t.Errorf("%d - Expected PersistenceEvent %v, Received %v\n", index, expected, actual)
	}
}

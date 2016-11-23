package test

import (
	"net/http"
	"net"
)

func startDockerMock() {

	response := []byte(`{
  "status": "start",
  "id": "ede54ee1afda366ab42f824e8a5ffd195155d853ceaec74a927f249ea270c743",
  "from": "alpine",
  "Type": "container",
  "Action": "start",
  "Actor": {
    "ID": "ede54ee1afda366ab42f824e8a5ffd195155d853ceaec74a927f249ea270c743",
    "Attributes": {
      "common_name": "makeamericagreatagain.cloudframe.wtf",
      "image": "alpine",
      "name": "my-container"
    }
  },
  "time": 1461943101,
  "timeNano": 1461943101607533796
}{
  "status": "stop",
  "id": "ede54ee1afda366ab42f824e8a5ffd195155d853ceaec74a927f249ea270c743",
  "from": "alpine",
  "Type": "container",
  "Action": "stop",
  "Actor": {
    "ID": "ede54ee1afda366ab42f824e8a5ffd195155d853ceaec74a927f249ea270c743",
    "Attributes": {
      "common_name": "makeamericagreatagain.cloudframe.wtf",
      "image": "alpine",
      "name": "my-container"
    }
  },
  "time": 1461943108,
  "timeNano": 1461943108607533796
}`)

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})

	l, err := net.Listen("unix","//var/run/dockermock.sock")
	if err != nil {
		panic(err.Error())
	}
	http.Serve(l,nil)
}

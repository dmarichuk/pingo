package job

import (
	"log"
	"net/http"
	"time"
)

type EndpointExecutor struct {
	Client *http.Client
}

func NewEndpointExecutor() *EndpointExecutor {
	return &EndpointExecutor{
		Client: &http.Client{
			Timeout: 10 * time.Second, // TODO make it configurable with reasonable default
		},
	}
}

func (e *EndpointExecutor) Exec(j *Job) bool {

	resp, err := e.Client.Get(j.Endpoint)
	if err != nil {
		log.Printf("Endpoint ping to %s is failed: ", err)
		return false
	}
	defer resp.Body.Close()
	return true
}


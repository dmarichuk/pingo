package job

import (
	"fmt"
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

func (e *EndpointExecutor) Exec(j *Job) (bool, string) {
	var msg string

	resp, err := e.Client.Get(j.Endpoint)
	if err != nil {
		msg = fmt.Sprintf("Endpoint health check for job %s with %s is failed: %s", j.Name, j.Endpoint, err.Error())
		return false, msg
	}
	defer resp.Body.Close()

	ok := resp.StatusCode == http.StatusOK
	if !ok {
		msg = fmt.Sprintf("Endpoint health check for job %s with %s is failed. Status code: %d", j.Name, j.Endpoint, resp.StatusCode)
	}
	return ok, msg
}

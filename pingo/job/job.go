package job

import (
	"fmt"
	"log"
	"time"
)

const (
	// job types
	SERVICE_PING = "service-ping"

	// task types
	TELEGRAM_ALERT = "telegram-alert"

	// job status
	SUCCESS = "SUCCESS"
	FAILED  = "FAILED"
)

type Task interface {
	Run()
}

type Executor interface {
	Exec(*Job) bool
}

type Job struct {
	Name      string
	Type      string
	Endpoint  string
	Interval  time.Duration
	OnFailure []Task
	Status    string
	TS        time.Time
	PerfTime  time.Duration
}

func (j *Job) RunJob() {
	e, err := j.GetExecutor()
	if err != nil {
		log.Fatalln(err)
	}
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("Running job %s", j.Name)
		j.TS = time.Now()
		if ok := e.Exec(j); !ok {
			log.Printf("Job %s failed", j.Name)
			j.Status = FAILED
			go j.RunOnFailure()
		} else {
			log.Printf("Job %s succeed", j.Name)
			j.Status = SUCCESS
		}
		j.PerfTime = time.Now().Sub(j.TS)
	}
}

func (j *Job) GetExecutor() (Executor, error) {
	switch j.Type {
	case SERVICE_PING:
		return NewEndpointExecutor(), nil
	default:
		return nil, fmt.Errorf("Executor not implemented! %s", j.Type)
	}
}

func (j *Job) RunOnFailure() {
	for _, task := range j.OnFailure {
		task.Run() // TODO now it runs sequentially, think about runing tasks in coruotines
	}
}

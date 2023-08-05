package job

import (
	"fmt"
	"log"
	"time"
)

const (
	// job types
	ENDPOING_HEALTH = "endpoint-health"
	RAM_USAGE       = "ram-usage"
	DISK_USAGE      = "disk-usage"

	// task types
	TELEGRAM_ALERT = "telegram-alert"
	EMAIL_ALERT    = "email-alert"

	// task classes
	ON_FAILURE  = "on_failure"
	ON_RECOVERY = "on_recovery"

	// job status
	SUCCESS = "SUCCESS"
	FAILED  = "FAILED"
)

type Alert interface {
	Send([]byte)
	GenerateMessage(*Job) []byte
}

type Executor interface {
	Exec(*Job) bool
}

type Job struct {
	Name       string
	Type       string
	Interval   time.Duration
	OnFailure  []Alert
	OnRecovery []Alert
	Status     string
	TS         time.Time
	PerfTime   time.Duration

	// Endpoint health fields
	Endpoint string

	// RAM usage fields
	RamThreshold float64

	// Disk usage fields
	DiskThreshold float64
	DiskPath      string
}

func (j *Job) Run() {
	log.Printf("[INFO] Launching %s job", j.Name)
	e, err := j.GetExecutor()
	if err != nil {
		log.Fatalln(err)
	}
	ticker := time.NewTicker(j.Interval)
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("[INFO] Running job %s", j.Name)
		j.TS = time.Now()
		if ok := e.Exec(j); !ok {
			log.Printf("[INFO] Job %s failed", j.Name)
			if j.Status != FAILED { // Launch on Failure task only if Job was not failed previously
				go j.SendAlerts(ON_FAILURE)
			}
			j.Status = FAILED
		} else {
			log.Printf("[INFO] Job %s succeed", j.Name)
			if j.Status == FAILED { // Launch on Recovery task only if Job was failed previously
				go j.SendAlerts(ON_RECOVERY)
			}
			j.Status = SUCCESS
		}
		j.PerfTime = time.Now().Sub(j.TS)
		log.Printf("[INFO] Job %s performance time %d", j.PerfTime)
	}
}

func (j *Job) GetExecutor() (Executor, error) {
	switch j.Type {
	case ENDPOING_HEALTH:
		return NewEndpointExecutor(), nil
	case RAM_USAGE:
		return NewMemoryUsageExecutor(j.RamThreshold), nil
	case DISK_USAGE:
		return NewDiskUsageExecutor(j.DiskPath, j.DiskThreshold), nil
	default:
		return nil, fmt.Errorf("Executor not implemented! %s", j.Type)
	}
}

func (j *Job) SendAlerts(class string) {
	var alerts *[]Alert
	switch class {
	case ON_FAILURE:
		alerts = &j.OnFailure
	case ON_RECOVERY:
		alerts = &j.OnRecovery
	default:
		log.Fatalln("[ERR] Unknown Task class: ", class)
	}

	for _, alert := range *alerts {
		message := alert.GenerateMessage(j)
		alert.Send(message) // TODO now it runs sequentially, think about runing tasks in goroutines
	}
}

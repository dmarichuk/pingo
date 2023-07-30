package job

import (
	"fmt"
	"log"
	"time"
)

const (
	// job types
	SERVICE_PING = "service-ping"
    RAM_USAGE = "ram-usage"
    DISK_USAGE = "disk-usage"

	// task types
	TELEGRAM_ALERT = "telegram-alert"
    
    // task classes
    ON_FAILURE = "on_failure"
    ON_RECOVERY = "on_recovery"

	// job status
	SUCCESS = "SUCCESS"
	FAILED  = "FAILED"
)

type Task interface {
	Launch()
}

type Executor interface {
	Exec(*Job) bool
}

type Job struct {
	Name      string
	Type      string
	Interval  time.Duration
	OnFailure []Task
    OnRecovery []Task
	Status    string
	TS        time.Time
	PerfTime  time.Duration

    // Service ping fields
	Endpoint  string

    // RAM usage fields
    RamThreshold float64

    // Disk usage fields
    DiskThreshold float64
    DiskPath string
}

func (j *Job) Run() {
    log.Printf("Launching %s job", j.Name)
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
            if j.Status != FAILED { // Launch on Failure task only if Job was not failed previously
                go j.LaunchTasks(ON_FAILURE)
            }
			j.Status = FAILED
		} else {
			log.Printf("Job %s succeed", j.Name)
            if j.Status == FAILED { // Launch on Recovery task only if Job was failed previously
                go j.LaunchTasks(ON_RECOVERY)
            }
			j.Status = SUCCESS
		}
		j.PerfTime = time.Now().Sub(j.TS)
	}
}

func (j *Job) GetExecutor() (Executor, error) {
	switch j.Type {
	case SERVICE_PING:
		return NewEndpointExecutor(), nil
    case RAM_USAGE:
        return NewMemoryUsageExecutor(j.RamThreshold), nil
    case DISK_USAGE:
        return NewDiskUsageExecutor(j.DiskPath, j.DiskThreshold), nil
	default:
		return nil, fmt.Errorf("Executor not implemented! %s", j.Type)
	}
}

func (j *Job) LaunchTasks(class string) {
    var tasks *[]Task
    switch class {
    case ON_FAILURE:
        tasks = &j.OnFailure
    case ON_RECOVERY:
        tasks = &j.OnRecovery
    default:
        log.Fatalln("Unknown Task class: ", class)
    }

	for _, task := range *tasks {
		task.Launch() // TODO now it runs sequentially, think about runing tasks in goroutines
	}
}


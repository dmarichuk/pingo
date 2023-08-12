package job

import (
	"database/sql"
	"fmt"
	"log"
	db "pingo/database"
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
	Send(string)
}

type Executor interface {
	Exec(*Job) (bool, string)
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
	Message    string

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
		if ok, msg := e.Exec(j); !ok {
			log.Printf("[INFO] Job %s failed. Message: %s", j.Name, msg)
			if j.Status != FAILED { // Launch on Failure task only if Job was not failed previously
				go j.SendAlerts(ON_FAILURE, msg)
			}
			j.Status = FAILED
			j.Message = msg
		} else {
			log.Printf("[INFO] Job %s succeed. Message: %s", j.Name, msg)
			if j.Status == FAILED { // Launch on Recovery task only if Job was failed previously
				go j.SendAlerts(ON_RECOVERY, fmt.Sprintf("Job %s recovered!", j.Name))
			}
			j.Status = SUCCESS
			j.Message = msg
		}
		j.PerfTime = time.Now().Sub(j.TS)
		log.Printf("[INFO] Job %s performance time %d", j.Name, j.PerfTime)
		err = j.DumpLogToDB(db.DB)
		if err != nil {
			log.Println(err)
		}
	}
}

func (j *Job) DumpToDB(db *sql.DB) error {
	sqlStmt := `
	INSERT INTO jobs (
		name, type, endpoint, ram_threshold, disk_threshold, disk_path
	) VALUES (
		?, ?, ?, ?, ?, ?
	)
	`
	_, err := db.Exec(sqlStmt, j.Name, j.Type, j.Endpoint, j.RamThreshold, j.DiskThreshold, j.DiskPath)
	return err
}

func (j *Job) DumpLogToDB(db *sql.DB) error {
	sqlStmt := `
	INSERT INTO job_logs (
		job, ts, status, perf_time, message
	) VALUES (
		?, ?, ?, ?, ? 
	)
	`
	_, err := db.Exec(sqlStmt, j.Name, j.TS.Format(time.DateTime), j.Status, j.PerfTime.Seconds(), j.Message)
	return err

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

func (j *Job) SendAlerts(class, msg string) {
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
		alert.Send(msg) // TODO now it runs sequentially, think about runing tasks in goroutines
	}
}

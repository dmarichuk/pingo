package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const DB_NAME = "./sqlite3.db"

var (
	DB    *sql.DB
	dbErr error
)

func init() {
	os.Remove(DB_NAME)
	DB, dbErr = sql.Open("sqlite3", DB_NAME)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	dbErr = CreateJobTable(DB)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	dbErr = CreateJobLog(DB)
	if dbErr != nil {
		log.Fatal(dbErr)
	}

}

type DBJob struct {
	Id            int     `json:"id,omitempty"`
	Name          string  `json:"name,omitempty"`
	Type          string  `json:"type,omitempty"`
	Endpoint      string  `json:"endpoint,omitempty"`
	RamThreshold  float64 `json:"ram_threshold,omitempty"`
	DiskThreshold float64 `json:"disk_threshold,omitempty"`
	DiskPath      string  `json:"disk_path,omitempty"`
}

type DBJobLog struct {
	Id       int    `json:"id,omitempty"`
	Job      string `json:"job,omitempty"`
	TS       string `json:"ts,omitempty"`
	Status   string `json:"status,omitempty"`
	PerfTime string `json:"perf_time,omitempty"`
	Message  string `json:"message,omitempty"`
}

type DBPieJob struct {
	Status string `json:"status,omitempty"`
	Count  string `json:"count,omitempty"`
}

func CreateJobTable(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY,
		name TEXT,
		type TEXT,
		endpoint TEXT,
		ram_threshold REAL,
		disk_threshold REAL,
		disk_path TEXT
	)
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `
	CREATE UNIQUE INDEX IF NOT EXISTS name_unique_idx ON jobs (name)
	`
	_, err = db.Exec(sqlStmt)
	return err
}

func CreateJobLog(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS job_logs (
		id INTEGER PRIMARY KEY,
		job TEXT,
		ts TEXT,
		status TEXT,
		perf_time REAL,
		message TEXT,
		FOREIGN KEY(job) REFERENCES jobs(name) ON DELETE CASCADE 
	)
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `
	CREATE INDEX IF NOT EXISTS ts_idx ON job_logs (ts DESC)
	`

	_, err = db.Exec(sqlStmt)
	return err
}

func SelectLatestJobLogs(db *sql.DB, name string, limit int) ([]DBJobLog, error) {
	result := make([]DBJobLog, 0, limit)
	sqlStmt := `
	SELECT t1.id, ts, status, perf_time, message
	FROM job_logs AS t1
	INNER JOIN jobs AS t2 ON t1.job = t2.name 
	WHERE t2.name=?
	ORDER BY ts DESC 
	LIMIT ? 
	`
	rows, err := db.Query(sqlStmt, name, limit)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {

		var j DBJobLog
		err = rows.Scan(&j.Id, &j.TS, &j.Status, &j.PerfTime, &j.Message)
		if err != nil {
			return result, err
		}
		result = append(result, j)
	}

	return result, nil
}

func SelectJobsForPieChart(db *sql.DB, name string) ([]DBPieJob, error) {
	result := make([]DBPieJob, 0, 2)
	sqlStmt := `
	SELECT	status, COUNT(status)
	FROM job_logs AS t1
	INNER JOIN jobs AS t2 ON t1.job = t2.name 
	WHERE t2.name=?
	GROUP BY status
	ORDER BY status
	`
	rows, err := db.Query(sqlStmt, name)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {

		var j DBPieJob
		err = rows.Scan(&j.Status, &j.Count)
		if err != nil {
			return result, err

		}
		result = append(result, j)
	}

	return result, nil
}

func SelectJobsInfo(db *sql.DB) ([]DBJob, error) {
	var result []DBJob
	sqlStmt := `
	SELECT id, name, type, endpoint, ram_threshold, disk_threshold, disk_path
	FROM jobs
	ORDER BY id
	`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {

		var j DBJob
		err = rows.Scan(&j.Id, &j.Name, &j.Type, &j.Endpoint, &j.RamThreshold, &j.DiskThreshold, &j.DiskPath)
		if err != nil {
			return result, err

		}
		result = append(result, j)
	}
	if err != nil {
		return result, err
	}

	return result, nil
}

package database

import (
	"database/sql"
	"log"
	"os"
	"pingo/job"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const DB_NAME = "/var/lib/pingo.db"

var (
	DB *sql.DB
)

func init() {
	os.Remove(DB_NAME)
	DB, err := sql.Open("sqlite3", DB_NAME)
	if err != nil {
		log.Fatal(err)
	}
	err = CreateJobTable(DB)
	if err != nil {
		log.Fatal(err)
	}
}

type readJob struct {
	Name          string  `json:"name,omitempty"`
	Type          string  `json:"type,omitempty"`
	TS            string  `json:"ts,omitempty"`
	Status        string  `json:"status,omitempty"`
	PerfTime      float64 `json:"perf_time,omitempty"`
	Endpoint      string  `json:"endpoint,omitempty"`
	RamThreshold  float64 `json:"ram_threshold,omitempty"`
	DiskThreshold float64 `json:"disk_threshold,omitempty"`
	DiskPath      string  `json:"disk_path,omitempty"`
}

func CreateJobTable(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY,
		name TEXT,
		type TEXT,
		ts TEXT,
		status TEXT,
		perf_time REAL,
		endpoint TEXT,
		ram_threshold REAL,
		disk_threshold REAL,
		disk_path TEXT
	)
	`
	_, err := db.Exec(sqlStmt)
	return err
}

func InsertJobRow(db *sql.DB, j *job.Job) error {
	sqlStmt := `
	INSERT INTO jobs (
		name, type, ts, status, perf_time, endpoint, ram_threshold, disk_threshold, disk_path
	) VALUES (
		?, ?, ?, ?, ?, ?, ?, ?, ?
	)
	`
	_, err := db.Exec(sqlStmt, j.Name, j.Type, j.TS.Format(time.DateTime), j.Status, j.PerfTime.Seconds(), j.RamThreshold, j.DiskThreshold, j.DiskPath)
	return err
}

func SelectJobsWithLimit(db *sql.DB, name string, limit int) ([]readJob, error) {
	result := make([]readJob, limit)
	sqlStmt := `
	SELECT	name, type, ts, status, perf_time, endpoint, ram_threshold, disk_threshold, disk_path
	FROM jobs
	WHERE name=?
	ORDER BY ts 
	DESC LIMIT ? 
	`
	rows, err := db.Query(sqlStmt, limit, name)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {

		var j readJob
		err = rows.Scan(&j.Name, &j.Type, &j.TS, &j.Status, &j.Endpoint, &j.PerfTime, &j.RamThreshold, &j.RamThreshold, &j.DiskThreshold, &j.DiskPath)

		if err != nil {
			return result, err
		}
		result = append(result, j)
	}

	return result, nil
}

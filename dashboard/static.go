package dashboard

import (
	"embed"
	db "pingo/database"
	"pingo/job"
)

var (
	//go:embed static
	static embed.FS
)

func mapTarget(j db.DBJob) string {
	switch j.Type {
	case job.ENDPOING_HEALTH:
		return j.Endpoint
	case job.DISK_USAGE:
		return j.DiskPath
	default:
		return " --- "
	}
}

func mapType(j db.DBJob) string {
	switch j.Type {
	case job.ENDPOING_HEALTH:
		return "Endpoint health"
	case job.RAM_USAGE:
		return "Memory usage"
	case job.DISK_USAGE:
		return "Disk usage"
	default:
		return j.Type
	}
}

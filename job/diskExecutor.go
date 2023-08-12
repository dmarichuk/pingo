package job

import (
	"fmt"
	"syscall"
)

type DiskUsageExecutor struct {
	info      *syscall.Statfs_t
	Path      string
	Threshold float64
}

func NewDiskUsageExecutor(path string, threshold float64) *DiskUsageExecutor {
	return &DiskUsageExecutor{
		info:      &syscall.Statfs_t{},
		Threshold: threshold,
		Path:      path,
	}
}

func (e *DiskUsageExecutor) Exec(j *Job) (bool, string) {
	var msg string

	if err := syscall.Statfs(e.Path, e.info); err != nil {
		msg = fmt.Sprintf("Couldn't access stat fs: %s", err.Error())
		return false, msg
	}

	freeDisk := float64(e.info.Bfree * uint64(e.info.Bsize))
	totalDisk := float64(e.info.Blocks * uint64(e.info.Bsize))
	ratio := 1.0 - freeDisk/totalDisk

	ok := ratio < e.Threshold
	if !ok {
		msg = fmt.Sprintf("Disk usage exceeded threshold %.2f for %s\nFree Disk: %.2f\nTotal Disk: %.2f", e.Threshold, j.Name, freeDisk, totalDisk)
	}
	return ok, msg
}

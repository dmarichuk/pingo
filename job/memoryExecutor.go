package job

import (
	"fmt"
	"syscall"
)

type MemoryUsageExecutor struct {
	info      *syscall.Sysinfo_t
	Threshold float64
}

func NewMemoryUsageExecutor(threshold float64) *MemoryUsageExecutor {
	return &MemoryUsageExecutor{
		info:      &syscall.Sysinfo_t{},
		Threshold: threshold,
	}
}

func (e *MemoryUsageExecutor) Exec(j *Job) (bool, string) {
	var msg string

	if err := syscall.Sysinfo(e.info); err != nil {
		msg = fmt.Sprintf("Couldn't access system info: %s", err.Error())
		return false, msg
	}

	freeMemory := float64(e.info.Freeram)
	totalMemory := float64(e.info.Totalram)

	ratio := 1.0 - freeMemory/totalMemory
	ok := ratio < e.Threshold
	if !ok {
		msg = fmt.Sprintf("RAM usage exeeded threshold %.2f for %s\nFree Memory: %.2f bytes\nTotal Memory: %.2f bytes", e.Threshold, j.Name, freeMemory, totalMemory)
	}
	return ok, msg
}

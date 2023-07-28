package job

import (
	"log"
    "syscall"
)


type MemoryUsageExecutor struct {
    info *syscall.Sysinfo_t
    Threshold float64
}

func NewMemoryUsageExecutor(threshold float64) *MemoryUsageExecutor {
    return &MemoryUsageExecutor{
        info: &syscall.Sysinfo_t{},
        Threshold: threshold,

    }
}

func (e *MemoryUsageExecutor) Exec(j *Job) bool {
    if err := syscall.Sysinfo(e.info); err != nil {
        log.Println("Couldn't access system info: ", err)
        return false	
    }
    ratio := 1.0 - float64(e.info.Freeram) / float64(e.info.Totalram)
	return ratio > e.Threshold
}


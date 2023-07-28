package job

import (
	"log"
    "syscall"
)


type DiskUsageExecutor struct {
    info *syscall.Statfs_t
    Path string
    Threshold float64
}

func NewDiskUsageExecutor(path string, threshold float64) *DiskUsageExecutor {
    return &DiskUsageExecutor{
        info: &syscall.Statfs_t{},
        Threshold: threshold,
        Path: path,
    }
}

func (e *DiskUsageExecutor) Exec(j *Job) bool {
    if err := syscall.Statfs(e.Path, e.info); err != nil {
        log.Println("Couldn't access stat fs: ", err)
        return false	
    }
    ratio := 1.0 - float64(e.info.Bfree * uint64(e.info.Bsize)) / float64(e.info.Blocks * uint64(e.info.Bsize))
	return ratio > e.Threshold
}


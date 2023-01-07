package cubeslurmtask

import (
	"time"
)

type Task struct {
	Encoding         string
	Threshold        int
	Timeout          time.Duration
	MaxCubes         int
	MinRefutedLeaves int
}

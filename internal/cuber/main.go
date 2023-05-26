package cuber

import "errors"

var (
	ErrCubesetViolatedConstraints = errors.New("cubeset violated constraints")
)

type ThresholdType string

const (
	CutoffVars  = "cutoff_vars"
	CutoffDepth = "cutoff_depth"
)

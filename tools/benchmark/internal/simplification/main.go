package simplification

import "time"

type Simplification struct {
	FreeVariables int
	Simplifier    string
	ProcessTime   time.Duration
	Eliminaton    int
	Conflicts     int
	Clauses       int
	InstanceName  string
}

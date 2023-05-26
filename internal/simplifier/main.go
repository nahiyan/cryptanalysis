package simplifier

import "time"

const (
	Cadical  = "cadical"
	Satelite = "satelite"
)

type CadicalOutput struct {
	FreeVariables int
	Clauses       int
	Eliminations  int
	ProcessTime   time.Duration
}

type SateliteOutput struct {
	FreeVariables int
	Clauses       int
	ProcessTime   time.Duration
}

type Simplifier string

type Result struct {
	ProcessTime          time.Duration
	NumVars              int
	NumClauses           int
	NumEliminatedVars    int
	NumEliminatedClauses int
}

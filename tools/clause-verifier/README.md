# Clause verifier

This tool checks reason and blocking clauses from a solver log against a SAT
solution assignment. It prints every unsatisfied clause, clause counts, and
`Solved` when the solver log contains a `c exit` line.

```sh
go run . <solution-log> <solver-log>
```

The solution log must contain assignment lines beginning with `v`. The solver
log can contain clauses in these forms:

```text
Reason clause: 1 -2 3
Blocking clause: -1 4
```

# Standalone tools

Each directory here is an independent research utility with a narrow input and
output contract. Tools must not depend on the repository-wide Go framework.

- `conditions-2bit/` derives and propagates two-bit conditions.
- `clause-verifier/` checks clauses against assignments.
- `proof-analyzer/` analyzes solver proof output.
- `rules-generator/` generates propagation rules.

# Configuration

- Solvers: CRYPTOMINISAT, XNFSAT
- Steps: 16 - 39
- Dobbertin's Attack: 1
- Max Time: 5000
- Hashes: ffffffffffffffffffffffffffffffff, 00000000000000000000000000000000
- Adders: Counter chain, Dot matrix
- XOR: Enabled for all

Each instance ran in a Slurm job with 1 CPU core and 300 MBs of memory

# Notes

XNFSAT performed poorly, failing to solve problems past 27 steps. Surprisingly, both the solvers found the Dobbertin's attack unsatisfiable in a lot of cases, even though they were solved by other solvers in non-XOR form.

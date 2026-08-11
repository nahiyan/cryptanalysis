# Configuration

- Solvers: CRYPTOMINISAT, KISSAT, CADICAL, GLUCOSE, MAPLESAT
- Steps: 32 - 39
- Dobbertin's Attack: 1
- Max Time: 5000
- Hashes: ffffffffffffffffffffffffffffffff, 00000000000000000000000000000000
- Adders: Counter chain, Dot matrix
- XOR: Enabled for none except CRYPTOMINISAT

Each instance ran in a Slurm job with 1 CPU core and 300 MBs of memory

# Notes

CryptoMiniSAT performed much better with XOR encodings, often performing better than all the other solvers at some points.

# Configuration

- Solvers: CRYPTOMINISAT, KISSAT, CADICAL, GLUCOSE, MAPLESAT
- Steps: 32 - 39
- Dobbertin's Attack: 1
- Max Time: 5000
- Hashes: ffffffffffffffffffffffffffffffff, 00000000000000000000000000000000
- Adders: Counter chain, Dot matrix
- XOR: Disabled for all

Each instance ran in a Slurm job with 1 CPU core and 300 MBs of memory

# Notes

Dobbertin's attack made a significant improvement - allowing us to solve 32-39 step versions of MD4 in 5000 seconds. Just like before, Kissat and CaDiCaL performed the best, solving the most instances. Several solvers, including MapleSAT, could solve the 39-step MD4.

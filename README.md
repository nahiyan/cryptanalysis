# MD4 Inversion

This repository houses various SAT DIMACS CNF encoders, tools for benchmarks using numerous SAT solvers, and techniques such as the [Dobbertin's attack](https://link.springer.com/content/pdf/10.1007/3-540-69710-1_19.pdf) and [Cube and Conquer using a Lookahead Solver](https://www.cs.utexas.edu/~marijn/publications/cube.pdf).

# Progress

## Dobbertin's Attack

16-39 steps of MD4 on hashes with all set and unset message bits have been inverted for the targets hashes with all set and unset bits. SAT solvers, such as Kissat, CaDiCaL, MapleSAT, Glucose, and CryptoMiniSAT, XNFSAT with/without XOR clauses have been tried out and compared visually using a cactus plot.

## Cube and Conquer

On top of the Dobbertin's attack, cube and conquer is being experimented for encodings generated using both TRANSALG and Saeed's encoder with variation of adders, CNF simplification techniques, cutoff variables for cubing, random sampling for estimation, etc. March is used for generating the cubes, while CaDiCaL is used for the simplification.

# Tools

## Encodings Generator

### Saeed Nejati

To use the encodings generator written by Saeed Nejati, compile it first by running `make` in the `encoders/saeed/crypto` directory, which should produce an executable named `main`. The documentation for using the encoder can be found through the `--help` or `-h` flag.

So far, for the research, we'll be looking at the following variations of encodings:

- With and without XOR
- Target hash with all set and unset bits
- Adders such as counter chain, dot matrix, and espresso
- n-step versions of MD4 with n ranging from 16 to 48
- With and without Dobbertin's attack
- Relaxation of one Dobbertin's constraint out of the 12 by $W - 32$ bits, where $W$ is the word size that is always 32

## Benchmark Tool

The benchmark tool is the heart of the project for experimenting with MD4 inversion. It's in `tools/benchmark`, written in Go, and has the following features:

- Drive the encoder for generating the encodings/instances
- Drive lookahead SAT solver[s] for generating cubes
- Drive CDCL SAT solver[s] for solving the instances
- Generate Slurm jobs for each instance
- Control the spawning of the instances, limit the max. concurrent instances, and keep track of the progress 
- Maintain an aggregated log from all the instances in a CSV file

To build the tool, just run `go build`, assuming that you have Go installed in your system already. As with any Go source code, you can run the code using `go run main.go`. For further documentation, simply call with the `--help` flag.

## Solution Analyzer

Once an instance is solved, a solutions output is generated, which holds all the variables that must be true/false to satisfy the constraints provided through the CNF. However, since there may be thousands of variables, it can be hard to read. The solutions analyzer, written in Rust, takes the solution file and presents the variables within the given ranges as bytes in binary (or other preferred bases). Moreover, it can normalize the solution to a specific format - currently MapleSAT's format is supported as the normalization target. The normalization feature is being used by the benchmark tool for every solution before it gets fed into the validator.

### Build and Run

To build and run the analyzer, ensure that you have [Cargo](https://doc.rust-lang.org/cargo/) installed and invoke the `cargo run` command in the `scripts/solution_analyzer` directory, or build the binary as a release build using `cargo build --release` and run it using `./target/release/solution_analyzer`.

The rest of the documentation can be accessed through the `--help` flag.

#### Normalization

To normalize a solution, which resides in /tmp/solution.sol, just run `solution_analyzer /tmp/solution.sol normalize`, and you should get a dump of the normalized version. You can pipe it to a file (existent or not) to save the normalization version to a file, like this: `solution_analyzer solution.sol > solution-normalized.sol`.

#### Summarization

So, for example, if you have a solutions file in `scripts/solution_analyzer/md4.sol`, you may want to run `solution_analyzer scripts/solution_analyzer/md4.sol --variables 1-512,641-768 summarize` to run the analyzer and output the values of the variables 1 to 512 and 641 to 768. You can provide as many variable ranges as you like, all provided in the single argument "--variables" in a comma-separated format.  It prints something like this (output truncated for obvious reasons):

```
Range: 1 - 512
00000000
00000001
10111110
00011101
11111111
...

Range: 641 - 768
11111111
11111111
11111111
11111111
11111111
...
```

The rest of the documentation for summarization can be accessed through `solution_analyzer summarize --help`, just like any other standard program.

## Cactus Plot

A small Python script that uses matplotlib for generating cactus plot is housed in `tools/cactus_plot` that is in the process of refactoring to be compatible with the recent major changes of the benchmark tool. It simply takes an aggregated log through stdin and plots the results of various SAT solvers in a cactus/survival plot, like so: `python generate.py < benchmark.log`, resulting in a PNG image of the plot.

## Cube & Conquer Cutoff Threshold Finder

Generates cubesets using various cutoff thresholds (number of max. variables) within constraints. For each cubeset, it takes a random sample set of size N to estimate the time to solve the sub-problems generated from the cubes.

# Credits

- `encoders/saeed` is a modified and trimmed version of https://github.com/saeednj/SAT-encoding
- `tools/cnc_threshold_finder` is a modified version of the CNC threshold calculator in https://github.com/olegzaikin/MD4-CnC

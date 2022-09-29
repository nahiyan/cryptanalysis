# MD4 Inversion

This repository houses various SAT DIMACS CNF encoders, tools for benchmarks using numerous, possibly crypto-aware, SAT solvers, etc.

For now, it focuses mainly on Saeed Nejati and Oleg Zaikin's work.

# Scripts

## Encodings Generator

### Saeed Nejati

To use the encodings generator written by Saeed Nejati, invoke the python script, `scripts/encodings-generator-saeed.py`, like in the following command:

```bash
python scripts/encodings-generator-saeed.py
```

This will generate encodings (DIMACS CNF) of MD4 in the `encodings/saeed` directory with the following variations:

- With and without XOR
- Target hash with all set bits and all unset bits
- Counter chain and doot matrix adders
- n-step versions of MD4 with n ranging from 16 to 48

## Solutions Analyzer

Once an instance is solved, a solutions output is generated, which holds all the variables that must be true/false to satisfy the constraints. However, since there may be thousands of variables, which can be hard to read. To make it easier to read, the solutions analyzer, written in Rust, analyses the solutions file and presents the variables within given ranges as bytes.

### Build and Run

To build and run the analyzer, ensure that you have [Cargo](https://doc.rust-lang.org/cargo/) installed and invoke the `cargo run` command in the `scripts/solution_analyzer` directory, or build the binary as a release build using `cargo build --release` and run it using `./target/release/solution_analyzer`.

The analyzer takes the following arguments:
- A solution file
- Multiple ranges in the form `{start}-{end}`, where `{start}` and `{end}` are starting and ending numbers for the variables

You can put as many ranges as you like. Just ensure that they do not overlap.

So, for example, if you have a solutions file in `scripts/solution_analyzer/md4.sol`, you may want to run `cargo run -- scripts/solution_analyzer/md4.sol 1-512 641-768` to run the analyzer and output the values of the variables 1 to 512 and 641 to 768. It prints something like this (output truncated for obvious reasons):

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

# Credits

- `encoders/saeed` is a modified and trimmed version of https://github.com/saeednj/SAT-encoding

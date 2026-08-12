<!-- TODO: Instruct people to use the root Makefile to build all the dependencies along with the cryptanalysis tool -->

# MD4, MD5, and SHA-256 Cryptanalysis

This repository houses a cryptanalysis tool aimed at (pre-image and collision) attacks on hash functions, namely MD4, MD5, and SHA-256, using SAT solvers. Techniques such as the [Dobbertin's attack](https://link.springer.com/content/pdf/10.1007/3-540-69710-1_19.pdf) and [Cube and Conquer using a Lookahead Solver](https://www.cs.utexas.edu/~marijn/publications/cube.pdf) are used to exploit the weakness of hash functions and increase feasibility of the attacks.

The repository is organized around cryptanalytic work rather than around a
single implementation language:

- `attacks/` contains attack-specific code and material.
- `encoders/` contains SAT encoders and encoder inputs.
- `tools/` contains small, standalone research utilities.
- `runners/` contains the Go and Rust programs that execute or coordinate
  experiments.
- `experiments/` is reserved for reproducible experiment definitions.
- `results/` contains recorded benchmark and experiment results.

The existing Go command is under `runners/go-driver` while its reusable parts
are gradually extracted into standalone tools. The experimental Rust runner is
under `runners/rust-driver`.

The Rust encoder implementations share only the low-level SAT formula and
embedded Espresso tables in `encoders/nejati-common`. Attack construction and
command-line behavior remain in their standalone encoder directories.

# Prerequisities

To use the cryptanalysis tool, the following executables may be required (depending on your use-case):

- [Kissat](https://github.com/arminbiere/kissat) as `kissat`
- [CaDiCaL](https://github.com/arminbiere/cadical) as `cadical`
- [March](https://github.com/BrianLi009/PhysicsCheck/tree/master/gen_cubes/march_cu) as `march_cu_pc` (suffixed with "_pc" for "PhysicsCheck")
- [Transalg](https://gitlab.com/transalg/transalg) as `transalg`
- [CryptoMiniSAT](https://github.com/msoos/cryptominisat) as `cryptominisat`
- [MapleSAT](https://github.com/nahiyan/maplesat) as `maplesat`
- [Glucose](https://github.com/mi-ki/glucose-syrup) as `glucose`
- The modified Nejati preimage and collision encoders under `encoders/`
- [xnfSAT](https://github.com/Vtec234/xnfsat) as `xnfsat`

Other 3rd party dependencies may be required on use-case, such as `lstech_maple`, `kissat_cf`, `yalsat`, `palsat, etc.

# Building

> You need Go 1.18 or newer and a current stable Rust toolchain to build this
> repository.

Run `make` in the repository root to build the Go driver and the modified
Nejati encoders. To build only the Go driver, run:

```sh
cd runners/go-driver
go build -o ../../cryptanalysis
```

Available commands and arguments can be found via the `--help` flag. For example, `cryptanalysis run --help` will show the instructions for the "run" command.

# Running

All queries to the cryptanalysis tool are provided through a schema file written in [TOML](https://toml.io/en/). A pipeline can be declared in the schema with each pipe defining the operation and its configuration.

For example, the following schema file instructs the cryptanalysis tool to encode and solve a 43-step MD4 instance:

```toml
[[pipeline]]
type = "encode"
[pipeline.EncodeParams]
encoder = "transalg"
function = "md4"
xor = [0]
dobbertin = [1]
dobbertinBits = [32]
adders = ["espresso"]
hashes = ["ffffffffffffffffffffffffffffffff"]
steps = [43]

[[pipeline]]
type = "solve"
[pipeline.SolveParams]
solvers = ["kissat"]
timeout = 10000
workers = 16
```

The pipeline can be run by executing `cryptanalysis run <schema-file>`, where `<schema-fle>` is the placeholder to the file path, e.g. schema.toml. Note that the pipeline propagates top-down sequentially.

A much more complex pipeline for encode > simplify > cube > select (cubes) > solve can be defined like this:

```toml
[[pipeline]]
type = "encode"
[pipeline.EncodeParams]
encoder = "transalg"
function = "md4"
xor = [0]
dobbertin = [1]
dobbertinBits = [32]
adders = ["espresso"]
hashes = ["ffffffffffffffffffffffffffffffff"]
steps = [43]

[[pipeline]]
type = "simplify"
[pipeline.SimplifyParams]
name = "cadical"
conflicts = [100]
workers = 1

[[pipeline]]
type = "cube"
[pipeline.CubeParams]
# thresholds = [130]
initialThreshold = 10
stepChange = -10
maxCubes = 10000000
minRefutedLeaves = 500
workers = 12
timeout = 28800

[[pipeline]]
type = "cube_select"
[pipeline.CubeSelectParams]
type = "random"
quantity = 1000
seed = 1

[[pipeline]]
type = "solve"
[pipeline.SolveParams]
solvers = ["kissat"]
timeout = 10000
workers = 16
```

The operations of the above pipeline are as follows:
- Encode a 43-step MD4 using Transalg with all-one target hash and including the Dobbertin's constraints
- Simplify the instance with CaDiCaL till 100 conflicts
- Cube till reaching cubesets of 10M cubes while only keeping cubesets of at least 500 refuted leaves
- Select 100 cubes from each cubeset in random order with a seed of 1 (you can exclude the quantity to select all the cubes)
- Solve the instances with Kissat (with a 10000s timeout) in 16 workers (16 processes of Kissat will be spawned at a time)

You can explore all the possible parameters and pipe types in
`runners/go-driver/internal/pipeline/main.go`.

# Configuration

The cryptanalysis tool's configuration is defined in a (optional) `config.toml` TOML file. The tool will conventionally look for the file relative to the current working directory.

This is how it may look like:

```toml
[Solver.Cadical]
LocalSearchRounds = 3

[Solver.Kissat]
LocalSearch = true
LocalSearchEffort = 10

[Solver.CryptoMiniSat]
LocalSearch = true
LocalSearchType = "walksat"

[Paths.Bin]
NejatiEncoder = "/tmp/SAT-encodings/crypto/main"
```

You can check out all the possible parameters in
`runners/go-driver/internal/config/main.go`.

# Encoders

The following SAT encoders are integrated into the cryptanalysis tool for generating the SAT encodings.

## Transalg

[Transalg](https://gitlab.com/transalg/transalg), a SAT encoder that takes the problem definition as a high-level C-like code to generate DIMACS CNF, has been utilized to encode attacks on MD4, MD5, and SHA-256.

## Nejati

Saeed Nejati wrote his [own encoders and verifiers](https://github.com/saeednj/SAT-encoding) for MD4, SHA-256, etc. However, this repository holds a modified and trimmed-down version of his project.

### Building

Run `make` in either `encoders/nejati-preimage` or
`encoders/nejati-collision`. Both encoders now build from Rust; their old C++
implementations remain available through `make legacy` for parity comparisons.
Manual invocation is normally unnecessary because the cryptanalysis tool calls
the configured encoder executable directly.

The following set of features is a subset of all that are available:

- XOR clauses
- Specification of the target hash
- Embedded Espresso adders in the Rust encoders
- Counter-chain and dot-matrix adders in the legacy C++ encoders
- Trimmed n-step version of the hash function
- Dobbertin's attack in MD4
- Relaxation of one Dobbertin's constraint out of the 12 by $W - 32$ bits, where $W$ is the word size that is always 32

> Important: The cryptanalysis tool recognizes (by default) the binary as `nejati_encoder` in the system's environment.

# Techniques

## Dobbertin's Attack

Exploiting the majority function in MD4, Dobbertin's constraints are encoded into the SAT problem to reduce the search space by containing values of 3 registers into 1 and making some of the pre-image words derivable with BCP before the CDCL phase. This makes it feasible to invert MD4 up to 43 steps.

## Cube and Conquer

Cube and conquer is a popular technique for generating assumption cubes that can be solved in parallel by CDCL solvers. The lookahead solver, March, is used for generating the cubes, while CaDiCaL is used for the simplification beforehand.

> Important: Please note that the version of March used is a modified version housed in the [PhysicsCheck repository](https://github.com/BrianLi009/PhysicsCheck/tree/b1212848392673eac93ba437017ef6979e2775f0/gen_cubes/march_cu). By default, March is assumed to be available as `march_cu_pc` ("pc" for PhysicsCheck) in the system's environment.

## Benchmark Tool

The Go driver under `runners/go-driver` can coordinate MD4 inversion
experiments. Its responsibilities include:

- Drive the encoder for generating the encodings/instances
- Drive lookahead SAT solver[s] for generating cubes
- Drive CDCL SAT solver[s] for solving the instances
- Generate Slurm jobs for each instance
- Control the spawning of the instances, limit the max. concurrent instances, and keep track of the progress 
- Maintain an aggregated log from all the instances in a CSV file

Build the driver with the root Makefile or run it directly from its module:

```sh
cd runners/go-driver
go run . --help
```

# Credits

- `encoders/nejati-preimage` and `encoders/nejati-collision` are modified and trimmed versions of https://github.com/saeednj/SAT-encoding
- Transalg code templates for MD4, MD5, and SHA-256 were based on that housed in https://gitlab.com/satencodings/satencodings/
- The threshold finding algorithm is a modified version of that found in https://github.com/olegzaikin/MD4-CnC

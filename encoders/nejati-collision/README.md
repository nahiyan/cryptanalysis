# Nejati SHA-256 collision encoder

This standalone Rust tool encodes SHA-256 collision attacks as DIMACS CNF. It
combines two round-reduced SHA-256 compression blocks with a supplied
differential path and the embedded propagation rules from `prop_rules.db`.

## Build

```sh
make
```

The resulting executable is `./nejati_collision_encoder`. The old C++ source is
retained temporarily for parity comparisons and can be built with `make legacy`
as `./nejati_collision_encoder_cpp`.

## Use

```sh
./nejati_collision_encoder \
  --rounds 38 \
  --diff_desc \
  --diff_const_file tables/38_mendel.txt \
  --adder_type espresso > collision-38.cnf
```

The encoder supports 16 to 64 SHA-256 rounds, native XOR clauses, fixed or free
starting chaining values, and the bundled one-bit differential-path format.
The propagation-rule database is embedded in the executable, so invocation is
independent of the current working directory.

Only the Espresso multi-operand adder is supported. The legacy
`--rand_input_diff` option never affected the C++ encoding and is rejected
instead of silently producing an unchanged instance.

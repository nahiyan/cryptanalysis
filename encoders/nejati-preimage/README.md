# Nejati preimage encoder

This directory contains a standalone Rust SAT encoder for preimage attacks on
round-reduced MD4, SHA-1, and SHA-256. It emits DIMACS CNF to standard output.

The Rust implementation replaces the original C++ executable. The C++ source
is retained temporarily as a reference while the encoder migration continues.
Its Boolean formula primitives and embedded Espresso tables live in the small
`../nejati-common` Rust crate so the collision encoder can reuse them.

## Build

Build the Rust encoder with:

```sh
make
```

The resulting executable is `./nejati_preimage_encoder`. To compile the old C++
implementation for comparison, run `make legacy`; its executable is
`./nejati_preimage_encoder_cpp`.

## Use

For example, generate a 20-round SHA-256 preimage instance with a random target:

```sh
./nejati_preimage_encoder \
  --function sha256 \
  --rounds 20 \
  --target random > sha256-20-preimage.cnf
```

Run `./nejati_preimage_encoder --help` for all options. The migrated encoder
supports:

- MD4 with 1 to 48 rounds
- SHA-1 with 16 to 80 rounds
- SHA-256 with 16 to 64 rounds
- native CNF or XOR clauses
- random or explicit target values
- a fixed prefix of a generated message
- MD4 Dobbertin constraints and relaxed-bit selection

Explicit targets are hexadecimal strings with the full compression-state size:
32 characters for MD4, 40 for SHA-1, and 64 for SHA-256.

The Rust implementation embeds the minimized Espresso truth-table encodings,
so running the encoder does not require the external `espresso` executable.
Only the Espresso adder encoding has been migrated. The legacy counter-chain
and dot-matrix adders are not available in the Rust executable.

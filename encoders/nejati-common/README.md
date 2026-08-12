# Nejati SAT primitives

This internal Rust crate contains the Boolean formula operations and embedded
Espresso adder encodings shared by the standalone Nejati encoders.

It is deliberately limited to SAT construction. Hash functions, attacks, file
formats, and command-line interfaces belong to the encoder tools that use it.

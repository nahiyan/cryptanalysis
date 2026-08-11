# Runners

Runners execute or coordinate experiments. Attack-specific logic should move
out of them as standalone tools are extracted.

- `go-driver/` contains the existing Go orchestration command and its module.
- `rust-driver/` contains the experimental Rust driver.

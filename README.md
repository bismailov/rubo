# Rubo Compiler

Rubo is a "Ruby-syntax frontend with a Go-orchestrated, Rust-accelerated backend."

## Project Structure

- `src/cmd/rubo`: The CLI entry point (Go).
- `src/parser`: The Go-based syntax analyzer.
- `src/runtime`: The Rust-based performance engine.
- `src/bridge`: FFI/C-header glue.

## Getting Started

### Prerequisites

- Go
- Rust
- Docker (for DevContainer)

### Building the Runtime

```bash
cd src/runtime
cargo build
```

### Running the CLI

```bash
cd src/cmd/rubo
LD_LIBRARY_PATH=../../runtime/target/debug go run main.go
```

## Architecture

Rubo combines the developer-friendly syntax of Ruby with the orchestration capabilities of Go and the raw performance of Rust. The Go-based CLI manages the compilation pipeline and uses CGO to interface with the high-performance Rust runtime.

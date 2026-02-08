# Rubo Compiler

Rubo is a "Ruby-syntax frontend with a Go-orchestrated, Rust-accelerated backend."

## 🚀 Performance
Rubo's hybrid architecture allows for dynamic optimization. When a function is identified as "hot," it is automatically promoted from Go interpretation to native Rust execution.

**Benchmark Results:**
- **Go Interpretation:** ~12.3µs
- **Native Rust:** ~2.9µs
- **Speedup:** **4.11x faster**

## 🛠️ Getting Started

### Prerequisites

We use [mise](https://mise.jdx.sh/) to manage project dependencies. If you don't have it installed, follow the [installation guide](https://mise.jdx.sh/getting-started.html).

Once `mise` is installed, run the following to install all required tools (Go, Rust, Ruby):

```bash
mise install
```

### Building the Project

You can build the entire stack (Rust runtime and Go CLI) using `mise`:

```bash
mise run build
```

This will:
1. Compile the Rust runtime in release mode.
2. Copy the library to the `lib/` directory.
3. Build the `rubo` executable in the root.

### How to Run (macOS)

Once built, you can run Rubo scripts using the `run` task:

```bash
mise run run -- examples/hello.rb
```

The `mise` environment automatically sets `DYLD_LIBRARY_PATH` so the executable can find the native library.


## 🏗️ Project Structure

- `src/cmd/rubo`: The CLI entry point (Go).
- `src/parser`: The Go-based syntax analyzer.
- `src/runtime`: The Rust-based performance engine (Rust).
- `src/bridge`: FFI/C-header glue.
- `src/orchestrator`: Manages the promotion of hot functions to Rust.

## Architecture

Rubo combines the developer-friendly syntax of Ruby with the orchestration capabilities of Go and the raw performance of Rust. The Go-based CLI manages the compilation pipeline and uses CGO to interface with the high-performance Rust runtime.

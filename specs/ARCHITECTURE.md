# Rubo Architecture Map

## Current State: Phase 7 (Complete)
- **Frontend:** Go-based Lexer and Parser generating a standard AST.
- **The Bridge:** `CGO` facilitates the call. `DYLD_LIBRARY_PATH` points to `./lib`.
- **Hot-Swap Mechanism:** - Profiler tracks call frequency.
    - Threshold: 1,000 calls.
    - Action: AST -> Rust Source -> `rustc --crate-type dylib` -> `dlopen` (or CGO equivalent).

## Interface Definition
- Rust functions exported via `#[no_mangle] pub extern "C"`.
- Data types currently supported: `int64`, `float64`.
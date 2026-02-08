# Rubo Constitution

## 1. Core Mandate
Rubo must maintain Ruby-like ergonomics while achieving near-native performance through dynamic Rust transpilation.

## 2. Technical Constraints
- **Orchestrator:** Must remain in Go for fast iteration and CLI stability.
- **Runtime:** Performance-critical paths must be offloaded to Rust.
- **Linkage:** Use CGO with dynamic linking (.dylib on macOS).
- **Memory:** Go owns the high-level AST; Rust owns the execution-time memory. Crossing the boundary must be explicit.

## 3. Performance Baselines
- **Math-heavy loop (Go):** ~2.35µs
- **Math-heavy loop (Rust):** ~573ns
- *Rule:* Any architectural change that regresses Rust execution speed by >5% must be rejected.
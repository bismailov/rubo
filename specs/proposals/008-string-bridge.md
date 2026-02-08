# Proposal 008: The String Bridge

**Status:** Accepted ✅
**Date:** 2026-02-08
**Selected Strategy:** Option B (Go-Side Cleanup / "The Loan")

## Problem
Go strings are garbage-collected headers (ptr + len). Rust strings are UTF-8 owned buffers or slices. Passing text between them via CGO requires a strategy to prevent memory leaks.

## Proposed Solution: The "Slab" Arena
1.  **Allocation:** Go allocates a `*C.char` buffer using `C.CString`.
2.  **Handoff:** Rust receives the pointer, wraps it in a `CStr`, and processes it.
3.  **Cleanup:** - *Option A:* Rust frees the memory (Requires `libc` crate).
    - *Option B:* Go defers `C.free` after the Rust call returns.
    
## Decision Log
**Decision:** We are proceeding with **Option B**.
**Reasoning:** 1. **Safety:** By using Go's `defer C.free`, we guarantee memory cleanup even if the Rust runtime panics. 
2. **Architecture:** Since Go is the "Orchestrator," it should manage the lifecycle of the buffers it allocates.
3. **Simplicity:** Avoids the need for custom deallocators or `libc` linkage inside the Rust kernels.

## Implementation Details (Option B)
- Use `C.CString(goStr)` on the Go side.
- Use `defer C.free(unsafe.Pointer(ptr))` immediately after allocation.
- Rust receives `*const c_char` and uses `CStr::from_ptr(ptr).to_str()`.

## Hardening Measures (Security & Stability)
- **Panic Guard:** All `extern "C"` functions must be wrapped in `panic::catch_unwind`.
- **Null Safety:** Explicit null pointer checks at the entry point of every Rust kernel.
- **Error Mapping:** Rust panics must be mapped to a specific Rubo error code (-3) so the Go orchestrator can handle it gracefully.


## Verification
- Must handle a "Hello World" string concatenation in 100ns or less.
- Must not show memory growth in `top` during a 1-million iteration string loop.
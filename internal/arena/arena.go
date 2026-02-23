package arena

/*
#include <stdlib.h>
#include <string.h>

void* malloc_arena(size_t size) {
    return malloc(size);
}

void free_arena(void* ptr) {
    free(ptr);
}
*/
import "C"
import (
	"unsafe"
)

// Arena manages an off-heap buffer for passing strings to C/Rust.
type Arena struct {
	start  unsafe.Pointer
	offset uintptr
	size   uintptr
}

// NewArena creates a new Arena with the specified capacity in bytes.
func NewArena(capacity int) *Arena {
	ptr := C.malloc_arena(C.size_t(capacity))
	return &Arena{
		start:  ptr,
		size:   uintptr(capacity),
		offset: 0,
	}
}

// Free releases the allocated memory.
func (a *Arena) Free() {
	if a.start != nil {
		C.free_arena(a.start)
		a.start = nil
	}
}

// Reset clears the arena, allowing reuse of the memory.
func (a *Arena) Reset() {
	a.offset = 0
}

// Copy allocates space for s in the arena, copies the string data,
// null-terminates it, and returns a C pointer.
func (a *Arena) Copy(s string) unsafe.Pointer {
	lenS := len(s)
	// Need space for string + null terminator
	required := uintptr(lenS + 1)

	if a.offset+required > a.size {
		// Simple panic on overflow for this implementation.
		// In a production system, we might allocate a new chunk.
		panic("Arena overflow: buffer execution cycle exceeded capacity")
	}

	// Calculate destination address
	dst := unsafe.Pointer(uintptr(a.start) + a.offset)

	// Copy string contents to the destination
	// Standard approach: convert Go string to bytes and copy
	// Or use C.memcpy if we wanted, but Go's copy is fine.

	// Create a slice view of the destination
	// This is safe because we know the bounds.
	// We construct a byte slice header carefully.
	dstSlice := unsafe.Slice((*byte)(dst), lenS+1)
	copy(dstSlice, s)
	dstSlice[lenS] = 0 // Null terminator

	// Advance offset
	a.offset += required

	return dst
}

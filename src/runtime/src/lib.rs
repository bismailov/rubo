pub mod hot_func;
use std::ffi::CStr;
use std::os::raw::c_char;
use std::panic;

#[no_mangle]
pub extern "C" fn rubo_add(a: i32, b: i32) -> i32 {
    a + b
}

/* #[no_mangle]
pub extern "C" fn rubo_string_len(ptr: *const c_char) -> i32 {
    if ptr.is_null() {
        return 0;
    }
    let c_str = unsafe { CStr::from_ptr(ptr) };
    let r_str = c_str.to_str().unwrap_or("");
    r_str.len() as i32
}
 */

#[no_mangle]
pub extern "C" fn rubo_string_len(ptr: *const c_char) -> i32 {
    // 1. Immediate Null Check
    if ptr.is_null() { return -1; }

    // 2. The Panic Boundary
    let result = panic::catch_unwind(|| {
        let c_str = unsafe { CStr::from_ptr(ptr) };
        
        // We use .to_str().expect() here because catch_unwind 
        // will now catch the panic if the string isn't valid UTF-8.
        let r_str = c_str.to_str().expect("Invalid UTF-8 passed to Rubo");
        
        r_str.len() as i32
    });

    // 3. Graceful Exit
    match result {
        Ok(len) => len,
        Err(_) => -3, // Custom error code for "Rust Runtime Panic"
    }
}

#[no_mangle]
pub extern "C" fn rubo_trigger_panic(_ptr: *const c_char) -> i32 {
    let result = panic::catch_unwind(|| {
        panic!("Intentional Rubo Panic for testing");
    });

    match result {
        Ok(_) => 0,
        Err(_) => -3,
    }
}

pub mod hot_func;

#[no_mangle]
pub extern "C" fn rubo_add(a: i32, b: i32) -> i32 {
    a + b
}

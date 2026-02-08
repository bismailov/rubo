#[no_mangle]
pub extern "C" fn heavy_math(x: i64, y: i64) -> i64 {
    ((((x + y) * (x - y)) + ((x * y) * (x + y))) - (y * y))
}
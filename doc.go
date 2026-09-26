// Package memless bridges the memless C ABI through purego, without cgo.
//
// The cdylib is found, in this order, through MEMLESS_LIB, a trusted
// (ideally absolute) path; the library a released module bundles for this
// platform, extracted once into the user cache directory and checked against
// its SHA-256 before every load; or, inside a checked-out workspace, under
// target/release then target/debug. MEMLESS_LIB loads arbitrary native code,
// like any FFI library path.
package memless

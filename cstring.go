package memless

import "unsafe"

// cString copies a NUL-terminated C string into a Go string. purego hands back
// a *byte so the caller keeps ownership and can free it; a string return would
// copy and lose the pointer to release.
func cString(message *byte) string {
	if message == nil {
		return ""
	}
	length := 0
	for *(*byte)(unsafe.Add(unsafe.Pointer(message), length)) != 0 {
		length++
	}
	return string(unsafe.Slice(message, length))
}

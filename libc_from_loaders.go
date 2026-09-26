package memless

import "slices"

var (
	muslLoaders  = []string{"/lib/ld-musl-x86_64.so.1", "/lib/ld-musl-aarch64.so.1"}
	glibcLoaders = []string{"/lib64/ld-linux-x86-64.so.2", "/lib/ld-linux-aarch64.so.1"}
)

// libcFromLoaders names the libc by its dynamic loader; the musl loader wins,
// since gcompat installs a glibc one on Alpine.
func libcFromLoaders(exists func(string) bool) string {
	if slices.ContainsFunc(muslLoaders, exists) {
		return "musl"
	}
	if slices.ContainsFunc(glibcLoaders, exists) {
		return "glibc"
	}
	return ""
}

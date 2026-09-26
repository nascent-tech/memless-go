package memless

const lddPath = "/usr/bin/ldd"

// detectLibc returns "glibc", "musl" or "": /usr/bin/ldd first, then the
// dynamic loader under /lib, for images that ship no ldd.
func detectLibc(readFile func(string) ([]byte, error), exists func(string) bool) string {
	data, err := readFile(lddPath)
	if libc := libcFromLdd(data, err); libc != "" {
		return libc
	}
	return libcFromLoaders(exists)
}

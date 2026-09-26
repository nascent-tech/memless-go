package memless

// bundleFor returns the library bundled for an OS and, on Linux, a libc. The
// Linux libraries target glibc: on musl, or when the libc cannot be told,
// there is none, as in the other bridges.
func bundleFor(goos, libc string) embeddedLibrary {
	if goos == "linux" && libc != "glibc" {
		return embeddedLibrary{}
	}
	return embeddedLibrary{files: embeddedFiles, file: embeddedFile, version: moduleVersion()}
}

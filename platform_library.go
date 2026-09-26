package memless

import (
	"os"
	"runtime"
)

// platformLibrary describes the library bundled for this platform.
func platformLibrary() embeddedLibrary {
	libc := ""
	if runtime.GOOS == "linux" {
		libc = detectLibc(os.ReadFile, fileExists)
	}
	return bundleFor(runtime.GOOS, libc)
}

package memless

import (
	"bytes"
	"os"
)

const lddPath = "/usr/bin/ldd"

// platformLibrary describes the library bundled for this platform. The Linux
// libraries target glibc: on musl there is none, as in the other bridges.
func platformLibrary() embeddedLibrary {
	if isMusl() {
		return embeddedLibrary{}
	}
	return embeddedLibrary{files: embeddedFiles, file: embeddedFile, version: moduleVersion()}
}

func isMusl() bool {
	data, err := os.ReadFile(lddPath)
	return err == nil && bytes.Contains(data, []byte("musl"))
}

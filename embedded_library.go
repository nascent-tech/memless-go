package memless

import (
	"io/fs"
	"path"
)

type embeddedLibrary struct {
	files   fs.FS
	file    string
	version string
}

// extract copies the bundled library into the cache and returns its path, or
// "" when there is none for this platform or the cache is unusable.
func (lib embeddedLibrary) extract(cacheDir string) string {
	if lib.files == nil || lib.file == "" || cacheDir == "" {
		return ""
	}
	data, err := fs.ReadFile(lib.files, lib.file)
	if err != nil {
		return ""
	}
	target, err := cachedCopy(cacheDir, lib.version, path.Base(lib.file), data)
	if err != nil {
		return ""
	}
	return target
}

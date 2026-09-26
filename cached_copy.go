package memless

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

const (
	shortSumBytes  = 6
	privateDirMode = 0o700
	libraryMode    = 0o755
)

// cachedCopy returns the path of a copy of data under
// cacheDir/<version>-<short sha256>/name. The copy on disk is checked against
// the full SHA-256 of data on every call and rewritten when it differs.
func cachedCopy(cacheDir, version, name string, data []byte) (string, error) {
	sum := sha256.Sum256(data)
	dir := filepath.Join(cacheDir, version+"-"+hex.EncodeToString(sum[:shortSumBytes]))
	if err := os.MkdirAll(dir, privateDirMode); err != nil {
		return "", err
	}
	target := filepath.Join(dir, name)
	if sameContent(target, sum) {
		return target, nil
	}
	return target, atomicWrite(dir, target, data)
}

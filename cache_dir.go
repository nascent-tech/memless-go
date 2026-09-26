package memless

import (
	"os"
	"path/filepath"
)

func cacheDir() string {
	base, err := os.UserCacheDir()
	if err != nil || base == "" {
		return ""
	}
	return filepath.Join(base, "memless")
}

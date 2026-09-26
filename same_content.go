package memless

import (
	"crypto/sha256"
	"os"
)

func sameContent(file string, sum [sha256.Size]byte) bool {
	data, err := os.ReadFile(file)
	return err == nil && sha256.Sum256(data) == sum
}

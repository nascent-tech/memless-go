package memless

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

const notFoundMessage = "memless cdylib not found: none bundled with this package, none under target/; set MEMLESS_LIB"

func libraryPath() (string, error) {
	return resolveLibrary(currentPlan())
}

func envLibrary(env string) (string, error) {
	if !exists(env) {
		return "", errors.New("memless cdylib not found at " + env + "; set MEMLESS_LIB")
	}
	return env, nil
}

func searchLibrary(root string) (string, error) {
	found := existingLibrary(root)
	if found == "" {
		return "", errors.New(notFoundMessage)
	}
	return found, nil
}

func existingLibrary(root string) string {
	found := ""
	for _, name := range libraryNames() {
		found = orExisting(found, filepath.Join(root, name))
	}
	return found
}

func orExisting(found, candidate string) string {
	if found != "" {
		return found
	}
	if !exists(candidate) {
		return ""
	}
	return candidate
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Same order as the PHP and Node bridges: release before debug, .dylib
// before .so.
func libraryNames() []string {
	return []string{
		"target/release/libmemless_capi.dylib",
		"target/debug/libmemless_capi.dylib",
		"target/release/libmemless_capi.so",
		"target/debug/libmemless_capi.so",
	}
}

func workspaceRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..")
}

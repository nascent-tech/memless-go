package memless

import (
	"fmt"
	"sync"

	"github.com/ebitengine/purego"
)

var (
	loadOnce          sync.Once
	loadErr           error
	memlessAbiVersion func() uint32
	memlessLoad       func(path string, outHandle *uint64, outMessage **byte) int32
	memlessRelease    func(handle uint64)
	memlessFreeString func(message *byte)

	memlessExecute           func(handle uint64, sql string, outAffected *uint64, outMessage **byte) int32
	memlessReload            func(handle uint64, outMessage **byte) int32
	memlessQuery             func(handle uint64, sql string, outResult *uint64, outMessage **byte) int32
	memlessResultColumnCount func(result uint64) uint64
	memlessResultRowCount    func(result uint64) uint64
	memlessResultColumn      func(result uint64, index uint64) string
	memlessResultCell        func(
		result, row, column uint64,
		outInteger *int64, outDecimal *float64, outBoolean *int32, outText **byte,
	) int32
	memlessResultRelease func(result uint64)
)

func ensureLoaded() error {
	loadOnce.Do(func() {
		loadErr = openLibrary()
	})
	return loadErr
}

func openLibrary() (err error) {
	defer func() { err = recovered(recover(), err) }()
	path, err := libraryPath()
	if err != nil {
		return err
	}
	handle, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_LOCAL)
	if err != nil {
		return err
	}
	registerFunctions(handle)
	return checkAbi()
}

// recovered turns a purego RegisterLibFunc panic (a missing symbol) into an
// error instead of crashing Load.
func recovered(panicked any, err error) error {
	if panicked != nil {
		return fmt.Errorf("memless library load failed: %v", panicked)
	}
	return err
}

func registerFunctions(handle uintptr) {
	registerCore(handle)
	registerResult(handle)
}

func registerCore(handle uintptr) {
	purego.RegisterLibFunc(&memlessAbiVersion, handle, "memless_abi_version")
	purego.RegisterLibFunc(&memlessLoad, handle, "memless_load")
	purego.RegisterLibFunc(&memlessRelease, handle, "memless_release")
	purego.RegisterLibFunc(&memlessFreeString, handle, "memless_free_string")
	purego.RegisterLibFunc(&memlessExecute, handle, "memless_execute")
	purego.RegisterLibFunc(&memlessReload, handle, "memless_reload")
	purego.RegisterLibFunc(&memlessQuery, handle, "memless_query")
}

func registerResult(handle uintptr) {
	purego.RegisterLibFunc(&memlessResultColumnCount, handle, "memless_result_column_count")
	purego.RegisterLibFunc(&memlessResultRowCount, handle, "memless_result_row_count")
	purego.RegisterLibFunc(&memlessResultColumn, handle, "memless_result_column")
	purego.RegisterLibFunc(&memlessResultCell, handle, "memless_result_cell")
	purego.RegisterLibFunc(&memlessResultRelease, handle, "memless_result_release")
}

func checkAbi() error {
	return checkVersion(memlessAbiVersion())
}

func checkVersion(version uint32) error {
	if version != abiVersion {
		return fmt.Errorf("memless ABI mismatch: expected %d, got %d", abiVersion, version)
	}
	return nil
}

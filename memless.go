package memless

import "strings"

const (
	abiVersion            uint32 = 5
	statusOk              int32  = 0
	statusRefused         int32  = 1
	statusInvalidArgument int32  = 2
)

// Instance is a live memless instance behind an opaque handle.
type Instance struct {
	handle   uint64
	released bool
}

// Load loads the file at path. It returns a RefusalError when the file is
// refused, and a FaultError on a boundary or internal fault. No instance is
// returned on any error.
func Load(path string) (*Instance, error) {
	if strings.IndexByte(path, 0) >= 0 {
		return nil, &FaultError{Status: statusInvalidArgument, Message: "path contains a NUL byte"}
	}
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	var handle uint64
	var message *byte
	status := memlessLoad(path, &handle, &message)
	return interpret(status, handle, message)
}

// Release releases the instance once. A second call, a zero-value instance, or
// an unknown handle is ignored.
func (i *Instance) Release() {
	if i.released || i.handle == 0 {
		return
	}
	i.released = true
	memlessRelease(i.handle)
}

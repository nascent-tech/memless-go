package memless

import "strings"

// Execute runs a write sql against the instance and returns the affected row
// count. It returns a RefusalError on a validation or disk failure, and a
// FaultError on a boundary or internal fault.
func (i *Instance) Execute(sql string) (uint64, error) {
	if strings.IndexByte(sql, 0) >= 0 {
		return 0, &FaultError{Status: statusInvalidArgument, Message: "sql contains a NUL byte"}
	}
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	var affected uint64
	var message *byte
	status := memlessExecute(i.handle, sql, &affected, &message)
	return interpretExecute(status, affected, message)
}

func interpretExecute(status int32, affected uint64, message *byte) (uint64, error) {
	text := takeMessage(message)
	if status == statusOk {
		return affected, nil
	}
	if status == statusRefused {
		return 0, &RefusalError{Message: text}
	}
	return 0, &FaultError{Status: status, Message: text}
}

package memless

import "fmt"

// FaultError is a boundary or internal fault (InvalidArgument or Internal), not
// a file refusal. errors.As tells it apart from a RefusalError.
type FaultError struct {
	Status  int32
	Message string
}

func (e *FaultError) Error() string {
	return fmt.Sprintf("memless fault (%d): %s", e.Status, e.Message)
}

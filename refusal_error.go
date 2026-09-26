package memless

// RefusalError carries the domain message verbatim (D13); errors.As tells it
// apart from a boundary fault.
type RefusalError struct {
	Message string
}

func (e *RefusalError) Error() string {
	return e.Message
}

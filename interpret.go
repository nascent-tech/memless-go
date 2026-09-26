package memless

func interpret(status int32, handle uint64, message *byte) (*Instance, error) {
	text := takeMessage(message)
	if status == statusOk {
		return &Instance{handle: handle}, nil
	}
	if status == statusRefused {
		return nil, &RefusalError{Message: text}
	}
	return nil, &FaultError{Status: status, Message: text}
}

func takeMessage(message *byte) string {
	if message == nil {
		return ""
	}
	text := cString(message)
	memlessFreeString(message)
	return text
}

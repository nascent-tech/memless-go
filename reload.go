package memless

// Reload re-reads the file the instance was loaded on and replaces the
// in-memory state with it, exactly as Load would build it. It returns a
// RefusalError while a transaction is open, or when the file would be refused
// by Load; the instance's state is left intact either way and stays usable.
func (i *Instance) Reload() error {
	if err := ensureLoaded(); err != nil {
		return err
	}
	var message *byte
	status := memlessReload(i.handle, &message)
	_, err := interpretExecute(status, 0, message)
	return err
}

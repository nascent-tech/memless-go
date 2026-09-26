package memless

// Begin opens a transaction on the instance. Writes then apply to a working
// state the file does not yet reflect, and a query sees them. A second Begin
// returns a RefusalError.
func (i *Instance) Begin() error {
	_, err := i.Execute("BEGIN")
	return err
}

// Commit validates the open transaction and rewrites the file once when the
// working state differs. With no open transaction it returns a RefusalError; a
// failed validation leaves memory and the file at the before-state.
func (i *Instance) Commit() error {
	_, err := i.Execute("COMMIT")
	return err
}

// Rollback discards the open transaction without touching the disk. With no
// open transaction it returns a RefusalError.
func (i *Instance) Rollback() error {
	_, err := i.Execute("ROLLBACK")
	return err
}

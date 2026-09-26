package memless

import "os"

// atomicWrite writes data to a temporary file in dir, then renames it over
// target, so a reader never sees a partial library.
func atomicWrite(dir, target string, data []byte) error {
	temp, err := os.CreateTemp(dir, ".memless-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(temp.Name()) }()
	if err := writeAndClose(temp, data); err != nil {
		return err
	}
	if err := os.Chmod(temp.Name(), libraryMode); err != nil {
		return err
	}
	return os.Rename(temp.Name(), target)
}

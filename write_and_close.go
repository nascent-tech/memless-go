package memless

import "os"

func writeAndClose(file *os.File, data []byte) error {
	_, err := file.Write(data)
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	return err
}

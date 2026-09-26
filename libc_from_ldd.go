package memless

import "bytes"

func libcFromLdd(data []byte, err error) string {
	if err != nil {
		return ""
	}
	if bytes.Contains(data, []byte("musl")) {
		return "musl"
	}
	if bytes.Contains(data, []byte("GNU C Library")) {
		return "glibc"
	}
	return ""
}

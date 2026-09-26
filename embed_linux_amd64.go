//go:build linux && amd64

package memless

import "embed"

//go:embed lib/linux-x64-gnu
var embeddedFiles embed.FS

const embeddedFile = "lib/linux-x64-gnu/libmemless_capi.so"

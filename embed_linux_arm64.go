//go:build linux && arm64

package memless

import "embed"

//go:embed lib/linux-arm64-gnu
var embeddedFiles embed.FS

const embeddedFile = "lib/linux-arm64-gnu/libmemless_capi.so"

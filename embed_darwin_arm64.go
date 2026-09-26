//go:build darwin && arm64

package memless

import "embed"

//go:embed lib/darwin-arm64
var embeddedFiles embed.FS

const embeddedFile = "lib/darwin-arm64/libmemless_capi.dylib"

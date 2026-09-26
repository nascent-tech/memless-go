//go:build darwin && amd64

package memless

import "embed"

//go:embed lib/darwin-x64
var embeddedFiles embed.FS

const embeddedFile = "lib/darwin-x64/libmemless_capi.dylib"

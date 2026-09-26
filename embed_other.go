//go:build !(darwin && (arm64 || amd64)) && !(linux && (arm64 || amd64))

package memless

import "embed"

var embeddedFiles embed.FS

const embeddedFile = ""

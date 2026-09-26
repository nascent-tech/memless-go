package memless

import (
	"runtime/debug"
	"slices"
	"strings"
)

const (
	modulePath   = "github.com/nascent-tech/memless-go"
	develVersion = "devel"
)

func moduleVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return develVersion
	}
	return versionIn(append(info.Deps, &info.Main))
}

func versionIn(modules []*debug.Module) string {
	index := slices.IndexFunc(modules, isThisModule)
	if index < 0 {
		return develVersion
	}
	return modules[index].Version
}

func isThisModule(module *debug.Module) bool {
	return module.Path == modulePath && strings.HasPrefix(module.Version, "v")
}

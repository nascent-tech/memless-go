package memless

type libraryPlan struct {
	env      string
	embedded embeddedLibrary
	cacheDir string
	root     string
}

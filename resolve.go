package memless

// resolveLibrary applies the search order shared by the three bridges:
// MEMLESS_LIB, then the library bundled for this platform, then target/.
func resolveLibrary(plan libraryPlan) (string, error) {
	if plan.env != "" {
		return envLibrary(plan.env)
	}
	if extracted := plan.embedded.extract(plan.cacheDir); extracted != "" {
		return extracted, nil
	}
	return searchLibrary(plan.root)
}

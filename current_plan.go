package memless

import "os"

func currentPlan() libraryPlan {
	plan := libraryPlan{env: os.Getenv("MEMLESS_LIB"), root: workspaceRoot()}
	plan.embedded = platformLibrary()
	plan.cacheDir = cacheDir()
	return plan
}

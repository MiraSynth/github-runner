package version

import "fmt"

const (
	Major           = "0"
	Minor           = "0"
	Patch           = "0"
	PreRelease      = "alpha"
	PreReleasePatch = "0"

	Tag = ""
)

func GetVersion() string {
	if PreRelease == "" {
		return fmt.Sprintf("%s.%s.%s", Major, Minor, Patch)
	}

	return fmt.Sprintf("%s.%s.%s-%s.%s", Major, Minor, Patch, PreRelease, PreReleasePatch)
}

func GetTag() string {
	return Tag
}

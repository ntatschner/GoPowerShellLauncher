package utils

import (
	"strings"
)

func SplitProfiles(profiles string) []string {
	return strings.Split(profiles, ",")
}

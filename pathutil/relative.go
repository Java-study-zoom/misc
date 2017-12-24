package pathutil

import (
	"strings"
)

// Relative returns the relative path of full to base.
func Relative(base, full string) string {
	if base == full {
		return "."
	}
	if strings.HasPrefix(full, base+"/") {
		return "./" + strings.TrimPrefix(full, base+"/")
	}
	return full
}

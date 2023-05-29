package fsutil

import (
	"os"
	"strings"
)

func IsHidden(dir string) bool {
	parts := strings.Split(dir, string(os.PathSeparator))
	if len(parts) == 0 {
		return false
	}
	return isHidden(parts[len(parts)-1])
}

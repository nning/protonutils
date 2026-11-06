package steam

import (
	"strconv"
)

// IsShortcut returns whether app id is a shortcut (non-Steam game). Shortcut
// app IDs have the most significant bit of the 32 bit ID set.
func IsShortcut(id string) bool {
	n, err := strconv.ParseUint(id, 10, 32)
	return err == nil && n > 1<<31-1
}

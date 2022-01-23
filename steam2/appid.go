package steam2

import (
	"strconv"
)

func IsShortcut(id string) bool {
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		_, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			return true
		}
	}

	return false
}

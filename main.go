package list_proton_versions

import "log"

func PanicOnError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

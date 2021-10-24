package list_proton_versions

import (
	"fmt"
	"log"
	"os"
)

func ExitOnError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

func PanicOnError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

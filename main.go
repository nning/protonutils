package list_proton_versions

import (
	"fmt"
	"os"
)

func ExitOnError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

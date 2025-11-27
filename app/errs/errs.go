package errs

import (
	"fmt"
	"os"
)

var HadError = false

func Error(message string) {
	fmt.Fprintf(os.Stderr, "[line 1] Error: %s\n", message)
	HadError = true
}

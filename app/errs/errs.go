package errs

import (
	"fmt"
	"os"
)

var HadError = false

func Error(line int, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, message)
	HadError = true
}

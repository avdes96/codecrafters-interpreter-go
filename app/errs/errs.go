package errs

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/token"
)

var HadError = false

func ErrorOnLine(line int, message string) {
	report(line, "", message)
}

func ErrorAtToken(t *token.Token, message string) {
	if t.Type == token.EOF {
		report(t.Line, " at end", message)
	} else {
		report(t.Line, fmt.Sprintf(" at '%s'", t.Lexeme), message)
	}
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	HadError = true
}

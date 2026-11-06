// Package debug controls diagnostic output enabled by the global --debug flag.
package debug

import (
	"fmt"
	"io"
)

// Enabled is set by the global --debug flag.
var Enabled bool

// Printf writes to w only when debug output is enabled.
func Printf(w io.Writer, format string, args ...any) {
	if !Enabled {
		return
	}
	fmt.Fprintf(w, format, args...)
}

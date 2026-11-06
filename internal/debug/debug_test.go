package debug

import (
	"bytes"
	"testing"
)

func TestPrintf(t *testing.T) {
	defer func() { Enabled = false }()

	var buf bytes.Buffer
	Enabled = false
	Printf(&buf, "hidden %d\n", 1)
	if buf.Len() != 0 {
		t.Errorf("expected no output when disabled, got %q", buf.String())
	}

	Enabled = true
	Printf(&buf, "shown %d\n", 2)
	if got := buf.String(); got != "shown 2\n" {
		t.Errorf("expected %q, got %q", "shown 2\n", got)
	}
}

package writer_test

import (
	"os"
	"testing"

	"github.com/gummiboll/forgetful/writer"
)

func TestgetEditor(t *testing.T) {
	var err error
	e := writer.GetEditor()

	if e != "vim" {
		t.Errorf("Expected 'vim', got: %s", e)
	}

	if os.Setenv("EDITOR", "pico"); err != nil {
		t.Errorf("Faied to set $EDITOR (%s)", err)
	}
	if e != "pico" {
		t.Errorf("Expected 'pico', got: %s", e)
	}
}

package writer_test

import (
	"testing"

	"github.com/gummiboll/forgetful/writer"
)

func TestGetEditor(t *testing.T) {
	e := writer.GetEditor()

	if e != "vim" {
		t.Errorf("Expected 'vim', got: %s", e)
	}
}

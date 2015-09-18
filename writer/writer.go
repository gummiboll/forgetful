package writer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/gummiboll/forgetful/storage"
)

// GetEditor returns a editor. $EDITOR if its set, vim otherwise.
func GetEditor() (e string) {
	envEditor := os.Getenv("EDITOR")
	if envEditor == "" {
		return "vim"
	}

	return envEditor
}

// WriteNote opens up a editor for Note-input
func WriteNote(n *storage.Note) (err error) {
	tmpfile := fmt.Sprintf("%sforgetful-%d.tmp", os.TempDir(), time.Now().Unix())
	editor := GetEditor()
	if err := ioutil.WriteFile(tmpfile, []byte(n.Text), 0600); err != nil {
		return err
	}

	cmd := exec.Command(editor, tmpfile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Read the note
	data, err := ioutil.ReadFile(tmpfile)
	if err != nil {
		return err
	}

	n.Text = string(data)
	// Cleanup
	os.Remove(tmpfile)

	return nil
}

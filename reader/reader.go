package reader

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/gummiboll/forgetful/storage"
)

// ReadNote pipes a note to less
func ReadNote(n storage.Note) (err error) {
	cmd := exec.Command("less")
	r, stdin := io.Pipe()
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if cmd.Start(); err != nil {
		return err
	}
	fmt.Fprintf(stdin, n.Text)
	if stdin.Close(); err != nil {
		return err
	}
	if cmd.Wait(); err != nil {
		return err
	}
	return nil
}

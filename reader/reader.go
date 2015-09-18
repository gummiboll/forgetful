package reader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/gummiboll/forgetful/storage"
)

// HastebinResponse represents a response from hastebin
type HastebinResponse struct {
	Key string `json:"key"`
}

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

// ShareNote sends a note to hastebin.com
func ShareNote(n storage.Note, url string) (purl string, err error) {
	resp, err := http.Post(url, "text", bytes.NewBuffer([]byte(n.Text)))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Share failed, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Share failed: %s", err)
	}

	hr := HastebinResponse{}
	if json.Unmarshal(body, &hr); err != nil {
		return "", fmt.Errorf("Share failed: %s", err)
	}

	return fmt.Sprintf("http://hastebin.com/%s", hr.Key), nil
}

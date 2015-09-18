package reader_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gummiboll/forgetful/reader"
	"github.com/gummiboll/forgetful/storage"
)

func TestShareNote(t *testing.T) {
	n := storage.Note{Name: "A test", Text: "Some example text"}

	echoHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"key":"aTeStStRiNg"}`)
	}

	ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()

	url, err := reader.ShareNote(n, ts.URL)
	if err != nil {
		t.Errorf("Share note failed (%s)", err)
	}

	if url != "http://hastebin.com/aTeStStRiNg" {
		t.Errorf("Tried sharing a note but didn't get the expected url. Url was: %s", url)
	}
}

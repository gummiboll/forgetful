package storage_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gummiboll/forgetful/storage"
	"github.com/jinzhu/gorm"
)

var (
	i storage.Impl
)

func init() {
	var err error
	testdb := "test.db"
	// Check if test.db exists, delete it if so
	if _, err := os.Stat(testdb); err == nil {
		if os.Remove(testdb); err != nil {
			panic(fmt.Sprintf("Failed to remove old testdb (%s)", err))
		}
	}

	i.DB, err = gorm.Open("sqlite3", testdb)
	if err != nil {
		panic(fmt.Sprintf("Failed to create test database (%s)", err))
	}

	i.DB.LogMode(false)
	i.InitSchema()
}

func TestSaveNote(t *testing.T) {
	n := storage.Note{Name: "a test", Text: "Some example text", Temporary: false}
	if err := i.SaveNote(&n); err != nil {
		t.Errorf("Failed to save note: %s", err)
	}

	// Create another note that is temporary
	n2 := storage.Note{Name: "Another test", Text: "foobar", Temporary: true}
	if err := i.SaveNote(&n2); err != nil {
		t.Errorf("Failed to save note: %s", err)
	}
}

func TestLoadNote(t *testing.T) {
	n, err := i.LoadNote("a test")
	if err != nil {
		t.Errorf("Failed to load note (%s)", err)
	}

	if n.Text != "Some example text" {
		t.Errorf("Text doesnt match expected text")
	}

	// Load a note case insensitive
	n2, err := i.LoadNote("another TEST")
	if err != nil {
		t.Errorf("Failed to load note (%s)", err)
	}

	if n2.Text != "foobar" {
		t.Errorf("Failed to load a note case insensitive")
	}

}

func TestListNotes(t *testing.T) {
	notes := i.ListNotes("")
	if len(notes) < 2 {
		t.Errorf("Two notes expected, found %d", len(notes))
	}

	notesWithFilter := i.ListNotes("another")
	if len(notesWithFilter) != 1 {
		t.Errorf("Expected one note, found %d", len(notesWithFilter))
	}
}

func TestSearchNotes(t *testing.T) {
	notes := i.SearchNotes("foobar")
	if len(notes) != 1 {
		t.Errorf("Expected one note, found %d", len(notes))
	}
}

/*
Figure out a better way to do this test
func TestRemoveExpiredNotes(t *testing.T) {
	n, err := i.LoadNote("another test")
	if err != nil {
		t.Errorf("Failed to load note (%s)", err)
	}

	n.UpdatedAt = time.Now().Add(-25 * time.Hour)
	if i.SaveNote(&n); err != nil {
		t.Errorf("Failed to set CreatedAt in the past (%s)", err)
	}

	if i.RemoveExpiredNotes(); err != nil {
		t.Errorf("Failed to remove expired notes (%s)", err)
	}

	nExists := i.NoteExists("another TEST")
	if nExists == true {
		t.Errorf("Tried to load a expired note efter RemoveExpiredNotes() was run and was successful")
	}

}
*/

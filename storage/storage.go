package storage

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/jinzhu/gorm"
	// For gorm
	_ "github.com/mattn/go-sqlite3"
)

// Note represents a note
type Note struct {
	gorm.Model
	Name      string `sql:"not null;primary_key;index"`
	Text      string
	Temporary bool
}

// Impl represents a gorm.DB
type Impl struct {
	DB *gorm.DB
}

// forgetfulDir returns /path/to/.forgetful and creates it
// if it doesnt exist yet
func forgetfulDir() (d string) {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	fDir := fmt.Sprintf("%s/.forgetful", usr.HomeDir)
	_, err = os.Stat(fDir)
	if err != nil {
		os.Mkdir(fDir, 0700)
	}

	return fDir
}

// InitDB tries to open db and creates one if needed
func (i *Impl) InitDB() (err error) {
	fDir := forgetfulDir()
	i.DB, err = gorm.Open("sqlite3", fDir+"/forgetful.db")
	if err != nil {
		return fmt.Errorf("Failed to connect to db, error was: %v", err)
	}

	i.DB.LogMode(false)

	return nil
}

// InitSchema initializes the schema
func (i *Impl) InitSchema() {
	i.DB.AutoMigrate(&Note{})
}

// NoteExists returns true if note exist
func (i *Impl) NoteExists(n string) bool {
	_, err := i.LoadNote(n)
	if err != nil {
		return false
	}

	return true
}

// RemoveExpiredNotes searches for temporary notes older then 24 hours
// and removes them
func (i *Impl) RemoveExpiredNotes() (err error) {
	notes := []Note{}
	timeFilter := time.Now().Add(-24 * time.Hour)
	if err := i.DB.Where("updated_at < ? and temporary == ?", timeFilter, true).Find(&notes).Error; err != nil {
		return fmt.Errorf("Failed to search for expired notes: %v", err)
	}

	for _, n := range notes {
		if err := i.DeleteNote(n); err != nil {
			return fmt.Errorf("Failed to delete expired note (%s): %v", n.Name, err)
		}
	}

	return nil
}

// ListNotes returns a list of all notes. If f is specified, filter results
func (i *Impl) ListNotes(f string) (n []Note) {
	notes := []Note{}
	if f == "" {
		i.DB.Find(&notes)
	} else {
		i.DB.Where("name LIKE ?", fmt.Sprintf("%s%s%s", "%", f, "%")).Find(&notes)
	}

	return notes
}

// SearchNotes returns a list of notes matching f
func (i *Impl) SearchNotes(f string) (n []Note) {
	notes := []Note{}
	like := "%" + f + "%"
	i.DB.Where("name LIKE ?", like).Or("text LIKE ?", like).Find(&notes)

	return notes
}

// SaveNote saves a note and returns it
func (i *Impl) SaveNote(n *Note) (err error) {
	if err := i.DB.Save(&n).Error; err != nil {
		return fmt.Errorf("Failed to save note: %s", err)
	}

	return nil
}

// LoadNote loads a note and returns it if successful
func (i *Impl) LoadNote(name string) (n Note, err error) {
	n = Note{}
	if i.DB.Where("name = ? COLLATE NOCASE", name).First(&n).Error != nil {
		return n, fmt.Errorf("Note: %s not found", name)
	}

	return n, nil
}

// DeleteNote deletes a note
func (i *Impl) DeleteNote(n Note) (err error) {
	if err := i.DB.Delete(&n).Error; err != nil {
		return fmt.Errorf("Failed to delete note: %s (id: %d)", n.Name, n.ID)
	}

	return nil
}

// RenameNote renames a notes
func (i *Impl) RenameNote(nID uint, newName string) (err error) {
	if err := i.DB.Model(&Note{}).Where("id = ?", nID).Update("name", newName).Error; err != nil {
		return err
	}

	return nil
}

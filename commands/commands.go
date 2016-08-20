package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/urfave/cli"

	"github.com/gummiboll/forgetful/reader"
	"github.com/gummiboll/forgetful/storage"
	"github.com/gummiboll/forgetful/writer"
)

type byLetterNocase []string

func (s byLetterNocase) Len() int {
	return len(s)
}
func (s byLetterNocase) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLetterNocase) Less(i, j int) bool {
	return strings.ToLower(s[i]) < strings.ToLower(s[j])
}

// FormatNoteList returns a formatted list of notes
func FormatNoteList(notes []storage.Note) (rnotes []string) {
	for _, n := range notes {
		nStr := fmt.Sprintf("\U0001f539  %s", n.Name)
		if n.Temporary {
			validTo := n.UpdatedAt.Add(24 * time.Hour)
			dur := validTo.Sub(time.Now())
			nStr += fmt.Sprintf(" (\U0001f4a5 in %s) ", dur)
		}
		rnotes = append(rnotes, nStr)
	}

	sort.Sort(byLetterNocase(rnotes))

	return rnotes
}

// NoteName returns name of note and error if note isnt present
func NoteName(c *cli.Context) (n string, err error) {
	if c.Args().Present() != true {
		return "", errors.New("Missing argument: name")
	}

	return strings.Join(c.Args(), " "), nil
}

// AddCommand adds a Note
func AddCommand(c *cli.Context, i storage.Impl) (n storage.Note, err error) {
	nName, err := NoteName(c)
	if err != nil {
		return n, err
	}

	if exists := i.NoteExists(nName); exists == true {
		return n, fmt.Errorf("Note already exists")
	}

	n.Name = nName
	n.Temporary = c.Bool("t")

	// Only open editor if -p (read from clipboard) isnt set
	if c.IsSet("p") {
		nText, err := clipboard.ReadAll()
		if err != nil {
			return n, err
		}
		n.Text = nText
	} else {
		if err := writer.WriteNote(&n); err != nil {
			return n, err
		}
	}

	if err := i.SaveNote(&n); err != nil {
		return n, err
	}

	return n, nil
}

// DeleteCommand deletes a Note
func DeleteCommand(c *cli.Context, i storage.Impl) (n storage.Note, err error) {
	nName, err := NoteName(c)
	if err != nil {
		return n, err
	}

	n, err = i.LoadNote(nName)
	if err != nil {
		return n, err
	}

	if i.DeleteNote(n) != nil {
		return n, err
	}

	return n, nil
}

// EditCommand edits a Note
func EditCommand(c *cli.Context, i storage.Impl) (n storage.Note, err error) {
	nName, err := NoteName(c)
	if err != nil {
		return n, err
	}

	n, err = i.LoadNote(nName)
	if err != nil {
		return n, err
	}

	if err := writer.WriteNote(&n); err != nil {
		return n, err
	}

	if err := i.SaveNote(&n); err != nil {
		return n, err
	}

	return n, nil

}

// ReadCommand reads a Note
func ReadCommand(c *cli.Context, i storage.Impl) (err error) {
	nName, err := NoteName(c)
	if err != nil {
		return err
	}

	n, err := i.LoadNote(nName)
	if err != nil {
		return err
	}

	if err := reader.ReadNote(n); err != nil {
		return err
	}

	return nil
}

// RenameCommand renames a Note
func RenameCommand(c *cli.Context, i storage.Impl) (nName string, newName string, err error) {
	nName, err = NoteName(c)
	if err != nil {
		return nName, newName, err
	}

	n, err := i.LoadNote(nName)
	if err != nil {
		return nName, newName, err
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(fmt.Sprintf("Rename note '%s' to: ", n.Name))
	newName, err = reader.ReadString('\n')
	if err != nil {
		return nName, newName, err
	}

	newName = strings.Trim(newName, "\n")

	if newName == "" {
		return nName, newName, errors.New("Note name can't be blank")
	}

	if i.NoteExists(newName) == true {
		return nName, newName, fmt.Errorf("Note '%s' already exists", newName)
	}

	if err = i.RenameNote(n.ID, newName); err != nil {
		return nName, newName, fmt.Errorf("Failed to rename note '%s' to '%s'.", nName, newName)
	}

	return nName, newName, nil
}

// ListCommand lists Notes
func ListCommand(c *cli.Context, i storage.Impl) (rnotes []string) {
	nName := strings.Join(c.Args(), " ")
	notes := i.ListNotes(nName)

	return FormatNoteList(notes)

}

// SearchCommand searches for Notes
func SearchCommand(c *cli.Context, i storage.Impl) (rnotes []string, err error) {
	nName, err := NoteName(c)
	if err != nil {
		return rnotes, err
	}

	notes := i.SearchNotes(nName)

	return FormatNoteList(notes), nil
}

// ShareCommand shares a Note
func ShareCommand(c *cli.Context, i storage.Impl) (n storage.Note, url string, err error) {
	nName, err := NoteName(c)
	if err != nil {
		return n, url, err
	}

	n, err = i.LoadNote(nName)
	if err != nil {
		return n, url, err
	}

	url, err = reader.ShareNote(n, "http://hastebin.com/documents")

	if err != nil {
		return n, url, err
	}

	return n, url, nil
}

// KeepCommand keeps/unkeeps a Note
func KeepCommand(c *cli.Context, i storage.Impl, k bool) (n storage.Note, err error) {
	nName, err := NoteName(c)
	if err != nil {
		return n, err
	}

	n, err = i.LoadNote(nName)
	if err != nil {
		return n, err
	}

	n.Temporary = !k
	if err := i.SaveNote(&n); err != nil {
		return n, err
	}

	return n, err
}

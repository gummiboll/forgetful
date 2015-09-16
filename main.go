package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/codegangsta/cli"
	"github.com/gummiboll/forgetful/reader"
	"github.com/gummiboll/forgetful/storage"
	"github.com/gummiboll/forgetful/writer"
)

const version string = "0.9"

// NoteName returns name of note and error if note isnt present
func NoteName(c *cli.Context) (n string, err error) {
	if c.Args().Present() != true {
		return "", errors.New("Missing argument: name")
	}

	return strings.Join(c.Args(), " "), nil
}

func main() {
	// Init
	i := storage.Impl{}
	if err := i.InitDB(); err != nil {
		panic(err)
	}

	i.InitSchema()
	if err := i.RemoveExpiredNotes(); err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "forgetful"
	app.Usage = "For your notes/cheat sheets"
	app.Version = version
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a note",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "t",
					Usage: "Mark as temporary (expires after 24 hours)",
				},
				cli.BoolFlag{
					Name:  "p",
					Usage: "Create note with contents from clipboard",
				},
			},
			Action: func(c *cli.Context) {
				nName, err := NoteName(c)
				if err != nil {
					fmt.Println(err)
					return
				}

				if exists := i.NoteExists(nName); exists == true {
					fmt.Println("Note already exists")
					return
				}

				n := storage.Note{Name: nName, Temporary: c.Bool("t")}

				// Only open editor if -p (read from clipboard) isnt set
				if c.IsSet("p") {
					nText, err := clipboard.ReadAll()
					if err != nil {
						fmt.Println(err)
						return
					}
					n.Text = nText
				} else {
					if err := writer.WriteNote(&n); err != nil {
						fmt.Println(err)
						return
					}
				}

				if err := i.SaveNote(&n); err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(fmt.Sprintf("Added note: %s", n.Name))
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete a note",
			Action: func(c *cli.Context) {
				nName, err := NoteName(c)
				if err != nil {
					fmt.Println(err)
					return
				}

				n, err := i.LoadNote(nName)
				if err != nil {
					fmt.Println(err)
					return
				}

				if i.DeleteNote(n) != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(fmt.Sprintf("Deleted note: %s", nName))
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "Edit/read a note",
			Action: func(c *cli.Context) {
				nName, err := NoteName(c)
				if err != nil {
					fmt.Println(err)
					return
				}

				n, err := i.LoadNote(nName)
				if err != nil {
					fmt.Println(err)
					return
				}

				if err := writer.WriteNote(&n); err != nil {
					fmt.Println(err)
				}

				if err := i.SaveNote(&n); err != nil {
					fmt.Println(err)
				}

				fmt.Println(fmt.Sprintf("Updated note: %s", n.Name))
			},
		},
		{
			Name:    "read",
			Aliases: []string{"r"},
			Usage:   "Read a note",
			Action: func(c *cli.Context) {
				nName, err := NoteName(c)
				if err != nil {
					fmt.Println(err)
					return
				}

				n, err := i.LoadNote(nName)
				if err != nil {
					fmt.Println(err)
					return
				}

				if err := reader.ReadNote(n); err != nil {
					fmt.Println(err)
					return
				}
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all notes, filter result if argument i present",
			Action: func(c *cli.Context) {
				nName := strings.Join(c.Args(), " ")
				notes := i.ListNotes(nName)
				if len(notes) > 0 {
					fmt.Println("Matching notes:")
				}

				for _, n := range notes {
					nStr := fmt.Sprintf("* %s", n.Name)
					if n.Temporary {
						validTo := n.CreatedAt.Add(24 * time.Hour)
						dur := validTo.Sub(time.Now())
						nStr += fmt.Sprintf(" (valid for %s)", dur)
					}

					fmt.Println(nStr)
				}
			},
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search notes for argument",
			Action: func(c *cli.Context) {
				nName, err := NoteName(c)
				if err != nil {
					fmt.Println(err)
					return
				}

				notes := i.SearchNotes(nName)
				if len(notes) > 0 {
					fmt.Println("Matching notes:")
				}

				for _, n := range notes {
					nStr := fmt.Sprintf("* %s", n.Name)
					if n.Temporary {
						validTo := n.CreatedAt.Add(24 * time.Hour)
						dur := validTo.Sub(time.Now())
						nStr += fmt.Sprintf(" (valid for %s)", dur)
					}

					fmt.Println(nStr)
				}
			},
		},
		{
			Name:  "share",
			Usage: "Share a note (publicly) on hastebin.com",
			Action: func(c *cli.Context) {
				nName, err := NoteName(c)
				if err != nil {
					fmt.Println(err)
					return
				}

				n, err := i.LoadNote(nName)
				if err != nil {
					fmt.Println(err)
					return
				}

				url, err := reader.ShareNote(n)

				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(fmt.Sprintf("Shared note '%s': %s", n.Name, url))
			},
		},
		{
			Name:    "keep",
			Aliases: []string{"k"},
			Usage:   "Sets a temporary note as permanent",
			Action: func(c *cli.Context) {
				nName, err := NoteName(c)
				if err != nil {
					fmt.Println(err)
					return
				}

				n, err := i.LoadNote(nName)
				if err != nil {
					fmt.Println(err)
					return
				}

				n.Temporary = false
				if err := i.SaveNote(&n); err != nil {
					fmt.Println(err)
				}

				fmt.Println(fmt.Sprintf("Keeping note: %s", n.Name))
			},
		},
	}

	app.Run(os.Args)

}

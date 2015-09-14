package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gummiboll/forgetful/reader"
	"github.com/gummiboll/forgetful/storage"
	"github.com/gummiboll/forgetful/writer"
)

const version string = "0.9"

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
			},
			Action: func(c *cli.Context) {
				if c.Args().Present() != true {
					fmt.Println("Missing argument: name")
					return
				}
				nName := strings.Join(c.Args(), " ")
				if exists := i.NoteExists(nName); exists == true {
					fmt.Println("Note already exists")
					return
				}
				n := storage.Note{Name: nName, Temporary: c.Bool("t")}
				if err := writer.WriteNote(&n); err != nil {
					fmt.Println(err)
					return
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
				if c.Args().Present() != true {
					fmt.Println("Missing name")
					return
				}
				nName := strings.Join(c.Args(), " ")
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
				if c.Args().Present() != true {
					fmt.Println("Missing argument: name")
					return
				}
				nName := strings.Join(c.Args(), " ")
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
				if c.Args().Present() != true {
					fmt.Println("Missing argument: name")
					return
				}
				nName := strings.Join(c.Args(), " ")
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
				if c.Args().Present() != true {
					fmt.Println("Nothing to search for")
					return
				}
				nName := strings.Join(c.Args(), " ")
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
			Name:    "keep",
			Aliases: []string{"k"},
			Usage:   "Sets a temporary note as permanent",
			Action: func(c *cli.Context) {
				if c.Args().Present() != true {
					fmt.Println("Missing argument: name")
					return
				}
				nName := strings.Join(c.Args(), " ")
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

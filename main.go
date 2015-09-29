package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gummiboll/forgetful/commands"
	"github.com/gummiboll/forgetful/storage"
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
				cli.BoolFlag{
					Name:  "p",
					Usage: "Create note with contents from clipboard",
				},
			},
			Action: func(c *cli.Context) {
				n, err := commands.AddCommand(c, i)
				if err != nil {
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
				n, err := commands.DeleteCommand(c, i)
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(fmt.Sprintf("Deleted note: %s", n.Name))
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "Edit/read a note",
			Action: func(c *cli.Context) {
				n, err := commands.EditCommand(c, i)
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(fmt.Sprintf("Updated note: %s", n.Name))
			},
		},
		{
			Name:    "read",
			Aliases: []string{"r"},
			Usage:   "Read a note",
			Action: func(c *cli.Context) {
				if err := commands.ReadCommand(c, i); err != nil {
					fmt.Println(err)
				}
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all notes, filter result if argument i present",
			Action: func(c *cli.Context) {
				notes := commands.ListCommand(c, i)

				if len(notes) > 0 {
					fmt.Println("Matching notes:")
				} else {
					fmt.Println("No matches")
				}

				for _, n := range notes {
					fmt.Println(n)
				}
			},
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search notes for argument",
			Action: func(c *cli.Context) {
				notes, err := commands.SearchCommand(c, i)
				if err != nil {
					fmt.Println(err)
					return
				}

				if len(notes) > 0 {
					fmt.Println("Matching notes:")
				} else {
					fmt.Println("No matches")
				}

				for _, n := range notes {
					fmt.Println(n)
				}
			},
		},
		{
			Name:  "share",
			Usage: "Share a note (publicly) on hastebin.com",
			Action: func(c *cli.Context) {
				n, url, err := commands.ShareCommand(c, i)
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
				n, err := commands.KeepCommand(c, i)
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println(fmt.Sprintf("Keeping note: %s", n.Name))
			},
		},
	}

	app.Run(os.Args)

}

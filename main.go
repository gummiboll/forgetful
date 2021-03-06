package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gummiboll/forgetful/commands"
	"github.com/gummiboll/forgetful/storage"
	"github.com/urfave/cli"
)

const version string = "1.1"

func printList(notes []string) {
	if len(notes) > 0 {
		fmt.Println(fmt.Sprintf("Found %d matching note(s):", len(notes)))

		for _, n := range notes {
			fmt.Println(n)
		}

	} else {
		fmt.Println("No matches")
	}
}

func printInfo(n storage.Note) {
	fmt.Println(fmt.Sprintf("%s:", n.Name))
	fmt.Println(fmt.Sprintf("Created at: %s, last updated at: %s", n.CreatedAt.Format("2006-01-02 15:04"), n.UpdatedAt.Format("2006-01-02 15:04")))
	if n.Temporary == true {
		validTo := n.UpdatedAt.Add(24 * time.Hour)
		dur := validTo.Sub(time.Now())
		fmt.Println(fmt.Sprintf("Is temporary, will expire in %s", dur-(dur%time.Second)))
	}
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
			Action: func(c *cli.Context) error {
				n, err := commands.AddCommand(c, i)
				if err != nil {
					return err
				}

				fmt.Println(fmt.Sprintf("Added note: %s", n.Name))
				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "Delete a note",
			Action: func(c *cli.Context) error {
				n, err := commands.DeleteCommand(c, i)
				if err != nil {
					return err
				}

				fmt.Println(fmt.Sprintf("Deleted note: %s", n.Name))
				return nil
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "Edit/read a note",
			Action: func(c *cli.Context) error {
				n, err := commands.EditCommand(c, i)
				if err != nil {
					return err
				}

				fmt.Println(fmt.Sprintf("Updated note: %s", n.Name))
				return nil
			},
		},
		{
			Name:    "info",
			Aliases: []string{"i"},
			Usage:   "Prints information about a note",
			Action: func(c *cli.Context) error {
				n, err := commands.InfoCommand(c, i)
				if err != nil {
					return err
				}

				printInfo(n)

				return nil
			},
		},
		{
			Name:    "read",
			Aliases: []string{"r"},
			Usage:   "Read a note",
			Action: func(c *cli.Context) error {
				if err := commands.ReadCommand(c, i); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:    "rename",
			Aliases: []string{"mv"},
			Usage:   "Rename a note",
			Action: func(c *cli.Context) error {
				nName, newName, err := commands.RenameCommand(c, i)
				if err != nil {
					return err
				}

				fmt.Println(fmt.Sprintf("Renamed %s to %s", nName, newName))

				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all notes, filter result if argument i present",
			Action: func(c *cli.Context) error {
				notes := commands.ListCommand(c, i)

				printList(notes)

				return nil
			},
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search notes for argument",
			Action: func(c *cli.Context) error {
				notes, err := commands.SearchCommand(c, i)
				if err != nil {
					return err
				}

				printList(notes)

				return nil
			},
		},
		{
			Name:  "share",
			Usage: "Share a note (publicly) on hastebin.com",
			Action: func(c *cli.Context) error {
				n, url, err := commands.ShareCommand(c, i)
				if err != nil {
					return err
				}

				fmt.Println(fmt.Sprintf("Shared note '%s': %s", n.Name, url))
				return nil
			},
		},
		{
			Name:    "keep",
			Aliases: []string{"k"},
			Usage:   "Sets a temporary note as permanent",
			Action: func(c *cli.Context) error {
				n, err := commands.KeepCommand(c, i, true)
				if err != nil {
					return err
				}

				fmt.Println(fmt.Sprintf("Keeping note: %s", n.Name))
				return nil
			},
		},
		{
			Name:    "unkeep",
			Aliases: []string{"u"},
			Usage:   "Sets a permanent note as temporary",
			Action: func(c *cli.Context) error {
				n, err := commands.KeepCommand(c, i, false)
				if err != nil {
					return err
				}

				fmt.Println(fmt.Sprintf("Unkeeping note: %s", n.Name))
				return nil
			},
		},
	}

	app.Run(os.Args)

}

forgetful
=======

`forgetful` is a small command line tool for your notes/cheat sheets. Inspiered by [Chris Lane's cheatsheet](https://github.com/chrisallenlane/cheat) but wanted something a little different and something that can automatically delete temporary notes that are old since I tend to create lots of textfiles with temporary notes from phone calls and whatnot.

# Install/setup
`go get github.com/gummiboll/forgetful` and copy bin/forgetful to your $PATH.

forgetful creates ~/.forgetful/forgetful.db (sqllite3) on the first run. Editor defaults to `vim` if $EDITOR is not set.

# Usage
`forgetful add (or 'a' for short) <note name>` - opens a editor with a blank note. If the flag `-t` is used the note is marked as temporary and will expire in 24 hours.

`forgetful edit (or 'e' for short) <note name>` - opens a editor with specified note.

`forgetful delete (or 'd' for short) <note name>` - deletes a note.

`forgetful list (or 'l' for short) <optional phrase>` - list notes, if phrase is specified filter results, e.g: `list foo` matches foobar, bigfoot and foo.

`forgetful search (or 's' for short) <phrase>` - searches notes for phrase.

`forgetful keep (or 'k' for short) <note name>` - mark a temporary note as permanent.

`forgetful help (or 'h' for short)` - show help

# Todo

 - Create note from copy/paste-buffer
 - ..?

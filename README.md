[![Build Status](https://travis-ci.org/gummiboll/forgetful.svg?branch=master)](https://travis-ci.org/gummiboll/forgetful)

forgetful
=======

`forgetful` is a small command line tool for your notes/cheat sheets. Inspiered by [Chris Lane's cheatsheet](https://github.com/chrisallenlane/cheat) but wanted something a little different and something that can automatically delete temporary notes that are old since I tend to create lots of textfiles with temporary notes from phone calls and whatnot.

# Install/setup
`go get github.com/gummiboll/forgetful` and copy bin/forgetful to your $PATH. Or download the [latest binary release](https://github.com/gummiboll/forgetful/releases/latest)

forgetful creates ~/.forgetful/forgetful.db (sqllite3) on the first run. Editor defaults to `vim` if $EDITOR is not set.

# Usage
command|short|description|optional arguments
-------|-----|-----------|------------------
add <name>|a|Adds a note|-p (from contents of clipboard) -t (temporary note, expires afer 24h)
edit <name>|e|Edit a note
delete <name>|d|Deletes a note
read <name>|r|Reads a note
list|l|List notes|If phrase is specified, filter results on that phrase
search|s|Search for notes
share||Shares a note publicaly on [hastebin](http://hastebin.com)
keep <name>|k|Mark a temporary note as permanent
help|h|Shows help

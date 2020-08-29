package main

import (
	"CopyPasteOnline/ports/app"
	"CopyPasteOnline/ports/mnemonics"
	"CopyPasteOnline/ports/sqlite"
)

func main() {
	// # Setup db manager
	db, err := sqlite.New("./data.sqlite")
	if err != nil {
		panic("unable to setup db: " + err.Error())
	}

	// # Setup mnemonics handler Will panic on error.
	mnemonics := mnemonics.New()

	// # Stat server.
	app.Start(db, mnemonics)
}

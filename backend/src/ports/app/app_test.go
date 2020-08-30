package app

import (
	"copypaste-api/ports/mnemonics"
	"copypaste-api/ports/sqlite"
	"os"
	"testing"
)

func tasker(t *testing.T, task func(*sqlite.SQLiteManager, *mnemonics.Poolhandler)) {
	path := "item.sqlite"
	dbManager, err := sqlite.New(path)
	mnemonicManager := mnemonics.New()
	if err != nil {
		t.Error("DB creation failure:", err)
		return
	}
	defer os.Remove(path)
	if err := dbManager.CreateItemTable(); err != nil {
		t.Error(err)
	}
	task(dbManager, mnemonicManager)
}

func TestSetup(t *testing.T) {
	// # Use with curl
	tasker(t, func(db *sqlite.SQLiteManager, mnemonics *mnemonics.Poolhandler) {
		Start(db, mnemonics)
	})
}

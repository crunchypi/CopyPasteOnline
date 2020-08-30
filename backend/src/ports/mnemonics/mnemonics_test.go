package mnemonics

import (
	"copypaste-api/ports/sqlite"
	"os"
	"testing"
)

func TestLoadFogleman(t *testing.T) {
	c := LoadAll()
	t.Log(c[0], "...", c[len(c)-1], "Total:", len(c))
}

func TestDrawEnsured(t *testing.T) {
	path := "mnemonic.sqlite"
	dbHandler, err := sqlite.New(path)
	if err != nil {
		t.Error(err)
	}
	dbHandler.CreateItemTable()
	defer os.Remove(path)

	// # Arbitrary values below.
	// # if defaultDrawN and defaultDrawTryLimit are low,
	// # then this should fail often. Converse should be true.
	mHandler := New()
	length := len(mHandler.corpus)
	defaultDrawN = 3         // # package var
	defaultDrawTryLimit = 10 // # package var
	for i := 0; i < length; i++ {
		menomic, ok := mHandler.DrawEnsured(dbHandler)
		if !ok {
			t.Error("Failed to draw on iter:", i)
		}
		dbHandler.CreateItem(menomic, []byte("testdata"))
	}
}

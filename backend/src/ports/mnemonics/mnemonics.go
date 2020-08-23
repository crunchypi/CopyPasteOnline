package mnemonics

import (
	"CopyPasteOnline/ports/sqlite"
	"fmt"
	"math/rand"
	"time"
)

var (
	defaultDrawN        = 3
	defaultDrawTryLimit = 10
)

// Poolhandler keeps track of a corpus pool consisting of
// mnemonic words and draws them in a ensured way.
type Poolhandler struct {
	corpus []string
}

// New returns a Poolhandler with a loaded corpus.
func New() *Poolhandler {
	return &Poolhandler{corpus: LoadAll()}
}

// draw draws defaultDrawN (package var) random mnemonic
// words from corpus and returns that as a string where
// each word is separated by a space (except first/last)
func (p *Poolhandler) draw() string {
	res := ""
	for i := 0; i < defaultDrawN; i++ {
		// # Choose random mnemonic.
		rand.Seed(time.Now().UnixNano())
		choose := rand.Intn(len(p.corpus)-0) + 0
		res += p.corpus[choose]
		// # Add trailing space, except for last.
		if i < defaultDrawN-1 {
			res += " "
		}
	}
	return res
}

// DrawEnsured draws n defaultDrawN (package var) mnemonic words while
// ensuring that a combination does not already exist in the database.
// Infinite loop prevented by defaultDrawTryLimit (package var), as that
// specifies the maximum amount of re-tries. Result is a string where
// each word is separated by a space(except to leading or trailing)
func (p *Poolhandler) DrawEnsured(s *sqlite.SQLiteManager) (string, bool) {
	if defaultDrawN < 1 {
		panic(fmt.Sprintf(
			"Implementation err. n must be > 1. Actual length:%d", defaultDrawN,
		))
	}

	for tryTimes := 0; tryTimes < defaultDrawTryLimit; tryTimes++ {
		res := p.draw()
		// # Check if mnemonic is taken.
		exists, err := s.ReadMnemonicExists(res)
		if err != nil {
			// # Not recoverable, has to be ensured during development.
			panic("Implementation issue:" + err.Error())
		}
		if !exists {
			return res, true
		}

	}
	return "", false
}

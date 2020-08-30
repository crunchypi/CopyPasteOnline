package app

import (
	"copypaste-api/ports/mnemonics"
	"copypaste-api/ports/sqlite"
	"io/ioutil"
	"log"
	"net/http"
)

// handler serves as a bridge between the app and
// other packages, mainly db and mnemonics
type handler struct {
	db        *sqlite.SQLiteManager
	mnemonics *mnemonics.Poolhandler
}

// copyData handles incoming data and returns mnemonics.
func (h *handler) copyData(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // # middleware err handled.

	mnemonic, ok := h.mnemonics.DrawEnsured(h.db)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err := h.db.CreateItem(mnemonic, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mnemonic))
}

// pasteData handles incoming mnemonics and returns stored data.
func (h *handler) pasteData(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // # middleware err handled.

	data, ok, err := h.db.ReadItemByMnemonic(string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data.Data)
}

// setRoutes sets up routes for this app.
func (h *handler) setRoutes() {
	http.Handle("/copy", midDOS(midBodyErr(http.HandlerFunc(h.copyData))))
	http.Handle("/paste", midDOS(midBodyErr(http.HandlerFunc(h.pasteData))))
}

// Start starts the app.
func Start(db *sqlite.SQLiteManager, mnemonics *mnemonics.Poolhandler) {
	db.CreateItemTable()
	handler := handler{db: db, mnemonics: mnemonics}

	handler.setRoutes()
	log.Println(http.ListenAndServe(":80", nil))
}

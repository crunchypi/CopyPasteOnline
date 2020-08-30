package app

// # TLS implemented as described in http://www.inanzzz.com/index.php/post/9ats/http2-and-tls-client-and-server-example-with-golang

import (
	"CopyPasteOnline/ports/mnemonics"
	"CopyPasteOnline/ports/sqlite"
	"io/ioutil"
	"log"
	"net/http"
)

type handler struct {
	db        *sqlite.SQLiteManager
	mnemonics *mnemonics.Poolhandler
}

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

func (h *handler) setRoutes() {
	http.Handle("/copy", midDOS(midBodyErr(http.HandlerFunc(h.copyData))))
	http.Handle("/paste", midDOS(midBodyErr(http.HandlerFunc(h.pasteData))))
}
func Start(db *sqlite.SQLiteManager, mnemonics *mnemonics.Poolhandler) {
	db.CreateItemTable()
	handler := handler{db: db, mnemonics: mnemonics}

	handler.setRoutes()

	if err := server().ListenAndServeTLS("", ""); err != nil {
		log.Fatal("shutdown:", err)
	}
}

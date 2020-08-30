package app

import (
	"io/ioutil"
	"net/http"
)

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

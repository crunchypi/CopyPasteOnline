package app

import (
	"copypaste-api/ports/app/dosguard"
	"io/ioutil"
	"net/http"
	"strings"
)

func midDOS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// # Only identify xx.xx.xx.xx:xxxx format (local can be [::1]:xxxx)
		ipport := strings.Split(r.RemoteAddr, ":")
		// # Guard indexing crash.
		if len(ipport) < 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ip := ipport[0]
		// # Control is a package var in doscheck which handles ip registration.
		// # Guard too many access attempts.
		if ok := dosguard.Control.RegisterCheck(ip); !ok {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func midBodyErr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// # Validate that body can be read.
		if _, err := ioutil.ReadAll(r.Body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

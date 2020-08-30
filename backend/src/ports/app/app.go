package app

import (
	"copypaste-api/config"
	"copypaste-api/ports/mnemonics"
	"copypaste-api/ports/sqlite"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ip   = config.IP
	port = config.Port
	// # Server IO time limitation.
	readTimeout  = config.ReadTimeout
	writeTimeout = config.WriteTimeout
	// # Path to index.html. Disabled if empty.
	staticFilePath = config.StaticFilePath
)

//
var (
	tlsEnabled         = config.TLSEnabled
	serverName         = config.ServerName
	certPublicKeyPath  = config.CertPublicKeyPath
	certPrivateKeyPath = config.CertPrivateKeyPath
)

// handler serves as a bridge between the app and
// other packages, mainly db and mnemonics
type handler struct {
	db        *sqlite.SQLiteManager
	mnemonics *mnemonics.Poolhandler
}

// setRoutes sets up routes for this app.
func (h *handler) setRoutes() {
	http.Handle("/copy", midDOS(midBodyErr(http.HandlerFunc(h.copyData))))
	http.Handle("/paste", midDOS(midBodyErr(http.HandlerFunc(h.pasteData))))

	if staticFilePath != "" {
		http.Handle("/", midDOS(http.FileServer(http.Dir(staticFilePath))))
	}
}

// tlsConfig is only called if tlsEnabled=true (pkg var)
func tlsConfig() *tls.Config {
	crt, err := ioutil.ReadFile(certPublicKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	key, err := ioutil.ReadFile(certPrivateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   serverName,
	}
}

// Start starts the app.
func Start(db *sqlite.SQLiteManager, mnemonics *mnemonics.Poolhandler) {
	// # Enable interface to other ports of this microservice.
	db.CreateItemTable()
	handler := handler{db: db, mnemonics: mnemonics}
	handler.setRoutes()

	// # Server configs.
	server := http.Server{
		Addr:         ip + ":" + port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	// # TLS on/off.
	switch tlsEnabled {
	case true:
		server.TLSConfig = tlsConfig()
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatal("shutdown:", err)
		}
	case false:
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("shutdown:", err)
		}
	}
}

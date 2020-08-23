package app

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	serverName         = "localhost"
	serverPort         = "8080"
	certPublicKeyPath  = "./cert/public.crt"
	certPrivateKeyPath = "./cert/private.key"

	readTimeout  = 3 * time.Second
	writeTimeout = 3 * time.Second
)

func server() *http.Server {
	return &http.Server{
		Addr:         ":" + serverPort,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		TLSConfig:    tlsConfig(),
	}
}

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

package main

import (
	"net/http"
	"time"

	"github.com/ministryofjustice/cloud-platform-label-pods/init_app"
)

func main() {
	ginMode := init_app.InitEnvVars()

	r := init_app.InitGin(ginMode)

	server := &http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	// to run this locally provide a self signed cert
	server.ListenAndServeTLS("/app/certs/tls.crt", "/app/certs/tls.key")
}

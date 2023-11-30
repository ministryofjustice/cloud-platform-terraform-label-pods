package main

import (
	"log"

	"github.com/ministryofjustice/cloud-platform-label-pods/init_app"
)

func main() {
	ginMode := init_app.InitEnvVars()

	r := init_app.InitGin(ginMode)

	// to run this locally provide a self signed cert
	err := r.RunTLS(":3000", "/app/certs/tls.crt", "/app/certs/tls.key")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

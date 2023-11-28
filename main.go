package main

import (
	"log"

	"github.com/ministryofjustice/cloud-platform-label-pods/init_app"
)

func main() {
	ginMode := init_app.InitEnvVars()

	r := init_app.InitGin(ginMode)

	err := r.RunTLS(":3000", "/etc/ssl/certs/tls.crt", "/etc/ssl/certs/tls.key")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

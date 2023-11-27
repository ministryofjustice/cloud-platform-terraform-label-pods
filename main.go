package main

import (
	"log"

	"github.com/ministryofjustice/cloud-platform-label-pods/init_app"
)

func main() {
	ginMode := init_app.InitEnvVars()

	r := init_app.InitGin(ginMode)

	err := r.Run(":3000")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

package main

import (
	"log"

	"github.com/33-token/model-api-platform/services/controlplane/internal/app"
)

func main() {
	server := app.NewServer()

	log.Printf("controlplane api listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

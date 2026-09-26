// Command server runs the toolshed reference API on top of JSON-file
// storage. It is deliberately tiny: its only job is wiring a storage
// backend into a lending.Service into an HTTP handler. Swapping storage
// backends means changing the two lines that construct resourceStore and
// loanStore — nothing else in this file, or in internal/lending or
// internal/api, needs to change.
package main

import (
	"log"
	"net/http"
	"os"

	"toolshed/internal/api"
	"toolshed/internal/lending"
	"toolshed/internal/storage/jsonfile"
)

func main() {
	dataDir := os.Getenv("TOOLSHED_DATA_DIR")
	if dataDir == "" {
		dataDir = "."
	}
	addr := os.Getenv("TOOLSHED_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	resources := jsonfile.NewResourceStore(dataDir + "/resources.json")
	loans := jsonfile.NewLoanStore(dataDir + "/loans.json")
	service := lending.NewService(resources, loans)
	handler := api.NewHandler(service)

	log.Printf("toolshed listening on %s (data dir: %s)", addr, dataDir)
	if err := http.ListenAndServe(addr, handler.Routes()); err != nil {
		log.Fatal(err)
	}
}

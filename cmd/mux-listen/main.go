package main

import (
	"log"

	"github.com/codycollier/hum-mux/pkg/mux"
)

func main() {

	log.Printf("[ml] Prepare to listen...")

	// Setup the client
	addr := "127.0.0.1:8888"
	cl, conn := mux.GetNewMuxClient(addr)
	defer conn.Close()

	// Listen (blocks)
	mux.Listen(cl)

	log.Printf("[ml] Stream closed")

}

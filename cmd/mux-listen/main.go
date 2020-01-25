package main

import (
	"hum/pkg/mux"
	"log"
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

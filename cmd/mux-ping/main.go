package main

import (
	"hum/pkg/mux"
	"log"
)

func main() {

	addr := "127.0.0.1:8888"

	log.Printf("[mp] Creating mux client")
	cl, conn := mux.GetNewMuxClient(addr)
	defer conn.Close()

	log.Printf("[mp] Sending pings")
	mux.Ping(cl)
	mux.Ping(cl)
	mux.Ping(cl)

	log.Printf("[mp] Done")

}

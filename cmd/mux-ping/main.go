package main

import (
	"hum/pkg/mux"
	"log"
)

func main() {

	addr := "127.0.0.1:8888"

	log.Printf("[mtx] Creating mux client")
	cl, conn := mux.GetNewMuxClient(addr)
	defer conn.Close()

	mux.CallPing(cl)
	mux.CallPing(cl)
	mux.CallPing(cl)
	log.Printf("[mtc] Done")

}

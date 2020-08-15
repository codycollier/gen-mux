package main

import (
	"hum/pkg/mux"
	"log"
	"time"
)

func main() {

	addr := "127.0.0.1:8888"

	log.Printf("[mp] Creating mux client")
	cl, conn := mux.GetNewMuxClient(addr)
	defer conn.Close()

	log.Printf("[mp] Starting ping loop")
	for {
		mux.Ping(cl)
		time.Sleep(time.Millisecond * 500)
	}

}

package mux

import (
	pb "hum/proto"
	"log"
	"math/rand"
)

// Initialize the core for a muxServer instance
func initMux(ms *muxServer) error {

	log.Printf("[muxd] Initializing Mux structures")
	rand.Seed(14028)

	// Initialize the input / output channels
	ms.mux_in = make(chan pb.Datum)
	ms.mux_out = make(map[int64]chan pb.Datum)
	go multiplexer(ms)

	return nil
}

// Muxer - any msg that comes in, send it out on every outbound channel
func multiplexer(ms *muxServer) {
	for {
		msg := <-ms.mux_in
		for out_id, out_chan := range ms.mux_out {
			log.Printf("[muxd] copying message to listener: %v", out_id)
			// TODO(cmc) - race condition with delete(ms.mux_out, my_id) in server.go
			out_chan <- msg
		}
	}
}

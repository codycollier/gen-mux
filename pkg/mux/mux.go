package mux

import (
	pb "hum/proto"
	"log"
	"math/rand"
)

// Initialize the core for a muxServer instance
func initMux(mux *muxServer) error {

	log.Printf("[muxd] Initializing Mux structures")
	rand.Seed(14028)

	// Initialize the input / output channels
	mux.mux_in = make(chan pb.Datum)
	mux.mux_out = make(map[int64]chan pb.Datum)
	go muxer(mux)

	return nil
}

// Muxer - any msg that comes in, send it out on every outbound channel
func muxer(mux *muxServer) {
	for {
		msg := <-mux.mux_in
		for out_id, out_chan := range mux.mux_out {
			log.Printf("[muxd] copying message to listener: %v", out_id)
			out_chan <- msg
		}
	}

}

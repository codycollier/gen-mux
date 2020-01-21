package mux

import (
	pb "hum/proto"
	"log"
)

// Initialize the core for a muxServer instance
func initMux(mux *muxServer) error {

	log.Printf("[muxd] Initializing Mux structures")

	// Initialize the input / output channels
	mux.mux_in = make(chan pb.Datum)
	mux.mux_out = make(chan pb.Datum)

	return nil
}

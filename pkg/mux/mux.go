package mux

import (
	"log"
	"math/rand"

	pb "github.com/codycollier/whisper-mux/proto"
)

// Initialize the core for a muxServer instance
func initMux(ms *muxServer) error {

	rand.Seed(14028)

	// Setup the input channel and the listeners map
	log.Printf("[muxd] Initializing mux structures")
	ms.mux_in = make(chan pb.Datum)
	ms.mux_listeners = make(map[int64]muxListener)

	// Start the mux
	log.Printf("[muxd] Starting multiplexer")
	go multiplexer(ms)

	return nil
}

// Multiplexer - send any input message to all available listeners
func multiplexer(ms *muxServer) {

	// Start mux loop
	for {

		// Block, waiting for input
		msg := <-ms.mux_in

		// Send it to every listener
		for _, listener := range ms.mux_listeners {

			select {

			// Don't send if the client has closed
			case <-listener.stream.Context().Done():
				log.Printf("[muxd] mux: Removing listener [lid:%v]", listener.id)
				// Remove the listener from the mux
				delete(ms.mux_listeners, listener.id)

			// Otherwise, send the message out
			default:
				// log.Printf("[muxd] mux: sending message to listener: %v", listener.id)
				listener.input <- msg
			}

			// next listener
		}

		// done processing message. loop around and wait for another.
	}
}

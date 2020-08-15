package main

import (
	"log"

	"github.com/golang/protobuf/ptypes"

	"github.com/codycollier/hum-mux/pkg/mux"
	pb "github.com/codycollier/hum-mux/proto"
)

func main() {

	addr := "127.0.0.1:8888"

	log.Printf("[mst] Creating mux client")
	cl, conn := mux.GetNewMuxClient(addr)
	defer conn.Close()

	// Create input channel and done-signal channel
	input := make(chan pb.Datum)
	done := make(chan int)
	go mux.Inject(cl, input, done)

	// Send test messages
	// for i := 0; i < 1; i++ {
	for i := 0; i < 1000000; i++ {
		input <- pb.Datum{Data: "foo", UtsEvent: ptypes.TimestampNow()}
	}

	// Close the input to signal no more messages to client
	close(input)

	// Wait for the client to finish up the stream
	<-done
	log.Printf("[mst] Done")

}

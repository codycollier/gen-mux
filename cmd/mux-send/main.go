package main

import (
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/codycollier/gen-mux/pkg/mux"
	pb "github.com/codycollier/gen-mux/proto"
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
	for {
		input <- pb.Datum{Data: "foo", UtsEvent: ptypes.TimestampNow()}
		time.Sleep(800 * time.Millisecond)
	}

}

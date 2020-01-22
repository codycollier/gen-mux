package mux

import (
	"context"
	"google.golang.org/grpc"
	pb "hum/proto"
	"io"
	"log"
	"math/rand"
	"net"
)

// Core structure of the mux server
type muxServer struct {
	mux_in  chan pb.Datum
	mux_out map[int64]chan pb.Datum
}

// Accept incoming messages/events from the hum
// Handler for:
//  rpc Inject (stream SendRequest) returns (SendResponse);
func (ms *muxServer) Inject(stream pb.Mux_InjectServer) error {

	// TODO(cmc): add a timeout in listening for inject?
	// TODO(cmc): ...

	count := int32(0)
	for {

		// Block and wait for an input
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("[muxd] Error receiving on Inject stream: %v", err)
			break
		}
		log.Printf("[muxd] Inject recv: %v", *req)

		// Push the message to the mux
		ms.mux_in <- *req.Datum

		count += 1

	}
	resp := &pb.InjectResponse{MsgCount: count}
	log.Printf("[muxd] Inject send close: %v", resp)
	stream.SendAndClose(resp)
	return nil
}

// Stream the hum out to a client
// Handler for:
//  rpc Listen (ListenRequest) returns (stream ListenResponse);
func (ms *muxServer) Listen(req *pb.ListenRequest, stream pb.Mux_ListenServer) error {

	// TODO(cmc): add support for filters
	// TODO(cmc): ...

	// Create and add a new listener channel to the mux
	my_id := rand.Int63()
	my_listener := make(chan pb.Datum)
	ms.mux_out[my_id] = my_listener

	log.Printf("[muxd] Listen recv: %v [lid:%v]", req, my_id)

	// Start the listener
listen:
	for {
		select {

		// Stop listening if the client has closed
		case <-stream.Context().Done():
			log.Printf("[muxd] Listen client done [lid:%v]", my_id)
			break listen

		// Listen for messages from mux listener and send to client
		case msg := <-my_listener:
			resp := &pb.ListenResponse{Datum: &msg}
			log.Printf("[muxd] Listen send: %v [lid:%v]", resp, my_id)
			if err := stream.Send(resp); err != nil {
				log.Printf("[muxd] Listen send err: %v [lid:%v]", err, my_id)
				break listen
			}
		}
	}

	// cleanup
	log.Printf("[muxd] Listen cleanup [lid:%v]", my_id)
	delete(ms.mux_out, my_id)
	close(my_listener)

	return nil
}

// Ping debugging endpoint
// Handler for:
//  rpc Ping (PingRequest) returns (PingResponse);
func (ms *muxServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {

	log.Printf("[muxd] Ping recv: %v", req)

	resp := &pb.PingResponse{Pong: true}
	log.Printf("[muxd] Ping send: %v", resp)

	return resp, nil
}

// Server initialization and start up
func StartMuxServer(addr string) {

	// TODO(cmc): Add support for ssl
	// TODO(cmc): ...

	// Initialize mux
	mux_server := &muxServer{}
	initMux(mux_server)

	// Setup gRPC server & register service
	log.Printf("[muxd] Setting up gRPC service")
	log.Printf("[muxd] Will be listening on: %v", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[muxd] Error listening on %s", addr)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMuxServer(grpcServer, mux_server)

	// Listen forever
	log.Printf("[muxd] Starting server...")
	grpcServer.Serve(listener)

}

package mux

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net"

	"google.golang.org/grpc"

	pb "github.com/codycollier/gen-mux/proto"
)

// Core structure of the mux server, multiplexing input to many listeners
type muxServer struct {
	mux_in        chan pb.Datum
	mux_listeners map[int64]muxListener
}

// Definition of a single listener
type muxListener struct {
	id     int64
	input  chan pb.Datum
	stream pb.Mux_ListenServer
}

// Accept incoming messages and inject
// Handler for:
//  rpc Inject (stream SendRequest) returns (SendResponse);
func (ms *muxServer) Inject(stream pb.Mux_InjectServer) error {

	// TODO(cmc): add a timeout in listening for inject?
	// TODO(cmc): ...

	count := int32(0)
	for {

		// Block and wait for an input
		req, err := stream.Recv()

		// Graceful close
		if err == io.EOF {
			resp := &pb.InjectResponse{MsgCount: count}
			log.Printf("[muxd] Inject send close: %v", resp)
			stream.SendAndClose(resp)
			break
		}

		// Unexpected close
		if err != nil {
			log.Printf("[muxd] Error receiving on Inject stream: %v", err)
			break
		}

		// Push the message to the mux
		log.Printf("[muxd] Inject recv: %v", *req)
		ms.mux_in <- *req.Datum

		count += 1

	}
	return nil
}

// Stream messages out to a client
// Handler for:
//  rpc Listen (ListenRequest) returns (stream ListenResponse);
func (ms *muxServer) Listen(req *pb.ListenRequest, stream pb.Mux_ListenServer) error {

	// TODO(cmc): add support for filters
	// TODO(cmc): ...

	// Create a new listener for self
	listener := muxListener{}
	listener.id = rand.Int63()
	listener.input = make(chan pb.Datum)
	listener.stream = stream

	// Add the listener to the mux
	ms.mux_listeners[listener.id] = listener

	// Start listening
	log.Printf("[muxd] Listen recv: %v [lid:%v]", req, listener.id)
	for {

		// Block until listener recieves a message
		msg := <-listener.input
		resp := &pb.ListenResponse{Datum: &msg}

		// Send the message out
		log.Printf("[muxd] Listen send: %v [lid:%v]", resp, listener.id)
		if err := stream.Send(resp); err != nil {
			log.Printf("[muxd] Listen send err: %v [lid:%v]", err, listener.id)
			return nil
		}
	}
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

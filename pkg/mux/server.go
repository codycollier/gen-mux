package mux

import (
	"context"
	"google.golang.org/grpc"
	pb "hum/proto"
	"log"
	"net"
)

// Core structure of the mux server
type muxServer struct {
	mux_in  chan pb.Datum
	mux_out chan pb.Datum
}

// Accept incoming messages/events from the hum
// Handler for:
//  rpc Inject (stream SendRequest) returns (SendResponse);
func (s *muxServer) Inject(stream pb.Mux_InjectServer) error {

	// TODO(cmc): find the mux and push messages

	for {

		req, err := stream.Recv()
		if err != nil {
			log.Printf("[muxd] Error receiving on Inject stream: %v", err)
		}

		log.Printf("[muxd] Inject recv: %v", req)
		// resp := &pb.InjectResponse{}
		// log.Printf("[muxd] Inject send: %v", resp)
		// stream.SendAndClose(resp)

	}
	return nil
}

// Stream the hum out to a client
// Handler for:
//  rpc Listen (ListenRequest) returns (stream ListenResponse);
func (s *muxServer) Listen(req *pb.ListenRequest, stream pb.Mux_ListenServer) error {

	log.Printf("[muxd] Listen: recv: %v", req)

	// TODO(cmc): accept msg (with filters?) and start streaming from mux
	// req

	for {
		resp := &pb.ListenResponse{}
		log.Printf("[muxd] Listen: send: %v", resp)
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil

}

// Ping debugging endpoint
// Handler for:
//  rpc Ping (PingRequest) returns (PingResponse);
func (s *muxServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {

	log.Printf("[muxd] Ping: recv: %v", req)

	resp := &pb.PingResponse{Pong: true}
	log.Printf("[muxd] Ping: send: %v", resp)

	return resp, nil
}

// Server initialization and start up
func StartMuxServer(addr string) {

	// Initialize mux
	mux_server := &muxServer{}
	initMux(mux_server)

	// Setup gRPC server & register service
	// TODO(cmc): Add support for ssl
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

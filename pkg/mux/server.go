package mux

import (
	"context"
	"google.golang.org/grpc"
	pb "hum/proto"
	"log"
	"net"
)

// Core structure of the mux server
type muxServer struct{}

// Handler for Mux.Send()
func (s *muxServer) Send(ctx context.Context, req *pb.SendRequest) (*pb.SendResponse, error) {

	log.Printf("[muxd] Send: received: %v", req)

	// TODO(cmc): accept msg and send to mux core
	resp := &pb.SendResponse{}

	log.Printf("[muxd] Send: sending: %v", resp)
	return resp, nil
}

// Handler for Mux.Listen()
func (s *muxServer) Listen(ctx context.Context, req *pb.ListenRequest) (*pb.ListenResponse, error) {

	log.Printf("[muxd] Listen: received: %v", req)

	// TODO(cmc): accept msg and start streaming from mux
	resp := &pb.ListenResponse{}

	log.Printf("[muxd] Listen: sending: %v", resp)
	return resp, nil
}

// Handler for Mux.Ping() handler
func (s *muxServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	log.Printf("[muxd] Ping: received: %v", req)

	resp := &pb.PingResponse{Pong: true}

	log.Printf("[muxd] Ping: sending: %v", resp)
	return resp, nil
}

// Server initialization and start up
func StartMuxServer(addr string) {

	// Setup gRPC server & register service
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[muxd] Error listening on %s", addr)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMuxServer(grpcServer, &muxServer{})

	// Listen forever
	grpcServer.Serve(listener)

}

package mux

import (
	"context"
	"google.golang.org/grpc"
	pb "hum/proto"
	"io"
	"log"
	"time"
)

// Call Mux and inject one or more events
func CallInject(cl pb.MuxClient) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Printf("[muxc] Starting Inject stream")
	stream, err := cl.Inject(ctx)
	if err != nil {
		log.Printf("[muxc] Error creating Inject stream: %v", err)
	}

	// TODO(cmc) - take an iterable and stream it up

	for {
		req := &pb.InjectRequest{}
		log.Printf("[muxc] Inject stream send: %v", req)
	}
	resp, err := stream.CloseAndRecv()
	log.Printf("[muxc] Inject stream recv: %v", resp)

	return nil
}

// Call Mux.Listen() and print each msg received
func CallListen(cl pb.MuxClient) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Printf("[muxc] Starting Listen stream")
	req := &pb.ListenRequest{}

	log.Printf("[muxc] Listen send: %v", req)
	stream, err := cl.Listen(ctx, req)
	if err != nil {
		log.Printf("[muxc] Error calling Listen: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("[muxc] Error on Listen stream: %v", err)
		}
		log.Printf("[muxc] Listen recv: %v", resp)
	}

	return nil
}

// Ping the Mux server
func CallPing(cl pb.MuxClient) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &pb.PingRequest{Ping: true}
	log.Printf("[muxc] Ping send: %v", req)
	resp, err := cl.Ping(ctx, req)
	if err != nil {
		log.Printf("[muxc] Error calling Ping: %v", err)
	}
	log.Printf("[muxc] Ping recv: %v", resp)

	return nil
}

func GetNewMuxClient(addr string) (pb.MuxClient, *grpc.ClientConn) {

	// Setup grpc conn and client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("[muxc] Error calling Dial: %v", err)
	}
	cl := pb.NewMuxClient(conn)

	return cl, conn

}

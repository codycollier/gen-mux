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
func Inject(cl pb.MuxClient, input chan pb.Datum, done chan int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Printf("[muxc] Starting Inject stream")
	stream, err := cl.Inject(ctx)
	if err != nil {
		log.Printf("[muxc] Error creating Inject stream: %v", err)
	}

	for datum := range input {
		req := &pb.InjectRequest{Datum: &datum}
		log.Printf("[muxc] Inject stream send: %v", req)
		stream.Send(req)
	}
	resp, err := stream.CloseAndRecv()
	log.Printf("[muxc] Inject stream close recv: %v", resp)

	log.Printf("[muxc] Inject stream done")
	done <- 0

	return nil
}

// Call Mux.Listen() and print each msg received
func Listen(cl pb.MuxClient) error {

	ctx := context.Background()

	// Send initial request and start the stream
	req := &pb.ListenRequest{}
	log.Printf("[muxc] Listen send: %v", req)
	stream, err := cl.Listen(ctx, req)
	if err != nil {
		log.Printf("[muxc] Error calling Listen: %v", err)
	}

	// Listen to the stream of messages
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
func Ping(cl pb.MuxClient) error {

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

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "LoadBalancer/gen"

	"google.golang.org/grpc"
)

func main() {
	port := 8082
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterLoadBalancerServer(grpcServer, &LoadBalancer{})
	grpcServer.Serve(lis)
}

type LoadBalancer struct {
	pb.UnimplementedLoadBalancerServer
}

func (s *LoadBalancer) LoadBalanceRequest(ctx context.Context, request *pb.IncomingRequest) (*pb.OutgoingResponse, error) {
	start := time.Now()

	fmt.Printf("Endpoint 2: Received message: %s from %s at %s\n", request.Message, request.Sender, start.String())

	// NOTE: here we have to know what is the previously used index in the machines array.
	// Round-robin algorithm

	duration := time.Since(start)
	return &pb.OutgoingResponse{HandledByMachine: 2, ResponseTime: duration.Milliseconds()}, nil
}

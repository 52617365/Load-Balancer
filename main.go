package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

import pb "LoadBalancer/gen"

type Settings struct {
	algorithm string
	port      int
}

func getArguments() Settings {
	var allowedAlgorithms = []string{"round-robin", "least-connections"}
	flag.Parse()

	port := flag.Int("port", 50051, "the port to serve on")
	algorithm := flag.String("a", "round-robin", "algorithm used for load balancer")

	for _, allowed := range allowedAlgorithms {
		if *algorithm == allowed {
			return Settings{algorithm: *algorithm, port: *port}
		}
	}
	panic("Invalid algorithm specified, only round-robin and least-connections are allowed")
}

func main() {
	settings := getArguments()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", settings.port))

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
	message := request.Message
	sender := request.Sender
	println("Received message: " + message + " from " + sender + " at " + start.String())
	duration := time.Since(start)
	return &pb.OutgoingResponse{HandledByMachine: 2, ResponseTime: duration.Milliseconds()}, nil
}

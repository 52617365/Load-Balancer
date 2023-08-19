package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "LoadBalancer/gen"
)

//type Settings struct {
//	algorithm string
//	port      int
//}

var algorithm string
var port int

var last_used_machine_index int

var machines = []string{"127.0.0.1:8081", "127.0.0.1:8082", "127.0.0.1:8083"}

func getArguments() {
	var allowedAlgorithms = []string{"round-robin", "least-connections"}

	p := flag.Int("port", 50051, "the port to serve on")
	algo := flag.String("algo", "round-robin", "algorithm used for load balancer")

	flag.Parse()

	if !slices.Contains(allowedAlgorithms, *algo) {
		log.Fatalf("invalid algorithm: %s", *algo)
	} else {
		algorithm = *algo
		port = *p
	}
}

func main() {
	getArguments()

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

	if algorithm == "round-robin" {
		index := last_used_machine_index + 1
		if index >= len(machines) {
			index = 0
		}
		endpoint := machines[index]
		fmt.Printf("Request -> %s | %s\n", endpoint, algorithm)

		conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		client := pb.NewLoadBalancerClient(conn)
		test, err := client.LoadBalanceRequest(context.Background(), &pb.IncomingRequest{Message: "hello world", Sender: "john doe"})

		if err != nil {
			panic(err)
		}
		last_used_machine_index = index

		// } else if algorithm == "least-connections" {
		// 	// Least connections algorithm

		return &pb.OutgoingResponse{HandledByMachine: test.HandledByMachine, ResponseTime: test.ResponseTime}, nil
	}

	panic("unhandled")
}

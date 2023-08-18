package main

import (
	pb "LoadBalancer/gen"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestEndpoint(t *testing.T) {
	serverAddr := flag.String("serverAddr", "localhost:50051", "The server address in the format of host:port")
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewLoadBalancerClient(conn)
	test, err := client.LoadBalanceRequest(context.Background(), &pb.IncomingRequest{Message: "hello world", Sender: "john doe"})

	if err != nil {
		panic(err)
	}
	// format print
	fmt.Printf("handled by: %d, response time: %dms.", test.HandledByMachine, test.ResponseTime)
}

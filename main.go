package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "LoadBalancer/gen"
)

var algorithm string
var port int

type State struct {
	lock                    sync.Mutex
	last_used_machine_index int64
}

// Inc increments the counter for the given key.
func (c *State) Inc() int64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.last_used_machine_index >= int64(len(machines)) {
		c.last_used_machine_index = 0
	} else {
		c.last_used_machine_index++
	}
	return c.last_used_machine_index
}

var state State = State{last_used_machine_index: 0, lock: sync.Mutex{}}

type Queue struct {
	lock         sync.Mutex
	queue_amount int64
}

type Machine struct {
	address string
	queue   Queue
}

func (q *Queue) Get() int64 {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.queue_amount
}
func (q *Queue) Inc() {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue_amount++
}

func (q *Queue) Dec() {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue_amount--
}

func getMachineWithLeastConnections() *Machine {
	leastConnectionsMachine := &machines[0]
	for i := range machines {
		if machines[i].queue.Get() < leastConnectionsMachine.queue.Get() {
			leastConnectionsMachine = &machines[i]
		}
	}
	return leastConnectionsMachine
}

var machines = []Machine{{address: "127.0.0.1:8081", queue: Queue{lock: sync.Mutex{}, queue_amount: 0}}, {address: "127.0.0.1:8082", queue: Queue{lock: sync.Mutex{}, queue_amount: 0}}, {address: "127.0.0.1:8083", queue: Queue{lock: sync.Mutex{}, queue_amount: 0}}}

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

// ghz --insecure --async \
//   --proto /balancer.proto \
//   --call LoadBalancer.LoadBalanceRequest \
//   -c 10 -n 10000 --rps 200 \
//   -d '{"message":"Hello world", "sender":"{{.WorkedID}}"}' 127.0.0.1:50051

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
		index := state.Inc()

		machine := &machines[index]
		fmt.Printf("Request -> %s | %s\n", machine.address, algorithm)

		conn, err := grpc.Dial(machine.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		client := pb.NewLoadBalancerClient(conn)
		test, err := client.LoadBalanceRequest(context.Background(), &pb.IncomingRequest{Message: request.Message, Sender: request.Sender})

		if err != nil {
			panic(err)
		}
		return &pb.OutgoingResponse{HandledByMachine: test.HandledByMachine, ResponseTime: test.ResponseTime, RandomIndex: test.RandomIndex}, nil
	} else if algorithm == "least-connections" {
		least_connections := getMachineWithLeastConnections()
		least_connections.queue.Inc()

		fmt.Printf("Request -> %s | %s, Queue counts: machine 1: %d, machine 2: %d, machine3: %d\n", least_connections.address, algorithm, machines[0].queue.Get(), machines[1].queue.Get(), machines[2].queue.Get())

		conn, err := grpc.Dial(least_connections.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		defer least_connections.queue.Dec()

		client := pb.NewLoadBalancerClient(conn)
		test, err := client.LoadBalanceRequest(context.Background(), &pb.IncomingRequest{Message: request.Message, Sender: request.Sender})

		if err != nil {
			panic(err)
		}
		return &pb.OutgoingResponse{HandledByMachine: test.HandledByMachine, ResponseTime: test.ResponseTime, RandomIndex: test.RandomIndex}, nil
	}

	panic("unhandled")
}

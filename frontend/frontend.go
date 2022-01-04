package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	mockPackage "github.com/Troelshjarne/Disys_mock_exam/increment"
	"google.golang.org/grpc"
)

var port = flag.String("port", ":9080", "port for server to listen on")
var replicas []mockPackage.CommunicationClient
var ctx context.Context
var lamTime = 0
var ips []string
var valueMutex sync.Mutex

type server struct {
	mockPackage.UnimplementedCommunicationServer
}

func main() {
	fmt.Println("=== Welcome to increment beta ===")
	// set up logging functionality
	logSetup()
	// append replica ips
	getIps()
	// connect to replicas
	replicaConns()

	ctx = context.Background()

	// start frontend server
	serverSetup()
}

func logSetup() {
	file, er := os.OpenFile("../logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if er != nil {
		log.Fatal(er)
	}
	log.SetOutput(file)
}

func getIps() {
	// possibly read from file.
	ips = append(ips, ":9081", ":9082", ":9083")
}

func replicaConns() {

	for _, ip := range ips {
		conn, err := grpc.Dial(ip, grpc.WithInsecure())
		if err != nil {
			continue
		}
		replica := mockPackage.NewCommunicationClient(conn)
		replicas = append(replicas, replica)

	}

}

func serverSetup() {
	flag.Parse()
	fmt.Println("=== Frontend starting up ===")
	list, err := net.Listen("tcp", *port)

	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}

	server := server{}

	grpc := grpc.NewServer()
	mockPackage.RegisterCommunicationServer(grpc, &server)

	if err := grpc.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func ticker(in int, local int) {
	if in > local {
		lamTime = in + 1
	} else {
		lamTime++
	}
}

func (s *server) Increment(ctx context.Context, in *mockPackage.IncRequest) (*mockPackage.Reply, error) {
	//log.Println("recieved at lamportTime :", lamTime)
	valueMutex.Lock()
	ticker(int(in.Time), lamTime)
	valueMutex.Unlock()
	fmt.Println("recieved at lamportTime :", lamTime)
	response := int32(0)
	request := in.Inc

	x := mockPackage.IncRequest{
		Inc:  int32(request),
		Time: int32(lamTime),
	}
	// direct incoming requests from clients to replicas
	for _, client := range replicas {

		reply, err := client.Increment(ctx, &x)
		if err != nil {
			log.Println("detected one unvailable replica, at lamportTimes:", lamTime)
			continue
		}
		// send back the biggest value, in case replicas are out of sync (if biggest i not "best", implement some consensus function)
		if response <= reply.Counter {
			response = reply.Counter
		}
		ticker(int(reply.Time), lamTime)
	}

	return &mockPackage.Reply{Counter: response, Time: int32(lamTime)}, nil
}

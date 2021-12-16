package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	mockPackage "github.com/Troelshjarne/Disys_mock_exam/increment"
	"google.golang.org/grpc"
)

var port = flag.String("port", ":9080", "port for server to listen on")
var valueMutex sync.Mutex
var value int32 = -1

type server struct {
	mockPackage.UnimplementedCommunicationServer
}

func main() {
	// go run server.go -port 127.0.0.1:9081
	//port := flag.String("port", ":9080", "port for server to listen on")
	flag.Parse()
	fmt.Println("=== Server starting up ===")
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

func (s *server) Increment(ctx context.Context, in *mockPackage.IncRequest) (*mockPackage.Reply, error) {

	valueMutex.Lock()
	value += in.Inc
	valueMutex.Unlock()

	return &mockPackage.Reply{Counter: value}, nil
}

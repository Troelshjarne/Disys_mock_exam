package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	mockPackage "github.com/Troelshjarne/Disys_mock_exam/increment"
	"google.golang.org/grpc"
)

var Mutex sync.Mutex
var value = 0

type Server struct {
	mockPackage.UnimplementedCommunicationServer
}

func main() {

	fmt.Println("=== Server starting up ===")
	list, err := net.Listen("tcp", ":9080")

	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}

	var options []grpc.ServerOption
	grpcServer := grpc.NewServer(options...)

	grpcServer.Serve(list)
}

func (s *Server) Increment(requestStream mockPackage.Communication_IncrementServer) error {

}

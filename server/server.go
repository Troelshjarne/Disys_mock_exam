package main

import (
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

type Server struct {
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

	var options []grpc.ServerOption
	grpcServer := grpc.NewServer(options...)

	mockPackage.RegisterCommunicationServer(grpcServer, &Server{})

	grpcServer.Serve(list)
}
func (s *Server) Increment(requestStream mockPackage.Communication_IncrementServer) error {

	request, err := requestStream.Recv()
	if err != nil {
		log.Printf("Request error: %v \n", err)
	}

	fmt.Println(request.Inc)
	fmt.Println(*port)
	inc := request.Inc

	valueMutex.Lock()
	value += inc
	valueMutex.Unlock()

	requestStream.SendAndClose(&mockPackage.Reply{Counter: value})

	return nil
}

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

var lamTime = 0
var port = flag.String("port", ":9080", "port for server to listen on")
var valueMutex sync.Mutex
var value int32 = -1

type server struct {
	mockPackage.UnimplementedCommunicationServer
}

func main() {

	logSetup()

	replicaSetup()
}

func replicaSetup() {
	// go run server.go -port 127.0.0.1:9081
	//port := flag.String("port", ":9080", "port for server to listen on")
	flag.Parse()
	fmt.Println("=== Replica starting up ===")
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

func logSetup() {
	file, er := os.OpenFile("../logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if er != nil {
		log.Fatal(er)
	}
	log.SetOutput(file)
}

func (s *server) Increment(ctx context.Context, in *mockPackage.IncRequest) (*mockPackage.Reply, error) {
	//log.Println("recieved at lamportTime :", lamTime)
	valueMutex.Lock()
	ticker(int(in.Time), lamTime)
	value += in.Inc
	valueMutex.Unlock()
	fmt.Println("recieved at lamportTime :", lamTime)
	lamTime += 1
	return &mockPackage.Reply{Counter: value, Time: int32(lamTime)}, nil
}

func ticker(in int, local int) {
	if in > local {
		lamTime = in + 1
	} else {
		lamTime++
	}
}

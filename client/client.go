package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	mockPackage "github.com/Troelshjarne/Disys_mock_exam/increment"
	"google.golang.org/grpc"
)

var tcpServer = flag.String("server", ":9080", "Tcp server")

func main() {
	fmt.Println("=== Welcome to Chitty Chat - Beta 0.1.2 ===")
	var options []grpc.DialOption
	options = append(options, grpc.WithBlock(), grpc.WithInsecure())
	//connect to server
	conn, err := grpc.Dial(*tcpServer, options...)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()

	// client connection interface
	client := mockPackage.NewCommunicationClient(conn)

	fmt.Println(client.Increment(ctx, &mockPackage.Empty{}))
}

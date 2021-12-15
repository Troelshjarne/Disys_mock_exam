package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	mockPackage "github.com/Troelshjarne/Disys_mock_exam/increment"
	"google.golang.org/grpc"
)

var incrementer = 1
var tcpServer = flag.String("server", ":9080", "Tcp server")

func main() {
	fmt.Println("=== Welcome to increment beta")
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

	for {
		fmt.Println("im alive")
		time.Sleep(time.Second * 2)
		increment(ctx, client)
		//Fix ID sent with message
	}
	//fmt.Println(client.Increment(ctx, &mockPackage.Empty{}))
}

func increment(ctx context.Context, client mockPackage.CommunicationClient) {
	fmt.Println("sending request")

	stream, err := client.Increment(ctx)

	if err != nil {
		log.Printf("Failure sending increment request. Got this error: %v", err)
	}

	request := mockPackage.IncRequest{
		Inc: int32(incrementer),
	}
	fmt.Println(request.Inc, "test")

	stream.Send(&request)
	incrementer++
	acc, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(acc.Counter)

}

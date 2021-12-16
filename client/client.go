package main

import (
	"context"
	"fmt"
	"log"
	"time"

	mockPackage "github.com/Troelshjarne/Disys_mock_exam/increment"
	"google.golang.org/grpc"
)

var clients []mockPackage.CommunicationClient
var incrementer = 1

func main() {
	fmt.Println("=== Welcome to increment beta")
	//var options []grpc.DialOption
	//options = append(options, grpc.WithBlock(), grpc.WithInsecure())
	//connect to server
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	conn2, err := grpc.Dial(":9081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	conn3, err := grpc.Dial(":9082", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// client connection interface
	client := mockPackage.NewCommunicationClient(conn)
	client2 := mockPackage.NewCommunicationClient(conn2)
	client3 := mockPackage.NewCommunicationClient(conn3)

	clients = append(clients, client, client2, client3)

	defer conn.Close()

	ctx := context.Background()

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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	mockPackage "github.com/Troelshjarne/Disys_mock_exam/increment"
	"google.golang.org/grpc"
)

var lamTime = 0
var clients []mockPackage.CommunicationClient
var incrementer = 1
var ctx context.Context
var incrementerMutex sync.Mutex

func main() {
	fmt.Println("=== Welcome to increment beta")

	file, er := os.OpenFile("../logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if er != nil {
		log.Fatal(er)
	}
	log.SetOutput(file)
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

	ctx = context.Background()

	message := make(chan string)

	go func() {
		time.Sleep(time.Second * 3)
		message <- "hello"
	}()

	go func() {
		msg := <-message

		fmt.Println("received", msg)
	}()

	for {
		//fmt.Println("The current Value on the Server is : ")
		time.Sleep(time.Second * 2)
		increment()

		//Fix ID sent with message
	}
	//fmt.Println(client.Increment(ctx, &mockPackage.Empty{}))
}

func increment() {
	lamTime++
	fmt.Println("sending request at lamportTime :", lamTime)
	log.Println("sending request at lamportTime :", lamTime)
	response := int32(0)
	for _, client := range clients {
		//request := mockPackage.IncRequest{}
		//lam := mockPackage.IncRequest{}
		//request.Inc = int32(incrementer)
		//lam.Time = int32(lamTime)
		incrementer++
		x := mockPackage.IncRequest{
			Inc:  int32(incrementer),
			Time: int32(lamTime),
		}

		ack, err := client.Increment(ctx, &x)
		if err != nil {

			continue
		}

		if response <= ack.Counter {
			response = ack.Counter
		}

		lamTime = int(ack.Time)
	}
	fmt.Println(response)

}

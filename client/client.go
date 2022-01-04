package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	mockPackage "github.com/Troelshjarne/Disys_mock_exam/increment"
	"google.golang.org/grpc"
)

var lamTime = 0
var incrementer = 1
var ctx context.Context
var ips []string
var connfrontend mockPackage.CommunicationClient

func main() {
	fmt.Println("=== Welcome to increment beta")

	getIps()
	connect()
	file, er := os.OpenFile("../logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if er != nil {
		log.Fatal(er)
	}
	log.SetOutput(file)

	ctx = context.Background()

	increment()

}

func getIps() {
	// possibly read from file.
	ips = append(ips, ":9084", ":9085")
}

// make a connection to the first available frontend

// connects to one available interface
func connect() {

	for _, ip := range ips {
		conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
		if err != nil {
			//give error
			continue
		} else {
			fmt.Println("connecting to :", ip)
			connfrontend = mockPackage.NewCommunicationClient(conn)
			break
		}

	}
}

func increment() {

	//log.Println("sending request at lamportTime :", lamTime)
	//response := int32(0)

	for {
		// increment lamport when sending request
		lamTime++
		x := mockPackage.IncRequest{
			Inc:  int32(incrementer),
			Time: int32(lamTime),
		}
		ack, err := connfrontend.Increment(ctx, &x)
		if err != nil {
			fmt.Println("increment request failed, connecting to new interface")
			connect()
			time.Sleep(time.Second * 2)
			//increment()
			continue
		}
		fmt.Println("sending request at lamportTime :", lamTime)
		ticker(int(ack.Time), lamTime)

		fmt.Println(ack.Counter)
		time.Sleep(time.Second * 2)

	}
}

func ticker(in int, local int) {
	if in > local {
		lamTime = in + 1
	} else {
		lamTime++
	}
}

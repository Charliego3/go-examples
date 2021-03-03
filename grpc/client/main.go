package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/whimthen/temp/grpc/helloworld"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	ctx := context.Background()

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("Greeting: %s", r.GetMessage())

			r, err = c.SyaHelloAgain(ctx, &pb.HelloRequest{Name: "again"})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Printf("AgainGreeting: %s", r.GetMessage())
		}
	}
}

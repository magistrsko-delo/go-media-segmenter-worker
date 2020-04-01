package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "main/proto"
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

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{
		Name: name,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	log.Println("NADALJUJEM  PROGRAM")

}
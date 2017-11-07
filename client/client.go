package main

import (
	"log"

	pb "github.com/supunz/go-job-queue/mailservice"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMailServiceClient(conn)

	// Contact the server and print out its response.
	name := "defaultName"
	receiver := "receiver"
	title := "title"

	response, err := c.PutEmail(context.Background(), &pb.EmailRequest{Sender: &name, Receiver: &receiver, Title: &title})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}
	log.Printf("it was a %s", response.Status)
}

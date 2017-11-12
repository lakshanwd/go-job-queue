package main

import (
	"log"

	pb "github.com/supunz/go-job-queue/mail"
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

	//creating new mail client
	c := pb.NewMailClient(conn)

	// Contact the server and print out its response.
	name := "name"
	receiver := "receiver"
	title := "title"
	content := "content goes here"

	for i := 0; i < 10000; i++ {
		response, err := c.PutEmail(context.Background(), &pb.EmailRequest{Sender: name, Receiver: receiver, Title: title, Content: content})
		if err != nil {
			log.Fatalf("could not put email on server due to : %v", err)
			return
		}
		log.Printf("response status was %v", response.Status)
	}
}

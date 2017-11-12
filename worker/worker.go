package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/supunz/go-job-queue/mail"
	"google.golang.org/grpc"
)

const (
	address        = "localhost:50051"
	maxDeamonCount = 20
	maxQueueSize   = 1000
)

var queue chan *pb.EmailRequest
var quit chan error

type deamon struct {
	ID int
}

func newDeamon(id int) *deamon {
	return &deamon{ID: id}
}

func (w *deamon) start(queue chan *pb.EmailRequest) {
	go func() {
		for {
			select {
			case email := <-queue:
				send(email, w.ID)
			}
		}
	}()
}

func send(e *pb.EmailRequest, deamonID int) {
	time.Sleep(time.Millisecond * 100)
	log.Printf("email %v sent by deamon %v\n", e.GetTitle(), deamonID)
}

func main() {
	queue := make(chan *pb.EmailRequest, maxQueueSize)
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//creating new mail client
	c := pb.NewMailClient(conn)

	for i := 0; i < maxDeamonCount; i++ {
		deamon := newDeamon(i)
		deamon.start(queue)
	}

	worker := &pb.Worker{WorkerName: os.Args[1]}
	for {
		response, _ := c.GetEmail(context.Background(), worker)
		select {
		case queue <- response:
			log.Printf("took email %v\n", response.GetTitle())
		}
	}
}

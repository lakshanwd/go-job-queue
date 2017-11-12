package main

import (
	"log"
	"net"

	pb "github.com/supunz/go-job-queue/mail"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//Port - port to Listen
//MaxQueueSize - max size of queue
//MaxWorkerCount - max count of workers
const (
	Port           = ":50051"
	MaxQueueSize   = 10000
	MaxWorkerCount = 20
)

var queue chan *pb.EmailRequest

// server is used to implement mail.MailServer.
type server struct {
}

// PutEmail implements mail.PutEmail
func (s *server) PutEmail(ctx context.Context, in *pb.EmailRequest) (*pb.EmailResponse, error) {
	queue <- in
	defer log.Printf("email received %v\n", in.GetTitle())
	return &pb.EmailResponse{Status: true}, nil
}

func (s *server) GetEmail(ctx context.Context, in *pb.Worker) (*pb.EmailRequest, error) {
	defer log.Printf("email taken by %v\n", in.GetWorkerName())
	return <-queue, nil
}

func main() {
	//create channel for queuing email requests
	queue = make(chan *pb.EmailRequest, MaxQueueSize)

	//listen to tcp
	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//create new server
	s := grpc.NewServer()
	pb.RegisterMailServer(s, &server{})

	// Register reflection service on gRPC server.
	log.Println("starting server...")
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

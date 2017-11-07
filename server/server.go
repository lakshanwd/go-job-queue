package main

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/supunz/go-job-queue/mailservice"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	Port           = ":50051"
	MaxQueueSize   = 10000
	MaxWorkerCount = 20
)

var queue chan Email

// server is used to implement helloworld.GreeterServer.
type server struct {
}

//Email - email
type Email struct {
	Sender   string
	Receiver string
	Title    string
}

func (e *Email) send(workerId int) {
	time.Sleep(time.Millisecond * 100)
	fmt.Println("email send by worker#", workerId)
}

type Worker struct {
	ID int
}

func NewWorker(id int) *Worker {
	return &Worker{ID: id}
}

func (w *Worker) start() {
	go func() {
		for {
			select {
			case email := <-queue:
				email.send(w.ID)
			}
		}
	}()
}

// SayHello implements helloworld.GreeterServer
func (s *server) PutEmail(ctx context.Context, in *pb.EmailRequest) *pb.EmailResponse {
	email := Email{Sender: in.GetSender(), Receiver: in.GetReceiver(), Title: in.GetTitle()}
	queue <- email
	return &pb.EmailResponse{Status: true}
}

func main() {
	queue = make(chan Email, MaxQueueSize)
	for i := 0; i < MaxWorkerCount; i++ {
		worker := NewWorker(i)
		worker.start()
	}
	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMailServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

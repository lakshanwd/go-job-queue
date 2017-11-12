package main

import (
	"fmt"
	"log"
	"net"
	"time"

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

var queue chan Email

// server is used to implement mail.MailServer.
type server struct {
}

//Email - email
type Email struct {
	Sender   string
	Receiver string
	Title    string
	Content  string
}

func (e *Email) send(workerID int) {
	time.Sleep(time.Millisecond * 100)
	fmt.Println("email send by worker#", workerID)
}

type worker struct {
	ID int
}

func newWorker(id int) *worker {
	return &worker{ID: id}
}

func (w *worker) start() {
	go func() {
		for {
			select {
			case email := <-queue:
				email.send(w.ID)
			}
		}
	}()
}

// PutEmail implements mail.PutEmail
func (s *server) PutEmail(ctx context.Context, in *pb.EmailRequest) (*pb.EmailResponse, error) {
	email := Email{Sender: in.GetSender(), Receiver: in.GetReceiver(), Title: in.GetTitle(), Content: in.GetContent()}
	queue <- email
	return &pb.EmailResponse{Status: true}, nil
}

func main() {
	//create channel for queuing email requests
	queue = make(chan Email, MaxQueueSize)

	//define maximum workers
	for i := 0; i < MaxWorkerCount; i++ {
		worker := newWorker(i)
		worker.start()
	}

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

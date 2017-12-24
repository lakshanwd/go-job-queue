package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/supunz/go-job-queue/mail"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//Port - port to Listen
var Port string

//MaxQueueSize - max size of queue
var MaxQueueSize int

//MaxWorkerCount - max count of workers
var MaxWorkerCount int

var queue chan *pb.EmailRequest

//config for server
type serverConfig struct {
	Port           int `json:"port"`
	MaxWorkerCount int `json:"max-worker-count"`
	MaxQueueSize   int `json:"max-queue-size"`
}

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
	//setup configuration
	config := loadConfiguration()
	Port = fmt.Sprintf(":%d", config.Port)
	MaxQueueSize = config.MaxQueueSize
	MaxWorkerCount = config.MaxWorkerCount

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

func loadConfiguration() serverConfig {
	file := "./config.json"
	var config serverConfig
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/lakshanwd/go-job-queue/mail"
	"google.golang.org/grpc"
)

var address string
var maxDeamonCount int
var queue chan *pb.EmailRequest
var quit chan bool

//config for worker
type workerConfig struct {
	Address        string `json:"address"`
	MaxDeamonCount int    `json:"max-deamon-count"`
}

type deamon struct {
	ID int
}

func newDeamon(id int) *deamon {
	return &deamon{ID: id}
}

func (d *deamon) start(mailClient pb.MailClient, worker *pb.Worker) {
	go func() {
		for {
			select {
			case email, ok := <-queue:
				if !ok {
					log.Printf("shutting down deamon %v\n", d.ID)
					quit <- true
					return
				}
				send(email, d.ID)
			}
		}
	}()
}

func send(e *pb.EmailRequest, deamonID int) {
	time.Sleep(time.Millisecond * 100)
	log.Printf("email %v sent by deamon %v\n", e.GetTitle(), deamonID)
}

func main() {
	//setup configuration
	config := loadConfiguration()
	address = config.Address
	maxDeamonCount = config.MaxDeamonCount

	//initialize channels
	queue = make(chan *pb.EmailRequest, maxDeamonCount)
	quit = make(chan bool)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect: %v", err.Error())
	}
	defer conn.Close()
	log.Println("connected")

	//initializing server
	mailClient := pb.NewMailClient(conn)
	worker := &pb.Worker{WorkerName: os.Args[1]}

	//getting mail requests from server
	go func() {
		for {
			if email, err := mailClient.GetEmail(context.Background(), worker); err != nil {
				queue <- email
			} else {
				quit <- true
				return
			}
		}
	}()

	//create deamons
	for i := 0; i < maxDeamonCount; i++ {
		deamon := newDeamon(i)
		deamon.start(mailClient, worker)
	}

	<-quit
	close(queue)

	//waiting for deamons to do their remaining work
	for i := 0; i < maxDeamonCount; i++ {
		<-quit
	}
}

func loadConfiguration() workerConfig {
	file := "./config.json"
	var config workerConfig
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

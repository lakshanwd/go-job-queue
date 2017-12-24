package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/supunz/go-job-queue/mail"
	"google.golang.org/grpc"
)

var address string
var maxDeamonCount int
var maxQueueSize int

var queue chan *pb.EmailRequest
var quit chan bool
var start chan bool
var mailClient chan pb.MailClient

//config for worker
type workerConfig struct {
	Address        string `json:"address"`
	MaxDeamonCount int    `json:"max-deamon-count"`
	MaxQueueSize   int    `json:"max-queue-size"`
}

type deamon struct {
	ID int
}

func newDeamon(id int) *deamon {
	return &deamon{ID: id}
}

func (w *deamon) start() {
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
	//setup configuration
	config := loadConfiguration()
	address = config.Address
	maxDeamonCount = config.MaxDeamonCount
	maxQueueSize = config.MaxQueueSize

	//initialize channels
	queue = make(chan *pb.EmailRequest, maxQueueSize)
	quit = make(chan bool)
	start = make(chan bool)
	mailClient = make(chan pb.MailClient)

	//create deamons
	for i := 0; i < maxDeamonCount; i++ {
		deamon := newDeamon(i)
		deamon.start()
	}

	for {
		// setup a connection to the server.
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("unable to connect: %v", err.Error())
		}
		log.Println("connected")

		//initializing server
		c := pb.NewMailClient(conn)
		worker := &pb.Worker{WorkerName: os.Args[1]}

		go func() {
			//waiting for server to start
			<-start
			for {
				//taking mail from server
				email, err := c.GetEmail(context.Background(), worker)
				if err == nil {
					//put email in queue
					queue <- email
					log.Printf("took %v from server", email.GetTitle())
				} else {
					//signal server to terminate
					quit <- true
					return
				}
			}
		}()

		//signal server to read emails
		start <- true

		//waiting for signal to quit current connection and attempt to reconnect
		<-quit

		//existing connection is closing
		conn.Close()

		//retry connection in 3 seconds
		time.Sleep(time.Second * 3)
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

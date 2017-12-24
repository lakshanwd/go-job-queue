package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	pb "github.com/supunz/go-job-queue/mail"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var address string

//config for client
type clientConfig struct {
	Address string `json:"address"`
}

func main() {
	//setup configuration
	config := loadConfiguration()
	address = config.Address

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
	content := "content goes here"

	for i := 0; i < 8000; i++ {
		email := &pb.EmailRequest{Sender: name, Receiver: receiver, Title: fmt.Sprintf("title %v", i), Content: content}
		response, err := c.PutEmail(context.Background(), email)
		if err != nil {
			log.Fatalf("could not put email on server due to : %v", err)
			return
		}
		log.Printf("response status was %v", response.Status)
	}
}

func loadConfiguration() clientConfig {
	file := "./config.json"
	var config clientConfig
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

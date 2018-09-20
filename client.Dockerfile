FROM golang:latest

WORKDIR /go/src/github.com/lakshanwd/go-job-queue

COPY client/client.go ./client/client.go
COPY mail ./mail
RUN go get -d -v ./... && go install -v ./...
COPY client/config.prod.json ./config.json
CMD [ "client" ]
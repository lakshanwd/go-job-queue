FROM golang:latest

WORKDIR /go/src/github.com/lakshanwd/go-job-queue

COPY worker/worker.go ./worker/worker.go
COPY mail ./mail
RUN go get -d -v ./... && go install -v ./...
COPY worker/config.prod.json ./config.json
CMD [ "worker","test-worker" ]
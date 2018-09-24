FROM golang:latest

WORKDIR /go/src/github.com/lakshanwd/go-job-queue

COPY server/server.go ./server/server.go
COPY mail ./mail
RUN go get -d ./... && go install -v ./...
COPY server/config.prod.json ./config.json
CMD [ "server" ]
EXPOSE 50051
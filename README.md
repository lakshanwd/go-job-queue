# go-job-queue ![build_status][build_status]

An example implementation of job queue with load balancer mechanism. Used [grpc](https://grpc.io/) in order to gain maximum performance.

## theory
- central server application listens to thousands of requests per second and loads them in to a queue.
- decentralized workers connect to central server and takes job parameters from the queue and store them in their own queue.
- each worker has limited no of deamon processes (go routines) which will take job parameters from their own worker queue and process and executes job

[build_status]: https://travis-ci.org/lakshanwd/go-job-queue.svg?branch=master "Travis Build Status"

## containerized and scalable
run `$ docker-compose up --scale worker=2 --scale client=5` to execute this app in containerized environment
run `$ docker-compose scale worker=3 client=8` to scale app

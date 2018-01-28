# go-job-queue ![alt text][build_status]

An example implementation of job queue with load balancer mechanism. Used [grpc](https://grpc.io/) in order to gain maximum performance.

## theory
- central server application listens to thousands of requests per second and loads them in to a queue.
- decentralized worker programs connect to central server and takes job parameters from the queue and store it in their own queue.
- each worker has limiter no of deamon processes (go routines) which will take job parameters from their own worker queue and process and executes job

[build_status]: https://travis-ci.org/supunz/go-job-queue.svg?branch=master "Travis Build Status"

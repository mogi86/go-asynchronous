# go-asynchronous

Sample code for asynchronous and AWS SQS of Go.
In this repository, we will use AWS SQS.
Message is published to SQS via HTTP protocol.
Worker subscribe to messages from SQS.

## Directory

- TODO

## Requirement

- Go
- AWS account
- AWS root account or AWS IAM user to use SQS

## How to use

### Start server to publish a message for SQS

- Start http server. `QUEUE_NAME` need to be specified your SQS name. ex) MogiSQS

```bash
$ export QUEUE_NAME={YOUR_SQS_NAME}
$ go run cmd/http/main.go
```

### Start Worker to subscribe a message for SQS from another terminal

- Build and execute to worker. `main.QueueName` need to specify your SQS name. ex) MogiSQS

```bash
$ go build -o ./worker -ldflags '-X main.QueueName={YOUR_SQS_NAME}' cmd/worker/main.go
$ ./worker
```

### Publish a message from another terminal

```bash
$ curl -v GET "http://localhost:8000/publish?message=hello"
```



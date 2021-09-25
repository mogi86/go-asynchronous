package http

import (
	"fmt"
	"github.com/mogi86/go-asynchronous/config"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// ref: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/sns-example-publish.html
// ref: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/sqs-example-receive-message.html
// ref: https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/sqs/SendMessage/SendMessage.go

func PublishMessage(w http.ResponseWriter, r *http.Request) {
	m := r.URL.Query().Get("message")
	if m == "" {
		fmt.Println("published message isn't set")
		http.Error(w, "published message isn't set", 500)
	}

	// get session
	// It seems that credential file (~/.aws/credentials) is referred.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		//Config: aws.Config{
		//	Region: aws.String("us-east-1"),
		//},
	}))

	// Get URL of queue
	c := config.GetConfig()
	result, err := GetQueueURL(sess, &c.QueueName)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		http.Error(w, "Got an error getting the queue URL:", 500)
	}

	err = SendMsg(sess, result.QueueUrl, m)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		http.Error(w, "Got an error sending the message:", 500)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintf(w, "success")
	if err != nil {
		fmt.Println("fail to print")
	}
}

// GetQueueURL gets the URL of an Amazon SQS queue
func GetQueueURL(sess *session.Session, queue *string) (*sqs.GetQueueUrlOutput, error) {
	// Create an SQS service client
	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SendMsg sends a message to an Amazon SQS queue
func SendMsg(sess *session.Session, queueURL *string, message string) error {
	svc := sqs.New(sess)

	_, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(1),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Sent Message"),
			},
		},
		MessageBody: aws.String(message),
		QueueUrl:    queueURL,
	})
	if err != nil {
		return err
	}

	return nil
}

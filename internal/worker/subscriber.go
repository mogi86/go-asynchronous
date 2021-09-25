package worker

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// ref: receive and delete message: https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/sqs/README.md#deletemessage
//   It seems that just subscribing doesn't delete the message.
//   => https://docs.aws.amazon.com/ja_jp/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-visibility-timeout.html

func GetMessage(queueName interface{}) error {
	var queue string
	queue, ok := queueName.(string)
	if ok == false {
		return errors.New("failed interface covert to string")
	}

	// Create a session that gets credential values from ~/.aws/credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Get URL of queue
	urlResult, err := GetQueueURL(sess, &queue)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return nil
	}

	var timeout int64 = 30
	msgResult, err := GetMessages(sess, urlResult.QueueUrl, &timeout)
	if err != nil {
		fmt.Println("Got an error receiving messages:")
		fmt.Println(err)
		return nil
	}

	if len(msgResult.Messages) == 0 {
		fmt.Println("There is no message")
		return nil
	}

	fmt.Println("Message ID: " + *msgResult.Messages[0].MessageId)
	fmt.Println("Message Handle: " + *msgResult.Messages[0].ReceiptHandle)
	fmt.Println("Message Body: " + *msgResult.Messages[0].Body)

	// delete
	err = DeleteMessage(sess, urlResult.QueueUrl, msgResult.Messages[0].ReceiptHandle)
	if err != nil {
		fmt.Println("Got an error deleting the message:")
		fmt.Println(err)
		return nil
	}
	return nil
}

// GetQueueURL gets the URL of an Amazon SQS queue
func GetQueueURL(sess *session.Session, queue *string) (*sqs.GetQueueUrlOutput, error) {
	svc := sqs.New(sess)

	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}

	return urlResult, nil
}

// GetMessages gets the messages from an Amazon SQS queue
func GetMessages(sess *session.Session, queueURL *string, timeout *int64) (*sqs.ReceiveMessageOutput, error) {
	// Create an SQS service client
	svc := sqs.New(sess)

	msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   timeout,
	})

	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

func DeleteMessage(sess *session.Session, queueURL *string, messageHandle *string) error {
	svc := sqs.New(sess)

	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueURL,
		ReceiptHandle: messageHandle,
	})
	if err != nil {
		return err
	}

	return nil
}

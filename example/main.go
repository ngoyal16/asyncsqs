package main

import (
	"context"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/ngoyal16/asyncsqs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	queueURL = "https://sqs.us-east-1.amazonaws.com/xxxxxxxxxxxx/qqqqqqqqqqqq"
)

var (
	client *asyncsqs.BufferedClient
)

func AsyncSQS(maxMessages uint64) {
	// Create a SQS client with appropriate credentials/IAM role, region etc.
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("config.LoadDefaultConfig() failed: %v", err)
	}
	sqsClient := sqs.NewFromConfig(awsCfg)

	// Create a asyncsqs buffered client; you'd have one per SQS queue
	client, err = asyncsqs.NewBufferedClient(asyncsqs.Config{
		SQSClient:            sqsClient,
		QueueURL:             queueURL,
		SendBatchEnabled:     true,
		DeleteBatchEnabled:   true,
		ReceiveBatchEnabled:  true,
		ReceiveWaitTime:      int32(10),
		OnSendMessageBatch:   sendResponseHandler,
		OnDeleteMessageBatch: deleteResponseHandler,
		OnReceiveMessage:     receiveResponseHandler,
	})
	if err != nil {
		log.Fatalf("asyncsqs.NewBufferedClient() failed: %v", err)
	}
	// important! Stop() ensures that requests in memory are gracefully
	// flushed/dispatched and resources like goroutines are cleaned-up
	defer client.Stop()

	for i := 0; uint64(i) < maxMessages; i++ {
		_ = client.SendMessageAsync(types.SendMessageBatchRequestEntry{
			Id:          aws.String(strconv.Itoa(i)),
			MessageBody: aws.String(strconv.Itoa(i)),
		})
	}

	// receive via normal SQS client and delete via async SQS client
	for {
		stats := client.Stats()
		if stats.MessagesReceived == maxMessages {
			break
		}

		client.ReceiveMessages()
	}

	for {
		stats := client.Stats()

		if stats.MessagesDeleted == maxMessages {
			break
		}
	}
}

func main() {
	AsyncSQS(uint64(100))
}

func sendResponseHandler(output *sqs.SendMessageBatchOutput, err error) {
	if err != nil {
		log.Printf("send returned error: %v", err)
		return
	}
	for _, s := range output.Successful {
		log.Printf("message send successful: msg id = %s", *s.Id)
	}
	for _, f := range output.Failed {
		log.Printf("message send failed: msg id = %s", *f.Id)
	}
}

func receiveResponseHandler(output *sqs.ReceiveMessageOutput, err error) {
	if err != nil {
		log.Printf("send returned error: %v", err)
		return
	}

	if output.Messages != nil {
		for _, m := range output.Messages {
			log.Printf("message received successful: msg id = %s", *m.MessageId)

			_ = client.DeleteMessageAsync(types.DeleteMessageBatchRequestEntry{
				Id:            aws.String(uuid.New().String()),
				ReceiptHandle: m.ReceiptHandle,
			})
		}
	}
}

func deleteResponseHandler(output *sqs.DeleteMessageBatchOutput, err error) {
	if err != nil {
		log.Printf("send returned error: %v", err)
		return
	}
	for _, s := range output.Successful {
		log.Printf("message delete successful: msg id = %s", *s.Id)
	}
	for _, f := range output.Failed {
		log.Printf("message delete failed: msg id = %s", *f.Id)
	}
}

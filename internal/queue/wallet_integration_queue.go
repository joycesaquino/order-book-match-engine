package queue

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/caarlos0/env"
	"log"
)

type Config struct {
	WalletQueue string `env:"WALLET_INTEGRATION_QUEUE_URL,required"`
	Region      string `env:"AWS_REGION,required"`
}

type sqsQueue struct {
	awsSqs *sqs.SQS
	config *Config
}

type Queue interface {
	Send(ctx context.Context, event interface{}) error
}

func (queue sqsQueue) Send(ctx context.Context, event interface{}) error {
	bytes, e := json.Marshal(event)
	if e != nil {
		return e
	}

	msg := &sqs.SendMessageInput{
		MessageBody: aws.String(string(bytes)),
		QueueUrl:    aws.String(queue.config.WalletQueue),
	}

	_, e = queue.awsSqs.SendMessageWithContext(ctx, msg)
	return e
}

func NewSQSQueue(sess *session.Session) Queue {

	var config *Config
	err := env.Parse(&config)
	if err != nil {
		log.Fatalf("[ERROR] - Erro on configure Wallet Integration Queue client: %s", err)
	}

	return &sqsQueue{
		awsSqs: sqs.New(sess),
		config: config,
	}

}

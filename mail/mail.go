package mail

import (
	"context"
	"jobsheet-go-aws2/constant"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SendParam struct {
	Title  string
	Body   string
	MailTo string
}

// メッセージをSQSに送信する
func SendMessage(param SendParam) error {
	ctx := context.Background()
	sqsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	sqsClient := sqs.NewFromConfig(sqsConfig)
	input := &sqs.SendMessageInput{
		MessageBody:  aws.String(param.Body),
		QueueUrl:     aws.String(constant.QUEUE_URL),
		DelaySeconds: 1,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"title": {
				DataType:    aws.String("String"),
				StringValue: aws.String(param.Title),
			},
			"body": {
				DataType:    aws.String("String"),
				StringValue: aws.String(param.Body),
			},
			"mailto": {
				DataType:    aws.String("String"),
				StringValue: aws.String(param.MailTo),
			},
		},
	}
	sqsClient.SendMessage(ctx, input)
	return nil
}

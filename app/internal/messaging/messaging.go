package messaging

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ortisan/router-go/internal/config"
	errApp "github.com/ortisan/router-go/internal/error"
	"github.com/ortisan/router-go/internal/util"
	"github.com/rs/zerolog/log"
)

var (
	timeoutVisibility int64 = 12 * 60
)

func getMessagesWithSession(sess *session.Session, queueURL *string, timeout *int64) (*sqs.ReceiveMessageOutput, error) {
	// Create an SQS service client
	client := sqs.New(sess)

	result, err := client.ReceiveMessage(&sqs.ReceiveMessageInput{
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

	return result, nil
}

func getMessages(queueURL *string, timeout *int64) (*sqs.ReceiveMessageOutput, error) {
	return getMessagesWithSession(config.NewAWSSession(), queueURL, timeout)
}

func GetHealthMessage() (string, error) {

	result, err := getMessages(&config.ConfigObj.AWS.SQS.HealthQueueUrl, &timeoutVisibility)

	if err != nil {
		return "", errApp.NewIntegrationError("Error to get sqs messages.", err)
	}

	msgReturn, err := util.ObjectToJson(result)

	log.Debug().Msg(msgReturn)

	return msgReturn, nil
}

func SendHealthMessage(message string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *config.AwsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := sns.New(sess)

	_, err := client.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: &config.ConfigObj.AWS.SNS.HealthTopicArn,
	})

	return err
}

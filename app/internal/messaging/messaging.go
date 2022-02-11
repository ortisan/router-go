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

var awsConfig = *aws.NewConfig()

func Config() {
	awsConfig.WithRegion(config.ConfigObj.AWS.Region)
	awsConfig.WithEndpoint(config.ConfigObj.AWS.EndpointUrl)
}

func getMessages(sess *session.Session, queueURL *string, timeout *int64) (*sqs.ReceiveMessageOutput, error) {
	// Create an SQS service client
	svc := sqs.New(sess)

	// snippet-start:[sqs.go.receive_messages.call]
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
	// snippet-end:[sqs.go.receive_messages.call]
	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

func GetHealthMessage() (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}))

	var timeoutVisibility int64
	timeoutVisibility = 12 * 60 * 60

	msgResult, err := getMessages(sess, &config.ConfigObj.AWS.SQS.HealthQueueUrl, &timeoutVisibility)

	if err != nil {
		return "", errApp.NewIntegrationError("Error to get sqs messages.", err)
	}
	msgReturn := string(util.ObjectToJson(msgResult))

	log.Debug().Msg(msgReturn)

	return msgReturn, nil
}

func SendHealthMessage(message string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sns.New(sess)

	_, err := svc.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: &config.ConfigObj.AWS.SNS.HealthTopicArn,
	})

	return err
}

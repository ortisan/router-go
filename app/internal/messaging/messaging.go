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
	awsConfig               = *aws.NewConfig()
	timeoutVisibility int64 = 12 * 60
)

func Setup() {
	// Config region
	awsConfig.WithRegion(config.ConfigObj.AWS.Region)
	// Config endpoint url (local and docker env)
	if len(config.ConfigObj.AWS.EndpointUrl) > 0 {
		awsConfig.WithEndpoint(config.ConfigObj.AWS.EndpointUrl)
	}
}

func getMessages(sess *session.Session, queueURL *string, timeout *int64) (*sqs.ReceiveMessageOutput, error) {
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

func GetHealthMessage() (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}))

	result, err := getMessages(sess, &config.ConfigObj.AWS.SQS.HealthQueueUrl, &timeoutVisibility)

	if err != nil {
		return "", errApp.NewIntegrationError("Error to get sqs messages.", err)
	}
	msgReturn := string(util.ObjectToJson(result))

	log.Debug().Msg(msgReturn)

	return msgReturn, nil
}

func SendHealthMessage(message string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := sns.New(sess)

	_, err := client.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: &config.ConfigObj.AWS.SNS.HealthTopicArn,
	})

	return err
}

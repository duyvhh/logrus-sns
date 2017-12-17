package logrus_sns

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

// SNSHook is a logrus Hook for dispatching messages to the specified topics on AWS SNS
type SNSHook struct {
	// Messages with a log level not contained in this array
	// will not be dispatched. If nil, all messages will be dispatched.
	AcceptedLevels []logrus.Level
	Session        *sns.SNS
	QueueUrl       *string
	Extra          map[string]interface{}
}

// Levels define the level of logs which will be sent to SNS
func (sh *SNSHook) Levels() []logrus.Level {
	if sh.AcceptedLevels == nil {
		return logrus.AllLevels
	}

	return sh.AcceptedLevels
}

func NewSNSHook(topic string) (*SNSHook, error) {
	// Creates a SNS hook with a standard AWS session configured by AWS Environment Variables
	s := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	hook, err := NewSNSHookWithSession(topic, s)

	if err != nil {
		return nil, err
	}

	return hook, nil
}

func NewSNSHookWithSession(topic string, sess *session.Session) (*SNSHook, error) {
	// Creates a SNS hook with a custom AWS session
	hook := &SNSHook{}

	hook.Session = sqs.New(sess)

	resultURL, err := hook.Session.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return nil, err
	}

	hook.QueueUrl = resultURL.QueueUrl

	return hook, nil
}

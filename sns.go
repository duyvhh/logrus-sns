package logrus_sns

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/sirupsen/logrus"
)

// SNSHook is a logrus Hook for dispatching messages to the specified topics on AWS SNS
type SNSHook struct {
	// Messages with a log level not contained in this array
	// will not be dispatched. If nil, all messages will be dispatched.
	AcceptedLevels []logrus.Level
	Session        *sns.SNS
	TopicArn       *string
	Subject        *string
	Extra          map[string]interface{}
}

// Levels define the level of logs which will be sent to SNS
func (hook *SNSHook) Levels() []logrus.Level {
	if hook.AcceptedLevels == nil {
		return logrus.AllLevels
	}

	return hook.AcceptedLevels
}

func NewSNSHook(topicArn, subject, region string) (*SNSHook, error) {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	hook, err := NewSNSHookWithSession(topicArn, subject, s)

	if err != nil {
		return nil, err
	}

	return hook, nil
}

func NewSNSHookWithSession(topicArn, subject string, s *session.Session) (*SNSHook, error) {
	// Creates a SNS hook with a custom AWS session
	hook := &SNSHook{}

	hook.Session = sns.New(s)

	_, err := hook.Session.GetTopicAttributes(&sns.GetTopicAttributesInput{
		TopicArn: aws.String(topicArn),
	})

	if err != nil {
		return nil, err
	}

	hook.TopicArn = aws.String(topicArn)
	hook.Subject = aws.String(subject)

	return hook, nil
}

func (hook *SNSHook) Fire(entry *logrus.Entry) error {
	publishInput := sns.PublishInput{}

	publishInput.TopicArn = hook.TopicArn
	publishInput.Message = &entry.Message
	publishInput.Subject = hook.Subject

	data, err := json.Marshal(&entry.Data)

	if err != nil {
		return fmt.Errorf("failed to serialize log data into JSON: %s", err.Error())
	}

	publishInput.MessageAttributes = map[string]*sns.MessageAttributeValue{
		"Level": {
			DataType:    aws.String("String"),
			StringValue: aws.String(entry.Level.String()),
		},
		"Time": {
			DataType:    aws.String("String"),
			StringValue: aws.String(entry.Time.String()),
		},
		"Data": {
			DataType:    aws.String("String"),
			StringValue: aws.String(string(data)),
		},
	}
	_, err = hook.Session.Publish(&publishInput)

	if err != nil {
		return err
	}

	return nil
}

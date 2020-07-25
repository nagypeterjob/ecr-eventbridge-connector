package exporter

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/nagypeterjob/ecr-eventbridge-connector/internal/formatter"
	"github.com/nagypeterjob/ecr-scan-lambda/pkg/api"
)

// SNSExporter publishes message to SNS topic as json
type SNSExporter struct {
	client   *api.SNSService
	name     string
	topicARN string
}

// NewSNSExporter .
func NewSNSExporter(name string, client *api.SNSService, topicARN string) *SNSExporter {
	return &SNSExporter{
		client:   client,
		name:     name,
		topicARN: topicARN,
	}
}

// Name .
func (s SNSExporter) Name() string {
	return s.name
}

func (s SNSExporter) Send(msg *formatter.Message) error {

	input := sns.PublishInput{
		Message:  &msg.Body,
		TopicArn: &s.topicARN,
	}

	if _, err := s.client.Publish(&input); err != nil {
		return err
	}

	return nil
}

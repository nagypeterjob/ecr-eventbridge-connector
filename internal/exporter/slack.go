package exporter

import (
	"fmt"
	"time"

	"github.com/nagypeterjob/ecr-eventbridge-connector/internal/formatter"
	"github.com/nlopes/slack"
)

// SlackService data structure for storing slack client related data
type SlackService struct {
	client  *slack.Client
	channel string
	name    string
}

// NewSlackExporter populates a new SlackService instance
func NewSlackExporter(name string, token string, channel string) *SlackService {
	return &SlackService{
		client:  slack.New(token),
		channel: channel,
		name:    name,
	}
}

// Name .
func (s SlackService) Name() string {
	return s.name
}

func (s SlackService) Send(msg *formatter.Message) error {
	blockParts := s.BuildMessageBlock(msg)
	channelID, timestamp, err := s.PostMessage(blockParts...)
	if err != nil {
		return err
	}
	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)

	return nil
}

// BuildMessageBlock constructs severity related message body
func (s *SlackService) BuildMessageBlock(msg *formatter.Message) []slack.Block {
	return []slack.Block{
		s.GenerateTextBlock(msg.Title),
		s.GenerateTextBlock(msg.Body),
		s.GenerateTextBlock(msg.Link),
		slack.NewDividerBlock(),
	}
}

// GenerateTextBlock returns a slack SectionBlock for text input
func (s *SlackService) GenerateTextBlock(text string) slack.Block {
	textBlock := slack.NewTextBlockObject("mrkdwn", text, false, false)
	return slack.NewSectionBlock(textBlock, nil, nil)
}

// PostMessage sends provided slack MessageBlocks to the given slack channel
func (s *SlackService) PostMessage(blocks ...slack.Block) (string, string, error) {
	// Wait one second so posting doesn't exceed Slack's rate limit
	time.Sleep(1 * time.Second)
	return s.client.PostMessage(s.channel, slack.MsgOptionBlocks(blocks...))
}

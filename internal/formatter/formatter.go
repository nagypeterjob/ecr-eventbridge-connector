package formatter

import (
	"fmt"

	"github.com/nagypeterjob/ecr-eventbridge-connector/pkg/eventbridge"
)

// Message is a common message type used by exporters
type Message struct {
	Title          string
	Link           string
	Body           string
	Status         string
	RepositoryName string
}

func formatLink(event eventbridge.ScanEvent) string {
	return fmt.Sprintf("https://console.aws.amazon.com/ecr/repositories/%s/image/%s/scan-results?region=%s", event.Detail.RepositoryName, event.Detail.ImageDigest, event.Region)
}

// Formatter defines a common interface for different formatters
type Formatter interface {
	// Format provides common interface for formatting text and providing it as a Message type
	Format(eventbridge.ScanEvent) (*Message, error)
}

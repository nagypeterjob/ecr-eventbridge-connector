package formatter

import (
	"fmt"

	"github.com/nagypeterjob/ecr-eventbridge-connector/pkg/eventbridge"
)

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

type Formatter interface {
	Format(eventbridge.ScanEvent) (*Message, error)
}

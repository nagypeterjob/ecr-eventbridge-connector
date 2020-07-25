package formatter

import (
	"bytes"
	"fmt"

	"github.com/nagypeterjob/ecr-eventbridge-connector/internal/severity"
	"github.com/nagypeterjob/ecr-eventbridge-connector/pkg/eventbridge"
)

type SlackFormatter struct{}

func bold(message interface{}) string {
	return fmt.Sprintf("*%v*", message)
}

func (sf SlackFormatter) Format(event eventbridge.ScanEvent) (*Message, error) {

	title := fmt.Sprintf("Vulnerabilities found in %s:", bold(event.Detail.RepositoryName))
	link := fmt.Sprintf("View detailed scan results <%s| on ECR console>", formatLink(event))

	var buffer bytes.Buffer
	for _, key := range severity.SeverityList {
		if val, ok := event.Detail.FindingSeverityCounts[key]; ok {
			buffer.WriteString(fmt.Sprintf("%s %s\n", key, bold(*val)))
		}
	}
	body := buffer.String()

	return &Message{
		Title:          title,
		Body:           body,
		Link:           link,
		Status:         event.Detail.ScanStatus,
		RepositoryName: event.Detail.RepositoryName,
	}, nil
}

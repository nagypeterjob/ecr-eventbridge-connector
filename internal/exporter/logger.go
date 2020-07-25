package exporter

import (
	"fmt"

	"github.com/nagypeterjob/ecr-eventbridge-connector/internal/formatter"
)

// LogExporter is a dummy exporter which prints formatted message to stdout.
// For debug and educational purposes.
type LogExporter struct {
	name string
}

// NewLogExporter .
func NewLogExporter(name string) *LogExporter {
	return &LogExporter{
		name: name,
	}
}

// Name .
func (l LogExporter) Name() string {
	return l.name
}

// Format clousure formats scan results and returns a function that sends report on invocation
func (l LogExporter) Send(msg *formatter.Message) error {
	fmt.Println(msg.Title)
	fmt.Println(msg.Body)
	fmt.Println(msg.Link)
	return nil
}

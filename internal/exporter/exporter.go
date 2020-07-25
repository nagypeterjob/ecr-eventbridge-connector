package exporter

import (
	"github.com/nagypeterjob/ecr-eventbridge-connector/internal/formatter"
)

// Exporter defines a common interface for different exporters
type Exporter interface {
	// Formats message types then returns function which sends formatted messages on invocation
	Send(msg *formatter.Message) error
	// Retrun exporter name
	Name() string
}

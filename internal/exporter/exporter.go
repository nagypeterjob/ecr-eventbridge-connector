package exporter

import (
	"github.com/nagypeterjob/ecr-eventbridge-connector/internal/formatter"
)

// Exporter defines a common interface for different exporters
type Exporter interface {
	// Send implements payload transmission for each exporter
	Send(msg *formatter.Message) error
	// Name retruns exporter name
	Name() string
}

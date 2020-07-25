package formatter

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/nagypeterjob/ecr-eventbridge-connector/internal/severity"
	"github.com/nagypeterjob/ecr-eventbridge-connector/pkg/eventbridge"
)

type JsonFormatter struct {
}

type jsonData struct {
	Title      string     `json:"title"`
	Repository repository `json:"repository"`
	Default    string     `json:"default"`
}

type repository struct {
	Name     string         `json:"name"`
	Link     string         `json:"link"`
	Findings []vulnerablity `json:"findings"`
}

type vulnerablity struct {
	Severity string `json:"severity"`
	Count    string `json:"count"`
}

func marshal(data jsonData) ([]byte, error) {
	return json.Marshal(data)
}

func (jf JsonFormatter) Format(event eventbridge.ScanEvent) (*Message, error) {

	js := jsonData{
		Title: fmt.Sprintf("Vulnerabilities found in %s:", event.Detail.RepositoryName),
		Repository: repository{
			Name:     event.Detail.RepositoryName,
			Link:     formatLink(event),
			Findings: []vulnerablity{},
		},
	}

	for _, key := range severity.SeverityList {
		if val, ok := event.Detail.FindingSeverityCounts[key]; ok {
			js.Repository.Findings = append(js.Repository.Findings, vulnerablity{
				Severity: key,
				Count:    strconv.FormatInt(int64(*val), 10),
			})
		}
	}

	msg, err := marshal(js)
	if err != nil {
		return nil, err
	}

	return &Message{
		Body:           string(msg),
		Status:         event.Detail.ScanStatus,
		RepositoryName: event.Detail.RepositoryName,
	}, nil
}

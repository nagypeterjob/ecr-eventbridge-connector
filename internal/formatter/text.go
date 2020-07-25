package formatter

import (
	"bytes"
	"text/template"

	"github.com/nagypeterjob/ecr-eventbridge-connector/pkg/eventbridge"
)

type values struct {
	Name               string
	CountCritical      *int
	CountHigh          *int
	CountMedium        *int
	CountLow           *int
	CountInformational *int
	CountUndefined     *int
	Link               string
}

type BasicFormatter struct{}

const (
	titleTmpl = `Vulnerabilities found in {{ .Name }}:{{printf "%s" "\n"}}`
	bodyTmpl  = `
	{{- if .CountCritical }}     CRITICAL: {{ .CountCritical }}{{printf "%s" "\n"}}{{end}}
	{{- if .CountHigh }}         HIGH: {{ .CountHigh }}{{printf "%s" "\n"}}{{end}}
	{{- if .CountMedium }}       MEDIUM: {{ .CountMedium }}{{printf "%s" "\n"}}{{end}}
	{{- if .CountLow }}          LOW: {{ .CountLow }}{{printf "%s" "\n"}}{{end}}
	{{- if .CountInformational }}INFORMATIONAL: {{ .CountInformational }}{{printf "%s" "\n"}}{{end}}
	{{- if .CountUndefined }}    UNDEFINED: {{ .CountUndefined }}{{end}}`
	linkTmpl = `View detailed scan results on console ({{ .Link }})
--------------------------------------`
)

func execTmpl(data interface{}, raw string) (bytes.Buffer, error) {
	tmpl, err := template.New("text formatter").Parse(raw)
	if err != nil {
		return bytes.Buffer{}, err
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, &data)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return buffer, nil
}

func (bf BasicFormatter) Format(event eventbridge.ScanEvent) (*Message, error) {

	data := values{
		Name:               event.Detail.RepositoryName,
		CountCritical:      event.Detail.FindingSeverityCounts["CRITICAL"],
		CountHigh:          event.Detail.FindingSeverityCounts["HIGH"],
		CountMedium:        event.Detail.FindingSeverityCounts["MEDIUM"],
		CountLow:           event.Detail.FindingSeverityCounts["LOW"],
		CountInformational: event.Detail.FindingSeverityCounts["INFORMATIONAL"],
		CountUndefined:     event.Detail.FindingSeverityCounts["UNDEFINED"],
		Link:               formatLink(event),
	}

	title, err := execTmpl(data, titleTmpl)
	if err != nil {
		return nil, err
	}

	body, err := execTmpl(data, bodyTmpl)
	if err != nil {
		return nil, err
	}

	link, err := execTmpl(data, linkTmpl)
	if err != nil {
		return nil, err
	}

	return &Message{
		Title:          title.String(),
		Body:           body.String(),
		Link:           link.String(),
		Status:         event.Detail.ScanStatus,
		RepositoryName: event.Detail.RepositoryName,
	}, nil
}

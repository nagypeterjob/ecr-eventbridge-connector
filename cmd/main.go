package main

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	exp "github.com/nagypeterjob/ecr-eventbridge-connector/internal/exporter"
	"github.com/nagypeterjob/ecr-eventbridge-connector/internal/formatter"
	"github.com/nagypeterjob/ecr-eventbridge-connector/internal/logger"
	"github.com/nagypeterjob/ecr-eventbridge-connector/internal/severity"
	"github.com/nagypeterjob/ecr-eventbridge-connector/pkg/eventbridge"
	"github.com/nagypeterjob/ecr-scan-lambda/pkg/api"
)

type app struct {
	exporters       []exporterTuple
	logger          *logger.Logger
	minimumSeverity string
}

type exporterTuple struct {
	exporter  exp.Exporter
	formatter formatter.Formatter
}

func initExporters(config config, logger *logger.Logger) ([]exporterTuple, error) {
	var tuple []exporterTuple

	logger.Infof("Exporters enabled: %s", config.exporters)

	enabledExporters := strings.Split(config.exporters, ",")

	for _, e := range enabledExporters {

		if e == "log" {
			logger.Debug("Initializing log exporter...")
			logexp := exp.NewLogExporter(e)
			tuple = append(tuple, exporterTuple{
				exporter:  logexp,
				formatter: formatter.BasicFormatter{},
			})
		}

		if e == "slack" {
			logger.Debug("Initializing slack exporter...")
			slack := exp.NewSlackExporter(e, config.slack.token, config.slack.channel)
			tuple = append(tuple, exporterTuple{
				exporter:  slack,
				formatter: formatter.SlackFormatter{},
			})
		}

		if e == "sns" {
			logger.Debug("Initializing sns exporter...")
			sess, err := session.NewSession(&aws.Config{Region: &config.region})
			if err != nil {
				return nil, err
			}
			client := sns.New(sess)

			service := api.NewSNSService(client)

			sns := exp.NewSNSExporter(e, service, config.sns.topicARN)
			tuple = append(tuple, exporterTuple{
				exporter:  sns,
				formatter: formatter.JSONFormatter{},
			})
		}

		if e == "mailgun" {
			logger.Debug("Initializing Mailgun exporter...")
			mg := exp.NewMailgunExporter(e, config.mailgun.recipients, config.mailgun.from, config.mailgun.apiKey)
			tuple = append(tuple, exporterTuple{
				exporter:  mg,
				formatter: formatter.BasicFormatter{},
			})
		}
	}
	return tuple, nil
}

func (a *app) Handle(request events.APIGatewayProxyRequest) error {
	var event eventbridge.ScanEvent
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return err
	}

	if !severity.HitSeverityThreshold(event, a.minimumSeverity) {
		a.logger.Infof("Severity did not reach user defined threshold %s", a.minimumSeverity)
		return nil
	}

	for _, tuple := range a.exporters {
		msg, err := tuple.formatter.Format(event)
		if err != nil {
			return err
		}

		if err = tuple.exporter.Send(msg); err != nil {
			return err
		}
		a.logger.Infof("%s exporter has sucessfully sent message", tuple.exporter.Name())
	}

	return nil
}

// Handler glues the lambda logic together
func Handler(request events.APIGatewayProxyRequest) error {
	if err := printVersion(); err != nil {
		return err
	}

	config, err := initConfig()
	if err != nil {
		return err
	}

	logger, err := logger.NewLogger(config.logLevel)
	if err != nil {
		return err
	}

	exporters, err := initExporters(config, logger)
	if err != nil {
		return err
	}

	app := app{
		exporters:       exporters,
		logger:          logger,
		minimumSeverity: config.minimumSeverity,
	}
	return app.Handle(request)
}

func main() {
	lambda.Start(Handler)
}

package main

import (
	"fmt"
	"os"
)

// Config stores lambda configuration
type config struct {
	exporters       string
	logLevel        string
	minimumSeverity string
	region          string

	slack   slackConfig
	sns     snsConfig
	mailgun mailgunConfig
}

type slackConfig struct {
	token   string
	channel string
}

type snsConfig struct {
	topicARN string
}

type mailgunConfig struct {
	apiKey     string
	from       string
	recipients string
}

func retrive(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func required(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return value, fmt.Errorf("Required environtment variable %s is not set", key)
	}
	return value, nil
}

// initConfig populates lambda configuration
func initConfig() (config, error) {

	return config{
		exporters:       retrive("EXPORTERS", "log"),
		logLevel:        retrive("LOG_LEVEL", "INFO"),
		minimumSeverity: retrive("MINIMUM_SEVERITY", "CRITICAL"),
		region:          retrive("REGION", "us-east-1"),
		mailgun: mailgunConfig{
			apiKey:     retrive("MAILGUN_API_KEY", ""),
			from:       retrive("MAILGUN_FROM", ""),
			recipients: retrive("MAILGUN_RECIPIENTS", ""),
		},
		slack: slackConfig{
			token:   retrive("SLACK_TOKEN", ""),
			channel: retrive("SLACK_CHANNEL", ""),
		},

		sns: snsConfig{
			topicARN: retrive("SNS_TOPIC_ARN", ""),
		},
	}, nil
}

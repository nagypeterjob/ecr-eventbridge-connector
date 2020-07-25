package formatter

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/nagypeterjob/ecr-eventbridge-connector/pkg/eventbridge"
)

func TestSlackFormatter(t *testing.T) {

	sf := SlackFormatter{}

	tm, err := time.Parse("2006-01-02T15:04:05Z", "2019-10-29T02:36:48Z")
	if err != nil {
		t.Fatal("Error parsing time")
	}

	event := eventbridge.ScanEvent{
		Version:      "0",
		ID:           "85fc3613-e913-7fc4-a80c-a3753e4aa9ae",
		DetailedType: "ECR Image Scan",
		Source:       "aws.ecr",
		Account:      "123456789012",
		Time:         tm,
		Region:       "us-east-1",
		Resources: []string{
			"arn:aws:ecr:us-east-1:123456789012:repository/my-repo",
		},
		Detail: eventbridge.ScanDetail{
			ScanStatus:     "COMPLETE",
			RepositoryName: "my-repo",
			FindingSeverityCounts: map[string]*int{
				"CRITICAL": aws.Int(10),
				"MEDIUM":   aws.Int(9),
			},
			ImageDigest: "sha256:7f5b2640fe6fb4f46592dfd3410c4a79dac4f89e4782432e0378abcd1234",
			ImageTags:   []string{},
		},
	}

	expectedMsg := Message{
		Title: "Vulnerabilities found in *my-repo*:",
		Body: `CRITICAL *10*
MEDIUM *9*
`,
		Link:           "View detailed scan results on console <" + formatLink(event) + "| on ECR console>",
		Status:         "COMPLETE",
		RepositoryName: "my-repo",
	}

	msg, err := sf.Format(event)
	if err != nil {
		t.Fatal("Failed to exec template")
	}

	if !reflect.DeepEqual(*msg, expectedMsg) {
		t.Fatalf("values not equal: wanted => %s, got => %s", expectedMsg, *msg)
	}
}

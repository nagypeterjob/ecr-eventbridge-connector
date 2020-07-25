package severity

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/nagypeterjob/ecr-eventbridge-connector/pkg/eventbridge"
)

func TestCalculateScore(t *testing.T) {
	cases := []struct {
		Severity Matrix
		Expected int
	}{
		{
			Severity: Matrix{
				Count: map[string]*int{
					"CRITICAL": aws.Int(1),
					"HIGH":     aws.Int(2),
				},
			},
			Expected: 150,
		},
		{
			Severity: Matrix{
				Count: map[string]*int{
					"CRITICAL": aws.Int(1),
					"HIGH":     aws.Int(2),
					"MEDIUM":   aws.Int(3),
				},
			},
			Expected: 170,
		},
		{
			Severity: Matrix{
				Count: map[string]*int{
					"CRITICAL":  aws.Int(1),
					"LOW":       aws.Int(2),
					"UNDEFINED": aws.Int(3),
				},
			},
			Expected: 111,
		},
	}

	for i, c := range cases {
		score := c.Severity.CalculateScore()
		if score != c.Expected {
			t.Fatalf("[%d], values not equal, wanting: %d, got: %d", i, c.Expected, score)
		}
	}
}

func TestHitSeverityThreshold(t *testing.T) {
	cases := []struct {
		input    eventbridge.ScanEvent
		expected bool
	}{
		{
			input: eventbridge.ScanEvent{
				Detail: eventbridge.ScanDetail{
					FindingSeverityCounts: map[string]*int{
						"INFORMATIONAL": aws.Int(10),
						"UNDEFINED":     aws.Int(9),
					},
				},
			},
			expected: false,
		},
		{
			input: eventbridge.ScanEvent{
				Detail: eventbridge.ScanDetail{
					FindingSeverityCounts: map[string]*int{
						"HIGH":     aws.Int(10),
						"CRITICAL": aws.Int(9),
					},
				},
			},
			expected: true,
		},
	}
	for i, c := range cases {
		repos := HitSeverityThreshold(c.input, "MEDIUM")
		if !reflect.DeepEqual(repos, c.expected) {
			t.Fatalf("[%d], values not equal, wanting: %v, got: %v", i, c.expected, repos)
		}
	}
}

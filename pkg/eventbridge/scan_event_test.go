package eventbridge

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
)

func TestScanEventParse(t *testing.T) {

	raw := `{
		"version": "0",
		"id": "85fc3613-e913-7fc4-a80c-a3753e4aa9ae",
		"detail-type": "ECR Image Scan",
		"source": "aws.ecr",
		"account": "123456789012",
		"time": "2019-10-29T02:36:48Z",
		"region": "us-east-1",
		"resources": [
			"arn:aws:ecr:us-east-1:123456789012:repository/my-repo"
		],
		"detail": {
			"scan-status": "COMPLETE",
			"repository-name": "my-repo",
			"finding-severity-counts": {
				"CRITICAL": 10,
				"MEDIUM": 9
			},
			"image-digest": "sha256:7f5b2640fe6fb4f46592dfd3410c4a79dac4f89e4782432e0378abcd1234",
			"image-tags": []
		}
	}`

	tm, err := time.Parse("2006-01-02T15:04:05Z", "2019-10-29T02:36:48Z")
	if err != nil {
		t.Fatal("Error parsing time")
	}

	except := ScanEvent{
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
		Detail: ScanDetail{
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
	var event ScanEvent
	err = json.Unmarshal([]byte(raw), &event)
	if err != nil {
		t.Fatalf("Error parsing json %v", err)
	}

	if !reflect.DeepEqual(event, except) {
		t.Fatalf("Error unmarshaling event: wanted => %v, got => %v", except, event)
	}
}

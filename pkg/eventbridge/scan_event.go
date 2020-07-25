package eventbridge

import (
	"time"
)

type ScanEvent struct {
	Version      string     `json:"version"`
	ID           string     `json:"id"`
	DetailedType string     `json:"detail-type"`
	Source       string     `json:"source"`
	Account      string     `json:"account"`
	Time         time.Time  `json:"time"`
	Region       string     `json:"region"`
	Resources    []string   `json:"resources"`
	Detail       ScanDetail `json:"detail"`
}

type ScanDetail struct {
	ScanStatus            string          `json:"scan-status"`
	RepositoryName        string          `json:"repository-name"`
	FindingSeverityCounts map[string]*int `json:"finding-severity-counts"`
	ImageDigest           string          `json:"image-digest"`
	ImageTags             []string        `json:"image-tags"`
}
